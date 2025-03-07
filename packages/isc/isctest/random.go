package isctest

import (
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotago/iotatest"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/clients/iscmove/iscmovetest"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state/statetest"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

type RandomAnchorOption struct {
	ID            *iotago.ObjectID
	Assets        *iscmove.AssetsBag
	StateMetadata *transaction.StateMetadata
	StateIndex    *uint32
	ObjectRef     *iotago.ObjectRef
	Owner         *iotago.Address
}

func RandomStateAnchor(opts ...RandomAnchorOption) isc.StateAnchor {
	var anchorOpts iscmovetest.RandomAnchorOption
	anchorRef := iotatest.RandomObjectRef()
	owner := iotatest.RandomAddress()
	if len(opts) == 1 {
		var stateMetadata []byte
		if opts[0].StateMetadata != nil {
			stateMetadata = opts[0].StateMetadata.Bytes()
		}
		anchorOpts = iscmovetest.RandomAnchorOption{
			ID:            opts[0].ID,
			Assets:        opts[0].Assets,
			StateMetadata: &stateMetadata,
			StateIndex:    opts[0].StateIndex,
		}
		if opts[0].Owner != nil {
			owner = iotatest.RandomAddress()
		}
		if opts[0].ObjectRef != nil {
			anchorRef = opts[0].ObjectRef
		}
		if opts[0].ID != nil {
			anchorRef.ObjectID = opts[0].ID
		}
	}
	anchor := iscmovetest.RandomAnchor(anchorOpts)

	anchorRefWithObject := iscmove.RefWithObject[iscmove.Anchor]{
		ObjectRef: *anchorRef,
		Object:    &anchor,
		Owner:     owner,
	}
	return isc.NewStateAnchor(&anchorRefWithObject, *iotatest.RandomAddress())
}

// simualte how StateAnchor is updated (state transition)
// assume the AssetsBag keep unchanged
func UpdateStateAnchor(anchor *isc.StateAnchor, stateMetadata ...[]byte) *isc.StateAnchor {
	// a := anchor.Clone()
	a := iscmove.AnchorWithRef{
		ObjectRef: *anchor.GetObjectRef(),
		Object: &iscmove.Anchor{
			ID:         *anchor.Anchor().ObjectID,
			Assets:     anchor.Anchor().Object.Assets,
			StateIndex: anchor.GetStateIndex() + 1,
		},
	}
	if len(stateMetadata) == 1 {
		a.Object.StateMetadata = stateMetadata[0]
	} else {
		a.Object.StateMetadata = transaction.NewStateMetadata(
			allmigrations.LatestSchemaVersion,
			statetest.NewRandL1Commitment(),
			&iotago.ObjectID{},
			gas.DefaultFeePolicy(),
			isc.NewCallArguments([]byte{1, 2, 3}),
			0,
			"http://url",
		).Bytes()
	}
	newStatAnchor := isc.NewStateAnchor(&a, anchor.ISCPackage())
	return &newStatAnchor
}

// RandomChainID creates a random chain ID. Used for testing only
func RandomChainID(seed ...[]byte) isc.ChainID {
	var h hashing.HashValue
	if len(seed) > 0 {
		h = hashing.HashData(seed[0])
	} else {
		h = hashing.PseudoRandomHash(nil)
	}
	chainID, err := isc.ChainIDFromBytes(h[:isc.ChainIDLength])
	if err != nil {
		panic(err)
	}
	return chainID
}

func RandomRequestWithRef() *iscmove.RefWithObject[iscmove.Request] {
	ref := iotatest.RandomObjectRef()
	a := iscmove.AssetsBagWithBalances{
		AssetsBag: iscmove.AssetsBag{ID: *iotatest.RandomAddress(), Size: 1},
		Balances:  make(iscmove.AssetsBagBalances),
	}
	a.Balances[iotajsonrpc.IotaCoinType] = 1000
	return &iscmove.RefWithObject[iscmove.Request]{
		ObjectRef: *ref,
		Object: &iscmove.Request{
			ID:        *ref.ObjectID,
			Sender:    cryptolib.NewRandomAddress(),
			AssetsBag: a,
			Message: iscmove.Message{
				Contract: 123,
				Function: 456,
				Args:     [][]byte{[]byte("testarg1"), []byte("testarg2")},
			},
			Allowance: iscmove.Assets{Coins: iscmove.CoinBalances{
				iotajsonrpc.IotaCoinType:                                111,
				iotajsonrpc.MustCoinTypeFromString("0x1::coin::TEST_A"): 222,
			}},
			GasBudget: 1000,
		},
	}
}

func RandomOnLedgerRequest() isc.OnLedgerRequest {
	req := RandomRequestWithRef()
	onReq, err := isc.OnLedgerFromRequest(req, cryptolib.NewRandomAddress())
	if err != nil {
		panic(err)
	}
	return onReq
}

func RandomOnLedgerDepositRequest(senders ...*cryptolib.Address) isc.OnLedgerRequest {
	sender := cryptolib.NewRandomAddress()
	if len(senders) != 0 {
		sender = senders[0]
	}
	ref := iotatest.RandomObjectRef()
	a := iscmove.AssetsBagWithBalances{
		AssetsBag: iscmove.AssetsBag{ID: *iotatest.RandomAddress(), Size: 1},
		Balances:  make(iscmove.AssetsBagBalances),
	}
	a.Balances[iotajsonrpc.IotaCoinType] = 1000
	req := iscmove.RefWithObject[iscmove.Request]{
		ObjectRef: *ref,
		Object: &iscmove.Request{
			ID:        *ref.ObjectID,
			Sender:    sender,
			AssetsBag: a,
			Message: iscmove.Message{
				Contract: uint32(isc.Hn("accounts")),
				Function: uint32(isc.Hn("deposit")),
			},
			Allowance: iscmove.Assets{Coins: iscmove.CoinBalances{iotajsonrpc.IotaCoinType: 10000}},
			GasBudget: 100000,
		},
	}
	onReq, err := isc.OnLedgerFromRequest(&req, sender)
	if err != nil {
		panic(err)
	}
	return onReq
}

// NewRandomAgentID creates random AgentID
func NewRandomAgentID() isc.AgentID {
	return isc.NewContractAgentID(RandomChainID(), isc.Hn("testName"))
}
