package main

import (
	"fmt"
	"math/big"

	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	old_isc "github.com/nnikolash/wasp-types-exported/packages/isc"
	old_codec "github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	old_util "github.com/nnikolash/wasp-types-exported/packages/util"
)

func OldChainIDToNewChainID(oldChainID old_isc.ChainID) isc.ChainID {
	return isc.ChainID(oldChainID)
}

func OldHnameToNewHname(oldHname old_isc.Hname) isc.Hname {
	return isc.Hname(oldHname)
}

func OldAgentIDtoNewAgentID(oldAgentID old_isc.AgentID) isc.AgentID {
	switch oldAgentID.Kind() {
	case old_isc.AgentIDKindAddress:
		oldAddr := oldAgentID.(*old_isc.AddressAgentID).Address()
		newAdd := iotago.MustAddressFromHex(oldAddr.String())
		return isc.NewAddressAgentID(cryptolib.NewAddressFromIota(newAdd))

	case old_isc.AgentIDKindContract:
		oldAgentID := oldAgentID.(*old_isc.ContractAgentID)
		chID := OldChainIDToNewChainID(oldAgentID.ChainID())
		hname := OldHnameToNewHname(oldAgentID.Hname())
		return isc.NewContractAgentID(chID, hname)

	case old_isc.AgentIDKindEthereumAddress:
		oldAgentID := oldAgentID.(*old_isc.EthereumAddressAgentID)
		chID := OldChainIDToNewChainID(oldAgentID.ChainID())
		ethAddr := oldAgentID.EthAddress()
		return isc.NewEthereumAddressAgentID(chID, ethAddr)

	case old_isc.AgentIDIsNil:
		panic(fmt.Sprintf("Found agent ID with kind = AgentIDIsNil: %v", oldAgentID))

	case old_isc.AgentIDKindNil:
		panic(fmt.Sprintf("Found agent ID with kind = AgentIDKindNil: %v", oldAgentID))

	default:
		panic(fmt.Sprintf("Unknown agent ID kind: %v = %v", oldAgentID.Kind(), oldAgentID))
	}
}

func OldTokensCountToNewCoinValue(oldTokensCount uint64) coin.Value {
	// TODO: what is the conversion rate?
	return coin.Value(oldTokensCount)
}

func DecodeOldTokens(b []byte) uint64 {
	amount := old_codec.MustDecodeBigIntAbs(b, big.NewInt(0))
	convertedAmount, _ := old_util.EthereumDecimalsToBaseTokenDecimals(amount, oldBaseTokenDecimals)

	return convertedAmount
}
