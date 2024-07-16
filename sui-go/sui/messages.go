package sui

import "github.com/iotaledger/wasp/sui-go/sui/serialization"

type TransactionData struct {
	V1 *TransactionDataV1
}

func (t TransactionData) IsBcsEnum() {}

type TransactionDataV1 struct {
	Kind       TransactionKind
	Sender     Address
	GasData    GasData
	Expiration TransactionExpiration
}

type TransactionExpiration struct {
	None  *serialization.EmptyEnum
	Epoch *EpochId
}

func (t TransactionExpiration) IsBcsEnum() {}

type GasData struct {
	Payment []*ObjectRef
	Owner   *Address
	Price   uint64
	Budget  uint64
}

type TransactionKind struct {
	ProgrammableTransaction *ProgrammableTransaction
	ChangeEpoch             *ChangeEpoch
	Genesis                 *GenesisTransaction
	ConsensusCommitPrologue *ConsensusCommitPrologue
}

func (t TransactionKind) IsBcsEnum() {}

type ConsensusCommitPrologue struct {
	Epoch             uint64
	Round             uint64
	CommitTimestampMs CheckpointTimestamp
}

type ProgrammableTransaction struct {
	Inputs   []CallArg
	Commands []Command
}

type TransferObject struct {
	Recipient Address
	ObjectRef ObjectRef
}

type MoveModulePublish struct {
	Modules [][]byte
}

type MoveCall struct {
	Package       ObjectID
	Module        string
	Function      string
	TypeArguments []*TypeTag
	Arguments     []*CallArg
}

type TransferSui struct {
	Recipient Address
	Amount    *uint64 `bcs:"optional"`
}

type Pay struct {
	Coins      []*ObjectRef
	Recipients []*Address
	Amounts    []*uint64
}

type PaySui = Pay

type PayAllSui struct {
	Coins     []*ObjectRef
	Recipient Address
}

type ChangeEpoch struct {
	Epoch                   EpochId
	ProtocolVersion         ProtocolVersion
	StorageCharge           uint64
	ComputationCharge       uint64
	StorageRebate           uint64
	NonRefundableStorageFee uint64
	EpochStartTimestampMs   uint64
	SystemPackages          []*struct {
		SequenceNumber SequenceNumber
		Bytes          [][]uint8
		Objects        []*ObjectID
	}
}

type GenesisTransaction struct {
	Objects []GenesisObject
}

type GenesisObject struct {
	RawObject *struct {
		Data  Data
		Owner Owner
	}
}

type CallArg struct {
	Pure   *[]byte
	Object *ObjectArg
}

func (c CallArg) IsBcsEnum() {}

type ObjectArg struct {
	ImmOrOwnedObject *ObjectRef
	SharedObject     *SharedObjectArg
	Receiving        *ObjectRef
}

type SharedObjectArg struct {
	Id                   *ObjectID
	InitialSharedVersion SequenceNumber
	Mutable              bool
}

func (o ObjectArg) IsBcsEnum() {}

func (o ObjectArg) id() *ObjectID {
	switch {
	case o.ImmOrOwnedObject != nil:
		return o.ImmOrOwnedObject.ObjectID
	case o.SharedObject != nil:
		return o.SharedObject.Id
	default:
		return &ObjectID{}
	}
}
