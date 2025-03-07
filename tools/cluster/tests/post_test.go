package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
)

// executed in cluster_test.go
func testPost1Request(t *testing.T, e *ChainEnv) {
	userKeyPair, userAddr, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	userClient := e.Chain.Client(userKeyPair)
	balance1 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddr), isc.BaseTokenCoinInfo.CoinType)

	reqTx, err := userClient.PostRequest(context.Background(), accounts.FuncDeposit.Message(), chainclient.PostRequestParams{
		Transfer: isc.NewAssets(2 * iotaclient.DefaultGasBudget),
	})
	require.NoError(t, err)

	receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, reqTx, true, 30*time.Second)
	require.NoError(t, err)
	balance2 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddr), isc.BaseTokenCoinInfo.CoinType)

	gasFeeCharged, err := util.DecodeUint64(receipts[0].GasFeeCharged)
	require.NoError(t, err)
	require.Equal(t, balance1+2*coin.Value(iotaclient.DefaultGasBudget)-coin.Value(gasFeeCharged), balance2)
}

// executed in cluster_test.go
func testPost3Requests(t *testing.T, e *ChainEnv) {
	userKeyPair, userAddr, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	myClient := e.Chain.Client(userKeyPair)
	balance1 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddr), isc.BaseTokenCoinInfo.CoinType)

	const numRepeats int64 = 3
	var receipts [numRepeats]*apiclient.ReceiptResponse
	for i := 0; i < int(numRepeats); i++ {
		tx, err := myClient.PostRequest(context.Background(), accounts.FuncDeposit.Message(), chainclient.PostRequestParams{
			Transfer: isc.NewAssets(2 * iotaclient.DefaultGasBudget),
		})
		require.NoError(t, err)

		recs, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, true, 30*time.Second)
		require.NoError(t, err)
		receipts[i] = recs[0]
	}

	gasFeeChargedSum := coin.Value(0)
	for i := 0; i < int(numRepeats); i++ {
		gasFeeCharged, err := util.DecodeUint64(receipts[0].GasFeeCharged)
		require.NoError(t, err)
		gasFeeChargedSum += coin.Value(gasFeeCharged)
	}

	balance2 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddr), isc.BaseTokenCoinInfo.CoinType)
	require.Equal(t, balance1+3*2*coin.Value(iotaclient.DefaultGasBudget)-gasFeeChargedSum, balance2)
}

// executed in cluster_test.go
func testPost5AsyncRequests(t *testing.T, e *ChainEnv) {
	userWallet, userAddr, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	userClient := e.Chain.Client(userWallet)
	balance1 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddr), isc.BaseTokenCoinInfo.CoinType)

	tx := [5]*iotajsonrpc.IotaTransactionBlockResponse{}
	gasFeeChargedSum := coin.Value(0)
	baseTokesSent := coin.Value(2 * iotaclient.DefaultGasBudget)
	for i := 0; i < 5; i++ {
		tx[i], err = userClient.PostRequest(context.Background(), accounts.FuncDeposit.Message(), chainclient.PostRequestParams{
			Transfer: isc.NewAssets(baseTokesSent),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 5; i++ {
		receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx[i], false, 30*time.Second)
		require.NoError(t, err)

		gasFeeCharged, err := util.DecodeUint64(receipts[0].GasFeeCharged)
		require.NoError(t, err)

		gasFeeChargedSum += coin.Value(gasFeeCharged)
	}

	balance2 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddr), isc.BaseTokenCoinInfo.CoinType)
	require.Equal(t, balance1+5*2*coin.Value(iotaclient.DefaultGasBudget)-gasFeeChargedSum, balance2)
}
