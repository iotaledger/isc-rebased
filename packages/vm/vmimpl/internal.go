package vmimpl

import (
	"math"
	"math/big"

	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/util/panicutil"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/corecontracts"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/vmexceptions"
	"github.com/iotaledger/wasp/sui-go/sui"
)

// creditToAccount credits assets to the chain ledger
func (reqctx *requestContext) creditToAccount(agentID isc.AgentID, coins isc.CoinBalances) {
	reqctx.accountsStateWriter(false).CreditToAccount(agentID, coins, reqctx.ChainID())
}

// creditToAccountFullDecimals credits assets to the chain ledger
func (reqctx *requestContext) creditToAccountFullDecimals(agentID isc.AgentID, amount *big.Int, gasBurn bool) {
	reqctx.accountsStateWriter(gasBurn).CreditToAccountFullDecimals(agentID, amount, reqctx.ChainID())
}

func (reqctx *requestContext) creditObjectsToAccount(agentID isc.AgentID, objectIDs []sui.ObjectID) {
	for _, id := range objectIDs {
		panic("TODO: how to get the object contents?")
		rec := accounts.ObjectRecord{
			ID:  id,
			BCS: []byte{},
		}
		reqctx.accountsStateWriter(false).CreditObjectToAccount(agentID, &rec, reqctx.ChainID())
	}
}

// debitFromAccount subtracts tokens from account if there are enough.
func (reqctx *requestContext) debitFromAccount(agentID isc.AgentID, coins isc.CoinBalances, gasBurn bool) {
	reqctx.accountsStateWriter(gasBurn).DebitFromAccount(agentID, coins, reqctx.ChainID())
}

// debitFromAccountFullDecimals subtracts basetokens tokens from account if there are enough.
func (reqctx *requestContext) debitFromAccountFullDecimals(agentID isc.AgentID, amount *big.Int, gasBurn bool) {
	reqctx.accountsStateWriter(gasBurn).DebitFromAccountFullDecimals(agentID, amount, reqctx.ChainID())
}

// debitObjectFromAccount removes a Object from an account.
func (reqctx *requestContext) debitObjectFromAccount(agentID isc.AgentID, objectID sui.ObjectID, gasBurn bool) {
	reqctx.accountsStateWriter(gasBurn).DebitObjectFromAccount(agentID, objectID, reqctx.ChainID())
}

func (reqctx *requestContext) mustMoveBetweenAccounts(fromAgentID, toAgentID isc.AgentID, assets *isc.Assets, gasBurn bool) {
	lo.Must0(reqctx.accountsStateWriter(gasBurn).MoveBetweenAccounts(fromAgentID, toAgentID, assets, reqctx.ChainID()))
}

func findContractByHname(chainState kv.KVStore, contractHname isc.Hname) (ret *root.ContractRecord) {
	return root.NewStateReaderFromChainState(chainState).FindContract(contractHname)
}

func (reqctx *requestContext) GetBaseTokensBalance(agentID isc.AgentID) (bts coin.Value, remainder *big.Int) {
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		bts, remainder = s.GetBaseTokensBalance(agentID, reqctx.ChainID())
	})
	return
}

func (reqctx *requestContext) GetBaseTokensBalanceDiscardRemainder(agentID isc.AgentID) (bts coin.Value) {
	bal, _ := reqctx.GetBaseTokensBalance(agentID)
	return bal
}

func (reqctx *requestContext) HasEnoughForAllowance(agentID isc.AgentID, allowance *isc.Assets) bool {
	var ret bool
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		ret = s.HasEnoughForAllowance(agentID, allowance, reqctx.ChainID())
	})
	return ret
}

func (reqctx *requestContext) GetCoinBalance(agentID isc.AgentID, nativeTokenID coin.Type) coin.Value {
	var ret coin.Value
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		ret = s.GetCoinBalance(agentID, nativeTokenID, reqctx.ChainID())
	})
	return ret
}

func (reqctx *requestContext) GetCoinBalanceTotal(coinType coin.Type) coin.Value {
	var ret coin.Value
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		ret = s.GetCoinBalanceTotal(coinType)
	})
	return ret
}

func (reqctx *requestContext) GetCoinBalances(agentID isc.AgentID) isc.CoinBalances {
	var ret isc.CoinBalances
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		ret = s.GetCoins(agentID, reqctx.ChainID())
	})
	return ret
}

