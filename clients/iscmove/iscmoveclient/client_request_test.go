package iscmoveclient_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	testcommon "github.com/iotaledger/wasp/clients/iota-go/test_common"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient/iscmoveclienttest"
	"github.com/iotaledger/wasp/clients/iscmove/iscmovetest"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/testutil/l1starter"
	"github.com/iotaledger/wasp/packages/util/bcs"
)

func ensureSingleCoin(t *testing.T, cryptolibSigner cryptolib.Signer, client clients.L1Client) {
	coinType := iotajsonrpc.IotaCoinType.String()
	coinObjects, err := client.GetCoins(context.Background(), iotaclient.GetCoinsRequest{
		CoinType: &coinType,
		Owner:    cryptolibSigner.Address().AsIotaAddress(),
	})
	require.NoError(t, err)

	if len(coinObjects.Data) == 1 {
		return
	}

	txb := iotago.NewProgrammableTransactionBuilder()
	primaryCoin := coinObjects.Data[0]
	coinsToMerge := make([]iotago.Argument, 0)
	for i := 1; i < len(coinObjects.Data); i++ {
		coinToMerge := coinObjects.Data[i]
		coinsToMerge = append(coinsToMerge, txb.MustObj(iotago.ObjectArg{ImmOrOwnedObject: coinToMerge.Ref()}))
	}

	_ = txb.Command(
		iotago.Command{
			MergeCoins: &iotago.ProgrammableMergeCoins{
				Destination: iotago.GetArgumentGasCoin(),
				Sources:     coinsToMerge,
			},
		},
	)

	txData := iotago.NewProgrammable(
		cryptolibSigner.Address().AsIotaAddress(),
		txb.Finish(),
		[]*iotago.ObjectRef{primaryCoin.Ref()},
		iotaclient.DefaultGasBudget,
		parameters.L1Default.Protocol.ReferenceGasPrice.Uint64(),
	)

	txnBytes, err := bcs.Marshal(&txData)
	require.NoError(t, err)

	result, err := client.SignAndExecuteTransaction(
		context.Background(),
		&iotaclient.SignAndExecuteTransactionRequest{
			Signer:      cryptolib.SignerToIotaSigner(cryptolibSigner),
			TxDataBytes: txnBytes,
			Options: &iotajsonrpc.IotaTransactionBlockResponseOptions{
				ShowEffects:       true,
				ShowObjectChanges: true,
			},
		})
	require.NoError(t, err)
	t.Logf("SignAndExecuteTransaction, result: %+v", result)

	coinObjects, err = client.GetCoins(context.Background(), iotaclient.GetCoinsRequest{
		CoinType: &coinType,
		Owner:    cryptolibSigner.Address().AsIotaAddress(),
	})
	require.NoError(t, err)
	t.Logf("SignAndExecuteTransaction, contObjects: %+v", coinObjects)

	if len(coinObjects.Data) != 1 {
		t.Fatalf("Failed to merge all coins into one")
	}
}

func TestProperCoinUse(t *testing.T) {
	l1 := l1starter.Instance().L1Client()
	client := iscmoveclienttest.NewHTTPClient()
	chainOwnerSigner := iscmoveclienttest.NewRandomSignerWithFunds(t, 0)
	anchor := startNewChain(t, client, chainOwnerSigner)

	cryptolibSigner := iscmoveclienttest.NewRandomSignerWithFunds(t, 1)

	// Ensure we only have one actual gas coin. Merge all coins into one - if needed.
	ensureSingleCoin(t, cryptolibSigner, l1)

	createAndSendRequestRes, err := client.CreateAndSendRequestWithAssets(
		context.Background(),
		&iscmoveclient.CreateAndSendRequestWithAssetsRequest{
			Signer:        cryptolibSigner,
			PackageID:     l1starter.ISCPackageID(),
			AnchorAddress: anchor.ObjectID,
			Assets:        iscmove.NewAssets(100000),
			Message:       iscmovetest.RandomMessage(),
			Allowance:     iscmove.NewAssets(100000),
			GasPrice:      iotaclient.DefaultGasPrice,
			GasBudget:     iotaclient.DefaultGasBudget,
		},
	)
	require.NoError(t, err)
	t.Log(createAndSendRequestRes)
}

