package isc

import (
	"fmt"

	"github.com/ethereum/go-ethereum"

	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/util/bcs"
	"github.com/iotaledger/wasp/sui-go/sui"
)

type OnLedgerRequestData struct {
	requestRef      sui.ObjectRef      `bcs:"export"`
	senderAddress   *cryptolib.Address `bcs:"export"`
	targetAddress   *cryptolib.Address `bcs:"export"`
	assets          *Assets            `bcs:"export"`
	assetsBag       *iscmove.AssetsBag `bcs:"export"`
	requestMetadata *RequestMetadata   `bcs:"export"`
}

var (
	_ Request         = new(OnLedgerRequestData)
	_ OnLedgerRequest = new(OnLedgerRequestData)
	_ Calldata        = new(OnLedgerRequestData)
)

func OnLedgerFromRequest(request *iscmove.RefWithObject[iscmove.Request], anchorAddress *cryptolib.Address) (OnLedgerRequest, error) {
	r := &OnLedgerRequestData{
		requestRef:    request.ObjectRef,
		senderAddress: request.Object.Sender,
		targetAddress: anchorAddress,
		assetsBag:     &request.Object.AssetsBag.AssetsBag,
		requestMetadata: &RequestMetadata{
			SenderContract: ContractIdentity{},
			Message: Message{
				Target: CallTarget{
					Contract:   Hname(request.Object.Message.Contract),
					EntryPoint: Hname(request.Object.Message.Function),
				},
				Params: nil,
			},
			Allowance: NewEmptyAssets(),
			GasBudget: 0,
		},
		assets: AssetsFromAssetsBagWithBalances(request.Object.AssetsBag),
	}

	return r, nil
}

func (req *OnLedgerRequestData) Allowance() *Assets {
	if req.requestMetadata == nil {
		return NewEmptyAssets()
	}
	return req.requestMetadata.Allowance
}

func (req *OnLedgerRequestData) Assets() *Assets {
	return req.assets
}

func (req *OnLedgerRequestData) Bytes() []byte {
	return bcs.MustMarshal(req)
}

func (req *OnLedgerRequestData) Message() Message {
	if req.requestMetadata == nil {
		return Message{}
	}
	return req.requestMetadata.Message
}

func (req *OnLedgerRequestData) Clone() OnLedgerRequest {
	outputRef := sui.ObjectRefFromBytes(req.requestRef.Bytes())

	ret := &OnLedgerRequestData{
		requestRef:    *outputRef,
		senderAddress: req.senderAddress.Clone(),
		targetAddress: req.targetAddress.Clone(),
	}

	if req.requestMetadata != nil {
		ret.requestMetadata = req.requestMetadata.Clone()
	}

	return ret
}

func (req *OnLedgerRequestData) GasBudget() (gasBudget uint64, isEVM bool) {
	if req.requestMetadata == nil {
		return 0, false
	}
	return req.requestMetadata.GasBudget, false
}

func (req *OnLedgerRequestData) ID() RequestID {
	return RequestID(*req.requestRef.ObjectID)
}

func (req *OnLedgerRequestData) IsOffLedger() bool {
	return false
}

func (req *OnLedgerRequestData) RequestID() sui.ObjectID {
	return *req.requestRef.ObjectID
}

func (req *OnLedgerRequestData) SenderAccount() AgentID {
	sender := req.SenderAddress()
	if sender == nil {
		return nil
	}
	if req.requestMetadata != nil && !req.requestMetadata.SenderContract.Empty() {
		chainID := ChainIDFromAddress(sender)
		return req.requestMetadata.SenderContract.AgentID(chainID)
	}
	return NewAddressAgentID(sender)
}

func (req *OnLedgerRequestData) SenderAddress() *cryptolib.Address {
	return req.senderAddress
}

func (req *OnLedgerRequestData) String() string {
	metadata := req.requestMetadata
	if metadata == nil {
		return "onledger request without metadata"
	}
	return fmt.Sprintf("OnLedgerRequestData::{ ID: %s, sender: %s, target: %s, entrypoint: %s, Params: %s, Assets: %v, GasBudget: %d }",
		req.ID().String(),
		metadata.SenderContract.String(),
		metadata.Message.Target.Contract.String(),
		metadata.Message.Target.EntryPoint.String(),
		metadata.Message.Params.String(),
		req.assets,
		metadata.GasBudget,
	)
}

func (req *OnLedgerRequestData) RequestRef() sui.ObjectRef {
	return req.requestRef
}

func (req *OnLedgerRequestData) AssetsBag() *iscmove.AssetsBag {
	return req.assetsBag
}

func (req *OnLedgerRequestData) TargetAddress() *cryptolib.Address {
	return req.targetAddress
}

func (req *OnLedgerRequestData) EVMCallMsg() *ethereum.CallMsg {
	return nil
}
