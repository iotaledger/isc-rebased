package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/clients/apiextensions"
	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/isctest"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/contracts/inccounter"
)

// ensures a nodes resumes normal operation after rebooting
func TestReboot(t *testing.T) {
	t.Skip("Cluster tests currently disabled")

	env := setupNativeInccounterTest(t, 4, []int{0, 1, 2, 3})
	// env := setupNativeInccounterTest(t, 3, []int{0, 1, 2})
	client := env.createNewClient()

	// ------ TODO why does this make the test fail?
	_, er := env.Clu.WaspClient(0).ChainsAPI.DeactivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)
	_, er = env.Clu.WaspClient(0).ChainsAPI.ActivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)

	_, er = env.Clu.WaspClient(1).ChainsAPI.DeactivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)
	_, er = env.Clu.WaspClient(1).ChainsAPI.ActivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)
	//-------

	tx, err := client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(t, err)

	_, err = apiextensions.APIWaitUntilAllRequestsProcessed(context.Background(), env.Clu.WaspClient(0), env.Chain.ChainID, tx, true, 10*time.Second)
	require.NoError(t, err)

	env.expectCounter(1)

	req, err := client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	require.NoError(t, err)

	_, _, err = env.Clu.WaspClient(0).ChainsAPI.
		WaitForRequest(context.Background(), env.Chain.ChainID.String(), req.ID().String()).
		TimeoutSeconds(10).
		Execute()
	require.NoError(t, err)

	env.expectCounter(2)

	// // ------ TODO why does this make the test fail?
	// er = env.Clu.WaspClient(0).DeactivateChain(env.Chain.ChainID)
	// require.NoError(t, er)
	// er = env.Clu.WaspClient(0).ActivateChain(env.Chain.ChainID)
	// require.NoError(t, er)

	// tx, err = client.PostRequest(inccounter.FuncIncCounter.Name)
	// require.NoError(t, err)
	// _, err = env.Clu.WaspClient(0).WaitUntilAllRequestsProcessed(env.Chain.ChainID, tx, true,10*time.Second)
	// require.NoError(t, err)
	// env.expectCounter( 3)

	// reqx, err := client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Name)
	// require.NoError(t, err)
	// env.Clu.WaspClient(0).WaitUntilRequestProcessed(env.Chain.ChainID, reqx.ID(), 10*time.Second)
	// env.expectCounter( 4)
	// //-------

	// restart the nodes
	err = env.Clu.RestartNodes(true, 0, 1, 2, 3)
	require.NoError(t, err)

	// after rebooting, the chain should resume processing requests without issues
	tx, err = client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(t, err)

	_, err = apiextensions.APIWaitUntilAllRequestsProcessed(context.Background(), env.Clu.WaspClient(0), env.Chain.ChainID, tx, true, 10*time.Second)
	require.NoError(t, err)
	env.expectCounter(3)

	// ensure offledger requests are still working
	req, err = client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	require.NoError(t, err)

	_, _, err = env.Clu.WaspClient(0).ChainsAPI.
		WaitForRequest(context.Background(), env.Chain.ChainID.String(), req.ID().String()).
		TimeoutSeconds(10).
		Execute()
	require.NoError(t, err)
	env.expectCounter(4)
}

