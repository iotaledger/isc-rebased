package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/contracts/inccounter"
)

// executed in cluster_test.go
func testBasicAccounts(t *testing.T, env *ChainEnv) {
	testAccounts(env)
}

func TestBasicAccountsNLow(t *testing.T) {
	t.Skip("Cluster tests currently disabled")
	runTest := func(tt *testing.T, n, t int) {
		e := setupWithNoChain(tt)
		chainNodes := make([]int, n)
		for i := range chainNodes {
			chainNodes[i] = i
		}
		chain, err := e.Clu.DeployChainWithDKG(chainNodes, chainNodes, uint16(t))
		require.NoError(tt, err)
		env := newChainEnv(tt, e.Clu, chain)
		testAccounts(env)
	}
	t.Run("N=1", func(tt *testing.T) { runTest(tt, 1, 1) })
	t.Run("N=2", func(tt *testing.T) { runTest(tt, 2, 2) })
	t.Run("N=3", func(tt *testing.T) { runTest(tt, 3, 3) })
	t.Run("N=4", func(tt *testing.T) { runTest(tt, 4, 3) })
}

func testAccounts(e *ChainEnv) {
	e.t.Logf("   %s: %s", root.Contract.Name, root.Contract.Hname().String())
	e.t.Logf("   %s: %s", accounts.Contract.Name, accounts.Contract.Hname().String())

	e.checkCoreContracts()

	for i := range e.Chain.CommitteeNodes {
		blockIndex, err2 := e.Chain.BlockIndex(i)
		require.NoError(e.t, err2)
		require.Greater(e.t, blockIndex, uint32(2))

		contractRegistry, err2 := e.Chain.ContractRegistry(i)
		require.NoError(e.t, err2)

		cr, ok := lo.Find(contractRegistry, func(item apiclient.ContractInfoResponse) bool {
			return item.HName == inccounter.Contract.Hname().String()
		})
		require.True(e.t, ok)

		require.EqualValues(e.t, inccounter.Contract.Name, cr.Name)

		counterValue, err2 := e.Chain.GetCounterValue(i)
		require.NoError(e.t, err2)
		require.EqualValues(e.t, 42, counterValue)
	}

	myWallet, myAddress, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)

	transferBaseTokens := coin.Value(1 * isc.Million)
	chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Chain.OriginatorClient().IscPackageID, myWallet)

	par := chainclient.NewPostRequestParams().WithBaseTokens(transferBaseTokens)
	reqTx, err := chClient.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), *par)
	require.NoError(e.t, err)

	receipts, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, reqTx, false, 10*time.Second)
	require.NoError(e.t, err)

	fees, err := util.DecodeUint64(receipts[0].GasFeeCharged)
	require.NoError(e.t, err)

	e.checkBalanceOnChain(isc.NewAddressAgentID(myAddress), coin.BaseTokenType, transferBaseTokens-coin.Value(fees))

	for i := range e.Chain.CommitteeNodes {
		counterValue, err := e.Chain.GetCounterValue(i)
		require.NoError(e.t, err)
		require.EqualValues(e.t, 43, counterValue)
	}

	if !e.Clu.AssertAddressBalances(myAddress, isc.NewAssets(iotaclient.FundsFromFaucetAmount-transferBaseTokens)) {
		e.t.Fatal()
	}

	incCounterAgentID := isc.NewContractAgentID(e.Chain.ChainID, inccounter.Contract.Hname())
	e.checkBalanceOnChain(incCounterAgentID, coin.BaseTokenType, 0)
}

// executed in cluster_test.go
func testBasic2Accounts(t *testing.T, env *ChainEnv) {
	chain := env.Chain

	env.checkCoreContracts()

	for _, i := range chain.CommitteeNodes {
		blockIndex, err2 := chain.BlockIndex(i)
		require.NoError(t, err2)
		require.Greater(t, blockIndex, uint32(2))

		contractRegistry, err2 := chain.ContractRegistry(i)
		require.NoError(t, err2)

		t.Logf("%+v", contractRegistry)
		cr, ok := lo.Find(contractRegistry, func(item apiclient.ContractInfoResponse) bool {
			return item.HName == inccounter.Contract.Hname().String()
		})
		require.True(t, ok)
		require.NotNil(t, cr)

		require.EqualValues(t, inccounter.Contract.Name, cr.Name)

		counterValue, err2 := chain.GetCounterValue(i)
		require.NoError(t, err2)
		require.EqualValues(t, 42, counterValue)
	}

	originatorSigScheme := chain.OriginatorKeyPair
	originatorAddress := chain.OriginatorAddress()

	myWallet, myAddress, err := env.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	transferBaseTokens := coin.Value(1 * isc.Million)
	myWalletClient := chainclient.New(env.Clu.L1Client(), env.Clu.WaspClient(0), chain.ChainID, env.Clu.Config.ISCPackageID(), myWallet)

	par := chainclient.NewPostRequestParams().WithBaseTokens(transferBaseTokens)
	reqTx, err := myWalletClient.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), *par)
	require.NoError(t, err)

	_, err = chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), chain.ChainID, reqTx, false, 30*time.Second)
	require.NoError(t, err)

	for _, i := range chain.CommitteeNodes {
		counterValue, err2 := chain.GetCounterValue(i)
		require.NoError(t, err2)
		require.EqualValues(t, 43, counterValue)
	}
	if !env.Clu.AssertAddressBalances(myAddress, isc.NewAssets(iotaclient.FundsFromFaucetAmount-transferBaseTokens)) {
		t.Fatal()
	}

	// withdraw back 500 base tokens to originator address
	fmt.Printf("\norig address from sigsheme: %s\n", originatorAddress.String())
	origL1Balance := env.Clu.AddressBalances(originatorAddress).BaseTokens()
	originatorClient := chainclient.New(env.Clu.L1Client(), env.Clu.WaspClient(0), chain.ChainID, env.Clu.Config.ISCPackageID(), originatorSigScheme)
	allowanceBaseTokens := coin.Value(uint64(800_000))
	req2, err := originatorClient.PostOffLedgerRequest(context.Background(), accounts.FuncWithdraw.Message(),
		chainclient.PostRequestParams{
			Allowance: isc.NewAssets(allowanceBaseTokens),
		},
	)
	require.NoError(t, err)

	_, err = chain.CommitteeMultiClient().WaitUntilRequestProcessedSuccessfully(context.Background(), chain.ChainID, req2.ID(), true, 30*time.Second)
	require.NoError(t, err)

	require.Equal(t, env.Clu.AddressBalances(originatorAddress).BaseTokens, origL1Balance+allowanceBaseTokens)
}
