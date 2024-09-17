package isc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/big"
	"slices"
	"strings"

	"github.com/samber/lo"

	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/util/rwutil"
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
)

type CoinBalances map[coin.Type]coin.Value

func NewCoinBalances() CoinBalances {
	return make(CoinBalances)
}

func (c CoinBalances) ToDict() dict.Dict {
	ret := dict.New()
	for coinType, amount := range c {
		ret.Set(kv.Key(coinType.Bytes()), amount.Bytes())
	}
	return ret
}

func CoinBalancesFromDict(d dict.Dict) (CoinBalances, error) {
	ret := NewCoinBalances()
	for key, val := range d {
		coinType, err := coin.TypeFromBytes([]byte(key))
		if err != nil {
			return nil, fmt.Errorf("CoinBalancesFromDict: %w", err)
		}
		coinValue, err := coin.ValueFromBytes(val)
		if err != nil {
			return nil, fmt.Errorf("CoinBalancesFromDict: %w", err)
		}
		ret.Add(coinType, coinValue)
	}
	return ret, nil
}

func (c CoinBalances) IterateSorted(f func(coin.Type, coin.Value) bool) {
	types := lo.Keys(c)
	slices.Sort(types)
	for _, coinType := range types {
		if !f(coinType, c[coinType]) {
			return
		}
	}
}

func (c *CoinBalances) Read(r io.Reader) error {
	*c = NewCoinBalances()
	rr := rwutil.NewReader(r)
	n := rr.ReadSize32()
	for i := 0; i < n; i++ {
		var coinType coin.Type
		var coinValue coin.Value
		rr.Read(&coinType)
		rr.Read(&coinValue)
		c.Add(coinType, coinValue)
	}
	return rr.Err
}

func (c CoinBalances) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteSize32(len(c))
	c.IterateSorted(func(t coin.Type, v coin.Value) bool {
		ww.Write(t)
		ww.Write(v)
		return true
	})
	return ww.Err
}

func (c CoinBalances) Bytes() []byte {
	return rwutil.WriteToBytes(c)
}

func CoinBalancesFromBytes(b []byte) (CoinBalances, error) {
	var r CoinBalances
	_, err := rwutil.ReadFromBytes(b, &r)
	return r, err
}

func (c CoinBalances) Add(coinType coin.Type, amount coin.Value) CoinBalances {
	if amount == 0 {
		return c
	}
	c[coinType] = c.Get(coinType) + amount
	return c
}

func (c CoinBalances) Set(coinType coin.Type, amount coin.Value) CoinBalances {
	if amount == 0 {
		delete(c, coinType)
		return c
	}
	c[coinType] = amount
	return c
}

func (c CoinBalances) AddBaseTokens(amount coin.Value) CoinBalances {
	return c.Add(coin.BaseTokenType, amount)
}

func (c CoinBalances) Sub(coinType coin.Type, amount coin.Value) CoinBalances {
	v := c.Get(coinType)
	switch {
	case v < amount:
		panic("negative coin balance")
	case v == amount:
		delete(c, coinType)
	default:
		c[coinType] = v - amount
	}
	return c
}

func (c CoinBalances) ToAssets() *Assets {
	return &Assets{
		Coins:   c,
		Objects: NewObjectIDSet(),
	}
}

func (c CoinBalances) Get(coinType coin.Type) coin.Value {
	return c[coinType]
}

func (c CoinBalances) BaseTokens() coin.Value {
	return c[coin.BaseTokenType]
}

func (c CoinBalances) IsEmpty() bool {
	return len(c) == 0
}

type CoinJSON struct {
	CoinType coin.Type          `json:"coinType" swagger:"required"`
	Balance  *suijsonrpc.BigInt `json:"balance" swagger:"required"`
}

