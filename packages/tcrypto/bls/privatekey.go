package bls

import (
	`bytes`
	`io`

	"github.com/btcsuite/btcd/btcutil/base58"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/bdn"

	"github.com/iotaledger/hive.go/ierrors"
)

// PrivateKey is the type of BLS private keys.
type PrivateKey struct {
	Scalar kyber.Scalar
}

// PrivateKeyFromBytes creates a PrivateKey from the given bytes.
func PrivateKeyFromBytes(b []byte) (privateKey PrivateKey, err error) {
	marshalUtil := bytes.NewReader(b)
	if privateKey, err = PrivateKeyFromReader(marshalUtil); err != nil {
		err = ierrors.Wrap(err, "failed to parse PrivateKey from Reader")
	}

	return
}

// PrivateKeyFromBase58EncodedString creates a PrivateKey from a base58 encoded string.
func PrivateKeyFromBase58EncodedString(base58String string) (privateKey PrivateKey, err error) {
	bytes := base58.Decode(base58String)
	if len(bytes) == 0 {
		err = ierrors.Wrapf(ErrBase58DecodeFailed, "error while decoding base58 encoded PrivateKey: %s", base58String)

		return
	}

	if privateKey, err = PrivateKeyFromBytes(bytes); err != nil {
		err = ierrors.Wrap(err, "failed to parse PrivateKey from bytes")

		return
	}

	return
}

// PrivateKeyFromReader unmarshals a PrivateKey using a Reader (for easier unmarshalling).
func PrivateKeyFromReader(reader *bytes.Reader) (privateKey PrivateKey, err error) {
	privateKeyBytes := make([]byte, PrivateKeySize)
	readLength, err := io.ReadFull(reader, privateKeyBytes)

	if err != nil {
		err = ierrors.Wrapf(ErrParseBytesFailed, "failed to read PrivateKey bytes: %v", err)
		return
	}

	if readLength != PrivateKeySize {
		err = ierrors.Wrapf(ErrParseBytesFailed, "failed to read PrivateKey bytes: %v", err)
		return
	}

	if err = privateKey.Scalar.UnmarshalBinary(privateKeyBytes); err != nil {
		err = ierrors.Wrapf(ErrParseBytesFailed, "failed to unmarshal PrivateKey: %w", err)

		return
	}

	return
}

// PrivateKeyFromRandomness generates a new random PrivateKey.
func PrivateKeyFromRandomness() (privateKey PrivateKey) {
	privateKey.Scalar, _ = bdn.NewKeyPair(blsSuite, randomness)

	return
}

// PublicKey returns the PublicKey corresponding to the PrivateKey.
func (p PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Point: blsSuite.G2().Point().Mul(p.Scalar, nil),
	}
}

// Sign signs the message and returns a SignatureWithPublicKey.
func (p PrivateKey) Sign(data []byte) (signatureWithPublicKey SignatureWithPublicKey, err error) {
	sig, err := bdn.Sign(blsSuite, p.Scalar, data)
	if err != nil {
		err = ierrors.Wrapf(ErrBLSFailed, "failed to sign data: %w", err)

		return
	}

	signatureWithPublicKey.PublicKey = p.PublicKey()
	copy(signatureWithPublicKey.Signature[:], sig)

	return
}

// Bytes returns a marshaled version of the PrivateKey.
func (p PrivateKey) Bytes() (bytes []byte) {
	bytes, err := p.Scalar.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return
}

// Base58 returns a base58 encoded version of the PrivateKey.
func (p PrivateKey) Base58() string {
	return base58.Encode(p.Bytes())
}

// String returns a human-readable version of the PrivateKey (base58 encoded).
func (p PrivateKey) String() string {
	return p.Base58()
}
