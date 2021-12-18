package vmtxbuilder

import (
	"math/big"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/util"
)

func (n *nativeTokenBalance) clone() *nativeTokenBalance {
	return &nativeTokenBalance{
		tokenID:            n.tokenID,
		input:              n.input,
		dustDepositCharged: n.dustDepositCharged,
		in:                 cloneInternalExtendedOutput(n.in),
		out:                cloneInternalExtendedOutput(n.out),
	}
}

// producesOutput if value update produces UTXO of the corresponding total native token balance
func (n *nativeTokenBalance) producesOutput() bool {
	if n.identicalInOut() {
		// value didn't change
		return false
	}
	if util.IsZeroBigInt(n.getOutValue()) {
		// end value is 0
		return false
	}
	return true
}

// requiresInput returns if value change requires input in the transaction
func (n *nativeTokenBalance) requiresInput() bool {
	if n.identicalInOut() {
		// value didn't change
		return false
	}
	if n.in == nil {
		// there's no input
		return false
	}
	return true
}

func (n *nativeTokenBalance) getOutValue() *big.Int {
	return n.out.NativeTokens[0].Amount
}

func (n *nativeTokenBalance) setOutValue(v *big.Int) {
	n.out.NativeTokens[0].Amount = v
}

func (n *nativeTokenBalance) identicalInOut() bool {
	switch {
	case n.in == n.out:
		panic("identicalExtendedOutputs: internal inconsistency 1")
	case n.in == nil || n.out == nil:
		return false
	case !n.in.Address.Equal(n.out.Address):
		return false
	case n.in.Amount != n.out.Amount:
		return false
	case !n.in.NativeTokens.Equal(n.out.NativeTokens):
		return false
	case !n.in.Blocks.Equal(n.out.Blocks):
		return false
	case len(n.in.NativeTokens) != 1:
		panic("identicalExtendedOutputs: internal inconsistency 2")
	case len(n.out.NativeTokens) != 1:
		panic("identicalExtendedOutputs: internal inconsistency 3")
	case n.in.NativeTokens[0].ID != n.tokenID:
		panic("identicalExtendedOutputs: internal inconsistency 4")
	case n.out.NativeTokens[0].ID != n.tokenID:
		panic("identicalExtendedOutputs: internal inconsistency 5")
	}
	return true
}

func cloneInternalExtendedOutput(o *iotago.ExtendedOutput) *iotago.ExtendedOutput {
	if o == nil {
		return nil
	}
	return &iotago.ExtendedOutput{
		Address:      o.Address, // immutable
		Amount:       o.Amount,
		NativeTokens: o.NativeTokens.Clone(),
		Blocks:       o.Blocks, // immutable
	}
}

func newInternalTokenOutput(aliasID iotago.AliasID, nativeTokenID iotago.NativeTokenID) *iotago.ExtendedOutput {
	return &iotago.ExtendedOutput{
		Address: aliasID.ToAddress(),
		Amount:  0,
		NativeTokens: iotago.NativeTokens{{
			ID:     nativeTokenID,
			Amount: big.NewInt(0),
		}},
		Blocks: iotago.FeatureBlocks{
			&iotago.SenderFeatureBlock{
				Address: aliasID.ToAddress(),
			},
		},
	}
}
