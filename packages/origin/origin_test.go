package origin_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/testutil/testmisc"
	"github.com/iotaledger/wasp/packages/testutil/utxodb"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/sui-go/suiclient"
)

var baseTokenCoinInfo = &isc.SuiCoinInfo{CoinType: coin.BaseTokenType}

func TestOrigin(t *testing.T) {
	l1commitment := origin.L1Commitment(0, nil, 0, baseTokenCoinInfo)
	store := state.NewStoreWithUniqueWriteMutex(mapdb.NewMapDB())
	initBlock, _ := origin.InitChain(0, store, nil, 0, baseTokenCoinInfo)
	latestBlock, err := store.LatestBlock()
	require.NoError(t, err)
	require.True(t, l1commitment.Equals(initBlock.L1Commitment()))
	require.True(t, l1commitment.Equals(latestBlock.L1Commitment()))
}

func TestCreateOrigin(t *testing.T) {
	var u *UtxoDB
	var originTx *iotago.Transaction
	var userKey *cryptolib.KeyPair
	var userAddr, stateAddr *cryptolib.Address
	var err error
	var chainID isc.ChainID
	var originTxID iotago.TransactionID

	initTest := func() {
		u = utxodb.New()
		userKey = cryptolib.NewKeyPair()
		userAddr = userKey.GetPublicKey().AsAddress()
		_, err2 := u.GetFundsFromFaucet(userAddr)
		require.NoError(t, err2)

		stateKey := cryptolib.NewKeyPair()
		stateAddr = stateKey.GetPublicKey().AsAddress()

		require.EqualValues(t, suiclient.FundsFromFaucetAmount, u.GetAddressBalanceBaseTokens(userAddr))
		require.EqualValues(t, 0, u.GetAddressBalanceBaseTokens(stateAddr))
	}
	createOrigin := func() {
		allOutputs, ids := u.GetUnspentOutputs(userAddr)

		panic("refactor me: origin.NewChainOriginTransaction")
		var originTx *iotago.Transaction
		var chainID isc.ChainID
		err = errors.New("refactor me: testConsBasic")
		_ = allOutputs
		_ = ids

		require.NoError(t, err)

		err = u.AddToLedger(originTx)
		require.NoError(t, err)

		originTxID, err = originTx.ID()
		require.NoError(t, err)

		txBack, ok := u.GetTransaction(originTxID)
		require.True(t, ok)
		txidBack, err2 := txBack.ID()
		require.NoError(t, err2)
		require.EqualValues(t, originTxID, txidBack)

		t.Logf("New chain ID: %s", chainID.String())
	}

	t.Run("create origin", func(t *testing.T) {
		initTest()
		createOrigin()

		anchor, _, err := transaction.GetAnchorFromTransaction(originTx)
		require.NoError(t, err)
		require.True(t, anchor.IsOrigin)
		require.EqualValues(t, chainID, anchor.ChainID)
		require.EqualValues(t, 0, anchor.StateIndex)
		require.True(t, stateAddr.Equals(anchor.StateController))
		require.True(t, stateAddr.Equals(anchor.GovernanceController))

		originStateMetadata := transaction.NewStateMetadata(
			origin.L1Commitment(
				allmigrations.DefaultScheme.LatestSchemaVersion(),
				isc.NewCallArguments(isc.NewAddressAgentID(anchor.GovernanceController).Bytes()),
				governance.DefaultMinBaseTokensOnCommonAccount,
				baseTokenCoinInfo,
			),
			gas.DefaultFeePolicy(),
			allmigrations.DefaultScheme.LatestSchemaVersion(),
			"",
		)

		require.EqualValues(t, anchor.StateData, originStateMetadata.Bytes())

		// only one output is expected in the ledger under the address of chainID
		outs, ids := u.GetUnspentOutputs(chainID.AsAddress())
		require.EqualValues(t, 1, len(outs))
		require.EqualValues(t, 1, len(ids))

		out := u.GetOutput(anchor.OutputID)
		require.NotNil(t, out)
	})
	t.Run("create init chain originTx", func(t *testing.T) {
		initTest()
		createOrigin()

		chainBaseTokens := originTx.Essence.Outputs[0].Deposit()

		t.Logf("chainBaseTokens: %d", chainBaseTokens)

		require.EqualValues(t, suiclient.FundsFromFaucetAmount-chainBaseTokens, int(u.GetAddressBalanceBaseTokens(userAddr)))
		require.EqualValues(t, 0, u.GetAddressBalanceBaseTokens(stateAddr))
		allOutputs, ids := u.GetUnspentOutputs(chainID.AsAddress())
		require.EqualValues(t, 1, len(allOutputs))
		require.EqualValues(t, 1, len(ids))
	})
}

