package corecontracts

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/webapi/common"
	"github.com/iotaledger/wasp/sui-go/sui"
)

func GetTotalAssets(ch chain.Chain, blockIndexOrTrieRoot string) (isc.CoinBalances, error) {
	ret, err := common.CallView(ch, accounts.ViewTotalAssets.Message(), blockIndexOrTrieRoot)
	if err != nil {
		return nil, err
	}
	return accounts.ViewTotalAssets.DecodeOutput(ret)
}

func GetAccountBalance(ch chain.Chain, agentID isc.AgentID, blockIndexOrTrieRoot string) (isc.CoinBalances, error) {
	ret, err := common.CallView(ch, accounts.ViewBalance.Message(&agentID), blockIndexOrTrieRoot)
	if err != nil {
		return nil, err
	}
	return accounts.ViewTotalAssets.DecodeOutput(ret)
}

func GetAccountNFTs(ch chain.Chain, agentID isc.AgentID, blockIndexOrTrieRoot string) ([]sui.Address, error) {
	ret, err := common.CallView(ch, accounts.ViewAccountObjects.Message(&agentID), blockIndexOrTrieRoot)
	if err != nil {
		return nil, err
	}
	return accounts.ViewAccountObjects.DecodeOutput(ret)
}

func GetAccountFoundries(ch chain.Chain, agentID isc.AgentID, blockIndexOrTrieRoot string) ([]coin.Type, error) {
	ret, err := common.CallView(ch, accounts.ViewAccountTreasuries.Message(&agentID), blockIndexOrTrieRoot)
	if err != nil {
		return nil, err
	}
	sns, err := accounts.ViewAccountTreasuries.DecodeOutput(ret)
	if err != nil {
		return nil, err
	}
	return sns, nil
}

func GetAccountNonce(ch chain.Chain, agentID isc.AgentID, blockIndexOrTrieRoot string) (uint64, error) {
	ret, err := common.CallView(ch, accounts.ViewGetAccountNonce.Message(&agentID), blockIndexOrTrieRoot)
	if err != nil {
		return 0, err
	}
	return accounts.ViewGetAccountNonce.DecodeOutput(ret)
}

func GetNFTData(ch chain.Chain, nftID iotago.NFTID, blockIndexOrTrieRoot string) (*isc.NFT, error) {
	panic("TODO")
	// ret, err := common.CallView(ch, accounts.ViewNFTData.Message(nftID), blockIndexOrTrieRoot)
	// if err != nil {
	// 	return nil, err
	// }
	// return accounts.ViewNFTData.DecodeOutput(ret)
}

func GetNativeTokenIDRegistry(ch chain.Chain, blockIndexOrTrieRoot string) ([]iotago.NativeTokenID, error) {
	panic("TODO")
	// ret, err := common.CallView(ch, accounts.ViewGetNativeTokenIDRegistry.Message(), blockIndexOrTrieRoot)
	// if err != nil {
	// 	return nil, err
	// }
	// return accounts.ViewGetNativeTokenIDRegistry.DecodeOutput(ret)
}

func GetFoundryOutput(ch chain.Chain, serialNumber uint32, blockIndexOrTrieRoot string) (*iotago.FoundryOutput, error) {
	panic("TODO")
	// ret, err := common.CallView(ch, accounts.ViewNativeToken.Message(serialNumber), blockIndexOrTrieRoot)
	// if err != nil {
	// 	return nil, err
	// }
	// out, err := accounts.ViewNativeToken.DecodeOutput(ret)
	// if err != nil {
	// 	return nil, err
	// }
	// return out.(*iotago.FoundryOutput), nil
}
