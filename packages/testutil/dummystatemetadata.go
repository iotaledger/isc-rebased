package testutil

import (
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

func DummyStateMetadata(commitment *state.L1Commitment) *transaction.StateMetadata {
	return transaction.NewStateMetadata(
		0,
		commitment,
		gas.DefaultFeePolicy(),
		[][]byte{},
		"",
	)
}