func TestReboot2(t *testing.T) {
	t.Skip("Cluster tests currently disabled")

	env := setupNativeInccounterTest(t, 4, []int{0, 1, 2, 3})
	client := env.createNewClient()

	tx, err := client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(t, err)

	_, err = apiextensions.APIWaitUntilAllRequestsProcessed(context.Background(), env.Clu.WaspClient(0), env.Chain.ChainID, tx, true, 10*time.Second)
	require.NoError(t, err)

	env.expectCounter(1)

	req, err := client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	require.NoError(t, err)

	_, _, err = env.Clu.WaspClient(0).ChainsAPI.
		WaitForRequest(context.Background(), env.Chain.ChainID.String(), req.ID().String()).
		WaitForL1Confirmation(true).
		TimeoutSeconds(10).
		Execute()
	require.NoError(t, err)

	env.expectCounter(2)

	// ------ TODO why does this make the test fail?
	_, er := env.Clu.WaspClient(0).ChainsAPI.DeactivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)

	_, er = env.Clu.WaspClient(0).ChainsAPI.ActivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)

	_, er = env.Clu.WaspClient(1).ChainsAPI.DeactivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)
	_, er = env.Clu.WaspClient(1).ChainsAPI.ActivateChain(context.Background(), env.Chain.ChainID.String()).Execute()
	require.NoError(t, er)
	//-------

	// // ------ TODO why does this make the test fail?
	// er = env.Clu.WaspClient(0).DeactivateChain(env.Chain.ChainID)
	// require.NoError(t, er)
	// er = env.Clu.WaspClient(0).ActivateChain(env.Chain.ChainID)
	// require.NoError(t, er)

	// tx, err = client.PostRequest(inccounter.FuncIncCounter.Message(nil))
	// require.NoError(t, err)
	// _, err = env.Clu.WaspClient(0).WaitUntilAllRequestsProcessed(env.Chain.ChainID, tx, true,10*time.Second)
	// require.NoError(t, err)
	// env.expectCounter( 3)

	// reqx, err := client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	// require.NoError(t, err)
	// env.Clu.WaspClient(0).WaitUntilRequestProcessed(env.Chain.ChainID, reqx.ID(), 10*time.Second)
	// env.expectCounter( 4)
	// //-------

	// restart the nodes
	err = env.Clu.RestartNodes(true, 0, 1, 2, 3)
	require.NoError(t, err)

	// after rebooting, the chain should resume processing requests without issues
	tx, err = client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(t, err)

	_, err = apiextensions.APIWaitUntilAllRequestsProcessed(context.Background(), env.Clu.WaspClient(0), env.Chain.ChainID, tx, true, 10*time.Second)
	require.NoError(t, err)

	env.expectCounter(3)
	// ensure off-ledger requests are still working
	req, err = client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	require.NoError(t, err)

	_, _, err = env.Clu.WaspClient(0).ChainsAPI.
		WaitForRequest(context.Background(), env.Chain.ChainID.String(), req.ID().String()).
		TimeoutSeconds(10).
		Execute()
	require.NoError(t, err)

	env.expectCounter(4)
}

type incCounterClient struct {
	expected int64
	t        *testing.T
	env      *ChainEnv
	client   *chainclient.Client
}

func newIncCounterClient(t *testing.T, env *ChainEnv, client *chainclient.Client) *incCounterClient {
	return &incCounterClient{t: t, env: env, client: client}
}

func (icc *incCounterClient) MustIncOnLedger() {
	tx, err := icc.client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(icc.t, err)

	_, err = apiextensions.APIWaitUntilAllRequestsProcessed(context.Background(), icc.env.Clu.WaspClient(0), icc.env.Chain.ChainID, tx, true, 10*time.Second)
	require.NoError(icc.t, err)

	icc.expected++
	icc.env.expectCounter(icc.expected)
}

func (icc *incCounterClient) MustIncOffLedger() {
	req, err := icc.client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	require.NoError(icc.t, err)

	_, _, err = icc.env.Clu.WaspClient(0).ChainsAPI.
		WaitForRequest(context.Background(), icc.env.Chain.ChainID.String(), req.ID().String()).
		TimeoutSeconds(10).
		Execute()
	require.NoError(icc.t, err)

	icc.expected++
	icc.env.expectCounter(icc.expected)
}

func (icc *incCounterClient) MustIncBoth(onLedgerFirst bool) {
	if onLedgerFirst {
		icc.MustIncOnLedger()
		icc.MustIncOffLedger()
	} else {
		icc.MustIncOffLedger()
		icc.MustIncOnLedger()
	}
}

