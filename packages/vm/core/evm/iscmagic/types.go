// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscmagic

import (
	"math/big"
	"time"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/sui-go/sui"
)

// ISCChainID matches the type definition in ISCTypes.sol
type ISCChainID [isc.ChainIDLength]byte

func init() {
	if isc.ChainIDLength != 32 {
		panic("static check: ChainID length does not match bytes32 in ISCTypes.sol")
	}
}

func WrapISCChainID(c isc.ChainID) (ret ISCChainID) {
	copy(ret[:], c.Bytes())
	return
}

func (c ISCChainID) Unwrap() (isc.ChainID, error) {
	return isc.ChainIDFromBytes(c[:])
}

func (c ISCChainID) MustUnwrap() isc.ChainID {
	ret, err := c.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// NativeTokenID matches the struct definition in ISCTypes.sol
type CoinType struct {
	Data string
}

func WrapCoinType(coinType isc.CoinType) CoinType {
	return CoinType{Data: coinType.String()}
}

func (a CoinType) Unwrap() (ret isc.CoinType) {
	ret = isc.CoinType(a.Data)
	return
}

func (a CoinType) MustUnwrap() (ret isc.CoinType) {
	ret = isc.CoinType(a.Data)
	return
}

type CoinBalances map[CoinType]*big.Int

func WrapCoinBalances(coinBalances isc.CoinBalances) CoinBalances {
	newCoinBalances := make(CoinBalances)

	for k, v := range coinBalances {
		newCoinBalances[WrapCoinType(k)] = v
	}

	return newCoinBalances
}

func (c CoinBalances) Unwrap() (isc.CoinBalances, error) {
	newCoinBalances := make(isc.CoinBalances)

	for k, v := range c {
		newCoinBalances[k.Unwrap()] = v
	}

	return newCoinBalances, nil
}

type ObjectIDSet map[ObjectID]struct{}

func WrapObjectIDSet(set isc.ObjectIDSet) ObjectIDSet {
	objectIDSet := make(ObjectIDSet)

	for k, v := range set {
		objectIDSet[WrapObjectID(k)] = v
	}

	return objectIDSet
}

func (c ObjectIDSet) Unwrap() (isc.ObjectIDSet, error) {
	objectIDSet := make(isc.ObjectIDSet)

	for k, v := range c {
		objectIDSet[k.Unwrap()] = v
	}

	return objectIDSet, nil
}

// TODO: refactor me: refactor or remove?
/*
// NativeToken matches the struct definition in ISCTypes.sol
type NativeToken struct {
	ID     CoinType
	Amount *big.Int
}

func WrapNativeToken(nativeToken *isc.NativeToken) NativeToken {
	return NativeToken{
		ID:     WrapCoinType(nativeToken.ID),
		Amount: nativeToken.Amount,
	}
}

func (nt NativeToken) Unwrap() *isc.NativeToken {
	return &isc.NativeToken{
		ID:     nt.ID.Unwrap(),
		Amount: nt.Amount,
	}
}*/

// L1Address matches the struct definition in ISCTypes.sol
type L1Address struct {
	Data []byte
}

func WrapL1Address(a *cryptolib.Address) L1Address {
	if a == nil {
		return L1Address{Data: []byte{}}
	}
	return L1Address{Data: a.Bytes()}
}

func (a L1Address) Unwrap() (*cryptolib.Address, error) {
	return cryptolib.NewAddressFromBytes(a.Data)
}

func (a L1Address) MustUnwrap() *cryptolib.Address {
	ret, err := a.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// ISCAgentID matches the struct definition in ISCTypes.sol
type ISCAgentID struct {
	Data []byte
}

func WrapISCAgentID(a isc.AgentID) ISCAgentID {
	return ISCAgentID{Data: a.Bytes()}
}

func (a ISCAgentID) Unwrap() (isc.AgentID, error) {
	return isc.AgentIDFromBytes(a.Data)
}

func (a ISCAgentID) MustUnwrap() isc.AgentID {
	ret, err := a.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// ISCRequestID matches the struct definition in ISCTypes.sol
type ISCRequestID struct {
	Data []byte
}

func WrapISCRequestID(rid isc.RequestID) ISCRequestID {
	return ISCRequestID{Data: rid.Bytes()}
}

func (rid ISCRequestID) Unwrap() (isc.RequestID, error) {
	return isc.RequestIDFromBytes(rid.Data)
}

func (rid ISCRequestID) MustUnwrap() isc.RequestID {
	ret, err := rid.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// ObjectID matches the type definition in ISCTypes.sol
type ObjectID [sui.AddressLen]byte

func init() {
	if sui.AddressLen != 32 {
		panic("static check: ObjectID length does not match bytes32 in ISCTypes.sol")
	}
}

func WrapObjectID(c sui.ObjectID) (ret ObjectID) {
	copy(ret[:], c[:])
	return
}

func (n ObjectID) Unwrap() (ret sui.ObjectID) {
	copy(ret[:], n[:])
	return
}

// TokenID returns the uint256 tokenID for ERC721
func (n ObjectID) TokenID() *big.Int {
	return new(big.Int).SetBytes(n[:])
}

// ISCNFT matches the struct definition in ISCTypes.sol
type ISCNFT struct {
	ID       ObjectID
	Issuer   L1Address
	Metadata []byte
	Owner    ISCAgentID
}

func WrapISCNFT(n *isc.NFT) ISCNFT {
	r := ISCNFT{
		ID:       WrapObjectID(n.ID),
		Issuer:   WrapL1Address(n.Issuer),
		Metadata: n.Metadata,
	}
	if n.Owner != nil {
		r.Owner = WrapISCAgentID(n.Owner)
	}
	return r
}

func (n ISCNFT) Unwrap() (*isc.NFT, error) {
	issuer, err := n.Issuer.Unwrap()
	if err != nil {
		return nil, err
	}
	return &isc.NFT{
		ID:       n.ID.Unwrap(),
		Issuer:   issuer,
		Metadata: n.Metadata,
		Owner:    n.Owner.MustUnwrap(),
	}, nil
}

func (n ISCNFT) MustUnwrap() *isc.NFT {
	ret, err := n.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// IRC27NFTMetadata matches the struct definition in ISCTypes.sol
type IRC27NFTMetadata struct {
	Standard    string
	Version     string
	MimeType    string
	Uri         string //nolint:revive // false positive
	Name        string
	Description string
}

func WrapIRC27NFTMetadata(m *isc.IRC27NFTMetadata) IRC27NFTMetadata {
	return IRC27NFTMetadata{
		Standard:    m.Standard,
		Version:     m.Version,
		MimeType:    m.MIMEType,
		Uri:         m.URI,
		Name:        m.Name,
		Description: m.Description,
	}
}

// IRC27NFT matches the struct definition in ISCTypes.sol
type IRC27NFT struct {
	Nft      ISCNFT
	Metadata IRC27NFTMetadata
}

// ISCAssets matches the struct definition in ISCTypes.sol
type ISCAssets struct {
	Coins   CoinBalances
	Objects ObjectIDSet
}

func WrapISCAssets(a *isc.Assets) ISCAssets {
	if a == nil {
		return WrapISCAssets(isc.NewEmptyAssets())
	}

	coins := CoinBalances{}
	for i, value := range a.Coins {
		coins[WrapCoinType(i)] = value
	}

	objects := ObjectIDSet{}
	for i, value := range a.Objects {
		objects[WrapObjectID(i)] = value
	}

	return ISCAssets{
		Coins:   coins,
		Objects: objects,
	}
}

func (a ISCAssets) BaseToken() uint64 {
	if val, ok := a.Coins[WrapCoinType(isc.BaseTokenType)]; ok {
		panic("refactor me: return val when changed to Uint64")
		_ = val
		return 0
	}

	return 0
}

func (a ISCAssets) Unwrap() *isc.Assets {
	coins := isc.CoinBalances{}
	for k, v := range a.Coins {
		coins[k.Unwrap()] = v
	}

	objects := isc.ObjectIDSet{}
	for k, v := range a.Objects {
		objects[k.Unwrap()] = v
	}

	panic("refactor me: set base token here when switched to Uint64")
	return &isc.Assets{

		Coins:   coins,
		Objects: objects,
	}
}

// ISCDictItem matches the struct definition in ISCTypes.sol
type ISCDictItem struct {
	Key   []byte
	Value []byte
}

// ISCDict matches the struct definition in ISCTypes.sol
type ISCDict struct {
	Items []ISCDictItem
}

func WrapISCDict(d dict.Dict) ISCDict {
	items := make([]ISCDictItem, 0, len(d))
	for k, v := range d {
		items = append(items, ISCDictItem{Key: []byte(k), Value: v})
	}
	return ISCDict{Items: items}
}

func (d ISCDict) Unwrap() dict.Dict {
	ret := dict.Dict{}
	for _, item := range d.Items {
		ret[kv.Key(item.Key)] = item.Value
	}
	return ret
}

type CallArguments struct {
	Data [][]byte
}

func WrapCallArguments(a isc.CallArguments) CallArguments {
	return CallArguments{
		Data: a[:],
	}
}

func (c CallArguments) Unwrap() isc.CallArguments {
	args := isc.CallArguments{}
	copy(args, c.Data)
	return args
}

type ISCSendMetadata struct {
	TargetContract uint32
	Entrypoint     uint32
	Params         CallArguments
	Allowance      ISCAssets
	GasBudget      uint64
}

func WrapISCSendMetadata(metadata isc.SendMetadata) ISCSendMetadata {
	ret := ISCSendMetadata{
		GasBudget:      metadata.GasBudget,
		TargetContract: uint32(metadata.Message.Target.Contract),
		Entrypoint:     uint32(metadata.Message.Target.EntryPoint),
		Params:         WrapCallArguments(metadata.Message.Params),
		Allowance:      WrapISCAssets(metadata.Allowance),
	}

	return ret
}

func (i ISCSendMetadata) Unwrap() *isc.SendMetadata {
	ret := isc.SendMetadata{
		Message: isc.NewMessage(
			isc.Hname(i.TargetContract),
			isc.Hname(i.Entrypoint),
			i.Params.Unwrap(),
		),
		Allowance: i.Allowance.Unwrap(),
		GasBudget: i.GasBudget,
	}

	return &ret
}

type ISCExpiration struct {
	Time          int64
	ReturnAddress L1Address
}

func (i *ISCExpiration) Unwrap() *isc.Expiration {
	if i == nil {
		return nil
	}

	if i.Time == 0 {
		return nil
	}

	address := i.ReturnAddress.MustUnwrap()

	ret := isc.Expiration{
		ReturnAddress: address,
		Time:          time.UnixMilli(i.Time),
	}

	return &ret
}

type ISCSendOptions struct {
	Timelock   int64
	Expiration ISCExpiration
}

func (i *ISCSendOptions) Unwrap() isc.SendOptions {
	var timeLock time.Time

	if i.Timelock > 0 {
		timeLock = time.UnixMilli(i.Timelock)
	}

	ret := isc.SendOptions{
		Timelock:   timeLock,
		Expiration: i.Expiration.Unwrap(),
	}

	return ret
}

type ISCTokenProperties struct {
	Name         string
	TickerSymbol string
	Decimals     uint8
	TotalSupply  *big.Int
}
