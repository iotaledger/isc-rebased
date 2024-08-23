package wbf

import (
	"fmt"
	"strconv"
	"strings"
)

type TypeOptions struct {
	//IncludeUnexported bool
	LenBytes LenBytesCount
	Bytes    ValueBytesCount
}

func (o *TypeOptions) Validate() error {
	if err := o.LenBytes.Validate(); err != nil {
		return fmt.Errorf("array len size: %w", err)
	}

	return nil
}

func (o *TypeOptions) Update(other TypeOptions) {
	if other.LenBytes != 0 {
		o.LenBytes = other.LenBytes
	}
	if other.Bytes != 0 {
		o.Bytes = other.Bytes
	}
}

var DefaultTypeOptions = TypeOptions{
	LenBytes: Len4Bytes,
}

type FieldOptions struct {
	TypeOptions
	Skip     bool
	Optional bool
	//OmitEmpty bool
	//ByteOrder    binary.ByteOrder
}

func (o *FieldOptions) Validate() error {
	if err := o.TypeOptions.LenBytes.Validate(); err != nil {
		return fmt.Errorf("array len size: %w", err)
	}

	return nil
}

func FieldOptionsFromTag(a string, defTypOpts TypeOptions) (_ FieldOptions, _ error) {
	if a == "" {
		return FieldOptions{TypeOptions: defTypOpts}, nil
	}
	if a == "-" {
		return FieldOptions{Skip: true}, nil
	}

	opts := FieldOptions{TypeOptions: defTypOpts}

	parts := strings.Split(a, ",")

	for _, part := range parts {
		subparts := strings.Split(part, "=")

		if len(subparts) > 2 {
			return FieldOptions{}, fmt.Errorf("invalid field tag: %s", part)
		}

		key := subparts[0]
		val := ""
		if len(subparts) == 2 {
			val = subparts[1]
		}

		switch key {
		case "bytes":
			bytes, err := strconv.Atoi(val)
			if err != nil {
				return FieldOptions{}, fmt.Errorf("invalid bytes tag: %s", val)
			}

			opts.Bytes = ValueBytesCount(bytes)
		case "len_bytes":
			bytes, err := strconv.Atoi(val)
			if err != nil {
				return FieldOptions{}, fmt.Errorf("invalid len_bytes tag: %s", val)
			}

			opts.LenBytes = LenBytesCount(bytes)
		case "optional":
			opts.Optional = true
		case "":
			return FieldOptions{}, fmt.Errorf("empty field tag")
		default:
			return FieldOptions{}, fmt.Errorf("unknown field tag: %s", key)
		}
	}

	return opts, nil
}

type LenBytesCount uint8

const (
	Len2Bytes LenBytesCount = 2
	Len4Bytes LenBytesCount = 4
)

func (s LenBytesCount) Validate() error {
	switch s {
	case Len2Bytes, Len4Bytes:
		return nil
	default:
		return fmt.Errorf("invalid collection len size: %v", s)
	}
}

type ValueBytesCount uint8

const (
	Value1Byte  ValueBytesCount = 1
	Value2Bytes ValueBytesCount = 2
	Value4Bytes ValueBytesCount = 4
	Value8Bytes ValueBytesCount = 8
)

func (s ValueBytesCount) Validate() error {
	switch s {
	case Value1Byte, Value2Bytes, Value4Bytes, Value8Bytes:
		return nil
	default:
		return fmt.Errorf("invalid value size: %v", s)
	}
}

type WBFType interface {
	WBFOptions() TypeOptions
}
