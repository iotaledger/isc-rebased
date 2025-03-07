package cryptolib

import (
	"crypto/ed25519"
	"encoding/binary"

	"github.com/wollac/iota-crypto-demo/pkg/bip32path"
	"github.com/wollac/iota-crypto-demo/pkg/slip10"
	"github.com/wollac/iota-crypto-demo/pkg/slip10/eddsa"
	"golang.org/x/crypto/blake2b"

	hivecrypto "github.com/iotaledger/hive.go/crypto/ed25519"

	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
	"github.com/iotaledger/wasp/packages/cryptolib/byteutils"
)

// testnet/alphanet uses COIN_TYPE = 1
const TestnetCoinType = iotasigner.TestnetCoinType

// / IOTA coin type <https://github.com/satoshilabs/slips/blob/master/slip-0044.md>
const IotaCoinType = iotasigner.IotaCoinType

// SubSeed returns a Seed (ed25519 Seed) from a master seed (that has arbitrary length)
// note that the accountIndex is actually an uint31
func SubSeed(walletSeed []byte, accountIndex uint32, useLegacyDerivation ...bool) Seed {
	if len(useLegacyDerivation) > 0 && useLegacyDerivation[0] {
		seed := SeedFromBytes(walletSeed)
		return legacyDerivation(&seed, accountIndex)
	}

	bip32Path, err := iotasigner.BuildBip32Path(iotasigner.SignatureFlagEd25519, iotasigner.IotaCoinType, accountIndex)
	if err != nil {
		panic(err)
	}

	path, err := bip32path.ParsePath(bip32Path)
	if err != nil {
		panic(err)
	}

	key, err := slip10.DeriveKeyFromPath(walletSeed, eddsa.Ed25519(), path)
	if err != nil {
		panic(err)
	}

	_, prvKey := key.Key.(eddsa.Seed).Ed25519Key()
	return SeedFromBytes(prvKey)
}

// ---
const (
	SeedSize = ed25519.SeedSize
)

type Seed [SeedSize]byte

func NewSeed() (ret Seed) {
	copy(ret[:], hivecrypto.NewSeed().Bytes())
	return ret
}

func SeedFromBytes(data []byte) (ret Seed) {
	copy(ret[:], data)
	return ret
}

func legacyDerivation(seed *Seed, index uint32) Seed {
	subSeed := make([]byte, SeedSize)
	indexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint32(indexBytes[:4], index)
	hashOfIndexBytes := blake2b.Sum256(indexBytes)
	byteutils.XORBytes(subSeed, seed[:], hashOfIndexBytes[:])
	return SeedFromBytes(subSeed)
}
