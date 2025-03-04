package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
)

func TestDepositWithdraw(t *testing.T) { // passed
	e := SetupWithChain(t)

	myWallet, myAddress, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)

	require.True(t,
		e.Clu.AssertAddressBalances(myAddress, isc.NewAssets(iotaclient.FundsFromFaucetAmount)),
	)

	myAgentID := isc.NewAddressAgentID(myAddress)
	// origAgentID := isc.NewAddressAgentID(e.Chain.OriginatorAddress(), 0)

	// chEnv.checkBalanceOnChain(origAgentID, isc.BaseTokenID, 0)
	e.checkBalanceOnChain(myAgentID, isc.BaseTokenCoinInfo.CoinType, 0)

	// deposit some base tokens to the chain
	var depositBaseTokens coin.Value = 10 * isc.Million
	chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), myWallet)

	par := chainclient.NewPostRequestParams().WithBaseTokens(depositBaseTokens)
	reqTx, err := chClient.PostRequest(context.Background(), accounts.FuncDeposit.Message(), *par)
	require.NoError(t, err)

	receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, reqTx, true, 30*time.Second)
	require.NoError(t, err)

	// e.checkBalanceOnChain(origAgentID, isc.BaseTokenID, 0)
	gasFees1, err := util.DecodeUint64(receipts[0].GasFeeCharged)
	require.NoError(t, err)

	var onChainBalance coin.Value = depositBaseTokens - coin.Value(gasFees1)
	e.checkBalanceOnChain(myAgentID, isc.BaseTokenCoinInfo.CoinType, onChainBalance)
	require.True(t,
		e.Clu.AssertAddressBalances(myAddress, isc.NewAssets(iotaclient.FundsFromFaucetAmount-depositBaseTokens-coin.Value(reqTx.Effects.Data.GasFee()))),
	)

	// withdraw some base tokens back
	var baseTokensToWithdraw coin.Value = 1 * isc.Million
	req, err := chClient.PostOffLedgerRequest(context.Background(), accounts.FuncWithdraw.Message(),
		chainclient.PostRequestParams{
			Allowance: isc.NewAssets(baseTokensToWithdraw),
		},
	)
	require.NoError(t, err)
	receipt, err := e.Chain.CommitteeMultiClient().WaitUntilRequestProcessedSuccessfully(context.Background(), e.Chain.ChainID, req.ID(), true, 30*time.Second)
	require.NoError(t, err)

	gasFees2, err := util.DecodeUint64(receipt.GasFeeCharged)
	require.NoError(t, err)

	e.checkBalanceOnChain(myAgentID, isc.BaseTokenCoinInfo.CoinType, onChainBalance-baseTokensToWithdraw-coin.Value(gasFees2))
	require.True(t,
		e.Clu.AssertAddressBalances(myAddress, isc.NewAssets(iotaclient.FundsFromFaucetAmount-depositBaseTokens+baseTokensToWithdraw-coin.Value(reqTx.Effects.Data.GasFee()))),
	)

	// TODO use "withdraw all base tokens" entrypoint to withdraw all remaining base tokens
}
