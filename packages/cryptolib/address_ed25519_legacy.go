package cryptolib

import (
	"fmt"

	"golang.org/x/crypto/blake2b"
)

const (
	// Ed25519AddressLegacyBytesLength is the length of a legacy Ed25519 address.
	Ed25519AddressLegacyBytesLength = blake2b.Size256
	// Ed25519AddressLegacySerializedBytesSize is the size of a serialized legacy Ed25519 address with its type denoting byte.
	//Ed25519AddressLegacySerializedBytesSize = serializer.SmallTypeDenotationByteSize + Ed25519AddressBytesLength
)

// Ed25519AddressLegacy defines a legacy (IOTA 2.0) Ed25519 address.
// An Ed25519AddressLegacy is the Blake2b-256 hash of an Ed25519 public key.
type Ed25519AddressLegacy [Ed25519AddressLegacyBytesLength]byte

var _ Address = &Ed25519AddressLegacy{}

func NewEd25519AddressLegacyFromBytes(bytes []byte) (*Ed25519AddressLegacy, error) {
	if len(bytes) != Ed25519AddressLegacyBytesLength {
		return nil, fmt.Errorf("address must be %d bytes long but is %d", Ed25519AddressLegacyBytesLength, len(bytes))
	}
	result := Ed25519AddressLegacy{}
	copy(result[:], bytes[:])
	return &result, nil
}

func (edAddr *Ed25519AddressLegacy) String() string {
	return EncodeHex(edAddr[:])
}

func (edAddr *Ed25519AddressLegacy) Type() AddressType {
	return AddressEd25519Legacy
}

func (edAddr *Ed25519AddressLegacy) Size() int {
	return Ed25519AddressLegacyBytesLength
}

func (edAddr *Ed25519AddressLegacy) Bytes() []byte {
	result := make([]byte, Ed25519AddressLegacyBytesLength)
	copy(result[:], edAddr[:])
	return result
}

func (edAddr *Ed25519AddressLegacy) Equal(other Address) bool {
	otherAddr, is := other.(*Ed25519AddressLegacy)
	if !is {
		return false
	}
	return *edAddr == *otherAddr
}

/*func (edAddr *Ed25519AddressLegacy) Clone() Address {
	cpy := &Ed25519Address{}
	copy(cpy[:], edAddr[:])
	return cpy
}

func (edAddr *Ed25519Address) VBytes(rentStruct *RentStructure, _ VBytesFunc) VBytes {
	return rentStruct.VBFactorData.Multiply(Ed25519AddressSerializedBytesSize)
}

func (edAddr *Ed25519Address) Key() string {
	return string(append([]byte{byte(AddressEd25519)}, (*edAddr)[:]...))
}

func (edAddr *Ed25519Address) Unlock(msg []byte, sig Signature) error {
	edSig, isEdSig := sig.(*Ed25519Signature)
	if !isEdSig {
		return fmt.Errorf("%w: can not unlock Ed25519 address with signature of type %s", ErrSignatureAndAddrIncompatible, sig.Type())
	}
	return edSig.Valid(msg, edAddr)
}
*/

/*func (edAddr *Ed25519Address) Bech32(hrp NetworkPrefix) string {
	return bech32String(hrp, edAddr)
}

func (edAddr *Ed25519Address) Deserialize(data []byte, deSeriMode serializer.DeSerializationMode, deSeriCtx interface{}) (int, error) {
	if deSeriMode.HasMode(serializer.DeSeriModePerformValidation) {
		if err := serializer.CheckMinByteLength(Ed25519AddressSerializedBytesSize, len(data)); err != nil {
			return 0, fmt.Errorf("invalid Ed25519 address bytes: %w", err)
		}
		if err := serializer.CheckTypeByte(data, byte(AddressEd25519)); err != nil {
			return 0, fmt.Errorf("unable to deserialize Ed25519 address: %w", err)
		}
	}
	copy(edAddr[:], data[serializer.SmallTypeDenotationByteSize:])
	return Ed25519AddressSerializedBytesSize, nil
}

func (edAddr *Ed25519Address) Serialize(_ serializer.DeSerializationMode, deSeriCtx interface{}) (data []byte, err error) {
	var b [Ed25519AddressSerializedBytesSize]byte
	b[0] = byte(AddressEd25519)
	copy(b[serializer.SmallTypeDenotationByteSize:], edAddr[:])
	return b[:], nil
}

func (edAddr *Ed25519Address) MarshalJSON() ([]byte, error) {
	jEd25519Address := &jsonEd25519Address{}
	jEd25519Address.PubKeyHash = EncodeHex(edAddr[:])
	jEd25519Address.Type = int(AddressEd25519)
	return json.Marshal(jEd25519Address)
}

func (edAddr *Ed25519Address) UnmarshalJSON(bytes []byte) error {
	jEd25519Address := &jsonEd25519Address{}
	if err := json.Unmarshal(bytes, jEd25519Address); err != nil {
		return err
	}
	seri, err := jEd25519Address.ToSerializable()
	if err != nil {
		return err
	}
	*edAddr = *seri.(*Ed25519Address)
	return nil
}

// Ed25519AddressFromPubKey returns the address belonging to the given Ed25519 public key.
func Ed25519AddressFromPubKey(pubKey ed25519.PublicKey) Ed25519Address {
	return blake2b.Sum256(pubKey[:])
}

// jsonEd25519Address defines the json representation of an Ed25519Address.
type jsonEd25519Address struct {
	Type       int    `json:"type"`
	PubKeyHash string `json:"pubKeyHash"`
}

func (j *jsonEd25519Address) ToSerializable() (serializer.Serializable, error) {
	addrBytes, err := DecodeHex(j.PubKeyHash)
	if err != nil {
		return nil, fmt.Errorf("unable to decode address from JSON for Ed25519 address: %w", err)
	}
	if err := serializer.CheckExactByteLength(len(addrBytes), Ed25519AddressBytesLength); err != nil {
		return nil, fmt.Errorf("unable to decode address from JSON for Ed25519 address: %w", err)
	}
	addr := &Ed25519Address{}
	copy(addr[:], addrBytes)
	return addr, nil
}
*/
