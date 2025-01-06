// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package testchain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/util/bcs"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

////////////////////////////////////////////////////////////////////////////////
// TestChainLedger

type TestChainLedger struct {
	t           *testing.T
	l1client    clients.L1Client
	iscPackage  *iotago.PackageID
	chainOwner  *cryptolib.KeyPair
	chainID     isc.ChainID
	fetchedReqs map[cryptolib.AddressKey]map[iotago.ObjectID]bool
}

func NewTestChainLedger(
	t *testing.T,
	originator *cryptolib.KeyPair,
	iscPackage *iotago.PackageID,
	l1client clients.L1Client,
) *TestChainLedger {
	return &TestChainLedger{
		t:           t,
		chainOwner:  originator,
		l1client:    l1client,
		iscPackage:  iscPackage,
		fetchedReqs: map[cryptolib.AddressKey]map[iotago.ObjectID]bool{},
	}
}

// Only set after MakeTxChainOrigin.
func (tcl *TestChainLedger) ChainID() isc.ChainID {
	return tcl.chainID
}

func (tcl *TestChainLedger) MakeTxChainOrigin(committeeAddress *cryptolib.Address) (*isc.StateAnchor, coin.Value) {
	coinType := iotajsonrpc.IotaCoinType.String()
	resGetCoins, err := tcl.l1client.GetCoins(context.Background(), iotaclient.GetCoinsRequest{Owner: tcl.chainOwner.Address().AsIotaAddress(), CoinType: &coinType})
	require.NoError(tcl.t, err)

	originDeposit := resGetCoins.Data[1]
	schemaVersion := allmigrations.DefaultScheme.LatestSchemaVersion()
	initParams := origin.DefaultInitParams(isc.NewAddressAgentID(committeeAddress)).Encode()
	originDepositVal := coin.Value(originDeposit.Balance.Uint64())
	l1commitment := origin.L1Commitment(schemaVersion, initParams, iotago.ObjectID{}, originDepositVal, isc.BaseTokenCoinInfo)
	stateMetadata := transaction.NewStateMetadata(
		schemaVersion,
		l1commitment,
		&iotago.ObjectID{},
		&gas.FeePolicy{
			GasPerToken: util.Ratio32{
				A: 1,
				B: 2,
			},
			EVMGasRatio: util.Ratio32{
				A: 3,
				B: 4,
			},
			ValidatorFeeShare: 5,
		},
		initParams,
		"https://iota.org",
	)
	gasCoin := resGetCoins.Data[0].Ref()
	// FIXME this may refer to the ObjectRef with older version, and trigger panic
	anchorRef, err := tcl.l1client.L2().StartNewChain(
		context.Background(),
		&iscmoveclient.StartNewChainRequest{
			Signer:            tcl.chainOwner,
			ChainOwnerAddress: tcl.chainOwner.Address(),
			PackageID:         *tcl.iscPackage,
			StateMetadata:     stateMetadata.Bytes(),
			InitCoinRef:       originDeposit.Ref(),
			GasPayments:       []*iotago.ObjectRef{gasCoin},
			GasPrice:          iotaclient.DefaultGasPrice,
			GasBudget:         iotaclient.DefaultGasBudget,
		},
	)
	require.NoError(tcl.t, err)
	stateAnchor := isc.NewStateAnchor(anchorRef, *tcl.iscPackage)
	require.NotNil(tcl.t, stateAnchor)
	tcl.chainID = stateAnchor.ChainID()

	return &stateAnchor, originDepositVal
}

func (tcl *TestChainLedger) MakeTxAccountsDeposit(account *cryptolib.KeyPair) (isc.Request, error) {
	resp, err := tcl.l1client.L2().CreateAndSendRequestWithAssets(
		context.Background(),
		&iscmoveclient.CreateAndSendRequestWithAssetsRequest{
			Signer:        account,
			PackageID:     *tcl.iscPackage,
			AnchorAddress: tcl.chainID.AsAddress().AsIotaAddress(),
			Assets:        iscmove.NewAssets(100_000_00),
			Message: &iscmove.Message{
				Contract: uint32(isc.Hn("accounts")),
				Function: uint32(isc.Hn("deposit")),
			},
			Allowance:        iscmove.NewAssets(100_000_000),
			OnchainGasBudget: 1000,
			GasPrice:         iotaclient.DefaultGasPrice,
			GasBudget:        iotaclient.DefaultGasBudget,
		},
	)
	if err != nil {
		return nil, err
	}
	reqRef, err := resp.GetCreatedObjectInfo(iscmove.RequestModuleName, iscmove.RequestObjectName)
	if err != nil {
		return nil, err
	}
	req, err := tcl.l1client.L2().GetRequestFromObjectID(context.Background(), reqRef.ObjectID)
	if err != nil {
		return nil, err
	}
	return isc.OnLedgerFromRequest(req, tcl.chainID.AsAddress())
}

