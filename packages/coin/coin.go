package coin

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/util/bcs"
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
)

// Value is the balance of a given coin
type Value uint64

func (v Value) Uint64() uint64 {
	return uint64(v)
}

func (v *Value) MarshalBCS(e *bcs.Encoder) error {
	e.WriteCompactUint(uint64(*v))
	return e.Err()
}

func (v *Value) UnmarshalBCS(d *bcs.Decoder) error {
	*v = Value(d.ReadCompactUint())
	return d.Err()
}

func (v Value) Bytes() []byte {
	return bcs.MustMarshal(&v)
}

func ValueFromBytes(b []byte) (Value, error) {
	return bcs.Unmarshal[Value](b)
}

// TODO: maybe it is not ok to consider this constant?
const BaseTokenType = Type(suijsonrpc.SuiCoinType)

// Type is the string representation of a Sui coin type, e.g. `0x2::sui::SUI`
type Type string

func (t *Type) MarshalBCS(e *bcs.Encoder) error {
	rt, err := sui.NewResourceType(string(*t))
	if err != nil {
		return fmt.Errorf("invalid Type %q: %w", t, err)
	}
	e.Encode(rt)
	return e.Err()
}

func (t *Type) UnmarshalBCS(d *bcs.Decoder) error {
	var rt sui.ResourceType
	d.Decode(&rt)
	if d.Err() != nil {
		return d.Err()
	}
	*t = Type(rt.ShortString())
	return nil
}

func (t Type) String() string {
	return string(t)
}

func (t Type) TypeTag() sui.TypeTag {
	coinTypeTag, err := sui.TypeTagFromString(t.String())
	if err != nil {
		panic(err)
	}
	return *coinTypeTag
}

func (t Type) Bytes() []byte {
	return bcs.MustMarshal(&t)
}

func TypeFromBytes(b []byte) (Type, error) {
	var r Type
	r, err := bcs.Unmarshal[Type](b)
	return r, err
}
