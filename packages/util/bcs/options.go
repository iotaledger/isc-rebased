package bcs

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type TypeOptions struct {
	//IncludeUnexported bool

	// TODO: Is this needed? It is present in rwutil as Size16/Size32, but it is more of validation.
	LenSizeInBytes LenBytesCount

	// TODO: Isthis really useful? The engineer can just change type of int to indicate its size.
	SizeInBytes ValueBytesCount

	IsCompactInt       bool
	InterfaceIsNotEnum bool

	ArrayElement *ArrayElemOptions
	MapKey       *TypeOptions
	MapValue     *TypeOptions
}

func (o *TypeOptions) Validate() error {
	if err := o.LenSizeInBytes.Validate(); err != nil {
		return fmt.Errorf("array len size: %w", err)
	}

	return nil
}

func (o *TypeOptions) Update(other TypeOptions) {
	if other.LenSizeInBytes != 0 {
		o.LenSizeInBytes = other.LenSizeInBytes
	}
	if other.SizeInBytes != 0 {
		o.SizeInBytes = other.SizeInBytes
	}
	if other.IsCompactInt {
		o.IsCompactInt = true
	}
	if other.InterfaceIsNotEnum {
		o.InterfaceIsNotEnum = true
	}
	if other.ArrayElement != nil {
		if o.ArrayElement == nil {
			o.ArrayElement = other.ArrayElement
		} else {
			o.ArrayElement.Update(*other.ArrayElement)
		}
	}
	if other.MapKey != nil {
		if o.MapKey == nil {
			o.MapKey = other.MapKey
		} else {
			o.MapKey.Update(*other.MapKey)
		}
	}
	if other.MapValue != nil {
		if o.MapValue == nil {
			o.MapValue = other.MapValue
		} else {
			o.MapValue.Update(*other.MapValue)
		}
	}
}

type ArrayElemOptions struct {
	TypeOptions
	AsByteArray bool
}

func (o *ArrayElemOptions) Update(other ArrayElemOptions) {
	o.TypeOptions.Update(other.TypeOptions)
	if other.AsByteArray {
		o.AsByteArray = true
	}
}

type FieldOptions struct {
	TypeOptions
	Skip     bool
	Optional bool
	//OmitEmpty bool
	//ByteOrder    binary.ByteOrder
	AsByteArray bool
}

func (o *FieldOptions) Validate() error {
	if err := o.TypeOptions.LenSizeInBytes.Validate(); err != nil {
		return fmt.Errorf("array len size: %w", err)
	}

	return nil
}

func FieldOptionsFromField(fieldType reflect.StructField, tagName string) (FieldOptions, bool, error) {
	a, hasTag := fieldType.Tag.Lookup(tagName)

	fieldOpts, err := FieldOptionsFromTag(a)
	if err != nil {
		return FieldOptions{}, false, fmt.Errorf("parsing annotation: %w", err)
	}

	switch fieldType.Type.Kind() {
	case reflect.Slice, reflect.Array:
		a, hasElemTag := fieldType.Tag.Lookup(tagName + "_elem")
		elemOpts, err := FieldOptionsFromTag(a)
		if err != nil {
			return FieldOptions{}, false, fmt.Errorf("parsing elem annotation: %w", err)
		}

		fieldOpts.ArrayElement = &ArrayElemOptions{
			TypeOptions: elemOpts.TypeOptions,
			AsByteArray: elemOpts.AsByteArray,
		}

		hasTag = hasTag || hasElemTag
	case reflect.Map:
		a, hasKeyTag := fieldType.Tag.Lookup(tagName + "_key")
		keyOpts, err := FieldOptionsFromTag(a)
		if err != nil {
			return FieldOptions{}, false, fmt.Errorf("parsing key annotation: %w", err)
		}

		fieldOpts.MapKey = &keyOpts.TypeOptions

		a, hasValueTag := fieldType.Tag.Lookup(tagName + "_value")
		valueOpts, err := FieldOptionsFromTag(a)
		if err != nil {
			return FieldOptions{}, false, fmt.Errorf("parsing value annotation: %w", err)
		}

		fieldOpts.MapValue = &valueOpts.TypeOptions

		hasTag = hasTag || hasKeyTag || hasValueTag
	}

	return fieldOpts, hasTag, nil
}

func FieldOptionsFromTag(a string) (_ FieldOptions, _ error) {
	if a == "" {
		return FieldOptions{}, nil
	}
	if a == "-" {
		return FieldOptions{Skip: true}, nil
	}

	opts := FieldOptions{}

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
		case "compact":
			opts.IsCompactInt = true
		case "bytes":
			bytes, err := strconv.Atoi(val)
			if err != nil {
				return FieldOptions{}, fmt.Errorf("invalid bytes tag: %s", val)
			}

			opts.SizeInBytes = ValueBytesCount(bytes)
		case "len_bytes":
			bytes, err := strconv.Atoi(val)
			if err != nil {
				return FieldOptions{}, fmt.Errorf("invalid len_bytes tag: %s", val)
			}

			opts.LenSizeInBytes = LenBytesCount(bytes)
		case "optional":
			opts.Optional = true
		case "bytearr":
			opts.AsByteArray = true
		case "not_enum":
			opts.InterfaceIsNotEnum = true
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

func defaultValueSize(k reflect.Kind) ValueBytesCount {
	switch k {
	case reflect.Int8, reflect.Uint8:
		return Value1Byte
	case reflect.Int16, reflect.Uint16:
		return Value2Bytes
	case reflect.Int32, reflect.Uint32:
		return Value4Bytes
	case reflect.Int64, reflect.Uint64, reflect.Int, reflect.Uint:
		return Value8Bytes
	default:
		panic(fmt.Errorf("unexpected kind: %v", k))
	}
}

type BCSType interface {
	BCSOptions() TypeOptions
}

var bcsTypeT = reflect.TypeOf((*BCSType)(nil)).Elem()
