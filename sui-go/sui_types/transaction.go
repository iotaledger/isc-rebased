package sui_types

import "github.com/iotaledger/wasp/sui-go/sui_types/serialization"

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
	sender *SuiAddress,
	pt ProgrammableTransaction,
	gasPayment []*ObjectRef,
	gasBudget uint64, // TODO set this to bigint
	gasPrice uint64, // TODO set this to bigint
) TransactionData {
	return NewProgrammableAllowSponsor(*sender, pt, gasPayment, gasBudget, gasPrice, *sender)
}

func NewProgrammableAllowSponsor(
	sender SuiAddress,
	pt ProgrammableTransaction,
	gasPayment []*ObjectRef,
	gasBudget uint64, // TODO set this to bigint
	gasPrice uint64, // TODO set this to bigint
	sponsor SuiAddress,
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
	// return newWithGasCoinsAllowSponsor(kind, sender, gasPayment, gasBudge, gasPrice, sponsor)
}

// func newWithGasCoinsAllowSponsor(
// 	kind TransactionKind,
// 	sender SuiAddress,
// 	gasPayment []*ObjectRef,
// 	gasBudget uint64,
// 	gasPrice uint64,
// 	gasSponsor SuiAddress,
// ) TransactionData {
// 	return TransactionData{
// 		V1: &TransactionDataV1{
// 			Kind:   kind,
// 			Sender: sender,
// 			GasData: GasData{
// 				Price:   gasPrice,
// 				Owner:   gasSponsor,
// 				Payment: gasPayment,
// 				Budget:  gasBudget,
// 			},
// 			Expiration: TransactionExpiration{
// 				None: &serialization.EmptyEnum{},
// 			},
// 		},
// 	}
// }