func (c *CoinBalances) UnmarshalJSON(b []byte) error {
	var coins []CoinJSON
	err := json.Unmarshal(b, &coins)
	if err != nil {
		return err
	}
	*c = NewCoinBalances()
	for _, cc := range coins {
		c.Add(cc.CoinType, coin.Value(cc.Balance.Int.Uint64()))
	}
	return nil
}

func (c CoinBalances) MarshalJSON() ([]byte, error) {
	var coins []CoinJSON
	c.IterateSorted(func(t coin.Type, v coin.Value) bool {
		coins = append(coins, CoinJSON{
			CoinType: t,
			Balance:  &suijsonrpc.BigInt{Int: new(big.Int).SetUint64(uint64(v))},
		})
		return true
	})
	return json.Marshal(coins)
}

func (c CoinBalances) Equals(b CoinBalances) bool {
	if len(c) != len(b) {
		return false
	}
	for coinType, amount := range c {
		bal := b[coinType]
		if bal != amount {
			return false
		}
	}
	return true
}

func (c CoinBalances) String() string {
	s := lo.MapToSlice(c, func(coinType coin.Type, amount coin.Value) string {
		return fmt.Sprintf("%s: %d", coinType, amount)
	})
	return fmt.Sprintf("CoinBalances{%s}", strings.Join(s, ", "))
}

func (c CoinBalances) Clone() CoinBalances {
	r := NewCoinBalances()
	for coinType, amount := range c {
		r.Add(coinType, amount)
	}
	return r
}

type ObjectIDSet map[sui.ObjectID]struct{}

func NewObjectIDSet() ObjectIDSet {
	return make(map[sui.ObjectID]struct{})
}

func NewObjectIDSetFromArray(ids []sui.ObjectID) ObjectIDSet {
	set := NewObjectIDSet()

	for _, id := range ids {
		set.Add(id)
	}

	return set
}

func (o ObjectIDSet) Add(id sui.ObjectID) {
	o[id] = struct{}{}
}

func (o ObjectIDSet) Has(id sui.ObjectID) bool {
	_, ok := o[id]
	return ok
}

func (o ObjectIDSet) Sorted() []sui.ObjectID {
	ids := lo.Keys(o)
	slices.SortFunc(ids, func(a, b sui.ObjectID) int { return bytes.Compare(a[:], b[:]) })
	return ids
}

func (o ObjectIDSet) IterateSorted(f func(sui.ObjectID) bool) {
	for _, id := range o.Sorted() {
		if !f(id) {
			return
		}
	}
}

func (o *ObjectIDSet) Read(r io.Reader) error {
	*o = NewObjectIDSet()
	rr := rwutil.NewReader(r)
	n := rr.ReadSize32()
	for i := 0; i < n; i++ {
		var id sui.ObjectID
		rr.ReadN(id[:])
		o.Add(id)
	}
	return rr.Err
}

func (o ObjectIDSet) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteSize32(len(o))
	o.IterateSorted(func(id sui.ObjectID) bool {
		ww.WriteN(id[:])
		return true
	})
	return ww.Err
}

func (o *ObjectIDSet) UnmarshalJSON(b []byte) error {
	var ids []sui.ObjectID
	err := json.Unmarshal(b, &ids)
	if err != nil {
		return err
	}
	*o = NewObjectIDSet()
	for _, id := range ids {
		o.Add(id)
	}
	return nil
}

func (o ObjectIDSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Sorted())
}

func (o ObjectIDSet) Equals(b ObjectIDSet) bool {
	if len(o) != len(b) {
		return false
	}
	for id := range o {
		_, ok := b[id]
		if !ok {
			return false
		}
	}
	return true
}

type Assets struct {
	// Coins is a set of coin balances
	Coins CoinBalances `json:"coins" swagger:"required"`
	// Objects is a set of non-Coin object IDs (e.g. NFTs)
	Objects ObjectIDSet `json:"objects" swagger:"required"`
}

func NewEmptyAssets() *Assets {
	return &Assets{
		Coins:   NewCoinBalances(),
		Objects: NewObjectIDSet(),
	}
}

