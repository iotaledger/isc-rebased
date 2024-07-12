package iscmove

import (
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
)

const (
	AnchorModuleName  = "anchor"
	AnchorObjectName  = "Anchor"
	ReceiptObjectName = "Receipt"

	AssetsBagModuleName = "assets_bag"
	AssetsBagObjectName = "AssetsBag"
	AssetObjectName     = "Asset"

	RequestModuleName      = "request"
	RequestDataObjectName  = "RequestData"
	RequestObjectName      = "Request"
	RequestEventObjectName = "RequestEvent"
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

// Related to: https://github.com/iotaledger/kinesis/tree/isc-suijsonrpc/dapps/isc/sources
// Might change completely: https://github.com/iotaledger/iota/pull/370#discussion_r1617682560
type Allowance struct {
	CoinAmounts []uint64
	CoinTypes   []string
	NFTs        []sui.ObjectID
}

type Referent[T any] struct {
	ID    sui.ObjectID
	Value *T `bcs:"optional"`
}

type AssetsBag struct {
	ID       sui.ObjectID
	Size     uint64
	Balances map[suijsonrpc.CoinType]*suijsonrpc.Balance
}

func NewAssetsBag() *AssetsBag {
	return &AssetsBag{
		Balances: make(map[suijsonrpc.CoinType]*suijsonrpc.Balance),
	}
}

type MoveAssetsBag struct {
	ID   suijsonrpc.MoveUID
	Size *suijsonrpc.BigInt
}

func NewMoveAssetsBag() *MoveAssetsBag {
	return &MoveAssetsBag{Size: suijsonrpc.NewBigInt(0)}
}

type Anchor struct {
	Ref        *sui.ObjectRef
	Assets     Referent[AssetsBag]
	InitParams []byte
	StateRoot  sui.Bytes
	StateIndex uint32
}

type anchorJsonObject struct {
	Assets struct {
		Type   string `json:"type"`
		Fields struct {
			ID    *sui.ObjectID `json:"id"`
			Value struct {
				Type   string `json:"type"`
				Fields struct {
					ID   suijsonrpc.MoveUID `json:"id"`
					Size suijsonrpc.BigInt  `json:"size"`
				} `json:"fields"`
			} `json:"value"`
		} `json:"fields"`
	} `json:"assets"`
	ID         suijsonrpc.MoveUID `json:"id"`
	StateIndex uint64             `json:"state_index"`
	StateRoot  []byte             `json:"state_root"`
}

type Receipt struct {
	RequestID sui.ObjectID
}

type RequestData struct {
	Contract string // TODO: should be isc.Hname
	Function string
	Args     []sui.Bytes
}

type Request struct {
	ID        *sui.ObjectID
	Sender    sui.Address
	AssetsBag Referent[AssetsBag] // Need to decide if we want to use this Referent wrapper as well. Could probably be of *AssetsBag with `bcs:"optional`
	Data      *RequestData        `bcs:"optional"`
}

type requestJsonObject struct {
	AssetsBag struct {
		Type   string `json:"type"`
		Fields struct {
			ID    *sui.ObjectID `json:"id"`
			Value struct {
				Type   string `json:"type"`
				Fields struct {
					ID   *suijsonrpc.MoveUID `json:"id"`
					Size suijsonrpc.BigInt   `json:"size"`
				} `json:"fields"`
			} `json:"value"`
		} `json:"fields"`
	} `json:"assets_bag"`
	Data struct {
		Type   string `json:"type"`
		Fields struct {
			Args     [][]byte `json:"args"`
			Contract string   `json:"contract"`
			Function string   `json:"function"`
		} `json:"fields"`
	} `json:"data"`
	ID     suijsonrpc.MoveUID `json:"id"`
	Sender *sui.Address       `json:"sender"`
}

// Related to: https://github.com/iotaledger/kinesis/blob/isc-suijsonrpc/crates/sui-framework/packages/stardust/sources/nft/irc27.move
type IRC27MetaData struct {
	Version           string
	MediaType         string
	URI               string // Actually of type "Url" in SUI -> Create proper type?
	Name              string
	CollectionName    *string `bcs:"optional"`
	Royalties         Table[sui.Address, uint32]
	IssuerName        *string  `bcs:"optional"`
	Description       *string  `bcs:"optional"`
	Attributes        []string // This is actually of Type VecSet which guarantees no duplicates. Not sure if we want to create a separate type for it. But we need to filter it to ensure no duplicates eventually.
	NonStandardFields Table[string, string]
}

// Related to: https://github.com/iotaledger/kinesis/blob/isc-suijsonrpc/crates/sui-framework/packages/stardust/sources/nft/nft.move

type NFT struct {
	ID                sui.ObjectID
	LegacySender      *sui.Address `bcs:"optional"`
	Metadata          *[]uint8     `bcs:"optional"`
	Tag               *[]uint8     `bcs:"optional"`
	ImmutableIssuer   *sui.Address `bcs:"optional"`
	ImmutableMetadata IRC27MetaData
}
