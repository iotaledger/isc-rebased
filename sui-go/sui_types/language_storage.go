package sui_types

import (
	"fmt"
	"strings"

	"github.com/iotaledger/wasp/sui-go/sui_types/serialization"
)

type StructTag struct {
	Address    SuiAddress
	Module     Identifier
	Name       Identifier
	TypeParams []TypeTag
}

func (s *StructTag) UnmarshalJSON(data []byte) error {
	// Split the string based on "::"
	parts := strings.Split(string(data), "::")
	if len(parts) != 3 {
		return fmt.Errorf("invalid StructTag format: %s", string(data))
	}

	// FIXME
	s.Address = *MustSuiAddressFromHex(parts[0])
	s.Module = Identifier(parts[1])
	s.Name = Identifier(parts[2])

	// FIXME TypeParams is ignored temporarily
	return nil
}

// refer BCS doc https://github.com/diem/bcs/blob/master/README.md#externally-tagged-enumerations
// IMPORTANT! The order of the fields MATTERS! DON'T CHANGE!
// this is enum `TypeTag` in `external-crates/move/crates/move-core-types/src/language_storage.rs`
// each field should be the same as enum `TypeTag` there
type TypeTag struct {
	Bool    *serialization.EmptyEnum
	U8      *serialization.EmptyEnum
	U64     *serialization.EmptyEnum
	U128    *serialization.EmptyEnum
	Address *serialization.EmptyEnum
	Signer  *serialization.EmptyEnum
	Vector  *TypeTag
	Struct  *StructTag

	U16  *serialization.EmptyEnum
	U32  *serialization.EmptyEnum
	U256 *serialization.EmptyEnum
}

func (t TypeTag) IsBcsEnum() {}
