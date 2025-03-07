// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscmagic

import (
	"bytes"
	_ "embed"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/hashing"
)

type addressKind uint8

const (
	addressKindISCMagic        = addressKind(iota)
	addressKindERC20BaseTokens // deprecated
	addressKindERC20Coin
	addressKindERC721NFTs
	addressKindERC721NFTCollection
	addressKindInvalid
)

var (
	AddressPrefix = []byte{0x10, 0x74}
	Address       = packMagicAddress(addressKindISCMagic, nil)

	kindByteIndex    = len(AddressPrefix)
	headerLength     = len(AddressPrefix) + 1 // AddressPrefix + kind (byte)
	maxPayloadLength = common.AddressLength - headerLength
)

// ERC20CoinAddress returns the Ethereum address of the ERC20 contract for
// the given coin.
func ERC20CoinAddress(coinType coin.Type) common.Address {
	hash := hashing.HashKeccak([]byte(coinType.String()))
	return packMagicAddress(addressKindERC20Coin, hash[:maxPayloadLength])
}

func ERC721NFTCollectionAddress(collectionID iotago.ObjectID) common.Address {
	return packMagicAddress(addressKindERC721NFTCollection, collectionID[:maxPayloadLength])
}

func packMagicAddress(kind addressKind, payload []byte) common.Address {
	var ret common.Address
	copy(ret[:], AddressPrefix)
	ret[kindByteIndex] = byte(kind)
	if len(payload) > maxPayloadLength {
		panic("packMagicAddress: invalid payload length")
	}
	copy(ret[headerLength:], payload)
	return ret
}

func unpackMagicAddress(addr common.Address) (addressKind, []byte, error) {
	if !bytes.Equal(addr[0:len(AddressPrefix)], AddressPrefix) {
		return 0, nil, errors.New("unpackMagicAddress: expected magic address prefix")
	}
	kind := addressKind(addr[kindByteIndex])
	if kind >= addressKindInvalid {
		return 0, nil, errors.New("unpackMagicAddress: unknown address kind")
	}
	payload := addr[headerLength:]
	return kind, payload, nil
}
