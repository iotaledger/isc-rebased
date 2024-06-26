package sui_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/iotaledger/wasp/sui-go/models"
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/sui/conn"
	"github.com/iotaledger/wasp/sui-go/sui_signer"
	"github.com/iotaledger/wasp/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

const (
	ComingChatValidatorAddress = "0x520289e77c838bae8501ae92b151b99a54407288fdd20dee6e5416bfe943eb7a"
)

func TestRequestAddDelegation(t *testing.T) {
	client, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)

	coins, err := client.GetCoins(context.Background(), signer.Address, nil, nil, 10)
	require.NoError(t, err)

	amount := uint64(sui_types.UnitSui)
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(amount), 0, 0, 0)
	require.NoError(t, err)

	validatorAddress := ComingChatValidatorAddress
	validator, err := sui_types.SuiAddressFromHex(validatorAddress)
	require.NoError(t, err)

	txBytes, err := sui.BCS_RequestAddStake(
		signer.Address,
		pickedCoins.CoinRefs(),
		models.NewBigInt(amount),
		validator,
		sui.DefaultGasBudget,
		1000,
	)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.True(t, simulate.Effects.Data.IsSuccess())
}

func TestRequestWithdrawDelegation(t *testing.T) {
	client := sui.NewSuiClient(conn.TestnetEndpointUrl)

	signer, err := sui_types.SuiAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	stakes, err := client.GetStakes(context.Background(), signer)
	require.NoError(t, err)
	require.True(t, len(stakes) > 0)
	require.True(t, len(stakes[0].Stakes) > 0)

	coins, err := client.GetCoins(context.Background(), signer, nil, nil, 10)
	require.NoError(t, err)
	pickedCoins, err := models.PickupCoins(coins, new(big.Int), sui.DefaultGasBudget, 0, 0)
	require.NoError(t, err)

	stakeId := stakes[0].Stakes[0].Data.StakedSuiId
	detail, err := client.GetObject(context.Background(), &stakeId, nil)
	require.NoError(t, err)
	txBytes, err := sui.BCS_RequestWithdrawStake(signer, detail.Data.Ref(), pickedCoins.CoinRefs(), sui.DefaultGasBudget, 1000)
	require.NoError(t, err)

	simulate, err := client.DryRunTransaction(context.Background(), txBytes)
	require.NoError(t, err)
	require.True(t, simulate.Effects.Data.IsSuccess())
}
