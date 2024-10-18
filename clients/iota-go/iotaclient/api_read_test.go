package iotaclient_test

import (
	"context"
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaconn"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
)

var (
	testMnemonic = "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
	testSeed     = []byte{
		50,
		230,
		119,
		9,
		86,
		155,
		106,
		30,
		245,
		81,
		234,
		122,
		116,
		90,
		172,
		148,
		59,
		33,
		88,
		252,
		134,
		42,
		231,
		198,
		208,
		141,
		209,
		116,
		78,
		21,
		216,
		24,
	}
	testAddress = iotago.MustAddressFromHex("0x786dff8a4ee13d45b502c8f22f398e3517e6ec78aa4ae564c348acb07fad7f50")
)

func TestGetChainIdentifier(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	chainID, err := client.GetChainIdentifier(context.Background())
	require.NoError(t, err)
	require.Equal(t, iotaconn.ChainIdentifierAlphanet, chainID)
}

func TestGetCheckpoint(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	checkpoint, err := client.GetCheckpoint(context.Background(), iotajsonrpc.NewBigInt(1000))
	require.NoError(t, err)
	targetCheckpoint := &iotajsonrpc.Checkpoint{
		Epoch:                    iotajsonrpc.NewBigInt(0),
		SequenceNumber:           iotajsonrpc.NewBigInt(1000),
		Digest:                   *iotago.MustNewDigest("Eu7yhUZ1oma3fk8KhHW86usFvSmjZ7QPEhPsX7ZYfRg3"),
		NetworkTotalTransactions: iotajsonrpc.NewBigInt(1004),
		PreviousDigest:           iotago.MustNewDigest("AcrgtLsNQxZQRU1JK395vanZzSR6nTun6huJAxEJuk14"),
		EpochRollingGasCostSummary: iotajsonrpc.GasCostSummary{
			ComputationCost:         iotajsonrpc.NewBigInt(0),
			StorageCost:             iotajsonrpc.NewBigInt(0),
			StorageRebate:           iotajsonrpc.NewBigInt(0),
			NonRefundableStorageFee: iotajsonrpc.NewBigInt(0),
		},
		TimestampMs:           iotajsonrpc.NewBigInt(1725548499477),
		Transactions:          []*iotago.Digest{iotago.MustNewDigest("8iu72fMHEFHiJMfjrPDTKBPufQgMSRKfeh2idG5CoHvE")},
		CheckpointCommitments: []iotago.CheckpointCommitment{},
		ValidatorSignature:    *iotago.MustNewBase64Data("k0u7tZR87vS8glhPgmCzgKFm1UU1ikmPmO9nVzFXn9XY20kpftc6zxdBe0lmSAzs"),
	}

	require.Equal(t, targetCheckpoint, checkpoint)
}