// Ensures a nodes resumes normal operation after rebooting.
// In this case we have F=0 and N=3, thus any reboot violates the assumptions.
func TestRebootN3Single(t *testing.T) {
	t.Skip("Cluster tests currently disabled")

	tm := util.NewTimer()
	allNodes := []int{0, 1, 2}
	env := setupNativeInccounterTest(t, 3, allNodes)
	tm.Step("setupNativeInccounterTest")
	client := env.createNewClient()
	tm.Step("createNewClient")

	env.DepositFunds(1_000_000, client.KeyPair.(*cryptolib.KeyPair)) // For Off-ledger requests to pass.
	tm.Step("DepositFunds")

	icc := newIncCounterClient(t, env, client)
	icc.MustIncBoth(true)
	tm.Step("incCounter")

	// Restart all nodes, one by one.
	for _, nodeIndex := range allNodes {
		require.NoError(t, env.Clu.RestartNodes(true, nodeIndex))
		icc.MustIncBoth(nodeIndex%2 == 1)
		tm.Step(fmt.Sprintf("incCounter-%v", nodeIndex))
	}
	t.Logf("Timing: %v", tm.String())
}

// Ensures a nodes resumes normal operation after rebooting.
// In this case we have F=0 and N=3, thus any reboot violates the assumptions.
// We restart 2 nodes each iteration in this scenario..
func TestRebootN3TwoNodes(t *testing.T) {
	t.Skip("Cluster tests currently disabled")

	tm := util.NewTimer()
	allNodes := []int{0, 1, 2}
	env := setupNativeInccounterTest(t, 3, allNodes)
	tm.Step("setupNativeInccounterTest")
	client := env.createNewClient()
	tm.Step("createNewClient")

	env.DepositFunds(1_000_000, client.KeyPair.(*cryptolib.KeyPair)) // For Off-ledger requests to pass.
	tm.Step("DepositFunds")

	icc := newIncCounterClient(t, env, client)
	icc.MustIncBoth(true)
	tm.Step("incCounter")

	// Restart all nodes, one by one.
	for _, nodeIndex := range allNodes {
		otherTwo := lo.Filter(allNodes, func(ni int, _ int) bool { return ni != nodeIndex })
		require.NoError(t, env.Clu.RestartNodes(true, otherTwo...))
		icc.MustIncBoth(nodeIndex%2 == 1)
		tm.Step(fmt.Sprintf("incCounter-%v", nodeIndex))
	}
	t.Logf("Timing: %v", tm.String())
}

