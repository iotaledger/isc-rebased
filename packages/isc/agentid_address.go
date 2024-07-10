package isc

import (
	"io"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// AddressAgentID is an AgentID backed by a non-alias address.
type AddressAgentID struct {
	a *cryptolib.Address
}

var _ AgentIDWithL1Address = &AddressAgentID{}

func NewAddressAgentID(addr *cryptolib.Address) *AddressAgentID {
	return &AddressAgentID{a: addr}
}

func addressAgentIDFromString(s string) (*AddressAgentID, error) {
	addr, err := cryptolib.NewAddressFromHexString(s)
	if err != nil {
		return nil, err
	}
	return &AddressAgentID{a: addr}, nil
}

func (a *AddressAgentID) Address() *cryptolib.Address {
	return a.a
}

func (a *AddressAgentID) Bytes() []byte {
	return rwutil.WriteToBytes(a)
}

func (a *AddressAgentID) BelongsToChain(ChainID) bool {
	return false
}

func (a *AddressAgentID) BytesWithoutChainID() []byte {
	return a.Bytes()
}

func (a *AddressAgentID) Equals(other AgentID) bool {
	if other == nil {
		return false
	}
	if other.Kind() != a.Kind() {
		return false
	}
	return other.(*AddressAgentID).a.Equals(a.a)
}

func (a *AddressAgentID) Kind() AgentIDKind {
	return AgentIDKindAddress
}

func (a *AddressAgentID) String() string {
	return a.a.String()
}

func (a *AddressAgentID) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadKindAndVerify(rwutil.Kind(a.Kind()))
	a.a = cryptolib.NewEmptyAddress()
	rr.Read(a.a)
	return rr.Err
}

func (a *AddressAgentID) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteKind(rwutil.Kind(a.Kind()))
	ww.Write(a.a)
	return ww.Err
}
