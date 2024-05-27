package sui_test

import (
	"context"
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/iotaledger/wasp/sui-go/models"
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/sui/conn"
	"github.com/iotaledger/wasp/sui-go/sui_signer"
	"github.com/iotaledger/wasp/sui-go/sui_types"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/require"
)

func TestGetChainIdentifier(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	chainID, err := client.GetChainIdentifier(context.Background())
	require.NoError(t, err)
	require.Equal(t, conn.ChainIdentifierSuiMainnet, chainID)
}

func TestGetCheckpoint(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	checkpoint, err := client.GetCheckpoint(context.Background(), models.NewSafeSuiBigInt(uint64(1000)))
	require.NoError(t, err)
	targetCheckpoint := &models.Checkpoint{
		Epoch:                    models.NewSafeSuiBigInt(uint64(0)),
		SequenceNumber:           models.NewSafeSuiBigInt(uint64(1000)),
		Digest:                   *sui_types.MustNewDigest("BE4JixC94sDtCgHJZruyk7QffZnWDFvM2oFjC8XtChET"),
		NetworkTotalTransactions: models.NewSafeSuiBigInt(uint64(1001)),
		PreviousDigest:           sui_types.MustNewDigest("41nPNZWHvvajmBQjX3GbppsgGZDEB6DhN4UxPkjSYRRj"),
		EpochRollingGasCostSummary: models.GasCostSummary{
			ComputationCost:         models.NewSafeSuiBigInt(uint64(0)),
			StorageCost:             models.NewSafeSuiBigInt(uint64(0)),
			StorageRebate:           models.NewSafeSuiBigInt(uint64(0)),
			NonRefundableStorageFee: models.NewSafeSuiBigInt(uint64(0)),
		},
		TimestampMs:           models.NewSafeSuiBigInt(uint64(1681393657483)),
		Transactions:          []*sui_types.Digest{sui_types.MustNewDigest("9NnjyPG8V2TPCSbNE391KDyge42AwV3vUD7aNtQQ9eqS")},
		CheckpointCommitments: []sui_types.CheckpointCommitment{},
		ValidatorSignature:    *sui_types.MustNewBase64Data("r8/5+Rm7niIlndcnvjSJ/vZLPrH3xY/ePGYTvrVbTascoQSpS+wsGlC+bQBpzIwA"),
	}
	require.Equal(t, targetCheckpoint, checkpoint)
}

func TestGetCheckpoints(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	cursor := models.NewSafeSuiBigInt(uint64(999))
	limit := uint64(2)
	checkpointPage, err := client.GetCheckpoints(context.Background(), &cursor, &limit, false)
	require.NoError(t, err)
	targetCheckpoints := []*models.Checkpoint{
		{
			Epoch:                    models.NewSafeSuiBigInt(uint64(0)),
			SequenceNumber:           models.NewSafeSuiBigInt(uint64(1000)),
			Digest:                   *sui_types.MustNewDigest("BE4JixC94sDtCgHJZruyk7QffZnWDFvM2oFjC8XtChET"),
			NetworkTotalTransactions: models.NewSafeSuiBigInt(uint64(1001)),
			PreviousDigest:           sui_types.MustNewDigest("41nPNZWHvvajmBQjX3GbppsgGZDEB6DhN4UxPkjSYRRj"),
			EpochRollingGasCostSummary: models.GasCostSummary{
				ComputationCost:         models.NewSafeSuiBigInt(uint64(0)),
				StorageCost:             models.NewSafeSuiBigInt(uint64(0)),
				StorageRebate:           models.NewSafeSuiBigInt(uint64(0)),
				NonRefundableStorageFee: models.NewSafeSuiBigInt(uint64(0)),
			},
			TimestampMs:           models.NewSafeSuiBigInt(uint64(1681393657483)),
			Transactions:          []*sui_types.Digest{sui_types.MustNewDigest("9NnjyPG8V2TPCSbNE391KDyge42AwV3vUD7aNtQQ9eqS")},
			CheckpointCommitments: []sui_types.CheckpointCommitment{},
			ValidatorSignature:    *sui_types.MustNewBase64Data("r8/5+Rm7niIlndcnvjSJ/vZLPrH3xY/ePGYTvrVbTascoQSpS+wsGlC+bQBpzIwA"),
		},
		{
			Epoch:                    models.NewSafeSuiBigInt(uint64(0)),
			SequenceNumber:           models.NewSafeSuiBigInt(uint64(1001)),
			Digest:                   *sui_types.MustNewDigest("8umKe5Ae2TAH5ySw2zeEua8cTeeTFZV8F3GfFViZ5cq3"),
			NetworkTotalTransactions: models.NewSafeSuiBigInt(uint64(1002)),
			PreviousDigest:           sui_types.MustNewDigest("BE4JixC94sDtCgHJZruyk7QffZnWDFvM2oFjC8XtChET"),
			EpochRollingGasCostSummary: models.GasCostSummary{
				ComputationCost:         models.NewSafeSuiBigInt(uint64(0)),
				StorageCost:             models.NewSafeSuiBigInt(uint64(0)),
				StorageRebate:           models.NewSafeSuiBigInt(uint64(0)),
				NonRefundableStorageFee: models.NewSafeSuiBigInt(uint64(0)),
			},
			TimestampMs:           models.NewSafeSuiBigInt(uint64(1681393661034)),
			Transactions:          []*sui_types.Digest{sui_types.MustNewDigest("9muLz7ZTocpBTdSo5Ak7ZxzEpfzywr6Y12Hj3AdT8dvV")},
			CheckpointCommitments: []sui_types.CheckpointCommitment{},
			ValidatorSignature:    *sui_types.MustNewBase64Data("jG5ViKThziBpnJnOw9dVdjIrv2IHhCrn8ZhvI1gUS2X1t90aRqhnLF6+WbS1S2WT"),
		},
	}
	require.Len(t, checkpointPage.Data, 2)
	require.Equal(t, checkpointPage.Data, targetCheckpoints)
	require.Equal(t, true, checkpointPage.HasNextPage)
	require.Equal(t, models.NewSafeSuiBigInt(uint64(1001)), *checkpointPage.NextCursor)
}

