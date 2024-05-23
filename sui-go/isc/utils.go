package isc

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils"
	"github.com/stretchr/testify/require"
)

// test only
func BuildAndDeployIscContracts(t *testing.T, client *Client, signer *sui_signer.Signer) *sui_types.PackageID {
	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/sui-go/isc/contracts/isc/")
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		sui_signer.TEST_ADDRESS,
		modules.Modules,
		modules.Dependencies,
		nil,
		models.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(), signer, txnBytes.TxBytes, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, txnResponse.Effects.Data.V1.Status.Status)

	packageID := txnResponse.GetPublishedPackageID()

	return packageID
}

// test only
func BuildDeployMintTestcoin(t *testing.T, client *Client, signer *sui_signer.Signer) (
	*sui_types.PackageID,
	*sui_types.ObjectID,
) {
	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/sui-go/contracts/testcoin/")
	require.NoError(t, err)

	txnBytes, err := client.Publish(
		context.Background(),
		sui_signer.TEST_ADDRESS,
		modules.Modules,
		modules.Dependencies,
		nil,
		models.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.NoError(t, err)
	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(), signer, txnBytes.TxBytes, &models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, txnResponse.Effects.Data.V1.Status.Status)

	packageID := txnResponse.GetPublishedPackageID()

	treasuryCap, _, err := sui.GetCreatedObjectIdAndType(txnResponse, "coin", "TreasuryCap")
	require.NoError(t, err)

	mintAmount := uint64(1000000)
	txnRes, err := client.MintToken(
		context.Background(),
		signer,
		packageID,
		"testcoin",
		treasuryCap,
		mintAmount,
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, models.ExecutionStatusSuccess, txnRes.Effects.Data.V1.Status.Status)

	return packageID, treasuryCap
}