func TestCreateAndSendRequest(t *testing.T) {
	client := iscmoveclienttest.NewHTTPClient()
	anchorSigner := iscmoveclienttest.NewRandomSignerWithFunds(t, 0)
	anchor := startNewChain(t, client, anchorSigner)

	cryptolibSigner := iscmoveclienttest.NewRandomSignerWithFunds(t, 1)
	ensureSingleCoin(t, cryptolibSigner, l1starter.Instance().L1Client())
	t.Run("success", func(t *testing.T) {
		txnResponse, err := newAssetsBag(client, cryptolibSigner)
		require.NoError(t, err)
		assetsBagRef, err := txnResponse.GetCreatedObjectInfo(iscmove.AssetsBagModuleName, iscmove.AssetsBagObjectName)
		require.NoError(t, err)

		createAndSendRequestRes, err := client.CreateAndSendRequest(
			context.Background(),
			&iscmoveclient.CreateAndSendRequestRequest{
				Signer:        cryptolibSigner,
				PackageID:     l1starter.ISCPackageID(),
				AnchorAddress: anchor.ObjectID,
				AssetsBagRef:  assetsBagRef,
				Message:       iscmovetest.RandomMessage(),
				Allowance:     iscmove.NewAssets(100),
				GasPrice:      iotaclient.DefaultGasPrice,
				GasBudget:     iotaclient.DefaultGasBudget,
			},
		)
		require.NoError(t, err)

		_, err = createAndSendRequestRes.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
		require.NoError(t, err)
	})

	t.Run("max size AssetsBag", func(t *testing.T) {
		txnResponse, err := newAssetsBag(client, cryptolibSigner)
		require.NoError(t, err)
		assetsBagRef, err := txnResponse.GetCreatedObjectInfo(iscmove.AssetsBagModuleName, iscmove.AssetsBagObjectName)
		require.NoError(t, err)

		for i := 0; i < 25; i++ {
			coinRef, _ := buildDeployMintTestcoin(t, client, cryptolibSigner)
			getCoinRef, err := client.GetObject(
				context.Background(),
				iotaclient.GetObjectRequest{
					ObjectID: coinRef.ObjectID,
					Options:  &iotajsonrpc.IotaObjectDataOptions{ShowType: true},
				},
			)
			require.NoError(t, err)

			coinResource, err := iotago.NewResourceType(*getCoinRef.Data.Type)
			require.NoError(t, err)

			testCointype, err := iotajsonrpc.CoinTypeFromString(coinResource.SubType1.String())
			require.NoError(t, err)

			_, err = PTBTestWrapper(
				&PTBTestWrapperRequest{
					Client:    client,
					Signer:    cryptolibSigner,
					PackageID: l1starter.ISCPackageID(),
					GasPrice:  iotaclient.DefaultGasPrice,
					GasBudget: iotaclient.DefaultGasBudget,
				},
				func(ptb *iotago.ProgrammableTransactionBuilder) *iotago.ProgrammableTransactionBuilder {
					return iscmoveclient.PTBAssetsBagPlaceCoinWithAmount(
						ptb,
						l1starter.ISCPackageID(),
						ptb.MustObj(iotago.ObjectArg{ImmOrOwnedObject: assetsBagRef}),
						ptb.MustObj(iotago.ObjectArg{ImmOrOwnedObject: coinRef}),
						iotajsonrpc.CoinValue(100),
						testCointype,
					)
				},
			)
			require.NoError(t, err)

			assetsBagRef, err = client.UpdateObjectRef(context.Background(), assetsBagRef)
			require.NoError(t, err)
		}

		createAndSendRequestRes, err := client.CreateAndSendRequest(
			context.Background(),
			&iscmoveclient.CreateAndSendRequestRequest{
				Signer:        cryptolibSigner,
				PackageID:     l1starter.ISCPackageID(),
				AnchorAddress: anchor.ObjectID,
				AssetsBagRef:  assetsBagRef,
				Message:       iscmovetest.RandomMessage(),
				Allowance:     iscmove.NewAssets(100),
				GasPrice:      iotaclient.DefaultGasPrice,
				GasBudget:     iotaclient.DefaultGasBudget,
			},
		)
		require.NoError(t, err)

		_, err = createAndSendRequestRes.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
		require.NoError(t, err)
	})

	t.Run("oversized AssetsBag", func(t *testing.T) {
		txnResponse, err := newAssetsBag(client, cryptolibSigner)
		require.NoError(t, err)
		assetsBagRef, err := txnResponse.GetCreatedObjectInfo(iscmove.AssetsBagModuleName, iscmove.AssetsBagObjectName)
		require.NoError(t, err)

		for i := 0; i < 26; i++ {
			coinRef, _ := buildDeployMintTestcoin(t, client, cryptolibSigner)
			getCoinRef, err := client.GetObject(
				context.Background(),
				iotaclient.GetObjectRequest{
					ObjectID: coinRef.ObjectID,
					Options:  &iotajsonrpc.IotaObjectDataOptions{ShowType: true},
				},
			)
			require.NoError(t, err)

			coinResource, err := iotago.NewResourceType(*getCoinRef.Data.Type)
			require.NoError(t, err)

			testCointype, err := iotajsonrpc.CoinTypeFromString(coinResource.SubType1.String())
			require.NoError(t, err)

			_, err = PTBTestWrapper(
				&PTBTestWrapperRequest{
					Client:    client,
					Signer:    cryptolibSigner,
					PackageID: l1starter.ISCPackageID(),
					GasPrice:  iotaclient.DefaultGasPrice,
					GasBudget: iotaclient.DefaultGasBudget,
				},
				func(ptb *iotago.ProgrammableTransactionBuilder) *iotago.ProgrammableTransactionBuilder {
					return iscmoveclient.PTBAssetsBagPlaceCoinWithAmount(
						ptb,
						l1starter.ISCPackageID(),
						ptb.MustObj(iotago.ObjectArg{ImmOrOwnedObject: assetsBagRef}),
						ptb.MustObj(iotago.ObjectArg{ImmOrOwnedObject: coinRef}),
						iotajsonrpc.CoinValue(100),
						testCointype,
					)
				},
			)
			require.NoError(t, err)

			assetsBagRef, err = client.UpdateObjectRef(context.Background(), assetsBagRef)
			require.NoError(t, err)
		}

		_, err = client.CreateAndSendRequest(
			context.Background(),
			&iscmoveclient.CreateAndSendRequestRequest{
				Signer:        cryptolibSigner,
				PackageID:     l1starter.ISCPackageID(),
				AnchorAddress: anchor.ObjectID,
				AssetsBagRef:  assetsBagRef,
				Message:       iscmovetest.RandomMessage(),
				Allowance:     iscmove.NewAssets(100),
				GasPrice:      iotaclient.DefaultGasPrice,
				GasBudget:     iotaclient.DefaultGasBudget,
			},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), `function_name: Some("create_and_send_request") }, 1`)
	})
}

