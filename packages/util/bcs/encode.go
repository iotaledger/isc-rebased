package bcs

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/iotaledger/wasp/packages/util/rwutil"
	"github.com/samber/lo"
)

type Encodable interface {
	MarshalBCS(e *Encoder) error
}

var encodableT = reflect.TypeOf((*Encodable)(nil)).Elem()

type EncoderConfig struct {
	TagName string
	// IncludeUnexported bool
	// IncludeUntaggedUnexported bool
	// ExcludeUntagged           bool
	CustomEncoders map[reflect.Type]CustomEncoder
}

func (c *EncoderConfig) InitializeDefaults() {
	if c.TagName == "" {
		c.TagName = "bcs"
	}
	if c.CustomEncoders == nil {
		c.CustomEncoders = CustomEncoders
	}
}

func NewEncoder(dest io.Writer, cfg EncoderConfig) *Encoder {
	cfg.InitializeDefaults()

	return &Encoder{
		cfg: cfg,
		w:   *rwutil.NewWriter(dest),
	}
}

type Encoder struct {
	cfg EncoderConfig
	w   rwutil.Writer
}

func (e *Encoder) Encode(v any) error {
	if v == nil {
		return fmt.Errorf("cannot encode a nil value")
	}

	return e.encodeValue(reflect.ValueOf(v), nil)
}

