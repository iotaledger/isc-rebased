package vmtxbuilder_test

import (
	"context"
	"testing"

	"github.com/iotaledger/wasp/clients"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaconn"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient/iscmoveclienttest"
	"github.com/iotaledger/wasp/packages/testutil/l1starter"
	"github.com/iotaledger/wasp/packages/transaction/transactiontest"
	"github.com/iotaledger/wasp/packages/util/bcs"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/vmtxbuilder"
)

func TestTxBuilderBasic(t *testing.T) {
	client := newLocalnetClient()
	signer := newSignerWithFunds(t, testSeed, 0)
	iscPackage := l1starter.DeployISCContracts(client, cryptolib.SignerToIotaSigner(signer))

	anchor, gasObjRef := iscmoveclienttest.StartNewChain(t, client, signer, iscPackage, 100)

	stateAnchor := isc.NewStateAnchor(anchor, signer.Address(), iscPackage)
	txb := vmtxbuilder.NewAnchorTransactionBuilder(iscPackage, &stateAnchor, signer.Address())

	req1 := createIscmoveReq(t, client, signer, iscPackage, anchor)
	txb.ConsumeRequest(req1)
	req2 := createIscmoveReq(t, client, signer, iscPackage, anchor)
	txb.ConsumeRequest(req2)
	pt := txb.BuildTransactionEssence(transactiontest.RandomStateMetadata().Bytes())

	tx := iotago.NewProgrammable(
		signer.Address().AsIotaAddress(),
		pt,
		[]*iotago.ObjectRef{gasObjRef},
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	txnBytes, err := bcs.Marshal(&tx)
	require.NoError(t, err)

	txnResponse, err := client.SignAndExecuteTransaction(
		context.Background(),
		cryptolib.SignerToIotaSigner(signer),
		txnBytes,
		&iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true, ShowObjectChanges: true},
	)
	require.NoError(t, err)
	require.True(t, txnResponse.Effects.Data.IsSuccess())

	getObjReq1, _ := client.GetObject(context.Background(), iotaclient.GetObjectRequest{ObjectID: req1.RequestRef().ObjectID, Options: &iotajsonrpc.IotaObjectDataOptions{ShowContent: true}})
	require.NotNil(t, getObjReq1.Error.Data.Deleted)
	getObjReq2, _ := client.GetObject(context.Background(), iotaclient.GetObjectRequest{ObjectID: req2.RequestRef().ObjectID})
	require.NotNil(t, getObjReq2.Error.Data.Deleted)
}

func TestTxBuilderSendAssetsAndRequest(t *testing.T) {
	client := newLocalnetClient()
	signer := newSignerWithFunds(t, testSeed, 0)
	recipient := newSignerWithFunds(t, testSeed, 1)
	iscPackage := l1starter.DeployISCContracts(client, cryptolib.SignerToIotaSigner(signer))

	anchor, gasObjRef := iscmoveclienttest.StartNewChain(t, client, signer, iscPackage, 100)

	stateAnchor := isc.NewStateAnchor(anchor, signer.Address(), iscPackage)
	txb1 := vmtxbuilder.NewAnchorTransactionBuilder(iscPackage, &stateAnchor, signer.Address())

	req1 := createIscmoveReq(t, client, signer, iscPackage, anchor)
	txb1.ConsumeRequest(req1)

	ptb1 := txb1.BuildTransactionEssence(transactiontest.RandomStateMetadata().Bytes())

	tx1 := iotago.NewProgrammable(
		signer.Address().AsIotaAddress(),
		ptb1,
		[]*iotago.ObjectRef{gasObjRef},
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	txnBytes1, err := bcs.Marshal(&tx1)
	require.NoError(t, err)

	txnResponse1, err := client.SignAndExecuteTransaction(
		context.Background(),
		cryptolib.SignerToIotaSigner(signer),
		txnBytes1,
		&iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true, ShowObjectChanges: true},
	)
	require.NoError(t, err)
	require.True(t, txnResponse1.Effects.Data.IsSuccess())

	getObjReq1, _ := client.GetObject(context.Background(), iotaclient.GetObjectRequest{ObjectID: req1.RequestRef().ObjectID, Options: &iotajsonrpc.IotaObjectDataOptions{ShowContent: true}})
	require.NotNil(t, getObjReq1.Error.Data.Deleted)

	// reset
	tmp, err := client.UpdateObjectRef(context.Background(), &anchor.ObjectRef)
	require.NoError(t, err)
	anchor.ObjectRef = *tmp
	txb2 := vmtxbuilder.NewAnchorTransactionBuilder(iscPackage, &stateAnchor, signer.Address())

	txb2.SendAssets(recipient.Address().AsIotaAddress(), isc.NewAssets(1))

	req2 := createIscmoveReq(t, client, signer, iscPackage, anchor)
	txb2.ConsumeRequest(req2)

	pt2 := txb2.BuildTransactionEssence(transactiontest.RandomStateMetadata().Bytes())

	gasObjRef, err = client.UpdateObjectRef(context.Background(), gasObjRef)
	require.NoError(t, err)

	tx2 := iotago.NewProgrammable(
		signer.Address().AsIotaAddress(),
		pt2,
		[]*iotago.ObjectRef{gasObjRef},
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	txnBytes2, err := bcs.Marshal(&tx2)
	require.NoError(t, err)

	txnResponse2, err := client.SignAndExecuteTransaction(
		context.Background(),
		cryptolib.SignerToIotaSigner(signer),
		txnBytes2,
		&iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true, ShowObjectChanges: true},
	)
	require.NoError(t, err)
	require.True(t, txnResponse2.Effects.Data.IsSuccess())

	getObjReq2, _ := client.GetObject(context.Background(), iotaclient.GetObjectRequest{ObjectID: req2.RequestRef().ObjectID})
	require.NotNil(t, getObjReq2.Error.Data.Deleted)
}

