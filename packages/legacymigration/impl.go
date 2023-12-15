// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0
package legacymigration

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/iotaledger/iota.go/encoding/t5b1"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

var Processor = Contract.Processor(nil,
	// view
	ViewMigratableBalance.WithHandler(viewMigratableBalance),
	ViewTotalBalance.WithHandler(viewTotalBalance),

	// funcs
	FuncMigrate.WithHandler(migrate),
	FuncBurn.WithHandler(burn),
	FuncSetNextAdmin.WithHandler(setNextAdmin),
	FuncClaimAdmin.WithHandler(claimAdmin),
)

//go:embed migratable.csv
var migrationData []byte

//go:embed migratable_test.csv
var migrationDataTest []byte

func SetInitialState(state kv.KVStore, admin isc.AgentID) {
	// set the admin agentID
	setAdminAgentID(state, admin)
	// read migration map from the provided file
	migrationMap := accountsMigrationMap(state)
	var csvReader *csv.Reader
	if os.Getenv("GO_TESTING") == "" {
		csvReader = csv.NewReader(bytes.NewBuffer(migrationData))
	} else {
		csvReader = csv.NewReader(bytes.NewBuffer(migrationDataTest))
	}

	csvReader.Comma = ';'
	for {
		record, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(fmt.Errorf("error reading migration csv: %s", err.Error()))
		}
		if len(record) != 2 {
			panic(fmt.Errorf("invalid migration csv, found record with length %d", len(record)))
		}
		address := t5b1.EncodeTrytes(record[0])
		amount, err := strconv.ParseUint(record[1], 10, 64)
		if err != nil {
			panic(fmt.Errorf("error parsing amount from migration CSV: %s", err.Error()))
		}
		migrationMap.SetAt(address, codec.Encode(amount))
	}
	// set the migration totals
	setTotalAmount(state, calcTotalAmount(state))
}

func migrate(ctx isc.Sandbox) dict.Dict {
	// get a valid bundle
	bundleBytes := ctx.Params().Get(ParamBundle)
	ctx.Requiref(bundleBytes != nil, "missing bundle parameter")

	bndl, err := validBundleFromBytes(bundleBytes)
	ctx.RequireNoError(err, "invalid bundle")

	migratedAddresses, targetAddress, err := addressesFromBundle(bndl)
	ctx.RequireNoError(err, "invalid bundle")

	// collect the sum of the funds to send
	tokensToMigrate := uint64(0)
	// check the SC mapping for all these legacy addresses,
	for _, migratedAddr := range migratedAddresses {
		tokensToMigrate += migratableBalance(ctx.State(), migratedAddr)
		accountsMigrationMap(ctx.State()).DelAt(migratedAddr) // delete the **migrated** addresses from the SC state
	}

	// issue event with amount+bundle
	ww := rwutil.NewBytesWriter()
	ww.WriteUint64(tokensToMigrate)                  // tokens migrated
	ww.WriteBytes(isc.AddressToBytes(targetAddress)) // target address
	ww.WriteUint8(uint8(len(migratedAddresses)))     // list of migrated addresses
	for _, a := range migratedAddresses {
		ww.WriteBytes(a)
	}
	ww.WriteBytes(bundleBytes) // bundle bytes
	ctx.Event("migration", ww.Bytes())

	// - send the funds via L1 to the target address
	ctx.Send(isc.RequestParameters{
		TargetAddress: targetAddress,
		Assets:        isc.NewAssetsBaseTokens(tokensToMigrate),
	})

	// in theory the code below is not needed, but let's keep for now as a sanity check
	{
		totalAmount := getTotalAmount(ctx.State())
		totalAmount -= tokensToMigrate
		setTotalAmount(ctx.State(), totalAmount)
		calculatedAmount := calcTotalAmount(ctx.State())
		// assert the total migration funds is still correct
		ctx.Requiref(calculatedAmount == totalAmount, "inconsistency in migrated funds totals")
		ctx.Requiref(ctx.BalanceBaseTokens() >= totalAmount, "inconsistency in migrated funds balance")
	}

	return nil
}

