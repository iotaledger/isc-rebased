package sui_test

import (
	"context"
	"testing"

	"github.com/iotaledger/wasp/sui-go/models"
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/sui/conn"
	"github.com/iotaledger/wasp/sui-go/sui_signer"
	"github.com/iotaledger/wasp/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestGetAllBalances(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	balances, err := api.GetAllBalances(context.TODO(), sui_signer.TEST_ADDRESS)
	require.NoError(t, err)
	for _, balance := range balances {
		t.Logf(
			"Coin Name: %v, Count: %v, Total: %v, Locked: %v",
			balance.CoinType, balance.CoinObjectCount,
			balance.TotalBalance.String(), balance.LockedBalance,
		)
	}
}

func TestGetAllCoins(t *testing.T) {
	type args struct {
		ctx     context.Context
		address *sui_types.SuiAddress
		cursor  *sui_types.ObjectID
		limit   uint
	}

	tests := []struct {
		name    string
		a       *sui.ImplSuiAPI
		args    args
		want    *models.CoinPage
		wantErr bool
	}{
		{
			name: "successful with limit",
			a:    sui.NewSuiClient(conn.TestnetEndpointUrl),
			args: args{
				ctx:     context.TODO(),
				address: sui_signer.TEST_ADDRESS,
				cursor:  nil,
				limit:   3,
			},
			wantErr: false,
		},
		{
			name: "successful without limit",
			a:    sui.NewSuiClient(conn.TestnetEndpointUrl),
			args: args{
				ctx:     context.TODO(),
				address: sui_signer.TEST_ADDRESS,
				cursor:  nil,
				limit:   0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.a.GetAllCoins(tt.args.ctx, &models.GetAllCoinsRequest{
					Owner:  tt.args.address,
					Cursor: tt.args.cursor,
					Limit:  tt.args.limit,
				})
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAllCoins() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				// we have called multiple times RequestFundFromFaucet() on testnet, so the account have several SUI objects.
				require.GreaterOrEqual(t, len(got.Data), int(tt.args.limit))
				require.NotNil(t, got.NextCursor)
			},
		)
	}
}

func TestGetBalance(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	balance, err := api.GetBalance(context.TODO(), &models.GetBalanceRequest{Owner: sui_signer.TEST_ADDRESS})
	require.NoError(t, err)
	t.Logf(
		"Coin Name: %v, Count: %v, Total: %v, Locked: %v",
		balance.CoinType, balance.CoinObjectCount,
		balance.TotalBalance.String(), balance.LockedBalance,
	)
}

func TestGetCoinMetadata(t *testing.T) {
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	metadata, err := api.GetCoinMetadata(context.TODO(), models.SuiCoinType)
	require.NoError(t, err)
	testSuiMetadata := &models.SuiCoinMetadata{
		Decimals:    9,
		Description: "",
		IconUrl:     "",
		Id:          sui_types.MustObjectIDFromHex("0x587c29de216efd4219573e08a1f6964d4fa7cb714518c2c8a0f29abfa264327d"),
		Name:        "Sui",
		Symbol:      "SUI",
	}
	require.Equal(t, testSuiMetadata, metadata)
}

func TestGetCoins(t *testing.T) {
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	defaultCoinType := models.SuiCoinType
	coins, err := api.GetCoins(context.TODO(), &models.GetCoinsRequest{
		Owner:    sui_signer.TEST_ADDRESS,
		CoinType: &defaultCoinType,
		Limit:    3,
	})
	require.NoError(t, err)

	require.Len(t, coins.Data, 3)
	for _, data := range coins.Data {
		require.Equal(t, models.SuiCoinType, data.CoinType)
		require.Greater(t, data.Balance.Int64(), int64(0))
	}
}

func TestGetTotalSupply(t *testing.T) {
	type args struct {
		ctx      context.Context
		coinType string
	}

	tests := []struct {
		name    string
		api     *sui.ImplSuiAPI
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "get Sui supply",
			api:  sui.NewSuiClient(conn.DevnetEndpointUrl),
			args: args{
				context.TODO(),
				models.SuiCoinType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.api.GetTotalSupply(tt.args.ctx, tt.args.coinType)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetTotalSupply() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				targetSupply := &models.Supply{Value: models.NewBigInt(10000000000000000000)}
				require.Equal(t, targetSupply, got)
			},
		)
	}
}
