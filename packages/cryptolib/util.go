package cryptolib

import (
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func SignatureFromBytes(bytes []byte) (result [SignatureSize]byte, err error) {
	if len(bytes) < SignatureSize {
		err = errors.New("bytes too short")
		return
	}
	copy(result[:], bytes)
	return
}

func IsVariantKeyPairValid(variantKeyPair VariantKeyPair) bool {
	if variantKeyPair == nil {
		return false
	}

	return !variantKeyPair.IsNil()
}

// EncodeHex encodes the bytes string to a hex string. It always adds the 0x prefix if bytes are not empty.
func EncodeHex(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return hexutil.Encode(b)
}
