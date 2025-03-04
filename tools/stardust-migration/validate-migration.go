package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/samber/lo"

	"github.com/pmezard/go-difflib/difflib"
	cmd "github.com/urfave/cli/v2"

	old_iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/state/indexedstore"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/tools/stardust-migration/cli"
	"github.com/iotaledger/wasp/tools/stardust-migration/db"
	"github.com/iotaledger/wasp/tools/stardust-migration/stateaccess/newstate"
	"github.com/iotaledger/wasp/tools/stardust-migration/stateaccess/oldstate"
	old_isc "github.com/nnikolash/wasp-types-exported/packages/isc"
	old_kv "github.com/nnikolash/wasp-types-exported/packages/kv"
	old_parameters "github.com/nnikolash/wasp-types-exported/packages/parameters"
	old_state "github.com/nnikolash/wasp-types-exported/packages/state"
	old_indexedstore "github.com/nnikolash/wasp-types-exported/packages/state/indexedstore"
	old_accounts "github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
)

const (
	newSchema = allmigrations.SchemaVersionIotaRebased
)

func validateMigration(c *cmd.Context) error {
	srcChainDBDir := c.Args().Get(0)
	destChainDBDir := c.Args().Get(1)

	srcKVS := db.ConnectOld(srcChainDBDir)
	srcStore := old_indexedstore.New(old_state.NewStoreWithUniqueWriteMutex(srcKVS))
	oldChainID := old_isc.ChainID(GetAnchorOutput(lo.Must(srcStore.LatestState())).AliasID)

	destKVS := db.ConnectNew(destChainDBDir)
	destStore := indexedstore.New(state.NewStoreWithUniqueWriteMutex(destKVS))
	newChainID := isc.ChainID(GetAnchorObject(lo.Must(destStore.LatestState())).ID)

	// TODO:
	// 1. Check equality of states:
	// 	For each of (latest state and some of intermediate states):
	// 		For each entity (balances, blocklog etc):
	// 			For each entry of that entity:
	// 				Get it using contract call (either API or directly) from old and new node
	// 				What cannot be retrieved using contracts - retrieve directly from db or write additional getters
	// 				Print in a text format
	// 				Compare texts
	// 				Also maybe encode old data using BCS and compare bytes
	//
	// 2. Check results of tracing
	// 	For each of (latest block, some of intermediate blocks, and some manually chosen blocks):
	// 		For each of requests in that block:
	// 			Run tracing of that block on old and new node, and compare results
	// 	As for manually chosen blocks I think about those where GasFeePolicy was changed or where requests failed because of gas.
	//
	// 3. (not sure) Check for unexpected DB keys difference
	// 	For each of (latest state and some of intermediate states):
	// 		Find difference in presence of keys between old and new state
	// 		Filter out keys by prefix, which are expected to not be present
	// 		Analyze the result. Ideally we should eventually have no unexpected difference there
	// 	We could also do that not over states, but over mutations in blocks.

	newLatestState := lo.Must(destStore.LatestState())
	cli.DebugLogf("Latest new state index: %v", newLatestState.BlockIndex())

	cli.DebugLogf("Reading old latest state for index #%v...", newLatestState.BlockIndex())
	oldLatestState := lo.Must(srcStore.StateByIndex(newLatestState.BlockIndex()))

	old_parameters.InitL1(&old_parameters.L1Params{
		Protocol: &old_iotago.ProtocolParameters{
			Bech32HRP: old_iotago.PrefixMainnet,
		},
	})

	cli.DebugLoggingEnabled = true
	validateStatesEqual(oldLatestState, newLatestState, oldChainID, newChainID)

	return nil
}

func validateStatesEqual(oldState old_state.State, newState state.State, oldChainID old_isc.ChainID, newChainID isc.ChainID) {
	cli.DebugLogf("Validating states equality...\n")
	oldStateContentStr := oldStateContentToStr(oldState, oldChainID)
	newStateContentStr := newStateContentToStr(newState, newChainID)

	oldStateFilePath := os.TempDir() + "/stardust-migration-old-state.txt"
	newStateFilePath := os.TempDir() + "/stardust-migration-new-state.txt"
	cli.DebugLogf("Writing old and new states to files %v and %v\n", oldStateFilePath, newStateFilePath)

	os.WriteFile(oldStateFilePath, []byte(oldStateContentStr), 0644)
	os.WriteFile(newStateFilePath, []byte(newStateContentStr), 0644)

	if oldStateContentStr != newStateContentStr {
		diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:       difflib.SplitLines(oldStateContentStr),
			B:       difflib.SplitLines(newStateContentStr),
			Context: 2,
		})

		cli.DebugLogf("States are not equal:\n%v\n", diff)
		os.Exit(1)
	}

	cli.DebugLogf("States are equal\n")
}

func oldStateContentToStr(chainState old_state.State, chainID old_isc.ChainID) string {
	accountsContractStr := oldAccountsContractContentToStr(oldstate.GetContactStateReader(chainState, old_accounts.Contract.Hname()), chainID)

	return accountsContractStr
}

