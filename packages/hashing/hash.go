package hashing

import (
	"encoding/json"
	"fmt"
	"hash"
	"math/rand"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/util/bcs"
)

const HashSize = 32

type HashValue [HashSize]byte

var NilHash = HashValue{}

func (h HashValue) Bytes() []byte {
	return h[:]
}

func (h HashValue) String() string {
	return h.Hex()
}

func (h HashValue) Hex() string {
	return iotago.EncodeHex(h[:])
}

func (h *HashValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *HashValue) UnmarshalJSON(buf []byte) error {
	var s string
	err := json.Unmarshal(buf, &s)
	if err != nil {
		return err
	}
	ret, err := HashValueFromHex(s)
	if err != nil {
		return err
	}
	copy(h[:], ret[:])
	return nil
}

func HashValueFromBytes(b []byte) (HashValue, error) {
	return bcs.Unmarshal[HashValue](b)
}

func MustHashValueFromHex(s string) HashValue {
	ret, err := HashValueFromHex(s)
	if err != nil {
		panic(err)
	}
	return ret
}

func HashValueFromHex(s string) (HashValue, error) {
	b, err := iotago.DecodeHex(s)
	if err != nil {
		return NilHash, err
	}
	return HashValueFromBytes(b)
}

// HashData Blake2b
func HashData(data ...[]byte) HashValue {
	return HashDataBlake2b(data...)
}

func HashDataBlake2b(data ...[]byte) (ret HashValue) {
	h := hashBlake2b()
	for _, d := range data {
		_, err := h.Write(d)
		if err != nil {
			panic(err)
		}
	}
	copy(ret[:], h.Sum(nil))
	return
}

func hashBlake2b() hash.Hash {
	h, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	if h.Size() != HashSize {
		panic("blake2b: hash size != 32")
	}
	return h
}

func HashKeccak(data ...[]byte) (ret HashValue) {
	h := hashKeccak()
	for _, d := range data {
		_, err := h.Write(d)
		if err != nil {
			panic(err)
		}
	}
	copy(ret[:], h.Sum(nil))
	return
}

func hashKeccak() hash.Hash {
	h := sha3.NewLegacyKeccak256()
	if h.Size() != HashSize {
		panic("keccak: hash size != 32")
	}
	return h
}

func HashSha3(data ...[]byte) (ret HashValue) {
	h := hashSha3()
	for _, d := range data {
		_, err := h.Write(d)
		if err != nil {
			panic(err)
		}
	}
	copy(ret[:], h.Sum(nil))
	return
}

func hashSha3() hash.Hash {
	h := sha3.New256()
	if h.Size() != HashSize {
		panic("sha3: hash size != 32")
	}
	return h
}

func HashStrings(str ...string) HashValue {
	tarr := make([][]byte, len(str))
	for i, s := range str {
		tarr[i] = []byte(s)
	}
	return HashData(tarr...)
}

func PseudoRandomHash(rnd *rand.Rand) HashValue {
	var s string
	if rnd == nil {
		s = fmt.Sprintf("%d", rand.Int())
	} else {
		s = fmt.Sprintf("%d", rnd.Int())
	}
	ret := HashStrings(s, s, s)
	return ret
}