func TestGetCheckpoints(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	cursor := iotajsonrpc.NewBigInt(999)
	limit := uint64(2)
	checkpointPage, err := client.GetCheckpoints(
		context.Background(), iotaclient.GetCheckpointsRequest{
			Cursor: cursor,
			Limit:  &limit,
		},
	)
	require.NoError(t, err)
	targetCheckpoints := []*iotajsonrpc.Checkpoint{
		{
			Epoch:                    iotajsonrpc.NewBigInt(0),
			SequenceNumber:           iotajsonrpc.NewBigInt(1000),
			Digest:                   *iotago.MustNewDigest("Eu7yhUZ1oma3fk8KhHW86usFvSmjZ7QPEhPsX7ZYfRg3"),
			NetworkTotalTransactions: iotajsonrpc.NewBigInt(1004),
			PreviousDigest:           iotago.MustNewDigest("AcrgtLsNQxZQRU1JK395vanZzSR6nTun6huJAxEJuk14"),
			EpochRollingGasCostSummary: iotajsonrpc.GasCostSummary{
				ComputationCost:         iotajsonrpc.NewBigInt(0),
				StorageCost:             iotajsonrpc.NewBigInt(0),
				StorageRebate:           iotajsonrpc.NewBigInt(0),
				NonRefundableStorageFee: iotajsonrpc.NewBigInt(0),
			},
			TimestampMs:           iotajsonrpc.NewBigInt(1725548499477),
			Transactions:          []*iotago.Digest{iotago.MustNewDigest("8iu72fMHEFHiJMfjrPDTKBPufQgMSRKfeh2idG5CoHvE")},
			CheckpointCommitments: []iotago.CheckpointCommitment{},
			ValidatorSignature:    *iotago.MustNewBase64Data("k0u7tZR87vS8glhPgmCzgKFm1UU1ikmPmO9nVzFXn9XY20kpftc6zxdBe0lmSAzs"),
		},
		{
			Epoch:                    iotajsonrpc.NewBigInt(0),
			SequenceNumber:           iotajsonrpc.NewBigInt(1001),
			Digest:                   *iotago.MustNewDigest("EJtUUwsKXJR9C9JcJ31e3VZ5rPEsjRu4cSMUaGiTARyo"),
			NetworkTotalTransactions: iotajsonrpc.NewBigInt(1005),
			PreviousDigest:           iotago.MustNewDigest("Eu7yhUZ1oma3fk8KhHW86usFvSmjZ7QPEhPsX7ZYfRg3"),
			EpochRollingGasCostSummary: iotajsonrpc.GasCostSummary{
				ComputationCost:         iotajsonrpc.NewBigInt(0),
				StorageCost:             iotajsonrpc.NewBigInt(0),
				StorageRebate:           iotajsonrpc.NewBigInt(0),
				NonRefundableStorageFee: iotajsonrpc.NewBigInt(0),
			},
			TimestampMs:           iotajsonrpc.NewBigInt(1725548500033),
			Transactions:          []*iotago.Digest{iotago.MustNewDigest("X3QFYvZm5yAgg3nPVPox6jWskpd2cw57Xg8uXNtCTW5")},
			CheckpointCommitments: []iotago.CheckpointCommitment{},
			ValidatorSignature:    *iotago.MustNewBase64Data("jHdu/+su0PZ+93y7du1LH48p1+WAqVm2+5EpvMaFrRBnT0Y63EOTl6fMJFwHEizu"),
		},
	}
	require.Len(t, checkpointPage.Data, 2)
	t.Log(checkpointPage.Data[0].Transactions[0].String())
	t.Log(checkpointPage.Data[1].Transactions[0].String())
	require.Equal(t, checkpointPage.Data, targetCheckpoints)
	require.Equal(t, true, checkpointPage.HasNextPage)
	require.Equal(t, iotajsonrpc.NewBigInt(1001), checkpointPage.NextCursor)
}

