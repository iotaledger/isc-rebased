package accounts_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/isctest"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/sui-go/sui"
)

func TestAccounts(t *testing.T) {
	// execute tests on all schema versions
	for v := isc.SchemaVersion(0); v <= allmigrations.DefaultScheme.LatestSchemaVersion(); v++ {
		testCreditDebit1(t, v)
		testCreditDebit2(t, v)
		testCreditDebit3(t, v)
		testCreditDebit4(t, v)
		testCreditDebit5(t, v)
		testCreditDebit6(t, v)
		testCreditDebit7(t, v)
		testMoveAll(t, v)
		testDebitAll(t, v)
		testTransferObjects(t, v)
		testCreditDebitObject1(t, v)
	}
}

func knownAgentID(b byte, h uint32) isc.AgentID {
	var chainID isc.ChainID
	for i := range chainID {
		chainID[i] = b
	}
	return isc.NewContractAgentID(chainID, isc.Hname(h))
}

var dummyAssetID = coin.Type("0x1::foo::bar")

func checkLedgerT(t *testing.T, v isc.SchemaVersion, state dict.Dict) isc.CoinBalances {
	require.NoError(t, accounts.NewStateReader(v, state).CheckLedgerConsistency())
	return accounts.NewStateReader(v, state).GetTotalL2FungibleTokens()
}

func testCreditDebit1(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)

	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := knownAgentID(1, 2)
	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	total = checkLedgerT(t, v, state)

	require.NotNil(t, total)
	require.True(t, total.Equals(transfer))

	transfer[coin.BaseTokenType] = 1
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	total = checkLedgerT(t, v, state)

	expected := isc.NewCoinBalances().AddBaseTokens(43).Add(dummyAssetID, 4)
	require.True(t, expected.Equals(total))

	userAssets := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID1, isc.ChainID{})
	require.EqualValues(t, 43, userAssets.BaseTokens())
	require.EqualValues(t, 4, userAssets[dummyAssetID])
	checkLedgerT(t, v, state)

	accounts.NewStateWriter(v, state).DebitFromAccount(agentID1, expected, isc.ChainID{})
	total = checkLedgerT(t, v, state)
	expected = isc.NewCoinBalances()
	require.True(t, expected.Equals(total))
}

func testCreditDebit2(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)
	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := isctest.NewRandomAgentID()
	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	total = checkLedgerT(t, v, state)

	expected := transfer
	require.True(t, expected.Equals(total))

	transfer = isc.NewCoinBalances().Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).DebitFromAccount(agentID1, transfer, isc.ChainID{})
	total = checkLedgerT(t, v, state)
	expected = isc.NewCoinBalances().AddBaseTokens(42)
	require.True(t, expected.Equals(total))

	require.Zero(t, accounts.NewStateReader(v, state).GetCoinBalance(agentID1, dummyAssetID, isc.ChainID{}))
	bal1 := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID1, isc.ChainID{})
	require.False(t, bal1.IsEmpty())
	require.True(t, total.Equals(bal1))
}

func testCreditDebit3(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)
	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := isctest.NewRandomAgentID()
	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	total = checkLedgerT(t, v, state)

	expected := transfer
	require.True(t, expected.Equals(total))

	transfer = isc.NewCoinBalances().Add(dummyAssetID, 100)
	require.Panics(t,
		func() {
			accounts.NewStateWriter(v, state).DebitFromAccount(agentID1, transfer, isc.ChainID{})
		},
	)
	total = checkLedgerT(t, v, state)

	expected = isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	require.True(t, expected.Equals(total))
}

