package sui_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/iotaledger/isc-private/sui-go/models"
	"github.com/iotaledger/isc-private/sui-go/sui"
	"github.com/iotaledger/isc-private/sui-go/sui_signer"
	"github.com/iotaledger/isc-private/sui-go/sui_types"

	"github.com/stretchr/testify/require"
)

// @return types.DryRunTransactionBlockResponse
func dryRunTxn(
	t *testing.T,
	api *sui.ImplSuiAPI,
	txBytes sui_types.Base64Data,
	showJson bool,
) *models.DryRunTransactionBlockResponse {
	simulate, err := api.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.Equal(t, "", simulate.Effects.Data.V1.Status.Error)
	require.True(t, simulate.Effects.Data.IsSuccess())
	if showJson {
		data, err := json.Marshal(simulate)
		require.NoError(t, err)
		t.Log(string(data))
		t.Log("gasFee: ", simulate.Effects.Data.GasFee())
	}
	return simulate
}

func executeTxn(
	t *testing.T,
	api *sui.ImplSuiAPI,
	txBytes sui_types.Base64Data,
	signer *sui_signer.Signer,
) *models.SuiTransactionBlockResponse {
	// First of all, make sure that there are no problems with simulated trading.
	simulate, err := api.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// sign and send
	signature, err := signer.SignTransactionBlock(txBytes, sui_signer.DefaultIntent())
	require.NoError(t, err)
	options := models.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	}
	resp, err := api.ExecuteTransactionBlock(
		context.TODO(), txBytes, []*sui_signer.Signature{&signature}, &options,
		models.TxnRequestTypeWaitForLocalExecution,
	)
	require.NoError(t, err)
	t.Log(resp)
	return resp
}