func (reqctx *requestContext) GetAccountObjects(agentID isc.AgentID) (ret []sui.ObjectID) {
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		ret = s.GetAccountObjects(agentID)
	})
	return ret
}

func (reqctx *requestContext) GetObjectBCS(objectID sui.ObjectID) (ret []byte, ok bool) {
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		ret = s.GetObjectBCS(objectID)
	})
	return ret, ret != nil
}

func (reqctx *requestContext) GetSenderTokenBalanceForFees() coin.Value {
	sender := reqctx.req.SenderAccount()
	if sender == nil {
		return 0
	}
	return reqctx.GetBaseTokensBalanceDiscardRemainder(sender)
}

func (reqctx *requestContext) requestLookupKey() blocklog.RequestLookupKey {
	return blocklog.NewRequestLookupKey(reqctx.vm.stateDraft.BlockIndex(), reqctx.requestIndex)
}

func (reqctx *requestContext) eventLookupKey() *blocklog.EventLookupKey {
	return blocklog.NewEventLookupKey(reqctx.vm.stateDraft.BlockIndex(), reqctx.requestIndex, reqctx.requestEventIndex)
}

func (reqctx *requestContext) writeReceiptToBlockLog(vmError *isc.VMError) *blocklog.RequestReceipt {
	receipt := &blocklog.RequestReceipt{
		Request:       reqctx.req,
		GasBudget:     reqctx.gas.budgetAdjusted,
		GasBurned:     reqctx.gas.burned,
		GasFeeCharged: reqctx.gas.feeCharged,
		GasBurnLog:    reqctx.gas.burnLog,
	}

	if vmError != nil {
		b := vmError.Bytes()
		if len(b) > isc.VMErrorMessageLimit {
			vmError = coreerrors.ErrErrorMessageTooLong
		}
		receipt.Error = vmError.AsUnresolvedError()
	}

	reqctx.Debugf("writeReceiptToBlockLog - reqID:%s err: %v", reqctx.req.ID(), vmError)

	key := reqctx.requestLookupKey()
	var err error
	reqctx.callCore(blocklog.Contract, func(s kv.KVStore) {
		err = blocklog.NewStateWriter(s).SaveRequestReceipt(receipt, key)
	})
	if err != nil {
		panic(err)
	}
	for _, f := range reqctx.onWriteReceipt {
		reqctx.callCore(corecontracts.All[f.contract], func(s kv.KVStore) {
			f.callback(s, receipt.GasBurned)
		})
	}
	return receipt
}

func (reqctx *requestContext) mustSaveEvent(hContract isc.Hname, topic string, payload []byte) {
	if reqctx.requestEventIndex == math.MaxUint16 {
		panic(vm.ErrTooManyEvents)
	}
	reqctx.Debugf("MustSaveEvent/%s: topic: '%s'", hContract.String(), topic)

	event := &isc.Event{
		ContractID: hContract,
		Topic:      topic,
		Payload:    payload,
		Timestamp:  uint64(reqctx.Timestamp().UnixNano()),
	}
	eventKey := reqctx.eventLookupKey().Bytes()
	reqctx.callCore(blocklog.Contract, func(s kv.KVStore) {
		blocklog.NewStateWriter(s).SaveEvent(eventKey, event)
	})
	reqctx.requestEventIndex++
}

// updateOffLedgerRequestNonce updates stored nonce for ISC off ledger requests
func (reqctx *requestContext) updateOffLedgerRequestNonce() {
	reqctx.callAccounts(func(s *accounts.StateWriter) {
		s.IncrementNonce(reqctx.req.SenderAccount(), reqctx.ChainID())
	})
}

// adjustL2BaseTokensIfNeeded adjust L2 ledger for base tokens if the L1 changed because of storage deposit changes
func (reqctx *requestContext) adjustL2BaseTokensIfNeeded(adjustment coin.Value, account isc.AgentID) {
	if adjustment == 0 {
		return
	}
	err := panicutil.CatchPanicReturnError(func() {
		reqctx.callAccounts(func(s *accounts.StateWriter) {
			s.AdjustAccountBaseTokens(account, adjustment, reqctx.ChainID())
		})
	}, accounts.ErrNotEnoughFunds)
	if err != nil {
		panic(vmexceptions.ErrNotEnoughFundsForMinFee)
	}
}

func (reqctx *requestContext) GetCoinInfo(coinType coin.Type) (*isc.SuiCoinInfo, bool) {
	panic("TODO")
}