func testCreditDebit4(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)
	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := isctest.NewRandomAgentID()
	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	total = checkLedgerT(t, v, state)

	expected := transfer
	require.True(t, expected.Equals(total))

	keys := accounts.NewStateReader(v, state).AllAccountsAsDict().Keys()
	require.EqualValues(t, 1, len(keys))

	agentID2 := isctest.NewRandomAgentID()
	require.NotEqualValues(t, agentID1, agentID2)

	transfer = isc.NewCoinBalances().AddBaseTokens(20)
	err := accounts.NewStateWriter(v, state).MoveBetweenAccounts(agentID1, agentID2, transfer.ToAssets(), isc.ChainID{})
	require.NoError(t, err)
	total = checkLedgerT(t, v, state)

	keys = accounts.NewStateReader(v, state).AllAccountsAsDict().Keys()
	require.EqualValues(t, 2, len(keys))

	expected = isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	require.True(t, expected.Equals(total))

	bm1 := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID1, isc.ChainID{})
	require.False(t, bm1.IsEmpty())
	expected = isc.NewCoinBalances().AddBaseTokens(22).Add(dummyAssetID, 2)
	require.True(t, expected.Equals(bm1))

	bm2 := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID2, isc.ChainID{})
	require.False(t, bm2.IsEmpty())
	expected = isc.NewCoinBalances().AddBaseTokens(20)
	require.True(t, expected.Equals(bm2))
}

func testCreditDebit5(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)
	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := isctest.NewRandomAgentID()
	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	total = checkLedgerT(t, v, state)

	expected := transfer
	require.True(t, expected.Equals(total))

	keys := accounts.NewStateReader(v, state).AllAccountsAsDict().Keys()
	require.EqualValues(t, 1, len(keys))

	agentID2 := isctest.NewRandomAgentID()
	require.NotEqualValues(t, agentID1, agentID2)

	transfer = isc.NewCoinBalances().AddBaseTokens(50)
	require.Error(t, accounts.NewStateWriter(v, state).MoveBetweenAccounts(agentID1, agentID2, transfer.ToAssets(), isc.ChainID{}))
	total = checkLedgerT(t, v, state)

	keys = accounts.NewStateReader(v, state).AllAccountsAsDict().Keys()
	require.EqualValues(t, 1, len(keys))

	expected = isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	require.True(t, expected.Equals(total))

	bm1 := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID1, isc.ChainID{})
	require.False(t, bm1.IsEmpty())
	require.True(t, expected.Equals(bm1))

	bm2 := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID2, isc.ChainID{})
	require.True(t, bm2.IsEmpty())
}

func testCreditDebit6(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)
	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := isctest.NewRandomAgentID()
	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	checkLedgerT(t, v, state)

	agentID2 := isctest.NewRandomAgentID()
	require.NotEqualValues(t, agentID1, agentID2)

	err := accounts.NewStateWriter(v, state).MoveBetweenAccounts(agentID1, agentID2, transfer.ToAssets(), isc.ChainID{})
	require.NoError(t, err)
	total = checkLedgerT(t, v, state)

	keys := accounts.NewStateReader(v, state).AllAccountsAsDict().Keys()
	require.EqualValues(t, 2, len(keys))

	bal := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID1, isc.ChainID{})
	require.True(t, bal.IsEmpty())

	bal2 := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID2, isc.ChainID{})
	require.False(t, bal2.IsEmpty())
	require.True(t, total.Equals(bal2))
}

func testCreditDebit7(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)
	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := isctest.NewRandomAgentID()
	transfer := isc.NewCoinBalances().Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	checkLedgerT(t, v, state)

	debitTransfer := isc.NewCoinBalances().AddBaseTokens(1)
	// debit must fail
	require.Panics(t, func() {
		accounts.NewStateWriter(v, state).DebitFromAccount(agentID1, debitTransfer, isc.ChainID{})
	})

	total = checkLedgerT(t, v, state)
	require.True(t, transfer.Equals(total))
}

func testMoveAll(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	agentID1 := isctest.NewRandomAgentID()
	agentID2 := isctest.NewRandomAgentID()

	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	accs := accounts.NewStateReader(v, state).AllAccountsAsDict()
	require.Len(t, accs, 1)
	require.EqualValues(t, 1, len(accs))
	_, ok := accs[kv.Key(agentID1.Bytes())]
	require.True(t, ok)

	err := accounts.NewStateWriter(v, state).MoveBetweenAccounts(agentID1, agentID2, transfer.ToAssets(), isc.ChainID{})
	require.NoError(t, err)
	accs = accounts.NewStateReader(v, state).AllAccountsAsDict()
	require.Len(t, accs, 2)
	require.EqualValues(t, 2, len(accs))
	_, ok = accs[kv.Key(agentID2.Bytes())]
	require.True(t, ok)
}

