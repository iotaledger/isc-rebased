package sbtests

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/testutil/testmisc"

	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
)

func TestTooManyOutputsInASingleCall(t *testing.T) { run2(t, testTooManyOutputsInASingleCall) }
func testTooManyOutputsInASingleCall(t *testing.T) {
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	// send 1 tx will 1_000_000 BaseTokens which should result in too many outputs, so the request must fail
	wallet, _ := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(20))

	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncSplitFunds.Name).
		AddBaseTokens(1000 * isc.Million).
		AddAllowanceBaseTokens(999 * isc.Million). // contract is sending 1Mi per output
		WithGasBudget(math.MaxUint64)
	_, err := ch.PostRequestSync(req, wallet)
	require.Error(t, err)
	testmisc.RequireErrorToBe(t, err, vm.ErrExceededPostedOutputLimit)
	require.NotContains(t, err.Error(), "skipped")
}

func TestSeveralOutputsInASingleCall(t *testing.T) { run2(t, testSeveralOutputsInASingleCall) }
func testSeveralOutputsInASingleCall(t *testing.T) {
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, walletAddr := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(20))

	err := ch.DepositBaseTokensToL2(100_000, wallet)
	require.NoError(t, err)

	beforeWallet := ch.L1L2Funds(walletAddr)
	t.Logf("----- BEFORE wallet: %s", beforeWallet)

	// this will SUCCEED because it will result in 4 outputs in the single call
	const allowance = 4 * isc.Million
	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncSplitFunds.Name).
		AddAllowanceBaseTokens(allowance).
		AddBaseTokens(allowance + 1*isc.Million).
		WithGasBudget(math.MaxUint64)
	tx, _, err := ch.PostRequestSyncTx(req, wallet)
	require.NoError(t, err)

	storageDeposit := tx.Essence.Outputs[0].Deposit()
	ch.Env.AssertL1BaseTokens(walletAddr, beforeWallet.AssetsL1.BaseTokens+allowance-storageDeposit)
}

func TestSeveralOutputsInASingleCallFail(t *testing.T) { run2(t, testSeveralOutputsInASingleCallFail) }
func testSeveralOutputsInASingleCallFail(t *testing.T) {
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, walletAddr := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(20))

	err := ch.DepositBaseTokensToL2(100_000, wallet)
	require.NoError(t, err)

	beforeWallet := ch.L1L2Funds(walletAddr)
	t.Logf("----- BEFORE wallet: %s", beforeWallet)

	// this will FAIL because it will result in 5 outputs in the single call
	const allowance = 5 * isc.Million
	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncSplitFunds.Name).
		AddAllowanceBaseTokens(allowance).
		AddBaseTokens(allowance + 1*isc.Million).
		WithGasBudget(math.MaxUint64)

	_, err = ch.PostRequestSync(req, wallet)
	testmisc.RequireErrorToBe(t, err, vm.ErrExceededPostedOutputLimit)
	require.NotContains(t, err.Error(), "skipped")
}

func TestSplitTokensFail(t *testing.T) { run2(t, testSplitTokensFail) }
func testSplitTokensFail(t *testing.T) {
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, _ := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(20))

	err := ch.DepositBaseTokensToL2(2*isc.Million, wallet)
	require.NoError(t, err)

	sn, nativeTokenID, err := ch.NewNativeTokenParams(100).
		WithUser(wallet).
		CreateFoundry()
	require.NoError(t, err)
	err = ch.MintTokens(sn, 100, wallet)
	require.NoError(t, err)

	// this will FAIL because it will result in 100 outputs in the single call
	allowance := isc.NewAssets(100*isc.Million).AddCoin(nativeTokenID, 100)
	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncSplitFundsNativeTokens.Name).
		AddAllowance(allowance).
		AddBaseTokens(200 * isc.Million).
		WithGasBudget(math.MaxUint64)
	_, err = ch.PostRequestSync(req, wallet)
	testmisc.RequireErrorToBe(t, err, vm.ErrExceededPostedOutputLimit)
	require.NotContains(t, err.Error(), "skipped")
}

