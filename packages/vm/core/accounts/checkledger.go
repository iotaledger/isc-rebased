package accounts

import (
	"fmt"
	"math/big"

	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/util"
)

func (s *StateReader) CheckLedgerConsistency() error {
	t := s.GetTotalL2FungibleTokens()
	c := s.calcL2TotalFungibleTokens()
	if !t.Equals(c) {
		return fmt.Errorf(
			"inconsistent on-chain account ledger\n total assets: %s\ncalc total: %s\n",
			t, c,
		)
	}
	return nil
}

func (s *StateReader) calcL2TotalFungibleTokens() isc.CoinBalances {
	ret := isc.CoinBalances{}
	totalWeiRemainder := big.NewInt(0)

	s.allAccountsMapR().IterateKeys(func(accountKey []byte) bool {
		// add all native tokens owned by each account
		s.coinsMapR(kv.Key(accountKey)).Iterate(func(coinType []byte, val []byte) bool {
			ret.Add(
				codec.CoinType.MustDecode(coinType),
				codec.CoinValue.MustDecode(val),
			)
			return true
		})
		// use the full decimals for each account, so no dust balance is lost in the calculation
		totalWeiRemainder.Add(totalWeiRemainder, s.getWeiRemainder(kv.Key(accountKey)))
		return true
	})

	// convert total remainder from 18 decimals, must be exact
	ret.Add(
		coin.BaseTokenType,
		util.MustEthereumDecimalsToBaseTokenDecimalsExact(totalWeiRemainder, parameters.Decimals),
	)

	return ret
}
