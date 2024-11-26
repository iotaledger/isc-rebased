package iscmoveclient_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient/iscmoveclienttest"
	"github.com/iotaledger/wasp/clients/iscmove/iscmovetest"
	"github.com/iotaledger/wasp/packages/testutil/l1starter"
)

func TestCreateAndSendRequest(t *testing.T) {
	client := iscmoveclienttest.NewLocalnetClient()
	cryptolibSigner := iscmoveclienttest.NewSignerWithFunds(t, iscmoveclienttest.TestSeed, 0)

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
			Allowance:     iscmove.NewAssets(100),
			GasPrice:      iotaclient.DefaultGasPrice,
			GasBudget:     iotaclient.DefaultGasBudget,
		},
	)
	require.NoError(t, err)

	_, err = createAndSendRequestRes.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
	require.NoError(t, err)
}

func TestCreateAndSendRequestWithAssets(t *testing.T) {
	client := iscmoveclienttest.NewLocalnetClient()
	cryptolibSigner := iscmoveclienttest.NewSignerWithFunds(t, iscmoveclienttest.TestSeed, 0)

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
					"0x1::iota::IOTA":    11,
					"0xa::testa::TEST_A": 12,
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
	client := iscmoveclienttest.NewLocalnetClient()
	cryptolibSigner := iscmoveclienttest.NewSignerWithFunds(t, iscmoveclienttest.TestSeed, 0)

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
					"0x1::iota::IOTA":    21,
					"0xa::testa::TEST_A": 12,
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
	require.Equal(t, uint64(12), req.Object.Allowance.Coins["0xa::testa::TEST_A"].Uint64())
	require.Equal(t, uint64(21), req.Object.Allowance.Coins["0x1::iota::IOTA"].Uint64())
}
