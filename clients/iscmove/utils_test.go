package iscmove_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/sui-go/contracts"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
	"github.com/iotaledger/wasp/sui-go/suiclient"
	"github.com/iotaledger/wasp/sui-go/sui"
)

func buildAndDeployISCContracts(t *testing.T, client *iscmove.Client, signer cryptolib.Signer) *sui.PackageID {
	iscBytecode := contracts.ISC()

	txnBytes, err := client.Publish(context.Background(), suiclient.PublishRequest{
		Sender:          signer.Address().AsSuiAddress(),
		CompiledModules: iscBytecode.Modules,
		Dependencies:    iscBytecode.Dependencies,
		GasBudget:       suijsonrpc.NewBigInt(suiclient.DefaultGasBudget * 10),
	})
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(), cryptolib.SignerToSuiSigner(signer), txnBytes.TxBytes, &suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageID, err := txnResponse.GetPublishedPackageID()
	require.NoError(t, err)

	return packageID
}

func buildDeployMintTestcoin(t *testing.T, client *iscmove.Client, signer cryptolib.Signer) (
	*sui.PackageID,
	*sui.ObjectID,
) {
	testcoinBytecode := contracts.Testcoin()
	suiSigner := cryptolib.SignerToSuiSigner(signer)

	txnBytes, err := client.Publish(context.Background(), suiclient.PublishRequest{
		Sender:          signer.Address().AsSuiAddress(),
		CompiledModules: testcoinBytecode.Modules,
		Dependencies:    testcoinBytecode.Dependencies,
		GasBudget:       suijsonrpc.NewBigInt(suiclient.DefaultGasBudget * 10),
	})
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(), suiSigner, txnBytes.TxBytes, &suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	packageID, err := txnResponse.GetPublishedPackageID()
	require.NoError(t, err)

	treasuryCapRef, err := txnResponse.GetCreatedObjectInfo("coin", "TreasuryCap")
	require.NoError(t, err)

	mintAmount := uint64(1000000)
	txnRes, err := client.MintToken(
		context.Background(),
		suiSigner,
		packageID,
		"testcoin",
		treasuryCapRef.ObjectID,
		mintAmount,
		&suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.True(t, txnRes.Effects.Data.IsSuccess())

	return packageID, treasuryCapRef.ObjectID
}
