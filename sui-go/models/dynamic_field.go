package models

import (
	"github.com/iotaledger/isc-private/sui-go/sui_types"
	"github.com/iotaledger/isc-private/sui-go/sui_types/serialization"
)

type DynamicFieldInfo struct {
	Name sui_types.DynamicFieldName `json:"name"`
	//Base58
	BcsName    sui_types.Base58                                  `json:"bcsName"`
	Type       serialization.TagJson[sui_types.DynamicFieldType] `json:"type"`
	ObjectType string                                            `json:"objectType"`
	ObjectID   sui_types.ObjectID                                `json:"objectId"`
	Version    sui_types.SequenceNumber                          `json:"version"`
	Digest     sui_types.ObjectDigest                            `json:"digest"`
}

type DynamicFieldPage = Page[DynamicFieldInfo, sui_types.ObjectID]