func TestTxBuilderSendCrossChainRequest(t *testing.T) {
	t.Skip("we may suppress the SendCrossChainRequest behavior for now")
	client := newLocalnetClient()
	signer := newSignerWithFunds(t, testSeed, 0)
	iscPackage := l1starter.DeployISCContracts(client, cryptolib.SignerToIotaSigner(signer))

	anchor1, gasObjRef1 := iscmoveclienttest.StartNewChain(t, client, signer, iscPackage, 100)
	anchor2, gasObjRef2 := iscmoveclienttest.StartNewChain(t, client, signer, iscPackage, 100)

	stateAnchor1 := isc.NewStateAnchor(anchor1, signer.Address(), iscPackage)
	txb1 := vmtxbuilder.NewAnchorTransactionBuilder(iscPackage, &stateAnchor1, signer.Address())

	req1 := createIscmoveReq(t, client, signer, iscPackage, anchor1)
	txb1.ConsumeRequest(req1)

	pt1 := txb1.BuildTransactionEssence(transactiontest.RandomStateMetadata().Bytes())

	tx1 := iotago.NewProgrammable(
		signer.Address().AsIotaAddress(),
		pt1,
		[]*iotago.ObjectRef{gasObjRef1},
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	txnBytes1, err := bcs.Marshal(&tx1)
	require.NoError(t, err)

	txnResponse1, err := client.SignAndExecuteTransaction(
		context.Background(),
		cryptolib.SignerToIotaSigner(signer),
		txnBytes1,
		&iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true, ShowObjectChanges: true},
	)
	require.NoError(t, err)
	require.True(t, txnResponse1.Effects.Data.IsSuccess())

	getObjReq1, _ := client.GetObject(context.Background(), iotaclient.GetObjectRequest{ObjectID: req1.RequestRef().ObjectID, Options: &iotajsonrpc.IotaObjectDataOptions{ShowContent: true}})
	require.NotNil(t, getObjReq1.Error.Data.Deleted)

	// reset
	tmp, err := client.UpdateObjectRef(context.Background(), &anchor1.ObjectRef)
	require.NoError(t, err)
	anchor1.ObjectRef = *tmp
	txb2 := vmtxbuilder.NewAnchorTransactionBuilder(iscPackage, &stateAnchor1, signer.Address())

	txb2.SendCrossChainRequest(&iscPackage, anchor2.ObjectID, isc.NewAssets(1), &isc.SendMetadata{
		Message:   isc.NewMessage(isc.Hn("accounts"), isc.Hn("deposit")),
		Allowance: isc.NewAssets(1),
		GasBudget: 2,
	})

	pt2 := txb2.BuildTransactionEssence(transactiontest.RandomStateMetadata().Bytes())

	tx2 := iotago.NewProgrammable(
		signer.Address().AsIotaAddress(),
		pt2,
		[]*iotago.ObjectRef{gasObjRef2},
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	txnBytes2, err := bcs.Marshal(&tx2)
	require.NoError(t, err)

	txnResponse2, err := client.SignAndExecuteTransaction(
		context.Background(),
		cryptolib.SignerToIotaSigner(signer),
		txnBytes2,
		&iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true, ShowObjectChanges: true},
	)
	require.NoError(t, err)
	require.True(t, txnResponse2.Effects.Data.IsSuccess())
	crossChainRequestRef, err := txnResponse2.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
	require.NoError(t, err)

	stateAnchor2 := isc.NewStateAnchor(anchor2, signer.Address(), iscPackage)
	txb3 := vmtxbuilder.NewAnchorTransactionBuilder(iscPackage, &stateAnchor2, signer.Address())

	reqWithObj, err := client.L2().GetRequestFromObjectID(context.Background(), crossChainRequestRef.ObjectID)
	require.NoError(t, err)
	req3, err := isc.OnLedgerFromRequest(reqWithObj, cryptolib.NewAddressFromIota(anchor2.ObjectID))
	require.NoError(t, err)
	txb3.ConsumeRequest(req3)

	gasObjRef2, err = client.UpdateObjectRef(context.Background(), gasObjRef2)
	require.NoError(t, err)
	pt3 := txb3.BuildTransactionEssence(transactiontest.RandomStateMetadata().Bytes())

	tx3 := iotago.NewProgrammable(
		signer.Address().AsIotaAddress(),
		pt3,
		[]*iotago.ObjectRef{gasObjRef2},
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)

	txnBytes3, err := bcs.Marshal(&tx3)
	require.NoError(t, err)

	txnResponse3, err := client.SignAndExecuteTransaction(
		context.Background(),
		cryptolib.SignerToIotaSigner(signer),
		txnBytes3,
		&iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true, ShowObjectChanges: true},
	)
	require.NoError(t, err)
	require.True(t, txnResponse3.Effects.Data.IsSuccess())
}

