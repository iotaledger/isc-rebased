package iotago

import (
	"github.com/iotaledger/wasp/clients/iota-go/iotago/serialization"
)

var (
	SuiSystemMut = CallArg{
		Object: &SuiSystemMutObj,
	}

	SuiSystemMutObj = ObjectArg{
		SharedObject: &SharedObjectArg{
			Id:                   SuiObjectIdSystemState,
			InitialSharedVersion: SuiSystemStateObjectSharedVersion,
			Mutable:              true,
		},
	}
)

func NewProgrammable(
	sender *Address,
	pt ProgrammableTransaction,
	gasPayment []*ObjectRef,
	gasBudget uint64,
	gasPrice uint64,
) TransactionData {
	return NewProgrammableAllowSponsor(*sender, pt, gasPayment, gasBudget, gasPrice, sender)
}

func NewProgrammableAllowSponsor(
	sender Address,
	pt ProgrammableTransaction,
	gasPayment []*ObjectRef,
	gasBudget uint64,
	gasPrice uint64,
	sponsor *Address,
) TransactionData {
	kind := TransactionKind{
		ProgrammableTransaction: &pt,
	}
	return TransactionData{
		V1: &TransactionDataV1{
			Kind:   kind,
			Sender: sender,
			GasData: GasData{
				Payment: gasPayment,
				Owner:   sponsor,
				Price:   gasPrice,
				Budget:  gasBudget,
			},
			Expiration: TransactionExpiration{
				None: &serialization.EmptyEnum{},
			},
		},
	}
}