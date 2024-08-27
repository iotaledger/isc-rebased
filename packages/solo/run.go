// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package solo

import (
	"errors"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/vm/vmimpl"
)

func (ch *Chain) RunOffLedgerRequest(r isc.Request) (isc.CallArguments, error) {
	defer ch.logRequestLastBlock()
	results := ch.RunRequestsSync([]isc.Request{r}, "off-ledger")
	if len(results) == 0 {
		return nil, errors.New("request was skipped")
	}
	res := results[0]
	return res.Return, ch.ResolveVMError(res.Receipt.Error).AsGoError()
}

func (ch *Chain) RunOffLedgerRequests(reqs []isc.Request) []*vm.RequestResult {
	defer ch.logRequestLastBlock()
	return ch.RunRequestsSync(reqs, "off-ledger")
}

func (ch *Chain) RunRequestsSync(reqs []isc.Request, trace string) (results []*vm.RequestResult) {
	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()
	return ch.runRequestsNolock(reqs, trace)
}

func (ch *Chain) estimateGas(req isc.Request) (result *vm.RequestResult) {
	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()

	res := ch.runTaskNoLock([]isc.Request{req}, true)
	require.Len(ch.Env.T, res.RequestResults, 1, "cannot estimate gas: request was skipped")
	return res.RequestResults[0]
}

func (ch *Chain) runTaskNoLock(reqs []isc.Request, estimateGas bool) *vm.VMTaskResult {
	anchorOutput := ch.GetAnchorOutputFromL1()
	task := &vm.VMTask{
		Processors:         ch.proc,
		AnchorOutput:       anchorOutput.GetAliasOutput(),
		AnchorOutputID:     anchorOutput.OutputID(),
		Requests:           reqs,
		TimeAssumption:     ch.Env.GlobalTime(),
		Store:              ch.store,
		Entropy:            hashing.PseudoRandomHash(nil),
		ValidatorFeeTarget: ch.ValidatorFeeTarget,
		Log:                ch.Log().Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
		// state baseline is always valid in Solo
		EnableGasBurnLogging: ch.Env.enableGasBurnLogging,
		EstimateGasMode:      estimateGas,
		Migrations:           allmigrations.DefaultScheme,
	}

	res, err := vmimpl.Run(task)
	require.NoError(ch.Env.T, err)
	err = accounts.NewStateReaderFromChainState(res.StateDraft.SchemaVersion(), res.StateDraft).
		CheckLedgerConsistency()
	require.NoError(ch.Env.T, err)
	return res
}

func (ch *Chain) runRequestsNolock(reqs []isc.Request, trace string) (results []*vm.RequestResult) {
	ch.Log().Debugf("runRequestsNolock ('%s')", trace)

	res := ch.runTaskNoLock(reqs, false)

	var essence *iotago.TransactionEssence
	if res.RotationAddress == nil {
		essence = res.TransactionEssence
		copy(essence.InputsCommitment[:], res.InputsCommitment)
	} else {
		var err error = errors.New("refactor me: runRequestsNolock")
		panic("refactor me: rotate.MakeRotateStateControllerTransaction")
		require.NoError(ch.Env.T, err)
	}
	sig, err := ch.StateControllerKeyPair.Sign(essence.InputsCommitment[:])
	require.NoError(ch.Env.T, err)

	tx := transaction.MakeAnchorTransaction(essence, sig)

	if res.RotationAddress == nil {
		// normal state transition
		ch.settleStateTransition(tx, res.StateDraft)
	}

	err = ch.Env.AddToLedger(tx)
	require.NoError(ch.Env.T, err)

	anchor, _, err := transaction.GetAnchorFromTransaction(tx)
	require.NoError(ch.Env.T, err)

	if res.RotationAddress != nil {
		ch.Log().Infof("ROTATED STATE CONTROLLER to %s", anchor.StateController)
	}

	rootC := ch.GetRootCommitment()
	l1C := ch.GetL1Commitment()
	require.Equal(ch.Env.T, rootC, l1C.TrieRoot())

	ch.Env.EnqueueRequests(tx)

	return res.RequestResults
}

func (ch *Chain) settleStateTransition(stateTx *iotago.Transaction, stateDraft state.StateDraft) {
	block := ch.store.Commit(stateDraft)
	err := ch.store.SetLatest(block.TrieRoot())
	if err != nil {
		panic(err)
	}

	latestState, _ := ch.LatestState(chain.ActiveOrCommittedState)

	ch.Env.Publisher().BlockApplied(ch.ChainID, block, latestState)

	blockReceipts, err := blocklog.RequestReceiptsFromBlock(block)
	if err != nil {
		panic(err)
	}
	for _, rec := range blockReceipts {
		ch.mempool.RemoveRequest(rec.Request.ID())
	}
	unprocessableRequests, err := blocklog.UnprocessableRequestsAddedInBlock(block)
	if err != nil {
		panic(err)
	}
	for _, req := range unprocessableRequests {
		ch.mempool.RemoveRequest(req.ID())
	}
	ch.Log().Infof("state transition --> #%d. Requests in the block: %d. Outputs: %d",
		stateDraft.BlockIndex(), len(blockReceipts), len(stateTx.Essence.Outputs))
}

func (ch *Chain) logRequestLastBlock() {
	recs := ch.GetRequestReceiptsForBlock(ch.GetLatestBlockInfo().BlockIndex())
	for _, rec := range recs {
		ch.Log().Infof("REQ: '%s'", rec.Short())
	}
}