var testSeed = []byte{50, 230, 119, 9, 86, 155, 106, 30, 245, 81, 234, 122, 116, 90, 172, 148, 59, 33, 88, 252, 134, 42, 231, 198, 208, 141, 209, 116, 78, 21, 216, 24}

func newSignerWithFunds(t *testing.T, seed []byte, index int) cryptolib.Signer {
	seed[0] = seed[0] + byte(index)
	kp := cryptolib.KeyPairFromSeed(cryptolib.Seed(seed))
	err := iotaclient.RequestFundsFromFaucet(context.TODO(), kp.Address().AsIotaAddress(), iotaconn.LocalnetFaucetURL)
	require.NoError(t, err)
	return kp
}

func newLocalnetClient() clients.L1Client {
	return clients.NewL1Client(clients.L1Config{
		APIURL:    iotaconn.LocalnetEndpointURL,
		FaucetURL: iotaconn.LocalnetFaucetURL,
	})
}

func createIscmoveReq(
	t *testing.T,
	client clients.L1Client,
	signer cryptolib.Signer,
	iscPackage iotago.Address,
	anchor *iscmove.AnchorWithRef,
) isc.OnLedgerRequest {
	err := iotaclient.RequestFundsFromFaucet(context.Background(), signer.Address().AsIotaAddress(), iotaconn.LocalnetFaucetURL)
	require.NoError(t, err)
	getCoinsRes, err := client.GetCoins(
		context.Background(),
		iotaclient.GetCoinsRequest{
			Owner: signer.Address().AsIotaAddress(),
		},
	)
	require.NoError(t, err)
	_ = getCoinsRes
	assetsBagNewRes, err := client.L2().AssetsBagNew(
		context.Background(),
		signer,
		iscPackage,
		nil,
		iotaclient.DefaultGasPrice,
		iotaclient.DefaultGasBudget,
	)
	require.NoError(t, err)
	assetsBagRef, err := assetsBagNewRes.GetCreatedObjectInfo(iscmove.AssetsBagModuleName, iscmove.AssetsBagObjectName)
	require.NoError(t, err)
	_, err = client.L2().AssetsBagPlaceCoinAmount(
		context.Background(),
		signer,
		iscPackage,
		assetsBagRef,
		getCoinsRes.Data[len(getCoinsRes.Data)-1].Ref(),
		iotajsonrpc.IotaCoinType,
		111,
		nil,
		iotaclient.DefaultGasPrice,
		iotaclient.DefaultGasBudget,
	)
	require.NoError(t, err)
	assetsBagRef, err = client.UpdateObjectRef(context.Background(), assetsBagRef)
	require.NoError(t, err)

	createAndSendRequestRes, err := client.L2().CreateAndSendRequest(
		context.Background(),
		signer,
		iscPackage,
		anchor.ObjectID,
		assetsBagRef,
		&iscmove.Message{
			Contract: uint32(isc.Hn("test_isc_contract")),
			Function: uint32(isc.Hn("test_isc_func")),
			Args:     [][]byte{[]byte("one"), []byte("two"), []byte("three")},
		},
		iscmove.NewAssets(0),
		10,
		nil,
		iotaclient.DefaultGasPrice,
		iotaclient.DefaultGasBudget,
	)
	require.NoError(t, err)
	reqRef, err := createAndSendRequestRes.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
	require.NoError(t, err)
	reqWithObj, err := client.L2().GetRequestFromObjectID(context.Background(), reqRef.ObjectID)
	require.NoError(t, err)
	req, err := isc.OnLedgerFromRequest(reqWithObj, cryptolib.NewAddressFromIota(anchor.ObjectID))
	require.NoError(t, err)

	return req
}
