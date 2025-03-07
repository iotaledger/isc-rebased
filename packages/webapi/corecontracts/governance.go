package corecontracts

import (
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/webapi/common"
)

func GetAllowedStateControllerAddresses(ch chain.Chain, blockIndexOrTrieRoot string) ([]*cryptolib.Address, error) {
	ret, err := common.CallView(ch, governance.ViewGetAllowedStateControllerAddresses.Message(), blockIndexOrTrieRoot)
	if err != nil {
		return nil, err
	}

	return governance.ViewGetAllowedStateControllerAddresses.DecodeOutput(ret)
}

func GetChainOwner(ch chain.Chain, blockIndexOrTrieRoot string) (isc.AgentID, error) {
	ret, err := common.CallView(ch, governance.ViewGetChainOwner.Message(), blockIndexOrTrieRoot)
	if err != nil {
		return nil, err
	}
	return governance.ViewGetChainOwner.DecodeOutput(ret)
}

func GetChainInfo(ch chain.Chain, blockIndexOrTrieRoot string) (*isc.ChainInfo, error) {
	ret, err := common.CallView(ch, governance.ViewGetChainInfo.Message(), blockIndexOrTrieRoot)
	if err != nil {
		return nil, err
	}
	return governance.ViewGetChainInfo.DecodeOutput(ret)
}