// Was used to find proper deposit values for a specific metadata according to the existing hashes.
func TestMetadataBad(t *testing.T) {
	t.SkipNow()

	// This test was also skipped for wasp.
	// When it is enabled, it fails in both repos, so I have no easy way to decode that string of bytes
	// to be able then to re-implement it for isc.CallArguments instead of dict.Dict.

	// metadataHex := "0300000001006102000000e60701006204000000ffffffff01006322000000010024ed2ed9d3682c9c4b801dd15103f73d1fe877224cb51c8b3def6f91b67f5067"
	// metadataBin, err := hex.DecodeString(metadataHex)
	// require.NoError(t, err)
	// var initParams dict.Dict
	// initParams, err = dict.FromBytes(metadataBin)
	// require.NoError(t, err)
	// require.NotNil(t, initParams)
	// t.Logf("Args=%v", initParams)
	// initParams.Iterate(kv.EmptyPrefix, func(key kv.Key, value []byte) bool {
	// 	t.Logf("Args, %v ===> %v", key, value)
	// 	return true
	// })

	// for deposit := uint64(0); deposit <= 10*isc.Million; deposit++ {
	// 	db := mapdb.NewMapDB()
	// 	st := state.NewStoreWithUniqueWriteMutex(db)
	// 	block1A, _ := origin.InitChain(0, st, initParams, deposit, isc.BaseTokenCoinInfo)
	// 	block1B, _ := origin.InitChain(0, st, initParams, 10*isc.Million-deposit, isc.BaseTokenCoinInfo)
	// 	block1C, _ := origin.InitChain(0, st, initParams, 10*isc.Million+deposit, isc.BaseTokenCoinInfo)
	// 	block2A, _ := origin.InitChain(0, st, nil, deposit, isc.BaseTokenCoinInfo)
	// 	block2B, _ := origin.InitChain(0, st, nil, 10*isc.Million-deposit, isc.BaseTokenCoinInfo)
	// 	block2C, _ := origin.InitChain(0, st, nil, 10*isc.Million+deposit, isc.BaseTokenCoinInfo)
	// 	t.Logf("Block0, deposit=%v => %v %v %v / %v %v %v", deposit,
	// 		block1A.L1Commitment(), block1B.L1Commitment(), block1C.L1Commitment(),
	// 		block2A.L1Commitment(), block2B.L1Commitment(), block2C.L1Commitment(),
	// 	)
	// }
}

func TestDictBytes(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		key := rapid.SliceOfBytesMatching(".*").Draw(t, "key")
		val := rapid.SliceOfBytesMatching(".+").Draw(t, "val")
		d := dict.New()
		d.Set(kv.Key(key), val)
		b := d.Bytes()
		d2, err := dict.FromBytes(b)
		require.NoError(t, err)
		require.Equal(t, d, d2)
	})
}

// example values taken from a test on the testnet
func TestMismatchOriginCommitment(t *testing.T) {
	store := state.NewStoreWithUniqueWriteMutex(mapdb.NewMapDB())
	oid, err := iotago.OutputIDFromHex("0xcf72dd6a8c8cd76eab93c80ae192677a17c554b91334a41bed5079eff37effc40000")
	require.NoError(t, err)
	originMetadata, err := cryptolib.DecodeHex("0x03016102e607016204ffffffff016322010024ed2ed9d3682c9c4b801dd15103f73d1fe877224cb51c8b3def6f91b67f5067")
	require.NoError(t, err)
	aoStateMetadata, err := cryptolib.DecodeHex("0x01000000006e55672af085d73ea0ed646f280a26e0eba053df10f439378fe4e99e0fb8774600761da7c0402da8640000000100000000010000000100000000")
	require.NoError(t, err)
	_, sender, err := iotago.ParseBech32("rms1qqjw6tke6d5ze8ztsqwaz5gr7u73l6rhyfxt28yt8hhklydk0agxwgerk65")
	require.NoError(t, err)
	_, stateController, err := iotago.ParseBech32("rms1qrkrlggl2plwfvxyuuyj55gw48ws0xwtteydez8y8e03elm3xf38gf7eq5r")
	require.NoError(t, err)
	_, govController, err := iotago.ParseBech32("rms1qqjw6tke6d5ze8ztsqwaz5gr7u73l6rhyfxt28yt8hhklydk0agxwgerk65")
	require.NoError(t, err)
	_, chainAliasAddress, err := iotago.ParseBech32("rms1pr27d4mr9wgesv8je5j6zkequhw0ysx55ftxt04z55dm9hc9yxkauqtukfl")
	require.NoError(t, err)

	ao := isc.NewAliasOutputWithID(
		&iotago.AliasOutput{
			Amount:         10000000,
			NativeTokens:   []*iotago.NativeToken{},
			AliasID:        chainAliasAddress.(*iotago.AliasAddress).AliasID(),
			StateIndex:     0,
			StateMetadata:  aoStateMetadata,
			FoundryCounter: 0,
			Conditions: []iotago.UnlockCondition{
				&iotago.StateControllerAddressUnlockCondition{Address: stateController},
				&iotago.GovernorAddressUnlockCondition{Address: govController},
			},
			Features: []iotago.Feature{
				&iotago.SenderFeature{
					Address: sender,
				},
				&iotago.MetadataFeature{Data: originMetadata},
			},
		},
		oid,
	)

	_, err = origin.InitChainByAliasOutput(store, ao)
	testmisc.RequireErrorToBe(t, err, "l1Commitment mismatch between originAO / originBlock")
}