func TestGetEvents(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	digest, err := sui_types.NewDigest("3vVi8XZgNpzQ34PFgwJTQqWtPMU84njcBX1EUxUHhyDk")
	require.NoError(t, err)
	events, err := client.GetEvents(context.Background(), digest)
	require.NoError(t, err)
	require.Len(t, events, 1)
	for _, event := range events {
		require.Equal(t, digest, &event.Id.TxDigest)
		require.Equal(
			t,
			sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			&event.PackageId,
		)
		require.Equal(t, "clob_v2", event.TransactionModule)
		require.Equal(
			t,
			sui_types.MustSuiAddressFromHex("0xf0f13f7ef773c6246e87a8f059a684d60773f85e992e128b8272245c38c94076"),
			&event.Sender,
		)
		require.Equal(
			t,
			"0xdee9::clob_v2::OrderPlaced<0x2::sui::SUI, 0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf::coin::COIN>",
			event.Type,
		)
		targetBcsBase85 := base58.Decode("yNS5iDS3Gvdo3DhXdtFpuTS12RrSiNkrvjcm2rejntCuqWjF1DdwnHgjowdczAkR18LQHcBqbX2tWL76rys9rTCzG6vm7Tg34yqUkpFSMqNkcS6cfWbN8SdVsxn5g4ZEQotdBgEFn8yN7hVZ7P1MKvMwWf")
		require.Equal(t, targetBcsBase85, event.Bcs.Data())
		// TODO check ParsedJson map
	}
}

func TestGetLatestCheckpointSequenceNumber(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	sequenceNumber, err := client.GetLatestCheckpointSequenceNumber(context.Background())
	require.NoError(t, err)
	num, err := strconv.Atoi(sequenceNumber)
	require.NoError(t, err)
	require.Greater(t, num, 34317507)
}

