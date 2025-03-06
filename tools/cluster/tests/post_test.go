package tests

import (
	"context"
	`strconv`
	"testing"
	"time"

	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/apiextensions"
	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/contracts/inccounter"
)

func deployInccounter42(e *ChainEnv) *isc.ContractAgentID {
	e.checkCoreContracts()

	myWallet, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)

	myClient := e.Chain.Client(myWallet)

	var numRepeats int64 = 42
	tx, err := myClient.PostRequest(context.Background(), inccounter.FuncIncAndRepeatMany.Message(nil, &numRepeats), chainclient.PostRequestParams{
		Transfer:  isc.NewAssets(10 * isc.Million),
		Allowance: isc.NewAssets(9 * isc.Million),
	})
	require.NoError(e.t, err)

	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, false, 30*time.Second)
	require.NoError(e.t, err)

	for i := range e.Chain.CommitteeNodes {
		counterValue, err2 := e.Chain.GetCounterValue(i)
		require.NoError(e.t, err2)
		require.EqualValues(e.t, 42, counterValue)
	}

	result, err := apiextensions.CallView(
		context.Background(),
		e.Chain.Cluster.WaspClient(),
		e.Chain.ChainID.String(),
		apiextensions.CallViewReq(root.ViewFindContract.Message(inccounter.Contract.Hname())),
	)
	require.NoError(e.t, err)

	found, _, err := root.ViewFindContract.DecodeOutput(result)
	require.NoError(e.t, err)
	require.True(e.t, found)

	e.expectCounter(42)
	return isc.NewContractAgentID(e.Chain.ChainID, inccounter.Contract.Hname())
}

// executed in cluster_test.go
func testPostDeployInccounter(t *testing.T, e *ChainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounter.Contract.Name, contractID.String())
}

// executed in cluster_test.go
func testPost1Request(t *testing.T, e *ChainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounter.Contract.Name, contractID.String())

	myWallet, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	myClient := e.Chain.Client(myWallet)

	tx, err := myClient.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(t, err)

	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, false, 30*time.Second)
	require.NoError(t, err)

	e.expectCounter(43)
}

// executed in cluster_test.go
func testPost3Recursive(t *testing.T, e *ChainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounter.Contract.Name, contractID.String())

	myWallet, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	myClient := e.Chain.Client(myWallet)

	var numRepeats int64 = 3
	tx, err := myClient.PostRequest(context.Background(), inccounter.FuncIncAndRepeatMany.Message(nil, &numRepeats), chainclient.PostRequestParams{
		Transfer:  isc.NewAssets(10 * isc.Million),
		Allowance: isc.NewAssets(9 * isc.Million),
	})
	require.NoError(t, err)

	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, false, 30*time.Second)
	require.NoError(t, err)

	e.waitUntilCounterEquals(43+3, 10*time.Second)
}

// executed in cluster_test.go
func testPost5Requests(t *testing.T, e *ChainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounter.Contract.Name, contractID.String())

	myWallet, myAddress, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	myAgentID := isc.NewAddressAgentID(myAddress)
	myClient := e.Chain.Client(myWallet)

	e.checkBalanceOnChain(myAgentID, coin.BaseTokenType, 0)
	onChainBalance := coin.Value(0)
	for i := 0; i < 5; i++ {
		baseTokesSent := coin.Value(1 * isc.Million)
		tx, err := myClient.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{
			Transfer: isc.NewAssets(baseTokesSent),
		})
		require.NoError(t, err)

		receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, false, 30*time.Second)
		require.NoError(t, err)

		gasFeeCharged, err := strconv.ParseUint(receipts[0].GasFeeCharged, 10, 64)
		require.NoError(t, err)

		onChainBalance += baseTokesSent - coin.Value(gasFeeCharged)
	}

	e.expectCounter(42 + 5)
	e.checkBalanceOnChain(myAgentID, coin.BaseTokenType, onChainBalance)
}

// executed in cluster_test.go
func testPost5AsyncRequests(t *testing.T, e *ChainEnv) {
	contractID := deployInccounter42(e)
	t.Logf("-------------- deployed contract. Name: '%s' id: %s", inccounter.Contract.Name, contractID.String())

	myWallet, myAddress, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	myAgentID := isc.NewAddressAgentID(myAddress)

	myClient := e.Chain.Client(myWallet)

	tx := [5]*iotajsonrpc.IotaTransactionBlockResponse{}
	onChainBalance := coin.Value(0)
	baseTokesSent := coin.Value(1 * isc.Million)
	for i := 0; i < 5; i++ {
		tx[i], err = myClient.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{
			Transfer: isc.NewAssets(baseTokesSent),
		})
		require.NoError(t, err)
	}

	for i := 0; i < 5; i++ {
		receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx[i], false, 30*time.Second)
		require.NoError(t, err)

		gasFeeCharged, err := strconv.ParseUint(receipts[0].GasFeeCharged, 10, 64)
		require.NoError(t, err)

		onChainBalance += baseTokesSent - coin.Value(gasFeeCharged)
	}

	e.expectCounter(42 + 5)
	e.checkBalanceOnChain(myAgentID, coin.BaseTokenType, onChainBalance)

	if !e.Clu.AssertAddressBalances(myAddress,
		isc.NewAssets(iotaclient.FundsFromFaucetAmount-5*baseTokesSent)) {
		t.Fatal()
	}
}