func newStateContentToStr(chainState state.State, chainID isc.ChainID) string {
	accountsContractStr := newAccountsContractContentToStr(newstate.GetContactStateReader(chainState, accounts.Contract.Hname()), chainID)

	return accountsContractStr
}

func oldAccountsContractContentToStr(contractState old_kv.KVStoreReader, chainID old_isc.ChainID) string {
	accsStr := oldAccountsListToStr(contractState, chainID)
	accsLines := strings.Split(accsStr, "\n")
	cli.DebugLogf("Old accounts preview:\n%v\n...", strings.Join(lo.Slice(accsLines, 0, 10), "\n"))

	return accsStr
}

func newAccountsContractContentToStr(contractState kv.KVStoreReader, chainID isc.ChainID) string {
	accsStr := newAccountsListToStr(contractState, chainID)
	accsLines := strings.Split(accsStr, "\n")
	cli.DebugLogf("New accounts preview:\n%v\n...", strings.Join(lo.Slice(accsLines, 0, 10), "\n"))

	return accsStr
}

func oldAccountsListToStr(contractState old_kv.KVStoreReader, chainID old_isc.ChainID) string {
	cli.DebugLogf("Reading old accounts list...\n")
	accs := old_accounts.AllAccountsMapR(contractState)

	cli.DebugLogf("Found %v accounts\n", accs.Len())
	cli.DebugLogf("Reading accounts...\n")
	printProgress, clearProgress := cli.NewProgressPrinter("accounts", accs.Len())
	defer clearProgress()

	var accsStr strings.Builder
	accs.Iterate(func(accKey []byte, accValue []byte) bool {
		accID := lo.Must(old_accounts.AgentIDFromKey(old_kv.Key(accKey), chainID))
		accsStr.WriteString("\t")
		accsStr.WriteString(oldAgentIDToStr(accID))
		accsStr.WriteString("\n")
		printProgress()
		return true
	})

	cli.DebugLogf("Replacing chain ID with constant placeholder...\n")
	res := fmt.Sprintf("Found %v accounts:\n%v", accs.Len(), sortLines(accsStr.String()))
	res = strings.Replace(res, chainID.String(), "<chain-id>", -1)

	return res
}

func oldAgentIDToStr(agentID old_isc.AgentID) string {
	switch agentID := agentID.(type) {
	case *old_isc.AddressAgentID:
		return fmt.Sprintf("AddressAgentID(%v)", agentID.Address().String())
	case *old_isc.ContractAgentID:
		return fmt.Sprintf("ContractAgentID(%v, %v)", agentID.ChainID().String(), agentID.Hname())
	case *old_isc.EthereumAddressAgentID:
		return fmt.Sprintf("EthereumAddressAgentID(%v, %v)", agentID.ChainID().String(), agentID.EthAddress().String())
	case *old_isc.NilAgentID:
		panic(fmt.Sprintf("Found agent ID with kind = AgentIDIsNil: %v", agentID))
	default:
		panic(fmt.Sprintf("Unknown agent ID kind: %v/%T = %v", agentID.Kind(), agentID, agentID))
	}
}

func newAccountsListToStr(contractState kv.KVStoreReader, chainID isc.ChainID) string {
	cli.DebugLogf("Reading new accounts list...\n")
	accs := accounts.NewStateReader(newSchema, contractState).AllAccountsAsDict()

	cli.DebugLogf("Found %v accounts\n", len(accs))
	cli.DebugLogf("Reading accounts...\n")
	printProgress, clearProgress := cli.NewProgressPrinter("accounts", uint32(len(accs)))
	defer clearProgress()

	var accsStr strings.Builder
	accs.IterateSorted("", func(accKey kv.Key, accValue []byte) bool {
		accID := lo.Must(accounts.AgentIDFromKey(kv.Key(accKey), chainID))
		accsStr.WriteString("\t")
		accsStr.WriteString(newAgentIDToStr(accID))
		accsStr.WriteString("\n")
		printProgress()
		return true
	})

	cli.DebugLogf("Replacing chain ID with constant placeholder...\n")
	res := fmt.Sprintf("Found %v accounts:\n%v", len(accs), sortLines(accsStr.String()))
	res = strings.Replace(res, chainID.String(), "<chain-id>", -1)

	return res
}

func newAgentIDToStr(agentID isc.AgentID) string {
	switch agentID := agentID.(type) {
	case *isc.AddressAgentID:
		return fmt.Sprintf("AddressAgentID(%v)", agentID.Address().String())
	case *isc.ContractAgentID:
		return fmt.Sprintf("ContractAgentID(%v, %v)", agentID.ChainID().String(), agentID.Hname())
	case *isc.EthereumAddressAgentID:
		return fmt.Sprintf("EthereumAddressAgentID(%v, %v)", agentID.ChainID().String(), agentID.EthAddress().String())
	case *isc.NilAgentID:
		panic(fmt.Sprintf("Found agent ID with kind = AgentIDIsNil: %v", agentID))
	default:
		panic(fmt.Sprintf("Unknown agent ID kind: %v/%T = %v", agentID.Kind(), agentID, agentID))
	}
}

func sortLines(s string) string {
	lines := strings.Split(s, "\n")
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}
