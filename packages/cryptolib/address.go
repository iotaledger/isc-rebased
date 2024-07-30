package cryptolib

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/util/rwutil"
	"github.com/iotaledger/wasp/sui-go/sui"
)

const AddressSize = 32

type (
	Address    [AddressSize]byte
	AddressKey [AddressSize]byte
)

func NewEmptyAddress() *Address {
	return &Address{}
}

func NewRandomAddress() *Address {
	data := make([]byte, AddressSize)
	_, err := rand.Read(data)
	if err != nil {
		panic(fmt.Errorf("Failed to read random data, %s", err))
	}
	address, err := NewAddressFromBytes(data)
	if err != nil {
		panic(fmt.Errorf("Failed to create random address, %s", err))
	}
	return address
}

func newAddressFromArray(addr [AddressSize]byte) *Address {
	result := Address(addr)
	return &result
}

func NewAddressFromBytes(addr []byte) (*Address, error) {
	if len(addr) != AddressSize {
		return nil, fmt.Errorf("Array of size %v expected, size %v received", AddressSize, len(addr))
	}
	result := &Address{}
	copy(result[:], addr)
	return result, nil
}

func NewAddressFromHexString(addr string) (*Address, error) {
	addrBytes, err := DecodeHex(addr)
	if err != nil {
		return nil, fmt.Errorf("Error decoding hex: %w", err)
	}
	return NewAddressFromBytes(addrBytes)
}

func NewAddressFromKey(key AddressKey) *Address {
	result := Address(key)
	return &result
}

// TODO: remove when not needed
func NewAddressFromIotago(addr iotago.Address) *Address {
	addrBytes, err := addr.Serialize(0, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to obtain byte array from iotago address: %s", err))
	}
	address, err := NewAddressFromBytes(addrBytes[1:])
	if err != nil {
		panic(fmt.Sprintf("Failed to obtain address from byte array: %s", err))
	}
	return address
}

// TODO: remove when not needed
func (a *Address) AsSuiAddress() *sui.Address {
	result := sui.Address(a[:])
	return &result
}

func (a *Address) Equals(other *Address) bool {
	return *a == *other
}

func (a *Address) Bytes() []byte {
	return a[:]
}

func (a *Address) String() string {
	return EncodeHex(a.Bytes())
}

func (a *Address) Key() AddressKey {
	return AddressKey(*a)
}

func (a *Address) Clone() *Address {
	result := &Address{}
	copy(result[:], a[:])
	return result
}

func (a *Address) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	address := rr.ReadBytes()
	copy(a[:], address)
	return rr.Err
}

func (a *Address) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteBytes(a[:])
	return ww.Err
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func AddressFromHex(str string) (*Address, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}
	if len(str)%2 != 0 {
		str = "0" + str
	}
	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	if len(data) > sui.AddressLen {
		return nil, errors.New("the len is invalid")
	}
	var address Address
	copy(address[sui.AddressLen-len(data):], data[:])
	return &address, nil
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var str *string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	if str == nil {
		return errors.New("nil address")
	}
	tmp, err := AddressFromHex(*str)
	if err == nil {
		*a = *tmp
	}

	return err
}

// FIXME may need to be pointer
func (a Address) MarshalBCS() ([]byte, error) {
	return a[:], nil
}