func (tcl *TestChainLedger) RunOnChainStateTransition(anchor *isc.StateAnchor, pt iotago.ProgrammableTransaction) (*isc.StateAnchor, error) {
	signer := cryptolib.SignerToIotaSigner(tcl.chainOwner)

	coinPage, err := tcl.l1client.GetCoins(context.Background(), iotaclient.GetCoinsRequest{Owner: signer.Address()})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GasPayment object: %w", err)
	}
	var gasPayments []*iotago.ObjectRef
	for _, coin := range coinPage.Data {
		if !pt.IsInInputObjects(coin.CoinObjectID) {
			gasPayments = []*iotago.ObjectRef{coin.Ref()}
			break
		}
	}
	tx := iotago.NewProgrammable(
		signer.Address(),
		pt,
		gasPayments,
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)
	txBytes, err := bcs.Marshal(&tx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TransactionData: %w", err)
	}
	_, err = tcl.l1client.SignAndExecuteTransaction(
		context.Background(),
		&iotaclient.SignAndExecuteTransactionRequest{
			TxDataBytes: txBytes,
			Signer:      signer,
			Options:     &iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to SignAndExecuteTransaction: %w", err)
	}
	return tcl.UpdateAnchor(anchor)
}

func (tcl *TestChainLedger) UpdateAnchor(anchor *isc.StateAnchor) (*isc.StateAnchor, error) {
	anchorRef, err := tcl.l1client.UpdateObjectRef(context.Background(), anchor.GetObjectRef())
	if err != nil {
		return nil, err
	}
	anchorWithRef, err := tcl.l1client.L2().GetAnchorFromObjectID(context.Background(), anchorRef.ObjectID)
	if err != nil {
		return nil, err
	}
	stateAnchor := isc.NewStateAnchor(anchorWithRef, anchor.ISCPackage())
	return &stateAnchor, nil
}

func (tcl *TestChainLedger) FakeRotationTX(anchor *isc.StateAnchor, nextCommitteeAddr *cryptolib.Address) *isc.StateAnchor {
	// FIXME a temp impl before the decision of Rotation
	signer := cryptolib.SignerToIotaSigner(tcl.chainOwner)
	ptb := iotago.NewProgrammableTransactionBuilder()

	ptb.Command(iotago.Command{
		TransferObjects: &iotago.ProgrammableTransferObjects{
			Objects: []iotago.Argument{ptb.MustObj(iotago.ObjectArg{ImmOrOwnedObject: anchor.GetObjectRef()})},
			Address: ptb.MustPure(nextCommitteeAddr),
		},
	},
	)

	pt := ptb.Finish()
	coins, err := tcl.l1client.GetCoinObjsForTargetAmount(context.Background(), signer.Address(), iotaclient.DefaultGasBudget, iotaclient.DefaultGasBudget)
	require.NoError(tcl.t, err)
	gasPayments := coins.CoinRefs()

	tx := iotago.NewProgrammable(
		signer.Address(),
		pt,
		gasPayments,
		iotaclient.DefaultGasBudget,
		iotaclient.DefaultGasPrice,
	)

	txnBytes, err := bcs.Marshal(&tx)
	require.NoError(tcl.t, err)
	txnResponse, err := tcl.l1client.SignAndExecuteTransaction(
		context.Background(),
		&iotaclient.SignAndExecuteTransactionRequest{
			TxDataBytes: txnBytes,
			Signer:      signer,
			Options:     &iotajsonrpc.IotaTransactionBlockResponseOptions{ShowEffects: true, ShowObjectChanges: true},
		},
	)

	require.NoError(tcl.t, err)
	require.True(tcl.t, txnResponse.Effects.Data.IsSuccess())

	anchorRef, err := tcl.l1client.L2().GetAnchorFromObjectID(context.Background(), anchor.GetObjectID())
	require.NoError(tcl.t, err)

	tmp := isc.NewStateAnchor(anchorRef, *tcl.iscPackage)
	return &tmp
}
