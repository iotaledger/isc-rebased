package iotaclient_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaconn"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
)

func TestGetCommitteeInfo(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	epochId := iotajsonrpc.NewBigInt(0)
	committeeInfo, err := client.GetCommitteeInfo(context.Background(), epochId)
	require.NoError(t, err)
	require.Equal(t, epochId, committeeInfo.EpochId)
	// just use a arbitrary big number to ensure there are enough validator
	require.Greater(t, len(committeeInfo.Validators), 3)
}

func TestGetLatestIotaSystemState(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	state, err := client.GetLatestIotaSystemState(context.Background())
	require.NoError(t, err)
	require.NotNil(t, state)
}

func TestGetReferenceGasPrice(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	gasPrice, err := client.GetReferenceGasPrice(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, gasPrice.Int64(), int64(1000))
}

func TestGetStakes(t *testing.T) {
	// FIXME change the valid staking iotago address
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)

	// This address has been taken from https://explorer.iota.cafe/validator/0x02e1df479da7b51573248016db5f460586aad4d4c93315a1a8ed3c1a7fac1754
	address, err := iotago.AddressFromHex("0x02e1df479da7b51573248016db5f460586aad4d4c93315a1a8ed3c1a7fac1754")
	require.NoError(t, err)
	stakes, err := client.GetStakes(context.Background(), address)
	require.NoError(t, err)
	require.Greater(t, len(stakes), 0)
	for _, validator := range stakes {
		require.Equal(t, address, &validator.ValidatorAddress)
		for _, stake := range validator.Stakes {
			if stake.Data.StakeStatus.Data.Active != nil {
				t.Logf(
					"earned amount %10v at %v",
					stake.Data.StakeStatus.Data.Active.EstimatedReward.Uint64(),
					validator.ValidatorAddress,
				)
			}
		}
	}
}

func TestGetStakesByIds(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	// This address has been taken from https://explorer.iota.cafe/validator/0x02e1df479da7b51573248016db5f460586aad4d4c93315a1a8ed3c1a7fac1754
	owner, err := iotago.AddressFromHex("0x02e1df479da7b51573248016db5f460586aad4d4c93315a1a8ed3c1a7fac1754")
	require.NoError(t, err)
	stakes, err := api.GetStakes(context.Background(), owner)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(stakes), 1)

	stake1 := stakes[0].Stakes[0].Data
	stakeId := stake1.StakedIotaId
	stakesFromId, err := api.GetStakesByIds(context.Background(), []iotago.ObjectID{stakeId})
	require.NoError(t, err)
	require.Equal(t, len(stakesFromId), 1)

	queriedStake := stakesFromId[0].Stakes[0].Data
	require.Equal(t, stake1, queriedStake)
	t.Log(stakesFromId)
}

func TestGetValidatorsApy(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	apys, err := api.GetValidatorsApy(context.Background())
	require.NoError(t, err)
	t.Logf("current epoch %v", apys.Epoch)
	apyMap := apys.ApyMap()
	for _, apy := range apys.Apys {
		key := apy.Address
		t.Logf("%v apy: %v", key, apyMap[key])
	}
}