func TestCreateAndSendRequestWithAssets(t *testing.T) {
	client := iscmoveclienttest.NewHTTPClient()
	cryptolibSigner := iscmoveclienttest.NewSignerWithFunds(t, testcommon.TestSeed, 0)

	anchor := startNewChain(t, client, cryptolibSigner)

	createAndSendRequestRes, err := client.CreateAndSendRequestWithAssets(
		context.Background(),
		&iscmoveclient.CreateAndSendRequestWithAssetsRequest{
			Signer:        cryptolibSigner,
			PackageID:     l1starter.ISCPackageID(),
			AnchorAddress: anchor.ObjectID,
			Assets:        iscmove.NewAssets(100),
			Message:       iscmovetest.RandomMessage(),
			Allowance: &iscmove.Assets{
				Coins: iscmove.CoinBalances{
					iotajsonrpc.MustCoinTypeFromString("0x1::iota::IOTA"):    11,
					iotajsonrpc.MustCoinTypeFromString("0xa::testa::TEST_A"): 12,
				},
			},
			GasPrice:  iotaclient.DefaultGasPrice,
			GasBudget: iotaclient.DefaultGasBudget,
		},
	)
	require.NoError(t, err)

	_, err = createAndSendRequestRes.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
	require.NoError(t, err)
}

func TestGetRequestFromObjectID(t *testing.T) {
	client := iscmoveclienttest.NewHTTPClient()
	cryptolibSigner := iscmoveclienttest.NewSignerWithFunds(t, testcommon.TestSeed, 0)

	anchor := startNewChain(t, client, cryptolibSigner)

	txnResponse, err := newAssetsBag(client, cryptolibSigner)
	require.NoError(t, err)
	assetsBagRef, err := txnResponse.GetCreatedObjectInfo(iscmove.AssetsBagModuleName, iscmove.AssetsBagObjectName)
	require.NoError(t, err)

	createAndSendRequestRes, err := client.CreateAndSendRequest(
		context.Background(),
		&iscmoveclient.CreateAndSendRequestRequest{
			Signer:        cryptolibSigner,
			PackageID:     l1starter.ISCPackageID(),
			AnchorAddress: anchor.ObjectID,
			AssetsBagRef:  assetsBagRef,
			Message:       iscmovetest.RandomMessage(),
			Allowance: &iscmove.Assets{
				Coins: iscmove.CoinBalances{
					iotajsonrpc.MustCoinTypeFromString("0x1::iota::IOTA"):    21,
					iotajsonrpc.MustCoinTypeFromString("0xa::testa::TEST_A"): 12,
				},
			},
			GasPrice:  iotaclient.DefaultGasPrice,
			GasBudget: iotaclient.DefaultGasBudget,
		},
	)
	require.NoError(t, err)

	reqInfo, err := createAndSendRequestRes.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
	require.NoError(t, err)

	req, err := client.GetRequestFromObjectID(context.Background(), reqInfo.ObjectID)
	require.NoError(t, err)
	require.Equal(t, iotajsonrpc.CoinValue(12), req.Object.Allowance.Coins[iotajsonrpc.MustCoinTypeFromString("0xa::testa::TEST_A")])
	require.Equal(t, iotajsonrpc.CoinValue(21), req.Object.Allowance.Coins[iotajsonrpc.MustCoinTypeFromString("0x1::iota::IOTA")])
}
