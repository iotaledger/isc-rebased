package iotaclient_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaconn"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotago/serialization"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
)

func TestGetDynamicFieldObject(t *testing.T) {
	t.Skip("FIXME")
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	parentObjectID, err := iotago.AddressFromHex("0x1719957d7a2bf9d72459ff0eab8e600cbb1991ef41ddd5b4a8c531035933d256")
	require.NoError(t, err)
	type args struct {
		ctx            context.Context
		parentObjectID *iotago.ObjectID
		name           *iotago.DynamicFieldName
	}
	tests := []struct {
		name    string
		args    args
		want    *iotajsonrpc.IotaObjectResponse
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				ctx:            context.TODO(),
				parentObjectID: parentObjectID,
				name: &iotago.DynamicFieldName{
					Type:  "address",
					Value: "0xf9ed7d8de1a6c44d703b64318a1cc687c324fdec35454281035a53ea3ba1a95a",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.GetDynamicFieldObject(
					tt.args.ctx, iotaclient.GetDynamicFieldObjectRequest{
						ParentObjectID: tt.args.parentObjectID,
						Name:           tt.args.name,
					},
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetDynamicFieldObject() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}

func TestGetDynamicFields(t *testing.T) {
	t.Skip("refactor me: provide address of deepbook object")

	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	limit := 5
	type args struct {
		ctx            context.Context
		parentObjectID *iotago.ObjectID
		cursor         *iotago.ObjectID
		limit          *uint
	}
	tests := []struct {
		name    string
		args    args
		want    *iotajsonrpc.DynamicFieldPage
		wantErr error
	}{
		{
			name: "a deepbook shared object",
			args: args{
				ctx:            context.TODO(),
				parentObjectID: iotago.MustAddressFromHex("0xa9d09452bba939b3172c0242d022274845cfe4e58648b73dd33b3d5b823dc8ae"),
				cursor:         nil,
				limit:          func() *uint { tmpLimit := uint(limit); return &tmpLimit }(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := client.GetDynamicFields(
					tt.args.ctx, iotaclient.GetDynamicFieldsRequest{
						ParentObjectID: tt.args.parentObjectID,
						Cursor:         tt.args.cursor,
						Limit:          tt.args.limit,
					},
				)
				require.ErrorIs(t, err, tt.wantErr)
				// object ID is '0x4405b50d791fd3346754e8171aaab6bc2ed26c2c46efdd033c14b30ae507ac33'
				// it has 'internal_nodes' field in type '0x2::table::Table<u64, 0xdee9::critbit::InternalNode'
				require.Len(t, got.Data, limit)
				for _, field := range got.Data {
					require.Equal(t, "u64", field.Name.Type)
					require.Equal(t, "0xdee9::critbit::InternalNode", field.ObjectType)
				}
			},
		)
	}
}

func TestGetOwnedObjects(t *testing.T) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	signer := iotasigner.NewSignerByIndex(testSeed, iotasigner.KeySchemeFlagEd25519, 0)
	t.Run(
		"struct tag", func(t *testing.T) {
			structTag, err := iotago.StructTagFromString("0x2::coin::Coin<0x2::iota::IOTA>")
			require.NoError(t, err)
			query := iotajsonrpc.IotaObjectResponseQuery{
				Filter: &iotajsonrpc.IotaObjectDataFilter{
					StructType: structTag,
				},
				Options: &iotajsonrpc.IotaObjectDataOptions{
					ShowType:    true,
					ShowContent: true,
				},
			}
			limit := uint(10)
			objs, err := client.GetOwnedObjects(
				context.Background(), iotaclient.GetOwnedObjectsRequest{
					Address: signer.Address(),
					Query:   &query,
					Limit:   &limit,
				},
			)

			require.NoError(t, err)
			require.Greater(t, len(objs.Data), 1)
		},
	)

	t.Run(
		"move module", func(t *testing.T) {
			query := iotajsonrpc.IotaObjectResponseQuery{
				Filter: &iotajsonrpc.IotaObjectDataFilter{
					AddressOwner: signer.Address(),
				},
				Options: &iotajsonrpc.IotaObjectDataOptions{
					ShowType:    true,
					ShowContent: true,
				},
			}
			limit := uint(9)
			objs, err := client.GetOwnedObjects(
				context.Background(), iotaclient.GetOwnedObjectsRequest{
					Address: signer.Address(),
					Query:   &query,
					Limit:   &limit,
				},
			)
			require.NoError(t, err)
			require.Greater(t, len(objs.Data), 1)
		},
	)
	// query := iotajsonrpc.IotaObjectResponseQuery{
	// 	Filter: &iotajsonrpc.IotaObjectDataFilter{
	// 		StructType: "0x2::coin::Coin<0x2::iota::IOTA>",
	// 	},
	// 	Options: &iotajsonrpc.IotaObjectDataOptions{
	// 		ShowType:    true,
	// 		ShowContent: true,
	// 	},
	// }
	// limit := uint(2)
	// objs, err := client.GetOwnedObjects(
	// 	context.Background(), iotaclient.GetOwnedObjectsRequest{
	// 		Address: signer.Address(),
	// 		Query:   &query,
	// 		Cursor:  nil,
	// 		Limit:   &limit,
	// 	},
	// )
	// require.NoError(t, err)
	// require.GreaterOrEqual(t, len(objs.Data), int(limit))
	// require.NoError(t, err)
	// var fields iotajsonrpc.CoinFields
	// err = json.Unmarshal(objs.Data[1].Data.Content.Data.MoveObject.Fields, &fields)

	// require.NoError(t, err)
	// require.Equal(t, "1000000000", fields.Balance.String())
}

func TestQueryEvents(t *testing.T) {
	t.Skip("refactor me: provide address of deepbook object")
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	limit := 10

	type args struct {
		ctx             context.Context
		query           *iotajsonrpc.EventFilter
		cursor          *iotajsonrpc.EventId
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *iotajsonrpc.EventPage
		wantErr error
	}{
		{
			name: "event in deepbook.batch_cancel_order()",
			args: args{
				ctx: context.TODO(),
				query: &iotajsonrpc.EventFilter{
					Sender: iotago.MustAddressFromHex("0xf0f13f7ef773c6246e87a8f059a684d60773f85e992e128b8272245c38c94076"),
				},
				cursor:          nil,
				limit:           func() *uint { tmpLimit := uint(limit); return &tmpLimit }(),
				descendingOrder: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.QueryEvents(
					tt.args.ctx,
					iotaclient.QueryEventsRequest{
						Query:           tt.args.query,
						Cursor:          tt.args.cursor,
						Limit:           tt.args.limit,
						DescendingOrder: tt.args.descendingOrder,
					},
				)
				require.ErrorIs(t, err, tt.wantErr)
				require.Len(t, got.Data, int(limit))

				for _, event := range got.Data {
					// FIXME we should change other filter to, so we can verify each fields of event more detailed.
					require.Equal(
						t,
						iotago.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
						event.PackageId,
					)
					require.Equal(t, "clob_v2", event.TransactionModule)
					require.Equal(t, tt.args.query.Sender, event.Sender)
				}
			},
		)
	}
}

func TestQueryTransactionBlocks(t *testing.T) {
	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	limit := uint(10)
	type args struct {
		ctx             context.Context
		query           *iotajsonrpc.IotaTransactionBlockResponseQuery
		cursor          *iotago.TransactionDigest
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *iotajsonrpc.TransactionBlocksPage
		wantErr bool
	}{
		{
			name: "test for queryTransactionBlocks",
			args: args{
				ctx: context.TODO(),
				query: &iotajsonrpc.IotaTransactionBlockResponseQuery{
					Filter: &iotajsonrpc.TransactionFilter{
						FromAddress: testAddress,
					},
					Options: &iotajsonrpc.IotaTransactionBlockResponseOptions{
						ShowInput:   true,
						ShowEffects: true,
					},
				},
				cursor:          nil,
				limit:           &limit,
				descendingOrder: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := api.QueryTransactionBlocks(
					tt.args.ctx,
					iotaclient.QueryTransactionBlocksRequest{
						Query:           tt.args.query,
						Cursor:          tt.args.cursor,
						Limit:           tt.args.limit,
						DescendingOrder: tt.args.descendingOrder,
					},
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("QueryTransactionBlocks() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}

func TestResolveNameServiceAddress(t *testing.T) {
	t.Skip()

	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	addr, err := api.ResolveNameServiceAddress(context.Background(), "2222.iotax")
	require.NoError(t, err)
	require.Equal(t, "0x6174c5bd8ab9bf492e159a64e102de66429cfcde4fa883466db7b03af28b3ce9", addr.String())

	_, err = api.ResolveNameServiceAddress(context.Background(), "2222.iotajjzzww")
	require.ErrorContains(t, err, "not found")
}

func TestResolveNameServiceNames(t *testing.T) {
	t.Skip("Fails with 'Method not found'")

	api := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	owner := iotago.MustAddressFromHex("0x57188743983628b3474648d8aa4a9ee8abebe8f6816243773d7e8ed4fd833a28")
	namePage, err := api.ResolveNameServiceNames(
		context.Background(), iotaclient.ResolveNameServiceNamesRequest{
			Owner: owner,
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, namePage.Data)
	t.Log(namePage.Data)

	owner = iotago.MustAddressFromHex("0x57188743983628b3474648d8aa4a9ee8abebe8f681")
	namePage, err = api.ResolveNameServiceNames(
		context.Background(), iotaclient.ResolveNameServiceNamesRequest{
			Owner: owner,
		},
	)
	require.NoError(t, err)
	require.Empty(t, namePage.Data)
}

func TestSubscribeEvent(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testlogger.NewLogger(t)
	api, err := iotaclient.NewWebsocket(ctx, iotaconn.AlphanetWebsocketEndpointURL, log)
	require.NoError(t, err)

	type args struct {
		ctx      context.Context
		filter   *iotajsonrpc.EventFilter
		resultCh chan *iotajsonrpc.IotaEvent
	}
	tests := []struct {
		name    string
		args    args
		want    *iotajsonrpc.EventPage
		wantErr bool
	}{
		{
			name: "test for filter events",
			args: args{
				ctx: context.TODO(),
				filter: &iotajsonrpc.EventFilter{
					Package: iotago.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
				},
				resultCh: make(chan *iotajsonrpc.IotaEvent),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := api.SubscribeEvent(
					tt.args.ctx,
					tt.args.filter,
					tt.args.resultCh,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("SubscribeEvent() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				cnt := 0
				for results := range tt.args.resultCh {
					fmt.Println("results: ", results)
					// FIXME we need to check finite number request in details
					cnt++
					if cnt > 3 {
						break
					}
				}
			},
		)
	}
}

func TestSubscribeTransaction(t *testing.T) {
	t.Skip()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := testlogger.NewLogger(t)
	api, err := iotaclient.NewWebsocket(ctx, iotaconn.AlphanetWebsocketEndpointURL, log)
	require.NoError(t, err)

	type args struct {
		ctx      context.Context
		filter   *iotajsonrpc.TransactionFilter
		resultCh chan *serialization.TagJson[iotajsonrpc.IotaTransactionBlockEffects]
	}
	tests := []struct {
		name    string
		args    args
		want    *iotajsonrpc.IotaTransactionBlockEffects
		wantErr bool
	}{
		{
			name: "test for filter transaction",
			args: args{
				ctx: context.TODO(),
				filter: &iotajsonrpc.TransactionFilter{
					MoveFunction: &iotajsonrpc.TransactionFilterMoveFunction{
						Package: *iotago.MustPackageIDFromHex("0x2c68443db9e8c813b194010c11040a3ce59f47e4eb97a2ec805371505dad7459"),
					},
				},
				resultCh: make(chan *serialization.TagJson[iotajsonrpc.IotaTransactionBlockEffects]),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := api.SubscribeTransaction(
					tt.args.ctx,
					tt.args.filter,
					tt.args.resultCh,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("SubscribeTransaction() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				cnt := 0
				for results := range tt.args.resultCh {
					fmt.Println("results: ", results.Data.V1)
					// FIXME we need to check finite number request in details
					cnt++
					if cnt > 3 {
						break
					}
				}
			},
		)
	}
}
