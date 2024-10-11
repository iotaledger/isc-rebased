package coreutil

import (
	"fmt"
	"io"
	"math"

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

// Field is a CallArgsCodec that converts a single value of T
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
	Codec codec.Codec[T]
}

var _ codec.Codec[*any] = &OptionalCodec[any]{}

func (c *OptionalCodec[T]) Decode(b []byte, def ...*T) (r *T, err error) {
	if b[0] == 0 {
		if len(def) != 0 {
			err = fmt.Errorf("%T: unexpected default value", r)
			return
		}
		return nil, nil
	}
	if b[0] != 1 {
		return nil, fmt.Errorf("%T: invalid optional flag value '%v'", b[0], r)
	}

	v, err := c.Codec.Decode(b[1:])
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
		return []byte{0}
	}

	b := c.Codec.Encode(*v)
	return append([]byte{1}, b...)
}

func (c *OptionalCodec[T]) EncodeStream(w io.Writer, v *T) {
	if v == nil {
		w.Write([]byte{0})
		return
	}

	w.Write([]byte{1})
	c.Codec.EncodeStream(w, *v)
}

func (c *OptionalCodec[T]) DecodeStream(r io.Reader) (*T, error) {
	var flag [1]byte
	if _, err := r.Read(flag[:]); err != nil {
		return nil, err
	}

	if flag[0] == 0 {
		return nil, nil
	}
	if flag[0] != 1 {
		return nil, fmt.Errorf("%T: invalid optional flag value '%v'", new(T), flag[0])
	}

	v, err := c.Codec.DecodeStream(r)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (c *OptionalCodec[T]) MustDecodeStream(r io.Reader) *T {
	return lo.Must(c.DecodeStream(r))
}

// FieldWithCodecOptional returns a Field that accepts an optional value
func FieldWithCodecOptional[T any](c codec.Codec[T]) Field[*T] {
	return Field[*T]{Codec: &OptionalCodec[T]{Codec: c}}
}

type FieldArrayOf[T any] struct {
	codec codec.Codec[T]
}

// FieldArrayWithCodec returns a Field that encodes a slice of T
func FieldArrayWithCodec[T any](codec codec.Codec[T]) FieldArrayOf[T] {
	return FieldArrayOf[T]{codec: codec}
}

func (a FieldArrayOf[T]) Encode(slice []T) []byte {
	if len(slice) > math.MaxUint16 {
		panic("too many values")
	}
	return codec.SliceToArray(a.codec, slice)
}

func (a FieldArrayOf[T]) Decode(r []byte) ([]T, error) {
	return codec.SliceFromArray(a.codec, r)
}