func TestGetEvents(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	digest, err := iotago.NewDigest("3vVi8XZgNpzQ34PFgwJTQqWtPMU84njcBX1EUxUHhyDk")
	require.NoError(t, err)
	events, err := client.GetEvents(context.Background(), digest)
	require.NoError(t, err)
	require.Len(t, events, 1)
	for _, event := range events {
		require.Equal(t, digest, &event.Id.TxDigest)
		require.Equal(
			t,
			iotago.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			event.PackageId,
		)
		require.Equal(t, "clob_v2", event.TransactionModule)
		require.Equal(
			t,
			iotago.MustAddressFromHex("0xf0f13f7ef773c6246e87a8f059a684d60773f85e992e128b8272245c38c94076"),
			event.Sender,
		)
		targetStructTag := iotago.StructTag{
			Address: iotago.MustAddressFromHex("0xdee9"),
			Module:  iotago.Identifier("clob_v2"),
			Name:    iotago.Identifier("OrderPlaced"),
			TypeParams: []iotago.TypeTag{
				{
					Struct: &iotago.StructTag{
						Address: iotago.MustAddressFromHex("0x2"),
						Module:  iotago.Identifier("iota"),
						Name:    iotago.Identifier("IOTA"),
					},
				},
				{
					Struct: &iotago.StructTag{
						Address: iotago.MustAddressFromHex("0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf"),
						Module:  iotago.Identifier("coin"),
						Name:    iotago.Identifier("COIN"),
					},
				},
			},
		}
		require.Equal(t, targetStructTag.Address, event.Type.Address)
		require.Equal(t, targetStructTag.Module, event.Type.Module)
		require.Equal(t, targetStructTag.Name, event.Type.Name)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.Address, event.Type.TypeParams[0].Struct.Address)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.Module, event.Type.TypeParams[0].Struct.Module)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.Name, event.Type.TypeParams[0].Struct.Name)
		require.Equal(t, targetStructTag.TypeParams[0].Struct.TypeParams, event.Type.TypeParams[0].Struct.TypeParams)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.Address, event.Type.TypeParams[1].Struct.Address)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.Module, event.Type.TypeParams[1].Struct.Module)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.Name, event.Type.TypeParams[1].Struct.Name)
		require.Equal(t, targetStructTag.TypeParams[1].Struct.TypeParams, event.Type.TypeParams[1].Struct.TypeParams)
		targetBcsBase85 := base58.Decode("yNS5iDS3Gvdo3DhXdtFpuTS12RrSiNkrvjcm2rejntCuqWjF1DdwnHgjowdczAkR18LQHcBqbX2tWL76rys9rTCzG6vm7Tg34yqUkpFSMqNkcS6cfWbN8SdVsxn5g4ZEQotdBgEFn8yN7hVZ7P1MKvMwWf")
		require.Equal(t, targetBcsBase85, event.Bcs.Data())
		// TODO check ParsedJson map
	}
}

func TestGetLatestCheckpointSequenceNumber(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	sequenceNumber, err := client.GetLatestCheckpointSequenceNumber(context.Background())
	require.NoError(t, err)
	num, err := strconv.Atoi(sequenceNumber)
	require.NoError(t, err)
	require.Greater(t, num, 34317507)
}

func TestGetObject(t *testing.T) {
	type args struct {
		ctx   context.Context
		objID *iotago.ObjectID
	}
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	coins, err := api.GetCoins(
		context.TODO(), iotaclient.GetCoinsRequest{
			Owner: testAddress,
			Limit: 1,
		},
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		api     *iotaclient.Client
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test for devnet",
			api:  api,
			args: args{
				ctx:   context.TODO(),
				objID: coins.Data[0].CoinObjectID,
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.api.GetObject(
					tt.args.ctx, iotaclient.GetObjectRequest{
						ObjectID: tt.args.objID,
						Options: &iotajsonrpc.IotaObjectDataOptions{
							ShowType:                true,
							ShowOwner:               true,
							ShowContent:             true,
							ShowDisplay:             true,
							ShowBcs:                 true,
							ShowPreviousTransaction: true,
							ShowStorageRebate:       true,
						},
					},
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetObject() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%+v", got)
			},
		)
	}
}

func TestGetProtocolConfig(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	version := iotajsonrpc.NewBigInt(47)
	protocolConfig, err := api.GetProtocolConfig(context.Background(), version)
	require.NoError(t, err)
	require.Equal(t, uint64(47), protocolConfig.ProtocolVersion.Uint64())
}

func TestGetTotalTransactionBlocks(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	res, err := api.GetTotalTransactionBlocks(context.Background())
	require.NoError(t, err)
	t.Log(res)
}

