package cryptolib

import (
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
)

// VariantKeyPair originates from KeyPair
type Signer interface {
	// IsNil is a mandatory nil check. This includes the referenced keypair implementation pointer. `kp == nil` is not enough.
	// IsNil() bool

	Address() *Address
	Sign(msg []byte) (signature *Signature, err error)
	SignTransactionBlock(txnBytes []byte, intent iotasigner.Intent) (*Signature, error)
}

type iotaSigner struct {
	s Signer
}

// TODO: remove, when it is not needed
func SignerToIotaSigner(s Signer) iotasigner.Signer {
	return &iotaSigner{s}
}

func (is *iotaSigner) Address() *iotago.Address {
	return is.s.Address().AsIotaAddress()
}

func (is *iotaSigner) Sign(msg []byte) (signature *iotasigner.Signature, err error) {
	b, err := is.s.Sign(msg)
	if err != nil {
		return nil, err
	}
	return b.AsIotaSignature(), err
}

func (is *iotaSigner) SignTransactionBlock(txnBytes []byte, intent iotasigner.Intent) (*iotasigner.Signature, error) {
	signature, err := is.s.SignTransactionBlock(txnBytes, intent)
	if err != nil {
		return nil, err
	}
	return signature.AsIotaSignature(), nil
}
