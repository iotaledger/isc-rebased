// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package jsonrpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"

	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/chainutil"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/trie"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

// WaspEVMBackend is the implementation of [ChainBackend] for the production environment.
type WaspEVMBackend struct {
	chain      chain.Chain
	nodePubKey *cryptolib.PublicKey
	baseToken  *parameters.BaseToken
}

var _ ChainBackend = &WaspEVMBackend{}

func NewWaspEVMBackend(ch chain.Chain, nodePubKey *cryptolib.PublicKey, baseToken *parameters.BaseToken) *WaspEVMBackend {
	return &WaspEVMBackend{
		chain:      ch,
		nodePubKey: nodePubKey,
		baseToken:  baseToken,
	}
}

func (b *WaspEVMBackend) FeePolicy(blockIndex uint32) (*gas.FeePolicy, error) {
	state, err := b.ISCStateByBlockIndex(blockIndex)
	if err != nil {
		return nil, err
	}
	ret, err := b.ISCCallView(state, governance.ViewGetFeePolicy.Message())
	if err != nil {
		return nil, err
	}
	return governance.ViewGetFeePolicy.Output.Decode(ret)
}

func (b *WaspEVMBackend) EVMSendTransaction(tx *types.Transaction) error {
	// Ensure the transaction has more gas than the basic Ethereum tx fee.
	intrinsicGas, err := core.IntrinsicGas(tx.Data(), tx.AccessList(), tx.To() == nil, true, true, true)
	if err != nil {
		return err
	}
	if tx.Gas() < intrinsicGas {
		return core.ErrIntrinsicGas
	}

	req, err := isc.NewEVMOffLedgerTxRequest(b.chain.ID(), tx)
	if err != nil {
		return err
	}
	b.chain.Log().Debugf("EVMSendTransaction, evm.tx.nonce=%v, evm.tx.hash=%v => isc.req.id=%v", tx.Nonce(), tx.Hash().Hex(), req.ID())
	if err := b.chain.ReceiveOffLedgerRequest(req, b.nodePubKey); err != nil {
		return fmt.Errorf("tx not added to the mempool: %v", err.Error())
	}

	return nil
}

func (b *WaspEVMBackend) EVMCall(aliasOutput *isc.AliasOutputWithID, callMsg ethereum.CallMsg) ([]byte, error) {
	return chainutil.EVMCall(b.chain, aliasOutput, callMsg)
}

func (b *WaspEVMBackend) EVMEstimateGas(aliasOutput *isc.AliasOutputWithID, callMsg ethereum.CallMsg) (uint64, error) {
	return chainutil.EVMEstimateGas(b.chain, aliasOutput, callMsg)
}

func (b *WaspEVMBackend) EVMTraceTransaction(
	aliasOutput *isc.AliasOutputWithID,
	blockTime time.Time,
	iscRequestsInBlock []isc.Request,
	txIndex uint64,
	tracer *tracers.Tracer,
) error {
	return chainutil.EVMTraceTransaction(
		b.chain,
		aliasOutput,
		blockTime,
		iscRequestsInBlock,
		txIndex,
		tracer,
	)
}

func (b *WaspEVMBackend) ISCCallView(chainState state.State, msg isc.Message) (dict.Dict, error) {
	return chainutil.CallView(chainState, b.chain, msg)
}

func (b *WaspEVMBackend) BaseToken() *parameters.BaseToken {
	return b.baseToken
}

func (b *WaspEVMBackend) ISCLatestAliasOutput() (*isc.AliasOutputWithID, error) {
	latestAliasOutput, err := b.chain.LatestAliasOutput(chain.ActiveOrCommittedState)
	if err != nil {
		return nil, fmt.Errorf("could not get latest AliasOutput: %w", err)
	}
	return latestAliasOutput, nil
}

func (b *WaspEVMBackend) ISCLatestState() (state.State, error) {
	latestState, err := b.chain.LatestState(chain.ActiveOrCommittedState)
	if err != nil {
		return nil, fmt.Errorf("couldn't get latest block index: %w", err)
	}
	return latestState, nil
}

func (b *WaspEVMBackend) ISCStateByBlockIndex(blockIndex uint32) (state.State, error) {
	latestState, err := b.chain.LatestState(chain.ActiveOrCommittedState)
	if err != nil {
		return nil, fmt.Errorf("couldn't get latest state: %s", err.Error())
	}
	if latestState.BlockIndex() == blockIndex {
		return latestState, nil
	}
	return b.chain.Store().StateByIndex(blockIndex)
}

func (b *WaspEVMBackend) ISCStateByTrieRoot(trieRoot trie.Hash) (state.State, error) {
	return b.chain.Store().StateByTrieRoot(trieRoot)
}

func (b *WaspEVMBackend) ISCChainID() *isc.ChainID {
	chID := b.chain.ID()
	return &chID
}

var errNotImplemented = errors.New("method not implemented")

func (*WaspEVMBackend) RevertToSnapshot(int) error {
	return errNotImplemented
}

func (*WaspEVMBackend) TakeSnapshot() (int, error) {
	return 0, errNotImplemented
}