func testDebitAll(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	agentID1 := isctest.NewRandomAgentID()

	transfer := isc.NewCoinBalances().AddBaseTokens(42).Add(dummyAssetID, 2)
	accounts.NewStateWriter(v, state).CreditToAccount(agentID1, transfer, isc.ChainID{})
	accs := accounts.NewStateReader(v, state).AllAccountsAsDict()
	require.Len(t, accs, 1)
	require.EqualValues(t, 1, len(accs))
	_, ok := accs[kv.Key(agentID1.Bytes())]
	require.True(t, ok)

	accounts.NewStateWriter(v, state).DebitFromAccount(agentID1, transfer, isc.ChainID{})
	accs = accounts.NewStateReader(v, state).AllAccountsAsDict()
	require.Len(t, accs, 1)
	require.EqualValues(t, 1, len(accs))
	require.True(t, ok)

	assets := accounts.NewStateReader(v, state).GetAccountFungibleTokens(agentID1, isc.ChainID{})
	require.True(t, assets.IsEmpty())

	assets = accounts.NewStateReader(v, state).GetTotalL2FungibleTokens()
	require.True(t, assets.IsEmpty())
}

func testTransferObjects(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()
	total := checkLedgerT(t, v, state)

	require.True(t, total.Equals(isc.NewCoinBalances()))

	agentID1 := isctest.NewRandomAgentID()
	object1 := &accounts.ObjectRecord{
		ID:  sui.ObjectID{123},
		BCS: []byte("foobar"),
	}
	accounts.NewStateWriter(v, state).CreditObjectToAccount(agentID1, object1, isc.ChainID{})
	// object is credited
	user1Objects := accounts.NewStateReader(v, state).GetAccountObjects(agentID1)
	require.Len(t, user1Objects, 1)
	require.Equal(t, user1Objects[0], object1.ID)

	// object data is saved (accounts.SaveObject must be called)
	accounts.NewStateWriter(v, state).SaveObject(object1)

	objectData := accounts.NewStateReader(v, state).GetObjectBCS(object1.ID)
	require.Equal(t, object1.BCS, objectData)

	agentID2 := isctest.NewRandomAgentID()

	// cannot move an Object that is not owned
	require.Error(t, accounts.NewStateWriter(v, state).MoveBetweenAccounts(agentID1, agentID2, isc.NewAssets(0).AddObject(sui.ObjectID{111}), isc.ChainID{}))

	// moves successfully when the Object is owned
	err := accounts.NewStateWriter(v, state).MoveBetweenAccounts(agentID1, agentID2, isc.NewAssets(0).AddObject(object1.ID), isc.ChainID{})
	require.NoError(t, err)

	user1Objects = accounts.NewStateReader(v, state).GetAccountObjects(agentID1)
	require.Len(t, user1Objects, 0)
	user2Objects := accounts.NewStateReader(v, state).GetAccountObjects(agentID2)
	require.Len(t, user2Objects, 1)
	require.Equal(t, user2Objects[0], object1.ID)

	// remove the Object from the chain
	accounts.NewStateWriter(v, state).DebitObjectFromAccount(agentID2, object1.ID, isc.ChainID{})
	accounts.NewStateWriter(v, state).DeleteObject(object1.ID)
	require.Nil(t, accounts.NewStateReader(v, state).GetObjectBCS(object1.ID))
}

func testCreditDebitObject1(t *testing.T, v isc.SchemaVersion) {
	state := dict.New()

	agentID1 := knownAgentID(1, 2)
	object := &accounts.ObjectRecord{
		ID:  sui.ObjectID{123},
		BCS: []byte("foobar"),
	}
	accounts.NewStateWriter(v, state).CreditObjectToAccount(agentID1, object, isc.ChainID{})

	accObjects := accounts.NewStateReader(v, state).GetAccountObjects(agentID1)
	require.Len(t, accObjects, 1)
	require.Equal(t, accObjects[0], object.ID)

	accounts.NewStateWriter(v, state).DebitObjectFromAccount(agentID1, object.ID, isc.ChainID{})

	accObjects = accounts.NewStateReader(v, state).GetAccountObjects(agentID1)
	require.Len(t, accObjects, 0)
}