func NewAssets(baseTokens coin.Value) *Assets {
	return NewEmptyAssets().AddCoin(coin.BaseTokenType, baseTokens)
}

func AssetsFromAssetsBagWithBalances(assetsBag iscmove.AssetsBagWithBalances) *Assets {
	assets := NewEmptyAssets()
	for k, v := range assetsBag.Balances {
		assets.Coins.Add(coin.Type(k), coin.Value(v.TotalBalance))
	}
	return assets
}

func AssetsFromBytes(b []byte) (*Assets, error) {
	return rwutil.ReadFromBytes(b, NewEmptyAssets())
}

func (a *Assets) Clone() *Assets {
	r := NewEmptyAssets()
	r.Coins = a.Coins.Clone()
	r.Objects = maps.Clone(a.Objects)
	return r
}

func (a *Assets) AddCoin(coinType coin.Type, amount coin.Value) *Assets {
	a.Coins.Add(coinType, amount)
	return a
}

func (a *Assets) AddObject(id sui.ObjectID) *Assets {
	a.Objects.Add(id)
	return a
}

func (a *Assets) CoinBalance(coinType coin.Type) coin.Value {
	return a.Coins.Get(coinType)
}

func (a *Assets) String() string {
	s := lo.MapToSlice(a.Coins, func(coinType coin.Type, amount coin.Value) string {
		return fmt.Sprintf("%s: %d", coinType, amount)
	})
	s = append(s, lo.MapToSlice(a.Objects, func(id sui.ObjectID, _ struct{}) string {
		return id.ShortString()
	})...)
	return fmt.Sprintf("Assets{%s}", strings.Join(s, ", "))
}

func (a *Assets) Bytes() []byte {
	return rwutil.WriteToBytes(a)
}

func (a *Assets) Equals(b *Assets) bool {
	if (a == nil) || (b == nil) {
		return (a == nil) && (b == nil)
	}
	if a == b {
		return true
	}
	if !a.Coins.Equals(b.Coins) {
		return false
	}
	if !a.Objects.Equals(b.Objects) {
		return false
	}
	return true
}

// Spend subtracts assets from the current set, mutating the receiver.
// If the budget is not enough, returns false and leaves receiver untouched.
func (a *Assets) Spend(toSpend *Assets) bool {
	// check budget
	for coinType, spendAmount := range toSpend.Coins {
		available, ok := a.Coins[coinType]
		if !ok || available < spendAmount {
			return false
		}
	}
	for id := range toSpend.Objects {
		if !a.Objects.Has(id) {
			return false
		}
	}

	// budget is enough
	for coinType, spendAmount := range toSpend.Coins {
		a.Coins.Sub(coinType, spendAmount)
	}
	for id := range toSpend.Objects {
		delete(a.Objects, id)
	}
	return true
}

func (a *Assets) Add(b *Assets) *Assets {
	for coinType, amount := range b.Coins {
		a.Coins.Add(coinType, amount)
	}
	for id := range b.Objects {
		a.Objects.Add(id)
	}
	return a
}

func (a *Assets) IsEmpty() bool {
	return len(a.Coins) == 0 && len(a.Objects) == 0
}

func (a *Assets) AddBaseTokens(amount coin.Value) *Assets {
	a.Coins.Add(coin.BaseTokenType, amount)
	return a
}

func (a *Assets) SetBaseTokens(amount coin.Value) *Assets {
	a.Coins.Set(coin.BaseTokenType, amount)
	return a
}

func (a *Assets) BaseTokens() coin.Value {
	return a.Coins.Get(coin.BaseTokenType)
}

func (a *Assets) Read(r io.Reader) error {
	*a = Assets{}
	rr := rwutil.NewReader(r)
	rr.Read(&a.Coins)
	rr.Read(&a.Objects)
	return rr.Err
}

func (a *Assets) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.Write(a.Coins)
	ww.Write(a.Objects)
	return ww.Err
}
