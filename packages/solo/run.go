// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package solo

import (
	"errors"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/vm/vmimpl"
)

func (ch *Chain) RunOffLedgerRequest(r isc.Request) (
	*iotajsonrpc.IotaTransactionBlockResponse,
	isc.CallArguments,
	error,
) {
	defer ch.logRequestLastBlock()
	ptbRes, results := ch.RunRequestsSync([]isc.Request{r})
	if len(results) == 0 {
		return nil, nil, errors.New("request was skipped")
	}
	res := results[0]
	return ptbRes, res.Return, ch.ResolveVMError(res.Receipt.Error).AsGoError()
}

func (ch *Chain) RunOffLedgerRequests(reqs []isc.Request) (
	*iotajsonrpc.IotaTransactionBlockResponse,
	[]*vm.RequestResult,
) {
	defer ch.logRequestLastBlock()
	return ch.RunRequestsSync(reqs)
}

func (ch *Chain) RunRequestsSync(reqs []isc.Request) (
	*iotajsonrpc.IotaTransactionBlockResponse,
	[]*vm.RequestResult,
) {
	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()
	return ch.runRequestsNolock(reqs)
}

func (ch *Chain) estimateGas(req isc.Request) (result *vm.RequestResult) {
	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()

	res := ch.runTaskNoLock([]isc.Request{req}, true)
	require.Len(ch.Env.T, res.RequestResults, 1, "cannot estimate gas: request was skipped")
	return res.RequestResults[0]
}

func (ch *Chain) runTaskNoLock(reqs []isc.Request, estimateGas bool) *vm.VMTaskResult {
	task := &vm.VMTask{
		Processors:         ch.proc,
		Anchor:             ch.GetLatestAnchor(),
		Requests:           reqs,
		Timestamp:          ch.Env.GlobalTime(),
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

func (ch *Chain) runRequestsNolock(reqs []isc.Request) (
	*iotajsonrpc.IotaTransactionBlockResponse,
	[]*vm.RequestResult,
) {
	res := ch.runTaskNoLock(reqs, false)
	gasPayment := ch.Env.selectCoinsForGas(
		ch.OriginatorAddress,
		&res.UnsignedTransaction,
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	ptbRes := ch.Env.executePTB(
		res.UnsignedTransaction,
		ch.OriginatorPrivateKey,
		gasPayment,
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	if res.RotationAddress == nil {
		// normal state transition
		ch.settleStateTransition(res.StateDraft)
	}
	return ptbRes, res.RequestResults
}

func (ch *Chain) settleStateTransition(stateDraft state.StateDraft) {
	block := ch.store.Commit(stateDraft)
	err := ch.store.SetLatest(block.TrieRoot())
	if err != nil {
		panic(err)
	}

	latestState := lo.Must(ch.LatestState())

	ch.Env.Publisher().BlockApplied(ch.ChainID, block, latestState)

	blockReceipts, err := blocklog.RequestReceiptsFromBlock(block)
	if err != nil {
		panic(err)
	}
	for _, rec := range blockReceipts {
		ch.mempool.RemoveRequest(rec.Request.ID())
	}
	ch.Log().Infof("state transition --> #%d. Requests in the block: %d",
		stateDraft.BlockIndex(), len(blockReceipts))
}

func (ch *Chain) logRequestLastBlock() {
	recs := ch.GetRequestReceiptsForBlock(ch.GetLatestBlockInfo().BlockIndex)
	for _, rec := range recs {
		ch.Log().Infof("REQ: '%s'", rec.Short())
	}
}