func TestGetTransactionBlock(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	digest, err := iotago.NewDigest("D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL")
	require.NoError(t, err)
	resp, err := client.GetTransactionBlock(
		context.Background(), iotaclient.GetTransactionBlockRequest{
			Digest: digest,
			Options: &iotajsonrpc.IotaTransactionBlockResponseOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowRawEffects:     true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
				ShowEvents:         true,
			},
		},
	)
	require.NoError(t, err)

	require.NoError(t, err)
	targetGasCostSummary := iotajsonrpc.GasCostSummary{
		ComputationCost:         iotajsonrpc.NewBigInt(750000),
		StorageCost:             iotajsonrpc.NewBigInt(32383600),
		StorageRebate:           iotajsonrpc.NewBigInt(21955032),
		NonRefundableStorageFee: iotajsonrpc.NewBigInt(221768),
	}
	require.Equal(t, digest, &resp.Digest)
	targetRawTxBase64, err := base64.StdEncoding.DecodeString("AQAAAAAACgEBpqVCwrKBCI6PELxQWossTD9mgGbIy8W++ipS7CWatqOAVmEAAAAAAAEBAG85p+0UjVUsc5qkxhWSZ/qr2vghuqeSNiZr1gQzhCIAV3XJAQAAAAAgKEbgAIwWMBRZ1grRBFQ6qrSWLHa/AfKG8ubjmkxM/zoAIEnHBYEE/EtGK3r1lzrUU9QPAiTHLBd2+R8GS7k042UqAQF/3Yg8C3Qn8YzbSYxMh6SnnWvsR4PLPyGqOBa7xkzo7wDr5AEAAAAAAQEBbg3e/ArZiInAS6uWOeUSwhdmxeY2b4nmlpVtm+aVKHENAAAAAAAAAAEAERAyMjIyMjIyMjIyMjIuc3VpAQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABgEAAAAAAAAAAAAgVxiHQ5g2KLNHRkjYqkqe6Kvr6PaBYkN3PX6O1P2DOigAERAyMjIyMjIyMjIyMjIuc3VpACBXGIdDmDYos0dGSNiqSp7oq+vo9oFiQ3c9fo7U/YM6KAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIFa2lvc2sKYm9ycm93X3ZhbAEH7klqDMBNBqNFmCumaXyQxhkCDenidECMeBn3h/9m4aEIc3VpZnJlbnMHU3VpRnJlbgEHiJT6AvxvNsvEha6RRdBfJHp44iCBT7hBmrJhvYHwjzIJYnVsbHNoYXJrCUJ1bGxzaGFyawADAQAAAQEAAQIAAGpuoUDgld3YL3x0WQUFSzIDEp3QSgnQN1QWwxFhky0tC2ZyZWVfY2xhaW1zCmZyZWVfY2xhaW0BB+5JagzATQajRZgrpml8kMYZAg3p4nRAjHgZ94f/ZuGhCHN1aWZyZW5zB1N1aUZyZW4BB4iU+gL8bzbLxIWukUXQXyR6eOIggU+4QZqyYb2B8I8yCWJ1bGxzaGFyawlCdWxsc2hhcmsABQEDAAEEAAMAAAAAAQUAAQYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACBWtpb3NrCnJldHVybl92YWwBB+5JagzATQajRZgrpml8kMYZAg3p4nRAjHgZ94f/ZuGhCHN1aWZyZW5zB1N1aUZyZW4BB4iU+gL8bzbLxIWukUXQXyR6eOIggU+4QZqyYb2B8I8yCWJ1bGxzaGFyawlCdWxsc2hhcmsAAwEAAAMAAAAAAwAAAQAA2sImUutAC+sfXiEmRZyuju3BFrc7itYLcePo1/2zF+IMZGlyZWN0X3NldHVwEnNldF90YXJnZXRfYWRkcmVzcwAEAQQAAgEAAQcAAQYAANrCJlLrQAvrH14hJkWcro7twRa3O4rWC3Hj6Nf9sxfiDGRpcmVjdF9zZXR1cBJzZXRfcmV2ZXJzZV9sb29rdXAAAgEEAAEIAAEBAgEAAQkAVxiHQ5g2KLNHRkjYqkqe6Kvr6PaBYkN3PX6O1P2DOigBAIV+3vABgFUzNcciYyljcM6zXwvwuD9FeVw6JU3rDUD/YO8BAAAAACBmxGapu4poDXYNHxLCokFFdgFBwBhoQW8vcK8+XuklpFcYh0OYNiizR0ZI2KpKnuir6+j2gWJDdz1+jtT9gzoo7gIAAAAAAADA8MQAAAAAAAABYQBao7U4xuiDfVJM+YnHs7cBOs9VJJVriNBdHr7neIyT+M9tzPcRbANj2P9q2s21wtgIiNtayH6IAAhgFEhKsEANMFE7Y3jZzVZy0dJdgxaL8YB9JBE0745Io7/8t/XlJ3w=")
	require.NoError(t, err)
	require.Equal(t, targetRawTxBase64, resp.RawTransaction.Data())
	require.True(t, resp.Effects.Data.IsSuccess())
	require.Equal(t, int64(183), resp.Effects.Data.V1.ExecutedEpoch.Int64())
	require.Equal(t, targetGasCostSummary, resp.Effects.Data.V1.GasUsed)
	require.Equal(t, int64(11178568), resp.Effects.Data.GasFee())
	// TODO check all the fields
}