func TestSplitTokensSuccess(t *testing.T) { run2(t, testSplitTokensSuccess) }
func testSplitTokensSuccess(t *testing.T) {
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, addr := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(20))
	agentID := isc.NewAddressAgentID(addr)

	err := ch.DepositBaseTokensToL2(2*isc.Million, wallet)
	require.NoError(t, err)

	var amountMintedTokens coin.Value = 100
	sn, nativeTokenID, err := ch.NewNativeTokenParams(amountMintedTokens).
		WithUser(wallet).
		CreateFoundry()
	require.NoError(t, err)
	err = ch.MintTokens(sn, amountMintedTokens, wallet)
	require.NoError(t, err)

	var amountTokensToSend coin.Value = 3
	allowance := isc.NewAssets(2*isc.Million).AddCoin(nativeTokenID, amountTokensToSend)
	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncSplitFundsNativeTokens.Name).
		AddAllowance(allowance).
		AddBaseTokens(2 * isc.Million).
		WithGasBudget(math.MaxUint64)
	_, err = ch.PostRequestSync(req, wallet)
	require.NoError(t, err)
	require.Equal(t, ch.L2CoinBalance(agentID, nativeTokenID), amountMintedTokens-amountTokensToSend)
	require.Equal(t, ch.Env.L1CoinBalance(addr, nativeTokenID), amountTokensToSend)
}

func TestPingBaseTokens1(t *testing.T) { run2(t, testPingBaseTokens1) }
func testPingBaseTokens1(t *testing.T) {
	// TestPingBaseTokens1 sends some base tokens to SC and receives the whole allowance sent back to L1 as on-ledger request
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	user, userAddr := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(1))

	userFundsBefore := ch.L1L2Funds(userAddr)
	commonBefore := ch.L2CommonAccountAssets()
	t.Logf("----- BEFORE -----\nUser funds left: %s\nCommon account: %s", userFundsBefore, commonBefore)

	const expectedBack = 1 * isc.Million
	ch.Env.AssertL1BaseTokens(userAddr, iotaclient.FundsFromFaucetAmount)

	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncPingAllowanceBack.Name).
		AddBaseTokens(expectedBack + 500). // add extra base tokens besides allowance in order to estimate the gas fees
		AddAllowanceBaseTokens(expectedBack).
		WithGasBudget(100_000)

	_, estimate, err := ch.EstimateGasOnLedger(req, user)
	require.NoError(t, err)

	req.
		WithFungibleTokens(isc.NewAssets(expectedBack + estimate.GasFeeCharged).Coins).
		WithGasBudget(estimate.GasBurned)

	// re-estimate (it's possible the result is slightly different because we send less tokens (req is changed from  `exptected+500` above to `expected+estimate.GasFeeCharged`))
	_, estimate2, err := ch.EstimateGasOnLedger(req, user)
	require.NoError(t, err)
	req.
		WithFungibleTokens(isc.NewAssets(expectedBack + estimate2.GasFeeCharged).Coins).
		WithGasBudget(estimate2.GasBurned)

	_, err = ch.PostRequestSync(req, user)
	require.NoError(t, err)
	receipt := ch.LastReceipt()

	userFundsAfter := ch.L1L2Funds(userAddr)
	commonAfter := ch.L2CommonAccountAssets()
	t.Logf("------ AFTER ------\nReceipt: %s\nUser funds left: %s\nCommon account: %s", receipt, userFundsAfter, commonAfter)

	require.EqualValues(t, userFundsAfter.L1.BaseTokens(), iotaclient.FundsFromFaucetAmount-receipt.GasFeeCharged)
	require.EqualValues(t, commonBefore.BaseTokens(), commonAfter.BaseTokens())
	require.EqualValues(t, iotaclient.FundsFromFaucetAmount-receipt.GasFeeCharged, userFundsAfter.L1.BaseTokens())
	require.Zero(t, userFundsAfter.L2.BaseTokens())
}

func TestEstimateMinimumStorageDeposit(t *testing.T) { run2(t, testEstimateMinimumStorageDeposit) }
func testEstimateMinimumStorageDeposit(t *testing.T) {
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, _ := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(20))

	// should fail without enough base tokens to pay for a L1 transaction storage deposit
	allowance := isc.NewAssets(1)
	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncEstimateMinStorageDeposit.Name).
		AddAllowance(allowance).
		AddBaseTokens(100_000).
		WithGasBudget(100_000)

	_, err := ch.PostRequestSync(req, wallet)
	require.Error(t, err)

	// should succeed with enough base tokens to pay for a L1 transaction storage deposit
	allowance = isc.NewAssets(100_000)
	req = solo.NewCallParamsEx(ScName, sbtestsc.FuncEstimateMinStorageDeposit.Name).
		AddAllowance(allowance).
		AddBaseTokens(100_000).
		WithGasBudget(100_000)

	_, err = ch.PostRequestSync(req, wallet)
	require.NoError(t, err)
}

func mintDummyNFT(t *testing.T, ch *solo.Chain, issuer *cryptolib.KeyPair, owner *cryptolib.Address) *isc.NFT {
	nftMetadata := []byte("foobar")
	nft, err := ch.Env.MintNFTL1(issuer, owner, nftMetadata)
	require.NoError(t, err)
	return nft
}

