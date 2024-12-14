package common

import (
	"fmt"
	"strconv"
	"strings"

	chainpkg "github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/chainutil"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/trie"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
)

func ParseReceipt(chain chainpkg.Chain, receipt *blocklog.RequestReceipt) (*isc.Receipt, error) {
	// Using latest state instead state of request.BlockIndex to avoid
	// possibility of referencing a state that has been pruned.
	state, err := chain.LatestState(chainpkg.ActiveOrCommittedState)
	if err != nil {
		return nil, err
	}

	resolvedReceiptErr, err := chainutil.ResolveError(state, receipt.Error)
	if err != nil {
		return nil, err
	}

	iscReceipt := receipt.ToISCReceipt(resolvedReceiptErr)

	return iscReceipt, nil
}

func CallView(ch chainpkg.Chain, msg isc.Message, blockIndexOrHash string) (isc.CallArguments, error) {
	var chainState state.State
	var err error

	switch {
	case blockIndexOrHash == "":
		chainState, err = ch.LatestState(chainpkg.ActiveOrCommittedState)
		if err != nil {
			return nil, fmt.Errorf("error getting latest chain state: %w", err)
		}

	case strings.HasPrefix(blockIndexOrHash, "0x"):
		hashBytes, err := cryptolib.DecodeHex(blockIndexOrHash)
		if err != nil {
			return nil, fmt.Errorf("invalid block hash: %v", blockIndexOrHash)
		}

		trieRoot, err := trie.HashFromBytes(hashBytes)
		if err != nil {
			return nil, fmt.Errorf("invalid block hash: %v", blockIndexOrHash)
		}

		chainState, err = ch.Store().StateByTrieRoot(trieRoot)
		if err != nil {
			return nil, fmt.Errorf("error getting block by trie root: %w", err)
		}

	default:
		blockIndex, err := strconv.ParseUint(blockIndexOrHash, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid block number: %v", blockIndexOrHash)
		}

		chainState, err = ch.Store().StateByIndex(uint32(blockIndex))
		if err != nil {
			return nil, fmt.Errorf("error getting block by index: %w", err)
		}
	}

	return chainutil.CallView(ch.ID(), chainState, ch.Processors(), ch.Log(), msg)
}

func EstimateGas(ch chainpkg.Chain, gasCoin *coin.CoinWithRef, req isc.Request) (*isc.Receipt, error) {
	anchor, err := ch.LatestAnchor(chainpkg.ActiveOrCommittedState)
	if err != nil {
		return nil, fmt.Errorf("error getting latest anchor: %w", err)
	}

	rec, err := chainutil.SimulateRequest(anchor, gasCoin, ch.Store(), ch.Processors(), ch.Log(), req, true)
	if err != nil {
		return nil, err
	}
	parsedRec, err := ParseReceipt(ch, rec)
	return parsedRec, err
}
