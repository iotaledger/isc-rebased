package iscmove_test

import (
	"context"
	"testing"

	"github.com/iotaledger/wasp/sui-go/suiclient"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
	"github.com/iotaledger/wasp/sui-go/suisigner"
	"github.com/stretchr/testify/require"
)

func TestAssetsBagNewAndDestroyEmpty(t *testing.T) {
	cryptolibSigner := newSignerWithFunds(t, suisigner.TestSeed, 0)
	client := newLocalnetClient()

	iscPackageID := buildAndDeployISCContracts(t, client, cryptolibSigner)

	assetsBagRef, err := client.AssetsBagNew(
		context.Background(),
		cryptolibSigner,
		iscPackageID,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		&suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, assetsBagRef)

	res, err := client.AssetsDestroyEmpty(
		context.Background(),
		cryptolibSigner,
		iscPackageID,
		assetsBagRef,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		&suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	_, err = res.GetCreatedObjectInfo("assets_bag", "AssetsBag")
	require.Error(t, err, "not found")
}

func TestAssetsBagAddItems(t *testing.T) {
	cryptolibSigner := newSignerWithFunds(t, suisigner.TestSeed, 0)
	client := newLocalnetClient()

	iscPackageID := buildAndDeployISCContracts(t, client, cryptolibSigner)

	assetsBagMain, err := client.AssetsBagNew(
		context.Background(),
		cryptolibSigner,
		iscPackageID,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		&suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, assetsBagMain)

	_, coinRef := buildDeployMintTestcoin(t, client, cryptolibSigner)
	getCoinRef, err := client.GetObject(
		context.Background(),
		suiclient.GetObjectRequest{
			ObjectID: coinRef.ObjectID,
			Options:  &suijsonrpc.SuiObjectDataOptions{ShowType: true},
		},
	)
	require.NoError(t, err)

	assetsBagAddItemsRes, err := client.AssetsBagPlaceCoin(
		context.Background(),
		cryptolibSigner,
		iscPackageID,
		assetsBagMain,
		coinRef,
		*getCoinRef.Data.Type,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		&suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, assetsBagAddItemsRes.Effects.Data.IsSuccess())
}
