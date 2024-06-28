package iscmove

import (
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/sui-go/sui_types"
)

/*
  Structs are currently set up with the following assumptions:
  * Option<T> types are nullable types => (*T)
    => Option<T> may require the `bcs:"optional"` tag.

  * "ID" and "UID" are for now both typed as ObjectID, the actual typing maybe needs to be reconsidered. On our end it maybe not make a difference.
  * Type "Bag" is not available as a proper type in Sui-Go. It needs to be considered if we will need this, as
	=> Bag is a heterogeneous map, so it can hold key-value pairs of arbitrary types (map[any]any)
  * Type "Table" is a typed map: map[K]V
*/

// Related to: https://github.com/iotaledger/kinesis/tree/isc-models/dapps/isc/sources
// Might change completely: https://github.com/iotaledger/iota/pull/370#discussion_r1617682560
type Allowance struct {
	CoinAmounts []uint64
	CoinTypes   []string
	NFTs        []sui_types.ObjectID
}

type Referent[T any] struct {
	ID    sui_types.ObjectID
	Value *T `bcs:"optional"`
}

type AssetBag struct {
	ID   sui_types.ObjectID
	Size uint64
}

type Anchor struct {
	ID         sui_types.ObjectID
	Assets     Referent[AssetBag]
	StateRoot  sui_types.Bytes
	StateIndex uint32
}

type Receipt struct {
	RequestID sui_types.ObjectID
}

type RequestData struct {
	Contract isc.Hname
	Function isc.Hname
	Args     [][]uint8
}

type Request struct {
	ID        sui_types.ObjectID
	Sender    sui_types.SuiAddress
	AssetsBag Referent[AssetBag] // Need to decide if we want to use this Referent wrapper as well. Could probably be of *AssetBag with `bcs:"optional`
	Data      *RequestData       `bcs:"optional"`
}

// Related to: https://github.com/iotaledger/kinesis/blob/isc-models/crates/sui-framework/packages/stardust/sources/nft/irc27.move
type IRC27MetaData struct {
	Version           string
	MediaType         string
	URI               string // Actually of type "Url" in SUI -> Create proper type?
	Name              string
	CollectionName    *string `bcs:"optional"`
	Royalties         Table[sui_types.SuiAddress, uint32]
	IssuerName        *string  `bcs:"optional"`
	Description       *string  `bcs:"optional"`
	Attributes        []string // This is actually of Type VecSet which guarantees no duplicates. Not sure if we want to create a separate type for it. But we need to filter it to ensure no duplicates eventually.
	NonStandardFields Table[string, string]
}

// Related to: https://github.com/iotaledger/kinesis/blob/isc-models/crates/sui-framework/packages/stardust/sources/nft/nft.move

type NFT struct {
	ID                sui_types.ObjectID
	LegacySender      *sui_types.SuiAddress `bcs:"optional"`
	Metadata          *[]uint8              `bcs:"optional"`
	Tag               *[]uint8              `bcs:"optional"`
	ImmutableIssuer   *sui_types.SuiAddress `bcs:"optional"`
	ImmutableMetadata IRC27MetaData
}