func viewMigratableBalance(ctx isc.SandboxView) dict.Dict {
	legacyAddrTrytes := ctx.Params().MustGetString(ParamAddress)
	legacyAddr := t5b1.EncodeTrytes(legacyAddrTrytes)

	return dict.Dict{
		ParamBalance: codec.Encode(migratableBalance(ctx.StateR(), legacyAddr)),
	}
}

func migratableBalance(state kv.KVStoreReader, legacyAddr []byte) uint64 {
	migrationMap := accountsMigrationMapR(state)
	return codec.MustDecodeUint64(migrationMap.GetAt(legacyAddr), 0)
}

func mustAdmin(ctx isc.Sandbox) {
	ctx.Requiref(ctx.Request().SenderAccount().Equals(adminAgentID(ctx.State())), "only the admin can call the burn function")
}

func setNextAdmin(ctx isc.Sandbox) dict.Dict {
	mustAdmin(ctx)
	setNextAdminAgentID(ctx.State(), ctx.Params().MustGetAgentID(ParamNextAdminAgentID))
	return nil
}

func claimAdmin(ctx isc.Sandbox) dict.Dict {
	nextAdmin := nextAdminAgentID(ctx.State())
	if ctx.Request().SenderAccount().Equals(nextAdmin) {
		setAdminAgentID(ctx.State(), nextAdmin)
	}
	return nil
}

// TODO make this a "withdraw func" instead
func burn(ctx isc.Sandbox) dict.Dict {
	mustAdmin(ctx)
	ctx.Send(isc.RequestParameters{
		TargetAddress: &iotago.Ed25519Address{0}, // send to the Zero address (0x00000...)
		Assets: &isc.Assets{
			BaseTokens: ctx.BalanceBaseTokens(), // burn the entire balance of this contract
		},
	})
	return nil
}

func viewTotalBalance(ctx isc.SandboxView) dict.Dict {
	return dict.Dict{
		ParamTotalMigrationAmount: codec.Encode(calcTotalAmount(ctx.StateR())),
		ParamBalance:              codec.Encode(ctx.BalanceBaseTokens()),
	}
}

// --- contract state

const (
	keyLegacyAccounts   = "a"
	keyTotalAmount      = "t"
	keyAdminAgentID     = "d"
	keyNextAdminAgentID = "x"
)

func adminAgentID(state kv.KVReader) isc.AgentID {
	a, err := isc.AgentIDFromBytes(state.Get(keyAdminAgentID))
	if err != nil {
		panic("invalid admin agentID")
	}
	return a
}

func setAdminAgentID(state kv.KVStore, a isc.AgentID) {
	state.Set(keyAdminAgentID, a.Bytes())
}

func nextAdminAgentID(state kv.KVReader) isc.AgentID {
	nextAdminBytes := state.Get(keyNextAdminAgentID)
	if nextAdminBytes == nil {
		return nil
	}
	a, err := isc.AgentIDFromBytes(nextAdminBytes)
	if err != nil {
		panic("invalid agentID")
	}
	return a
}

func setNextAdminAgentID(state kv.KVStore, a isc.AgentID) {
	state.Set(keyNextAdminAgentID, a.Bytes())
}

func setTotalAmount(state kv.KVStore, amount uint64) {
	state.Set(keyTotalAmount, codec.Encode(amount))
}

func getTotalAmount(state kv.KVStoreReader) uint64 {
	return codec.MustDecodeUint64(state.Get(keyTotalAmount))
}

func calcTotalAmount(state kv.KVStoreReader) uint64 {
	migrationsMap := accountsMigrationMapR(state)
	accTotal := uint64(0)
	migrationsMap.Iterate(func(elemKey []byte, value []byte) bool {
		accTotal += codec.MustDecodeUint64(value)
		return true
	})
	return accTotal
}

func accountsMigrationMap(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, keyLegacyAccounts)
}

func accountsMigrationMapR(state kv.KVStoreReader) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, keyLegacyAccounts)
}