// Test rebooting nodes during operation.
func TestRebootDuringTasks(t *testing.T) {
	t.Skip("Cluster tests currently disabled")

	env := setupNativeInccounterTest(t, 4, []int{0, 1, 2, 3})
	restartDelay := 20 * time.Second
	restartCases := [][]int{
		{0},
		{0, 1},
		{2, 3},
		{1, 2, 3},
		{3},
		{0, 1, 2, 3},
	}
	postDelay := 200 * time.Millisecond
	postCount := int(restartDelay/postDelay) * len(restartCases) // To have enough posts for all restarts.

	// keep the nodes spammed
	{
		// deposit funds for off-ledger requests
		keyPair, _, err := env.Clu.NewKeyPairWithFunds()
		require.NoError(t, err)

		env.DepositFunds(iotaclient.FundsFromFaucetAmount, keyPair)
		client := env.Chain.Client(keyPair)

		go func() {
			for i := 0; i < postCount; i++ {
				// ignore any error (nodes might be down when sending)
				client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
				time.Sleep(postDelay)
			}
		}()

		go func() {
			keyPair, _, err := env.Clu.NewKeyPairWithFunds()
			require.NoError(t, err)
			client := env.Chain.Client(keyPair)
			for i := 0; i < postCount; i++ {
				_, err = client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
				require.NoError(t, err)
				time.Sleep(postDelay)
			}
		}()
	}

	lastCounter := int64(0)

	restart := func(indexes ...int) {
		t.Logf("restart, indexes=%v", indexes)
		// restart the nodes
		err := env.Clu.RestartNodes(true, indexes...)
		require.NoError(t, err)
		time.Sleep(restartDelay)

		// after rebooting, the chain should resume processing requests/views without issues
		ret, err := apiextensions.CallView(
			context.Background(),
			env.Clu.WaspClient(0),
			env.Chain.ChainID.String(),
			apiclient.ContractCallViewRequest{
				ContractHName: inccounter.Contract.Hname().String(),
				FunctionHName: inccounter.ViewGetCounter.Hname().String(),
			})
		require.NoError(t, err)

		counter, err := inccounter.ViewGetCounter.DecodeOutput(ret)
		require.NoError(t, err)
		require.Greater(t, counter, lastCounter)
		lastCounter = counter

		// assert the node still processes on/off-ledger requests
		keyPair2, _, err := env.Clu.NewKeyPairWithFunds()
		require.NoError(t, err)
		// deposit funds, then move them via off-ledger request
		env.DepositFunds(iotaclient.FundsFromFaucetAmount, keyPair2)
		client := env.Chain.Client(keyPair2)
		targetAgentID := isctest.NewRandomAgentID()
		req, err := client.PostOffLedgerRequest(
			context.Background(),
			accounts.FuncTransferAllowanceTo.Message(targetAgentID),
			chainclient.PostRequestParams{Allowance: isc.NewAssets(5000)},
		)
		require.NoError(t, err)
		_, err = env.Clu.MultiClient().WaitUntilRequestProcessed(context.Background(), env.Chain.ChainID, req.ID(), true, 10*time.Second)
		require.NoError(t, err)
		env.checkBalanceOnChain(targetAgentID, coin.BaseTokenType, 5000)
	}

	for _, restartIndexes := range restartCases {
		restart(restartIndexes...)
	}
}

func TestRebootRecoverFromWAL(t *testing.T) {
	t.Skip("Cluster tests currently disabled")

	env := setupNativeInccounterTest(t, 4, []int{0, 1, 2, 3})
	client := env.createNewClient()

	tx, err := client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(t, err)

	_, err = apiextensions.APIWaitUntilAllRequestsProcessed(context.Background(), env.Clu.WaspClient(0), env.Chain.ChainID, tx, true, 10*time.Second)
	require.NoError(t, err)

	env.expectCounter(1)

	req, err := client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	require.NoError(t, err)

	_, _, err = env.Clu.WaspClient(0).ChainsAPI.
		WaitForRequest(context.Background(), env.Chain.ChainID.String(), req.ID().String()).
		TimeoutSeconds(10).
		Execute()
	require.NoError(t, err)

	env.expectCounter(2)

	// restart the nodes, delete the DB
	err = env.Clu.RestartNodes(false, 0, 1, 2, 3)
	require.NoError(t, err)

	// after rebooting, the chain should resume processing requests without issues
	tx, err = client.PostRequest(context.Background(), inccounter.FuncIncCounter.Message(nil), chainclient.PostRequestParams{})
	require.NoError(t, err)

	_, err = apiextensions.APIWaitUntilAllRequestsProcessed(context.Background(), env.Clu.WaspClient(0), env.Chain.ChainID, tx, false, 10*time.Second)
	require.NoError(t, err)
	env.expectCounter(3)

	// ensure off-ledger requests are still working
	req, err = client.PostOffLedgerRequest(context.Background(), inccounter.FuncIncCounter.Message(nil))
	require.NoError(t, err)

	_, _, err = env.Clu.WaspClient(0).ChainsAPI.
		WaitForRequest(context.Background(), env.Chain.ChainID.String(), req.ID().String()).
		TimeoutSeconds(10).
		Execute()
	require.NoError(t, err)
	env.expectCounter(4)
}
