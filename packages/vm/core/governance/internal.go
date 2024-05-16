// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package governance

import (
	"math/big"

	"github.com/samber/lo"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

func (s *StateWriter) SetInitialState(chainOwner isc.AgentID, blockKeepAmount int32) {
	s.SetChainOwnerID(chainOwner)
	s.SetGasFeePolicy(gas.DefaultFeePolicy())
	s.SetGasLimits(gas.LimitsDefault)
	s.SetMaintenanceStatus(false)
	s.SetBlockKeepAmount(blockKeepAmount)
	s.SetMinCommonAccountBalance(DefaultMinBaseTokensOnCommonAccount)
	s.SetPayoutAgentID(chainOwner)
}

// GetRotationAddress tries to read the state of 'governance' and extract rotation address
// If succeeds, it means this block is fake.
// If fails, return nil
func (s *StateReader) GetRotationAddress() iotago.Address {
	ret, err := codec.Address.Decode(s.state.Get(varRotateToAddress), nil)
	if err != nil {
		return nil
	}
	return ret
}

func (s *StateWriter) SetRotationAddress(a iotago.Address) {
	s.state.Set(varRotateToAddress, codec.Address.Encode(a))
}

// GetChainInfo returns global variables of the chain
func (s *StateReader) GetChainInfo(chainID isc.ChainID) *isc.ChainInfo {
	ret := &isc.ChainInfo{
		ChainID:  chainID,
		Metadata: &isc.PublicChainMetadata{},
	}
	ret.ChainOwnerID = s.GetChainOwnerID()
	ret.GasFeePolicy = s.GetGasFeePolicy()
	ret.GasLimits = s.GetGasLimits()
	ret.BlockKeepAmount = s.GetBlockKeepAmount()
	ret.PublicURL = s.GetPublicURL()
	ret.Metadata = s.GetMetadata()
	return ret
}

func (s *StateReader) GetMinCommonAccountBalance() uint64 {
	return lo.Must(codec.Uint64.Decode(s.state.Get(varMinBaseTokensOnCommonAccount)))
}

func (s *StateWriter) SetMinCommonAccountBalance(m uint64) {
	s.state.Set(varMinBaseTokensOnCommonAccount, codec.Uint64.Encode(m))
}

func (s *StateReader) GetChainOwnerID() isc.AgentID {
	return lo.Must(codec.AgentID.Decode(s.state.Get(varChainOwnerID)))
}

func (s *StateWriter) SetChainOwnerID(a isc.AgentID) {
	s.state.Set(varChainOwnerID, codec.AgentID.Encode(a))
	if s.GetChainOwnerIDDelegated() != nil {
		s.state.Del(varChainOwnerIDDelegated)
	}
}

func (s *StateReader) GetChainOwnerIDDelegated() isc.AgentID {
	return lo.Must(codec.AgentID.Decode(s.state.Get(varChainOwnerIDDelegated), nil))
}

func (s *StateWriter) SetChainOwnerIDDelegated(a isc.AgentID) {
	s.state.Set(varChainOwnerIDDelegated, codec.AgentID.Encode(a))
}

func (s *StateReader) GetPayoutAgentID() isc.AgentID {
	return lo.Must(codec.AgentID.Decode(s.state.Get(varPayoutAgentID)))
}

func (s *StateWriter) SetPayoutAgentID(a isc.AgentID) {
	s.state.Set(varPayoutAgentID, codec.AgentID.Encode(a))
}

func (s *StateReader) GetGasFeePolicy() *gas.FeePolicy {
	return lo.Must(gas.FeePolicyFromBytes(s.state.Get(varGasFeePolicyBytes)))
}

func (s *StateWriter) SetGasFeePolicy(fp *gas.FeePolicy) {
	s.state.Set(varGasFeePolicyBytes, fp.Bytes())
}

func (s *StateReader) GetDefaultGasPrice() *big.Int {
	return s.GetGasFeePolicy().DefaultGasPriceFullDecimals(parameters.L1().BaseToken.Decimals)
}

func (s *StateReader) GetGasLimits() *gas.Limits {
	data := s.state.Get(varGasLimitsBytes)
	if data == nil {
		return gas.LimitsDefault
	}
	return lo.Must(gas.LimitsFromBytes(data))
}

func (s *StateWriter) SetGasLimits(gl *gas.Limits) {
	s.state.Set(varGasLimitsBytes, gl.Bytes())
}

func (s *StateReader) GetBlockKeepAmount() int32 {
	return lo.Must(codec.Int32.Decode(s.state.Get(varBlockKeepAmount), DefaultBlockKeepAmount))
}

func (s *StateWriter) SetBlockKeepAmount(n int32) {
	s.state.Set(varBlockKeepAmount, codec.Int32.Encode(n))
}

func (s *StateWriter) SetPublicURL(url string) {
	s.state.Set(varPublicURL, codec.String.Encode(url))
}

func (s *StateReader) GetPublicURL() string {
	return codec.String.MustDecode(s.state.Get(varPublicURL), "")
}

func (s *StateWriter) SetMetadata(metadata *isc.PublicChainMetadata) {
	s.state.Set(varMetadata, metadata.Bytes())
}

func (s *StateReader) GetMetadata() *isc.PublicChainMetadata {
	metadataBytes := s.state.Get(varMetadata)
	if metadataBytes == nil {
		return &isc.PublicChainMetadata{}
	}
	return lo.Must(isc.PublicChainMetadataFromBytes(metadataBytes))
}

func (s *StateWriter) AccessNodesMap() *collections.Map {
	return collections.NewMap(s.state, varAccessNodes)
}

func (s *StateReader) AccessNodesMap() *collections.ImmutableMap {
	return collections.NewMapReadOnly(s.state, varAccessNodes)
}

func (s *StateWriter) AccessNodeCandidatesMap() *collections.Map {
	return collections.NewMap(s.state, varAccessNodeCandidates)
}

func (s *StateReader) AccessNodeCandidatesMap() *collections.ImmutableMap {
	return collections.NewMapReadOnly(s.state, varAccessNodeCandidates)
}

func (s *StateWriter) AllowedStateControllerAddressesMap() *collections.Map {
	return collections.NewMap(s.state, varAllowedStateControllerAddresses)
}

func (s *StateReader) AllowedStateControllerAddressesMap() *collections.ImmutableMap {
	return collections.NewMapReadOnly(s.state, varAllowedStateControllerAddresses)
}

func (s *StateReader) GetMaintenanceStatus() bool {
	r := s.state.Get(varMaintenanceStatus)
	if r == nil {
		return false // chain is being initialized, governance has not been initialized yet
	}
	return lo.Must(codec.Bool.Decode(r))
}

func (s *StateWriter) SetMaintenanceStatus(status bool) {
	s.state.Set(varMaintenanceStatus, codec.Bool.Encode(status))
}

func (s *StateReader) AccessNodes() []*cryptolib.PublicKey {
	accessNodes := []*cryptolib.PublicKey{}
	s.AccessNodesMap().IterateKeys(func(pubKeyBytes []byte) bool {
		pubKey, err := cryptolib.PublicKeyFromBytes(pubKeyBytes)
		if err != nil {
			panic(err)
		}
		accessNodes = append(accessNodes, pubKey)
		return true
	})
	return accessNodes
}

func (s *StateReader) CandidateNodes() []*AccessNodeInfo {
	candidateNodes := []*AccessNodeInfo{}
	s.AccessNodeCandidatesMap().Iterate(func(pubKeyBytes, accessNodeInfoBytes []byte) bool {
		ani, err := AccessNodeInfoFromBytes(pubKeyBytes, accessNodeInfoBytes)
		if err != nil {
			panic(err)
		}
		candidateNodes = append(candidateNodes, ani)
		return true
	})
	return candidateNodes
}
