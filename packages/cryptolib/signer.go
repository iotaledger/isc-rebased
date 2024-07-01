package cryptolib

import (
	"github.com/iotaledger/wasp/sui-go/suisigner"
	"github.com/iotaledger/wasp/sui-go/sui"
)

// VariantKeyPair originates from KeyPair
type Signer interface {
	// IsNil is a mandatory nil check. This includes the referenced keypair implementation pointer. `kp == nil` is not enough.
	// IsNil() bool

	Address() *Address
	Sign(msg []byte) (signature *Signature, err error)
	SignTransactionBlock(txnBytes []byte, intent suisigner.Intent) (Signature, error)
}

type suiSigner struct {
	s Signer
}

// TODO: remove, when it is not needed
func SignerToSuiSigner(s Signer) suisigner.Signer {
	return &suiSigner{s}
}

func (is *suiSigner) Address() *sui.Address {
	return is.s.Address().AsSuiAddress()
}

func (is *suiSigner) Sign(msg []byte) (signature *suisigner.Signature, err error) {
	b, err := is.s.Sign(msg)
	if err != nil {
		return nil, err
	}
	return b.AsSuiSignature(), err
}

func (is *suiSigner) SignTransactionBlock(txnBytes []byte, intent suisigner.Intent) (*suisigner.Signature, error) {
	signature, err := is.s.SignTransactionBlock(txnBytes, intent)
	if err != nil {
		return nil, err
	}
	return signature.AsSuiSignature(), nil
}