func TestSendNFTsBack(t *testing.T) { run2(t, testSendNFTsBack) }
func testSendNFTsBack(t *testing.T) {
	// Send NFT and receive it back (on-ledger request)
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, addr := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(0))

	nft := mintDummyNFT(t, ch, wallet, addr)

	var baseTokensToSend coin.Value = 300_000
	var baseTokensForGas coin.Value = 100_000
	assetsToSend := isc.NewAssets(baseTokensToSend)
	assetsToAllow := isc.NewAssets(baseTokensToSend - baseTokensForGas)

	// receive an NFT back that is sent in the same request
	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncSendNFTsBack.Name).
		AddFungibleTokens(assetsToSend.Coins).
		WithObject(nft.ID).
		AddAllowance(assetsToAllow.AddObject(nft.ID)).
		WithMaxAffordableGasBudget()

	_, err := ch.PostRequestSync(req, wallet)
	require.NoError(t, err)
	require.True(t, ch.Env.HasL1NFT(addr, nft.ID))
}

func TestNFTOffledgerWithdraw(t *testing.T) { run2(t, testNFTOffledgerWithdraw) }

func testNFTOffledgerWithdraw(t *testing.T) {
	// Deposit an NFT, then claim it back via offleger-request
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, issuerAddr := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(0))

	nft := mintDummyNFT(t, ch, wallet, issuerAddr)

	require.True(t, ch.Env.HasL1NFT(issuerAddr, nft.ID))
	require.False(t, ch.Env.HasL1NFT(ch.ChainID.AsAddress(), nft.ID))
	require.False(t, ch.HasL2NFT(isc.NewAddressAgentID(issuerAddr), nft.ID))

	req := solo.NewCallParams(accounts.FuncDeposit.Message()).
		AddFungibleTokens(isc.NewAssets(1_000_000).Coins).
		WithObject(nft.ID).
		WithMaxAffordableGasBudget()

	_, err := ch.PostRequestSync(req, wallet)
	require.NoError(t, err)

	require.False(t, ch.Env.HasL1NFT(issuerAddr, nft.ID))
	require.True(t, ch.Env.HasL1NFT(ch.ChainID.AsAddress(), nft.ID))
	require.True(t, ch.HasL2NFT(isc.NewAddressAgentID(issuerAddr), nft.ID))

	wdReq := solo.NewCallParams(accounts.FuncWithdraw.Message()).
		WithAllowance(isc.NewAssets(10_000).AddObject(nft.ID)).
		WithMaxAffordableGasBudget()

	_, err = ch.PostRequestOffLedger(wdReq, wallet)
	require.NoError(t, err)

	require.True(t, ch.Env.HasL1NFT(issuerAddr, nft.ID))
	require.False(t, ch.Env.HasL1NFT(ch.ChainID.AsAddress(), nft.ID))
	require.False(t, ch.HasL2NFT(isc.NewAddressAgentID(issuerAddr), nft.ID))
}

func TestNFTMintToChain(t *testing.T) { run2(t, testNFTMintToChain) }
func testNFTMintToChain(t *testing.T) {
	// Mints an NFT as a request
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil)

	wallet, addr := ch.Env.NewKeyPairWithFunds(ch.Env.NewSeedFromIndex(0))

	nftToBeMinted := &isc.NFT{
		ID:       iotago.ObjectID{},
		Issuer:   addr,
		Metadata: []byte("foobar"),
	}

	var baseTokensToSend coin.Value = 300_000
	var baseTokensForGas coin.Value = 100_000
	assetsToSend := isc.NewAssets(baseTokensToSend)
	assetsToAllow := isc.NewAssets(baseTokensToSend - baseTokensForGas)

	// receive an NFT back that is sent in the same request
	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncClaimAllowance.Name).
		AddFungibleTokens(assetsToSend.Coins).
		WithObject(nftToBeMinted.ID).
		AddAllowance(assetsToAllow.AddObject(sui.Address{})). // empty NFTID
		WithMaxAffordableGasBudget()

	_, err := ch.PostRequestSync(req, wallet)
	require.NoError(t, err)
	// find out the NFTID
	receipt := ch.LastReceipt()
	nftID := isc.NFTIDFromOutputID(receipt.DeserializedRequest().ID().OutputID())

	// - Chain owns the NFT on L1
	require.True(t, ch.Env.HasL1NFT(ch.ChainID.AsAddress(), &nftID))
	// - The target contract owns the NFT on L2
	contractAgentID := isc.NewContractAgentID(ch.ChainID, sbtestsc.Contract.Hname())
	require.True(t, ch.HasL2NFT(contractAgentID, &nftID))
}
