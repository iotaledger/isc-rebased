package coreutil

import (
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

// field is a CallArgsCodec that converts a single value of T
type field[T any] struct{}

var _ CallArgsCodec[any] = (*field[any])(nil)

func (f field[T]) Encode(v T) []byte {
	return codec.Encode(v)
}

func (f field[T]) Decode(d []byte) (T, error) {
	return codec.Decode[T](d)
}

func Field[T any]() field[T] {
	return field[T]{}
}

// OptionalCodec is a Codec that converts to/from an optional value of type T.

func (c optionalField[T]) Decode(b []byte) (r *T, err error) {
	return codec.DecodeOptional[T](b)
}

func (c optionalField[T]) Encode(v *T) []byte {
	return codec.EncodeOptional(v)
}

// FieldWithCodecOptional returns a Field that accepts an optional value
func FieldOptional[T any]() optionalField[T] {
	return optionalField[T]{}
}

type optionalField[T any] struct{}

var _ CallArgsCodec[*any] = (*optionalField[any])(nil)
