package sui_signer

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
)

const (
	SignatureFlagEd25519   = 0x0
	SignatureFlagSecp256k1 = 0x1

	DerivationPathEd25519   = `m/44'/4218'/0'/0'/0'`
	DerivationPathSecp256k1 = `m/54'/4218'/0'/0/0`
)

var (
	TEST_MNEMONIC = "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
	TEST_ADDRESS  = sui_types.MustSuiAddressFromHex("0x786dff8a4ee13d45b502c8f22f398e3517e6ec78aa4ae564c348acb07fad7f50")
)

// FIXME support more than ed25519
type Signer struct {
	ed25519Keypair *KeypairEd25519
	// secp256k1Keypair *KeypairSecp256k1
	Address *sui_types.SuiAddress
}

func NewSigner(seed []byte) *Signer {
	prikey := ed25519.NewKeyFromSeed(seed[:])
	pubkey := prikey.Public().(ed25519.PublicKey)

	buf := append([]byte{FlagEd25519.Byte()}, pubkey...)
	addrBytes := blake2b.Sum256(buf)
	addr := "0x" + hex.EncodeToString(addrBytes[:])

	return &Signer{
		ed25519Keypair: &KeypairEd25519{
			PriKey: prikey,
			PubKey: pubkey,
		},
		Address: sui_types.MustSuiAddressFromHex(addr),
	}
}

// TODO add NewSignerWithFund

func NewSignerWithMnemonic(mnemonic string) (*Signer, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := DeriveForPath(DerivationPathEd25519, seed)
	if err != nil {
		return nil, err
	}
	return NewSigner(key.Key), nil
}

func (s *Signer) PrivateKey() []byte {
	switch {
	case s.ed25519Keypair != nil:
		return s.ed25519Keypair.PriKey
	default:
		return nil
	}
}

func (s *Signer) PublicKey() []byte {
	switch {
	case s.ed25519Keypair != nil:
		return s.ed25519Keypair.PubKey
	default:
		return nil
	}
}

func (s *Signer) Sign(data []byte) Signature {
	// FIXME support more than ed25519
	return Signature{
		Ed25519SuiSignature: NewEd25519SuiSignature(s, data),
	}
}

// FIXME support more than ed25519
func (a *Signer) SignTransactionBlock(txnBytes []byte, intent Intent) (Signature, error) {
	data := MessageWithIntent(intent, bcsBytes(txnBytes))
	hash := blake2b.Sum256(data)
	return a.Sign(hash[:]), nil
}

type bcsBytes []byte

func (b bcsBytes) MarshalBCS() ([]byte, error) {
	return b, nil
}