func TestGetObject(t *testing.T) {
	type args struct {
		ctx   context.Context
		objID *sui_types.ObjectID
	}
	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	coins, err := api.GetCoins(context.TODO(), sui_signer.TEST_ADDRESS, nil, nil, 1)
	require.NoError(t, err)

	tests := []struct {
		name    string
		api     *sui.ImplSuiAPI
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
					tt.args.ctx, tt.args.objID, &models.SuiObjectDataOptions{
						ShowType:                true,
						ShowOwner:               true,
						ShowContent:             true,
						ShowDisplay:             true,
						ShowBcs:                 true,
						ShowPreviousTransaction: true,
						ShowStorageRebate:       true,
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

func TestGetTotalTransactionBlocks(t *testing.T) {
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	res, err := api.GetTotalTransactionBlocks(context.Background())
	require.NoError(t, err)
	t.Log(res)
}

func TestGetTransactionBlock(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)
	digest, err := sui_types.NewDigest("D1TM8Esaj3G9xFEDirqMWt9S7HjJXFrAGYBah1zixWTL")
	require.NoError(t, err)
	resp, err := client.GetTransactionBlock(
		context.Background(), digest, &models.SuiTransactionBlockResponseOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowRawEffects:     true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
		},
	)
	require.NoError(t, err)

	require.NoError(t, err)
	targetGasCostSummary := models.GasCostSummary{
		ComputationCost:         models.NewSafeSuiBigInt(uint64(750000)),
		StorageCost:             models.NewSafeSuiBigInt(uint64(32383600)),
		StorageRebate:           models.NewSafeSuiBigInt(uint64(21955032)),
		NonRefundableStorageFee: models.NewSafeSuiBigInt(uint64(221768)),
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
	api := sui.NewSuiClient(conn.DevnetEndpointUrl)
	coins, err := api.GetCoins(context.TODO(), sui_signer.TEST_ADDRESS, nil, nil, 1)
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		t.Log("Warning: No Object Id for test.")
		return
	}

	obj := coins.Data[0].CoinObjectID
	objs := []*sui_types.ObjectID{obj, obj}
	resp, err := api.MultiGetObjects(
		context.Background(), objs, &models.SuiObjectDataOptions{
			ShowType:                true,
			ShowOwner:               true,
			ShowContent:             true,
			ShowDisplay:             true,
			ShowBcs:                 true,
			ShowPreviousTransaction: true,
			ShowStorageRebate:       true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, len(objs), len(resp))
	require.Equal(t, resp[0], resp[1])
}

func TestMultiGetTransactionBlocks(t *testing.T) {
	client := sui.NewSuiClient(conn.TestnetEndpointUrl)

	resp, err := client.MultiGetTransactionBlocks(
		context.Background(),
		[]*sui_types.Digest{
			sui_types.MustNewDigest("6A3ckipsEtBSEC5C53AipggQioWzVDbs9NE1SPvqrkJr"),
			sui_types.MustNewDigest("8AL88Qgk7p6ny3MkjzQboTvQg9SEoWZq4rknEPeXQdH5"),
		},
		&models.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		},
	)
	require.NoError(t, err)
	require.Len(t, resp, 2)
	require.Equal(t, "6A3ckipsEtBSEC5C53AipggQioWzVDbs9NE1SPvqrkJr", resp[0].Digest.String())
	require.Equal(t, "8AL88Qgk7p6ny3MkjzQboTvQg9SEoWZq4rknEPeXQdH5", resp[1].Digest.String())
}

func TestTryGetPastObject(t *testing.T) {
	api := sui.NewSuiClient(conn.MainnetEndpointUrl)
	objId, err := sui_types.ObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f")
	require.NoError(t, err)
	obj, err := api.TryGetPastObject(context.Background(), objId, 187584506, &models.SuiObjectDataOptions{
		ShowType:  true,
		ShowOwner: true,
	})
	require.NoError(t, err)
	require.Equal(t, objId, obj.Data.VersionFound.ObjectID)
	require.Equal(t, sui_types.MustNewDigest("61cY5vK1LXMo6QsdihMPKr5aXtRN31DA7pFLX2LzBQTB"), obj.Data.VersionFound.Digest)
	require.Equal(t, "0x1eabed72c53feb3805120a081dc15963c204dc8d091542592abaf7a35689b2fb::config::GlobalConfig", *obj.Data.VersionFound.Type)
}

func TestTryMultiGetPastObjects(t *testing.T) {
	api := sui.NewSuiClient(conn.MainnetEndpointUrl)
	req := []*models.SuiGetPastObjectRequest{
		{
			ObjectId: sui_types.MustObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
			Version:  models.NewSafeSuiBigInt(uint64(187584506)),
		},
		{
			ObjectId: sui_types.MustObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"),
			Version:  models.NewSafeSuiBigInt(uint64(187584500)),
		},
	}
	resp, err := api.TryMultiGetPastObjects(context.Background(), req, &models.SuiObjectDataOptions{
		ShowType:  true,
		ShowOwner: true,
	})
	require.NoError(t, err)
	require.Equal(t, sui_types.MustObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"), resp[0].Data.VersionFound.ObjectID)
	require.Equal(t, sui_types.MustNewDigest("61cY5vK1LXMo6QsdihMPKr5aXtRN31DA7pFLX2LzBQTB"), resp[0].Data.VersionFound.Digest)
	require.Equal(t, "0x1eabed72c53feb3805120a081dc15963c204dc8d091542592abaf7a35689b2fb::config::GlobalConfig", *resp[0].Data.VersionFound.Type)

	require.Equal(t, sui_types.MustObjectIDFromHex("0xdaa46292632c3c4d8f31f23ea0f9b36a28ff3677e9684980e4438403a67a3d8f"), resp[1].Data.VersionFound.ObjectID)
	require.Equal(t, sui_types.MustNewDigest("BeE8rwAHdUvgrFTiGaRXHKAPdFMucxKFaZDZNSGLQ2DW"), resp[1].Data.VersionFound.Digest)
	require.Equal(t, "0x1eabed72c53feb3805120a081dc15963c204dc8d091542592abaf7a35689b2fb::config::GlobalConfig", *resp[1].Data.VersionFound.Type)
}
