// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package jsonrpc

import (
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/trie"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

// ChainBackend provides access to the underlying ISC chain.
type ChainBackend interface {
	EVMSendTransaction(tx *types.Transaction) error
	EVMCall(anchor *isc.StateAnchor, callMsg ethereum.CallMsg) ([]byte, error)
	EVMEstimateGas(anchor *isc.StateAnchor, callMsg ethereum.CallMsg) (uint64, error)
	EVMTraceTransaction(anchor *isc.StateAnchor, blockTime time.Time, iscRequestsInBlock []isc.Request, txIndex uint64, tracer *tracers.Tracer) error
	FeePolicy(blockIndex uint32) (*gas.FeePolicy, error)
	ISCChainID() *isc.ChainID
	ISCCallView(chainState state.State, msg isc.Message) (isc.CallArguments, error)
	ISCLatestAnchor() (*isc.StateAnchor, error)
	ISCLatestState() (state.State, error)
	ISCStateByBlockIndex(blockIndex uint32) (state.State, error)
	ISCStateByTrieRoot(trieRoot trie.Hash) (state.State, error)
	BaseToken() *parameters.BaseToken
	TakeSnapshot() (int, error)
	RevertToSnapshot(int) error
}
