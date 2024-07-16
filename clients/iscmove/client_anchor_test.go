package iscmove_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/suiclient"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
	"github.com/iotaledger/wasp/sui-go/suisigner"
)

func TestStartNewChain(t *testing.T) {
	client := newLocalnetClient()
	signer := newSignerWithFunds(t, suisigner.TestSeed, 0)

	iscPackageID := buildAndDeployISCContracts(t, client, signer)

	anchor, err := client.StartNewChain(
		context.Background(),
		signer,
		iscPackageID,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		[]byte{},
		false,
	)
	require.NoError(t, err)
	t.Log("anchor: ", anchor)
}

func TestGetAnchorFromObjectID(t *testing.T) {
	client := newLocalnetClient()
	signer := newSignerWithFunds(t, suisigner.TestSeed, 0)

	iscPackageID := buildAndDeployISCContracts(t, client, signer)

	anchor1, err := client.StartNewChain(
		context.Background(),
		signer,
		iscPackageID,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		[]byte{},
		false,
	)
	require.NoError(t, err)
	t.Log("anchor1: ", anchor1)

	anchor2, err := client.GetAnchorFromObjectID(context.Background(), anchor1.Ref.ObjectID)
	require.NoError(t, err)
	require.Equal(t, anchor1, anchor2)
}

func TestReceiveAndUpdateStateRootRequest(t *testing.T) {
	client := newLocalnetClient()
	cryptolibSigner := newSignerWithFunds(t, suisigner.TestSeed, 0)
	chainSigner := newSignerWithFunds(t, suisigner.TestSeed, 1)

	iscPackageID := buildAndDeployISCContracts(t, client, cryptolibSigner)

	anchor := startNewChain(t, client, chainSigner, iscPackageID)

	txnResponse, err := client.AssetsBagNew(
		context.Background(),
		cryptolibSigner,
		iscPackageID,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		false,
	)
	require.NoError(t, err)
	sentAssetsBagRef, err := txnResponse.GetCreatedObjectInfo(iscmove.AssetsBagModuleName, iscmove.AssetsBagObjectName)
	require.NoError(t, err)

	createAndSendRequestRes, err := client.CreateAndSendRequest(
		context.Background(),
		cryptolibSigner,
		iscPackageID,
		anchor.Ref.ObjectID,
		sentAssetsBagRef,
		"test_isc_contract",
		"test_isc_func",
		[][]byte{[]byte("one"), []byte("two"), []byte("three")},
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		false,
	)

	require.NoError(t, err)

	requestRef, err := createAndSendRequestRes.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
	require.NoError(t, err)

	resGetObject, err := client.GetObject(context.Background(),
		suiclient.GetObjectRequest{ObjectID: anchor.Ref.ObjectID, Options: &suijsonrpc.SuiObjectDataOptions{ShowType: true}})
	require.NoError(t, err)
	anchorRef := resGetObject.Data.Ref()

	_, err = client.ReceiveAndUpdateStateRootRequest(
		context.Background(),
		chainSigner,
		iscPackageID,
		anchorRef,
		[]sui.ObjectRef{*requestRef},
		[]byte{1, 2, 3},
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		false,
	)
	require.NoError(t, err)
}

func startNewChain(t *testing.T, client *iscmove.Client, signer cryptolib.Signer, iscPackageID sui.PackageID) *iscmove.Anchor {
	anchor, err := client.StartNewChain(
		context.Background(),
		signer,
		iscPackageID,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		[]byte{},
		false,
	)
	require.NoError(t, err)
	return anchor
}
