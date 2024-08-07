package coreutil

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/kv/codec"
)

// Optional returns an optional value (type *T) from a variadic parameter
// (...T) which can be of length 0 or 1.
func Optional[T any](v ...T) *T {
	if len(v) > 0 {
		return &v[0]
	}
	return nil
}

// FromOptional extracts a value of type T from an optional (*T) and a default.
func FromOptional[T any](opt *T, def T) T {
	if opt == nil {
		return def
	}
	return *opt
}

// CallArgsCodec is the interface for any type that can be converted to/from dict.Dict
type CallArgsCodec[T any] interface {
	Encode(T) []byte
	Decode([]byte) (T, error)
}

// RawCallArgsCodec is a CallArgsCodec that performs no conversion
type RawCallArgsCodec struct{}

func (RawCallArgsCodec) Decode(d []byte) ([]byte, error) {
	return d, nil
}

func (RawCallArgsCodec) Encode(d []byte) []byte {
	return d
}

// Field is a CallArgsCodec that converts a single value into a single dict key
type Field[T any] struct {
	Codec codec.Codec[T]
}

func (f Field[T]) Encode(v T) []byte {
	b := f.Codec.Encode(v)
	if b == nil {
		return []byte{}
	}
	return b
}

func (f Field[T]) Decode(d []byte) (T, error) {
	return f.Codec.Decode(d)
}

func FieldWithCodec[T any](codec codec.Codec[T]) Field[T] {
	return Field[T]{Codec: codec}
}

// OptionalCodec is a Codec that converts to/from an optional value of type T.
type OptionalCodec[T any] struct {
	codec.Codec[T]
}

func (c *OptionalCodec[T]) Decode(b []byte, def ...*T) (r *T, err error) {
	if b == nil {
		if len(def) != 0 {
			err = fmt.Errorf("%T: unexpected default value", r)
			return
		}
		return nil, nil
	}
	v, err := c.Codec.Decode(b)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *OptionalCodec[T]) MustDecode(b []byte, def ...*T) *T {
	return lo.Must(c.Decode(b, def...))
}

func (c *OptionalCodec[T]) Encode(v *T) []byte {
	if v == nil {
		return nil
	}
	return c.Codec.Encode(*v)
}

// FieldWithCodecOptional returns a Field that accepts an optional value
func FieldWithCodecOptional[T any](c codec.Codec[T]) Field[*T] {
	return Field[*T]{Codec: &OptionalCodec[T]{Codec: c}}
}