func (e *Encoder) encodeValue(v reflect.Value, typeOptionsFromTag *TypeOptions) error {
	v, typeOptions, enumVariantIdx, customEncoder, err := e.dereferenceValue(v)
	if err != nil {
		return fmt.Errorf("%v: %w", v.Type(), err)
	}

	if customEncoder != nil {
		if err := customEncoder(e, v); err != nil {
			return fmt.Errorf("%v: custom encoder: %w", v.Type(), err)
		}

		return nil
	}

	if typeOptionsFromTag != nil {
		typeOptions.Update(*typeOptionsFromTag)
	}

	switch v.Kind() {
	case reflect.Bool:
		e.w.WriteBool(v.Bool())
	case reflect.Int8:
		if err := e.encodeInt(v, Value1Byte, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Uint8:
		if err := e.encodeUint(v, Value1Byte, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Int16:
		if err := e.encodeInt(v, Value2Bytes, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Uint16:
		if err := e.encodeUint(v, Value2Bytes, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Int32:
		if err := e.encodeInt(v, Value4Bytes, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Uint32:
		if err := e.encodeUint(v, Value4Bytes, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Int64:
		if err := e.encodeInt(v, Value8Bytes, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Uint64:
		if err := e.encodeUint(v, Value8Bytes, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Int:
		if err := e.encodeInt(v, Value8Bytes, typeOptions.Bytes); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.String:
		e.w.WriteString(v.String())
	case reflect.Slice:
		if err := e.encodeSlice(v, typeOptions); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Array:
		if err := e.encodeArray(v); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Map:
		if err := e.encodeMap(v, typeOptions); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	case reflect.Struct:
		if enumVariantIdx == -1 {
			if err := e.encodeStruct(v); err != nil {
				return fmt.Errorf("%v: %w", v.Type(), err)
			}
		} else {
			if err := e.encodeEnum(v.Field(enumVariantIdx), enumVariantIdx); err != nil {
				return fmt.Errorf("%v: %w", v.Type(), err)
			}
		}
	case reflect.Interface:
		if err := e.encodeEnum(v.Elem(), enumVariantIdx); err != nil {
			return fmt.Errorf("%v: %w", v.Type(), err)
		}
	default:
		return fmt.Errorf("%v: cannot encode unknown type type", v.Type())
	}

	if e.w.Err != nil {
		return fmt.Errorf("%v: %w", v.Type(), e.w.Err)
	}

	return nil
}

func (e *Encoder) dereferenceValue(v reflect.Value) (dereferenced reflect.Value, _ TypeOptions, enumVariantIdx int, _ CustomEncoder, _ error) {
	// Removing all redundant pointers

	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v, TypeOptions{}, -1, nil, fmt.Errorf("attempt to encode non-optinal nil value of type %v", v.Type())
		}

		// Before dereferencing pointer, we should check if maybe current type is already the type we should encode.
		typeOptions, enumVariantIdx, customEncoder, err := e.retrieveTypeInfo(v)
		if err != nil || typeOptions != nil || enumVariantIdx != -1 || customEncoder != nil {
			if typeOptions == nil {
				typeOptions = &TypeOptions{}
			}
			return v, *typeOptions, enumVariantIdx, customEncoder, err
		}

		v = v.Elem()
	}

	typeOptions, enumVariantIdx, customEncoder, err := e.retrieveTypeInfo(v)
	if typeOptions == nil {
		typeOptions = &TypeOptions{}
	}

	return v, *typeOptions, enumVariantIdx, customEncoder, err
}

func (e *Encoder) retrieveTypeInfo(v reflect.Value) (_ *TypeOptions, enumVariantIdx int, _ CustomEncoder, _ error) {
	// Detecting enum variant index might return error, so we
	// should first check for existance of custom encoder.
	if customEncoder := e.getCustomEncoder(v); customEncoder != nil {
		return nil, -1, customEncoder, nil
	}

	kind := v.Kind()

	switch {
	case kind == reflect.Interface:
		// Interface Enum
		enumVariantIdx, err := e.getInterfaceEnumVariantIdx(v)
		// // Rechecking for existance of custom encoder after we found type of enum variant.
		// // Maybe there was no custom encoder for enum type itself, but there could be for its variant.
		// customEncoder := e.getCustomEncoder(v.Elem())
		return nil, enumVariantIdx, nil, err
	case kind == reflect.Struct && v.Type().Implements(enumT):
		// Struct enum
		enumVariantIdx, err := e.getStructEnumVariantIdx(v)
		return nil, enumVariantIdx, nil, err
	default:
		vI := v.Interface()

		// This type does not have custom encoder, but it might provide encoding options.
		if bcsType, ok := vI.(BCSType); ok {
			typeOptions := bcsType.BCSOptions()

			return &typeOptions, -1, nil, nil
		}

		return nil, -1, nil, nil
	}
}

func (e *Encoder) getCustomEncoder(v reflect.Value) CustomEncoder {
	t := v.Type()

	// Check if this type has custom encoder func
	if customEncoder, ok := e.cfg.CustomEncoders[t]; ok {
		return customEncoder
	}

	// Check if this type implements custom encoding interface.
	// Although we could allow encoding of interfaces, which implement Encodable, still
	// we exclude them here to ensure symetric behaviour with decoding.
	if t.Kind() != reflect.Interface && t.Implements(encodableT) {
		encodable := v.Interface().(Encodable)

		customEncoder := func(e *Encoder, v reflect.Value) error {
			return encodable.MarshalBCS(e)
		}

		return customEncoder
	}

	return nil
}

func (e *Encoder) getInterfaceEnumVariantIdx(v reflect.Value) (enumVariantIdx int, _ error) {
	if v.IsNil() {
		return -1, fmt.Errorf("attemp to encode non-optional nil interface")
	}

	t := v.Type()

	enumVariants, registered := EnumTypes[t]
	if !registered {
		return -1, fmt.Errorf("interface %v is not registered as enum type", t)
	}

	valT := v.Elem().Type()

	for i, variant := range enumVariants {
		if valT == variant {
			enumVariantIdx = i
		}
	}

	if enumVariantIdx == -1 {
		return -1, fmt.Errorf("variant %v is not registered as part of enum type %v", valT, t)
	}

	return enumVariantIdx, nil
}

func (e *Encoder) getStructEnumVariantIdx(v reflect.Value) (enumVariantIdx int, _ error) {
	enumVariantIdx = -1

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		k := field.Kind()
		switch k {
		case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
			if field.IsNil() {
				continue
			}

			if enumVariantIdx != -1 {
				prevSetField := v.Type().Field(enumVariantIdx)
				currentField := v.Type().Field(i)
				return -1, fmt.Errorf("multiple options are set in enum struct %v: %v and %v", v.Type(), prevSetField.Name, currentField.Name)
			}

			enumVariantIdx = i
			// We do not break here to check if there are multiple options set
		default:
			fieldType := v.Type().Field(i)
			return -1, fmt.Errorf("field %v of enum %v is of non-nullable type %v", fieldType.Name, v.Type(), fieldType.Type)
		}
	}

	if enumVariantIdx == -1 {
		return -1, fmt.Errorf("no options are set in enum struct %v", v.Type())
	}

	return enumVariantIdx, nil
}

func (e *Encoder) encodeInt(v reflect.Value, origSize, customSize ValueBytesCount) error {
	size := lo.Ternary(customSize != 0, customSize, origSize)

	switch size {
	case Value1Byte:
		e.w.WriteInt8(int8(v.Int()))
	case Value2Bytes:
		e.w.WriteInt16(int16(v.Int()))
	case Value4Bytes:
		e.w.WriteInt32(int32(v.Int()))
	case Value8Bytes:
		e.w.WriteInt64(v.Int())
	default:
		return fmt.Errorf("invalid value size: %v", size)
	}

	return e.w.Err
}

func (e *Encoder) encodeUint(v reflect.Value, origSize, customSize ValueBytesCount) error {
	size := lo.Ternary(customSize != 0, customSize, origSize)

	switch size {
	case Value1Byte:
		e.w.WriteUint8(uint8(v.Uint()))
	case Value2Bytes:
		e.w.WriteUint16(uint16(v.Uint()))
	case Value4Bytes:
		e.w.WriteUint32(uint32(v.Uint()))
	case Value8Bytes:
		e.w.WriteUint64(v.Uint())
	default:
		return fmt.Errorf("invalid value size: %v", size)
	}

	return e.w.Err
}

func (e *Encoder) encodeSlice(v reflect.Value, typOpts TypeOptions) error {
	switch typOpts.LenBytes {
	case Len2Bytes:
		e.w.WriteSize16(v.Len())
	case Len4Bytes, 0:
		e.w.WriteSize32(v.Len())
	default:
		return fmt.Errorf("invalid collection size type: %v", typOpts.LenBytes)
	}

	for i := 0; i < v.Len(); i++ {
		if err := e.encodeValue(v.Index(i), nil); err != nil {
			return fmt.Errorf("[%v]: %w", i, err)
		}
	}

	return nil
}

func (e *Encoder) encodeArray(v reflect.Value) error {
	for i := 0; i < v.Len(); i++ {
		if err := e.encodeValue(v.Index(i), nil); err != nil {
			return fmt.Errorf("[%v]: %w", i, err)
		}
	}

	return nil
}

func (e *Encoder) encodeMap(v reflect.Value, typOpts TypeOptions) error {
	if v.IsNil() {
		return fmt.Errorf("attemp to encode non-optional nil-map")
	}

	switch typOpts.LenBytes {
	case Len2Bytes:
		e.w.WriteSize16(v.Len())
	case Len4Bytes, 0:
		e.w.WriteSize32(v.Len())
	default:
		return fmt.Errorf("invalid collection size type: %v", typOpts.LenBytes)
	}

	entries := make([]*lo.Entry[reflect.Value, reflect.Value], 0, v.Len())
	for elem := v.MapRange(); elem.Next(); {
		entries = append(entries, &lo.Entry[reflect.Value, reflect.Value]{Key: elem.Key(), Value: elem.Value()})
	}

	// Need to sort map entries to ensure deterministic encoding
	if err := sortMap(entries); err != nil {
		return fmt.Errorf("sorting map: %w", err)
	}

	for i := range entries {
		if err := e.encodeValue(entries[i].Key, nil); err != nil {
			return fmt.Errorf("key: %w", err)
		}

		if err := e.encodeValue(entries[i].Value, nil); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}

	return nil
}

func (e *Encoder) encodeStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)

		fieldOpts, hasTag, err := e.fieldOptsFromTag(fieldType)
		if err != nil {
			return fmt.Errorf("%v: parsing annotation: %w", fieldType.Name, err)
		}

		if fieldOpts.Skip {
			continue
		}

		fieldVal := v.Field(i)

		if !fieldType.IsExported() {
			if !hasTag {
				// Unexported fields without tags are skipped
				continue
			}

			if !fieldVal.CanAddr() {
				// Field is not addresable yet - making it addressable
				addressableV := reflect.New(t).Elem()
				addressableV.Set(v)
				v = addressableV
				fieldVal = v.Field(i)
			}

			// Accesing unexported field
			// Trick to access unexported fields: https://stackoverflow.com/questions/42664837/how-to-access-unexported-struct-fields/43918797#43918797
			fieldVal = reflect.NewAt(fieldVal.Type(), unsafe.Pointer(fieldVal.UnsafeAddr())).Elem()
		}

		fieldKind := fieldVal.Kind()

		if fieldKind == reflect.Ptr || fieldKind == reflect.Interface || fieldKind == reflect.Map {
			// The field is nullable

			isNil := fieldVal.IsNil()

			if isNil && !fieldOpts.Optional {
				return fmt.Errorf("%v: non-optional nil value", fieldType.Name)
			}

			if fieldOpts.Optional {
				e.w.WriteByte(lo.Ternary[byte](isNil, 0, 1))

				if isNil {
					continue
				}
			}
		}

		if err := e.encodeValue(fieldVal, &fieldOpts.TypeOptions); err != nil {
			return fmt.Errorf("%v: %w", fieldType.Name, err)
		}
	}

	return nil
}

func (e *Encoder) encodeEnum(v reflect.Value, variantIdx int) error {
	e.w.WriteSize32(variantIdx)

	if err := e.encodeValue(v, nil); err != nil {
		return fmt.Errorf("%v: %w", v.Type(), err)
	}

	return nil
}

func (e *Encoder) fieldOptsFromTag(fieldType reflect.StructField) (FieldOptions, bool, error) {
	a, hasTag := fieldType.Tag.Lookup(e.cfg.TagName)

	fieldOpts, err := FieldOptionsFromTag(a)
	if err != nil {
		return FieldOptions{}, false, fmt.Errorf("%v: parsing annotation: %w", fieldType.Name, err)
	}

	return fieldOpts, hasTag, nil
}

// func (e *Encoder) Writer() *rwutil.Writer {
// 	return &e.w
// }

func (e *Encoder) Write(b []byte) (n int, err error) {
	e.w.WriteN(b)
	return len(b), e.w.Err
}

func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer

	if err := NewEncoder(&buf, EncoderConfig{}).Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func MustMarshal(v any) []byte {
	b, err := Marshal(v)
	if err != nil {
		panic(fmt.Errorf("failed to marshal object of type %T into BCS: %w", v, err))
	}

	return b
}

type CustomEncoder func(e *Encoder, v reflect.Value) error

var CustomEncoders = make(map[reflect.Type]CustomEncoder)

func MakeCustomEncoder[V any](f func(e *Encoder, v V) error) func(e *Encoder, v reflect.Value) error {
	return func(e *Encoder, v reflect.Value) error {
		return f(e, v.Interface().(V))
	}
}

func AddCustomEncoder[V any](f func(e *Encoder, v V) error) {
	CustomEncoders[reflect.TypeOf(lo.Empty[V]())] = MakeCustomEncoder(f)
}