func TestMultiGetObjects(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	coins, err := api.GetCoins(
		context.TODO(), iotaclient.GetCoinsRequest{
			Owner: testAddress,
			Limit: 1,
		},
	)
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		t.Log("Warning: No Object Id for test.")
		return
	}

	obj := coins.Data[0].CoinObjectID
	objs := []*iotago.ObjectID{obj, obj}
	resp, err := api.MultiGetObjects(
		context.Background(), iotaclient.MultiGetObjectsRequest{
			ObjectIDs: objs,
			Options: &iotajsonrpc.IotaObjectDataOptions{
				ShowType:                true,
				ShowOwner:               true,
				ShowContent:             true,
				ShowDisplay:             true,
				ShowBcs:                 true,
				ShowPreviousTransaction: true,
				ShowStorageRebate:       true,
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, len(objs), len(resp))
	require.Equal(t, resp[0], resp[1])
}

func TestMultiGetTransactionBlocks(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)

	resp, err := client.MultiGetTransactionBlocks(
		context.Background(),
		iotaclient.MultiGetTransactionBlocksRequest{
			Digests: []*iotago.Digest{
				iotago.MustNewDigest("6A3ckipsEtBSEC5C53AipggQioWzVDbs9NE1SPvqrkJr"),
				iotago.MustNewDigest("8AL88Qgk7p6ny3MkjzQboTvQg9SEoWZq4rknEPeXQdH5"),
			},
			Options: &iotajsonrpc.IotaTransactionBlockResponseOptions{
				ShowEffects: true,
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, resp, 2)
	require.Equal(t, "6A3ckipsEtBSEC5C53AipggQioWzVDbs9NE1SPvqrkJr", resp[0].Digest.String())
	require.Equal(t, "8AL88Qgk7p6ny3MkjzQboTvQg9SEoWZq4rknEPeXQdH5", resp[1].Digest.String())
}

func TestTryGetPastObject(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	// there is no software-level guarantee/SLA that objects with past versions can be retrieved by this API
	resp, err := api.TryGetPastObject(
		context.Background(), iotaclient.TryGetPastObjectRequest{
			ObjectID: iotago.MustObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
			Version:  187584506,
			Options: &iotajsonrpc.IotaObjectDataOptions{
				ShowType:  true,
				ShowOwner: true,
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp.Data.VersionNotFound)
}

func TestTryMultiGetPastObjects(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	req := []*iotajsonrpc.IotaGetPastObjectRequest{
		{
			ObjectId: iotago.MustObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
			Version:  iotajsonrpc.NewBigInt(187584506),
		},
		{
			ObjectId: iotago.MustObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
			Version:  iotajsonrpc.NewBigInt(187584500),
		},
	}
	// there is no software-level guarantee/SLA that objects with past versions can be retrieved by this API
	resp, err := api.TryMultiGetPastObjects(
		context.Background(), iotaclient.TryMultiGetPastObjectsRequest{
			PastObjects: req,
			Options: &iotajsonrpc.IotaObjectDataOptions{
				ShowType:  true,
				ShowOwner: true,
			},
		},
	)
	require.NoError(t, err)
	for _, data := range resp {
		require.NotNil(t, data.Data.VersionNotFound)
	}
}
