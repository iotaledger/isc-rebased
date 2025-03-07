package tests

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/tools/cluster/templates"
)

// executed in cluster_test.go
func testBasicAccounts(t *testing.T, env *ChainEnv) {
	testAccounts(env)
}

func TestBasicAccountsNLow(t *testing.T) {
	runTest := func(tt *testing.T, n, t int) {
		blockKeepAmount := 10
		clu := newCluster(tt, waspClusterOpts{
			nNodes: 4,
			modifyConfig: func(nodeIndex int, configParams templates.WaspConfigParams) templates.WaspConfigParams {
				// set node 0 as an "archive node"
				if nodeIndex == 0 {
					configParams.PruningMinStatesToKeep = -1
				} else {
					// all other nodes will only keep 10 blocks
					configParams.PruningMinStatesToKeep = blockKeepAmount
				}

				return configParams
			},
		})

		// set blockKeepAmount (active state pruning) to 10 as well
		chain, err := clu.DeployChainWithDKG(clu.Config.AllNodes(), clu.Config.AllNodes(), 4, int32(blockKeepAmount))
		require.NoError(tt, err)
		env := newChainEnv(tt, clu, chain)

		testAccounts(env)
	}
	t.Run("N=1", func(tt *testing.T) { runTest(tt, 1, 1) }) // passed
	t.Run("N=2", func(tt *testing.T) { runTest(tt, 2, 2) }) // passed
	t.Run("N=3", func(tt *testing.T) { runTest(tt, 3, 3) }) // passed
	t.Run("N=4", func(tt *testing.T) { runTest(tt, 4, 3) }) // passed
}

func testAccounts(e *ChainEnv) {
	e.t.Logf("   %s: %s", root.Contract.Name, root.Contract.Hname().String())
	e.t.Logf("   %s: %s", accounts.Contract.Name, accounts.Contract.Hname().String())

	e.checkCoreContracts()

	keyPair, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)
	originatorClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
	_, err = originatorClient.DepositFunds(10 * isc.Million)
	require.NoError(e.t, err)
	time.Sleep(3 * time.Second)
	chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
	balance1, err := chClient.L1Client.GetBalance(context.TODO(), iotaclient.GetBalanceRequest{Owner: keyPair.Address().AsIotaAddress()})
	require.NoError(e.t, err)

	for i := range e.Chain.CommitteeNodes {
		chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(i), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
		balance2, err := chClient.L1Client.GetBalance(context.TODO(), iotaclient.GetBalanceRequest{Owner: keyPair.Address().AsIotaAddress()})
		require.NoError(e.t, err)
		require.Equal(e.t, balance1.TotalBalance.Int64(), balance2.TotalBalance.Int64())
	}

	_, err = originatorClient.PostOffLedgerRequest(context.Background(),
		accounts.FuncWithdraw.Message(),
		chainclient.PostRequestParams{
			Allowance: isc.NewAssets(10),
		},
	)
	require.NoError(e.t, err)
	time.Sleep(3 * time.Second)

	for i := range e.Chain.CommitteeNodes {
		chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(i), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
		balance3, err := chClient.L1Client.GetBalance(context.TODO(), iotaclient.GetBalanceRequest{Owner: keyPair.Address().AsIotaAddress()})
		require.NoError(e.t, err)
		require.Equal(e.t, balance1.TotalBalance.Int64()+10, balance3.TotalBalance.Int64())
	}
}

// executed in cluster_test.go
func testBasic2Accounts(t *testing.T, e *ChainEnv) {
	e.checkCoreContracts()

	keyPair, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)
	originatorClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
	_, err = originatorClient.DepositFunds(10 * isc.Million)
	require.NoError(e.t, err)
	time.Sleep(3 * time.Second)
	chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
	balance1, err := chClient.L1Client.GetBalance(context.TODO(), iotaclient.GetBalanceRequest{Owner: keyPair.Address().AsIotaAddress()})
	require.NoError(e.t, err)

	for i := range e.Chain.CommitteeNodes {
		chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(i), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
		balance2, err := chClient.L1Client.GetBalance(context.TODO(), iotaclient.GetBalanceRequest{Owner: keyPair.Address().AsIotaAddress()})
		require.NoError(e.t, err)
		require.Equal(e.t, balance1.TotalBalance.Int64(), balance2.TotalBalance.Int64())
	}

	_, err = originatorClient.PostOffLedgerRequest(context.Background(),
		accounts.FuncWithdraw.Message(),
		chainclient.PostRequestParams{
			Allowance: isc.NewAssets(10),
		},
	)
	require.NoError(e.t, err)
	time.Sleep(3 * time.Second)

	for i := range e.Chain.CommitteeNodes {
		chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(i), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), keyPair)
		balance3, err := chClient.L1Client.GetBalance(context.TODO(), iotaclient.GetBalanceRequest{Owner: keyPair.Address().AsIotaAddress()})
		require.NoError(e.t, err)
		require.Equal(e.t, balance1.TotalBalance.Int64()+10, balance3.TotalBalance.Int64())
	}

	userWallet1, userAddress1, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	e.DepositFunds(iotaclient.DefaultGasBudget, userWallet1)
	_, userAddress2, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	time.Sleep(3 * time.Second)

	userWalletClient1 := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), userWallet1)

	user1L2Bal1 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddress1), isc.BaseTokenCoinInfo.CoinType)
	require.NoError(e.t, err)
	user2L2Bal1 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddress2), isc.BaseTokenCoinInfo.CoinType)
	require.NoError(e.t, err)

	var transferAmount coin.Value = 10
	req, err := userWalletClient1.PostOffLedgerRequest(context.Background(),
		accounts.FuncTransferAllowanceTo.Message(isc.NewAddressAgentID(userAddress2)),
		chainclient.PostRequestParams{
			Allowance: isc.NewAssets(transferAmount),
		},
	)
	require.NoError(e.t, err)
	time.Sleep(3 * time.Second)

	reqceipt, err := e.Chain.CommitteeMultiClient().WaitUntilRequestProcessedSuccessfully(context.Background(), e.Chain.ChainID, req.ID(), false, 30*time.Second)
	require.NoError(e.t, err)

	user1L2Bal2 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddress1), isc.BaseTokenCoinInfo.CoinType)
	require.NoError(e.t, err)
	user2L2Bal2 := e.getBalanceOnChain(isc.NewAddressAgentID(userAddress2), isc.BaseTokenCoinInfo.CoinType)
	require.NoError(e.t, err)
	gasFeeCharged, err := strconv.ParseUint(reqceipt.GasFeeCharged, 10, 64)
	require.NoError(e.t, err)
	require.Equal(t, user1L2Bal1-coin.Value(gasFeeCharged)-transferAmount, user1L2Bal2)
	require.Equal(t, user2L2Bal1+transferAmount, user2L2Bal2)
}
