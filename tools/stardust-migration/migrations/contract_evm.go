package migrations

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmimpl"
	"github.com/iotaledger/wasp/tools/stardust-migration/utils"
	"github.com/iotaledger/wasp/tools/stardust-migration/utils/cli"
	old_isc "github.com/nnikolash/wasp-types-exported/packages/isc"
	old_kv "github.com/nnikolash/wasp-types-exported/packages/kv"
	old_evm "github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	old_evmimpl "github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/evmimpl"
)

func MigrateEVMContract(oldChainState old_kv.KVStoreReader, newChainState kv.KVStore) {
	cli.DebugLog("Migrating evm contract...")

	oldContractState := old_evm.ContractPartitionR(oldChainState)
	newContractState := evm.ContractPartition(newChainState)

	migrateEVMEmulator(oldContractState, newContractState)

	oldMagicState := old_evm.ISCMagicSubrealmR(oldContractState)
	newMagicState := evm.ISCMagicSubrealm(newContractState)

	migrateISCMagicPrivileged(oldMagicState, newMagicState)
	migrateISCMagicAllowance(oldMagicState, newMagicState)
	migrateISCMagicERC20ExternalNativeTokens(oldMagicState, newMagicState)

	cli.DebugLog("Migrated evm contract")
}

func migrateEVMEmulator(oldContractState old_kv.KVStoreReader, newContractState kv.KVStore) {
	cli.DebugLog("Migrating evm/emulator...")

	oldEmulatorState := old_evm.EmulatorStateSubrealmR(oldContractState)
	newEmulatorState := evm.EmulatorStateSubrealm(newContractState)

	progress := NewProgressPrinter(500)

	// We simply copy all bytes
	oldEmulatorState.Iterate("", func(key old_kv.Key, value []byte) bool {
		newEmulatorState.Set(kv.Key(key), value)
		progress.Print()
		return true
	})

	cli.DebugLogf("Migrated %v keys for evm/emulator", progress.Count)
}

func migrateISCMagicPrivileged(oldMagicState old_kv.KVStoreReader, newMagicState kv.KVStore) {
	cli.DebugLog("Migrating iscmagic/privileged...")

	count := 0

	// Simply copy all bytes
	oldMagicState.IterateSorted(old_evmimpl.PrefixPrivileged, func(k old_kv.Key, v []byte) bool {
		newMagicState.Set(kv.Key(k), v)
		count++
		return true
	})

	cli.DebugLogf("Migrated %v keys for iscmagic/privileged", count)
}

func migrateISCMagicAllowance(oldMagicState old_kv.KVStoreReader, newMagicState kv.KVStore) {
	cli.DebugLog("Migrating iscmagic/allowance...")

	progress := NewProgressPrinter()

	oldMagicState.IterateSorted(old_evmimpl.PrefixAllowance, func(k old_kv.Key, v []byte) bool {
		k = utils.MustRemovePrefix(k, old_evmimpl.PrefixAllowance)
		if len(k) != 2*common.AddressLength {
			log.Panicf("unexpected key length: %v", len(k))
		}

		oldFromBytes := []byte(k[:common.AddressLength])
		oldToBytes := []byte(k[common.AddressLength:])

		from := common.BytesToAddress(oldFromBytes)
		to := common.BytesToAddress(oldToBytes)

		if v == nil {
			newMagicState.Del(kv.Key(k))
		} else {
			oldAllowance := old_isc.MustAssetsFromBytes(v)
			newAllowance := OldAssetsToNewAssets(oldAllowance)

			newMagicState.Set(evmimpl.KeyAllowance(from, to), newAllowance.Bytes())
		}

		progress.Print()
		return true
	})

	cli.DebugLogf("Migrated %v keys for iscmagic/allowance", progress.Count)
}

func migrateISCMagicERC20ExternalNativeTokens(oldMagicState old_kv.KVStoreReader, newMagicState kv.KVStore) {
	cli.DebugLog("Migrating iscmagic/erc20_external_native_tokens...")

	count := 0

	// Simply copying all bytes, because for now not sure what to do with it, plus according to the information about keys usage
	// this feature seems not even being used.
	// TODO: revisit this before doing actual migration.
	oldMagicState.IterateSorted(old_evmimpl.PrefixERC20ExternalNativeTokens, func(k old_kv.Key, v []byte) bool {
		newMagicState.Set(kv.Key(k), v)
		count++
		return true
	})

	cli.DebugLogf("Migrated %v keys for iscmagic/erc20_external_native_tokens", count)
}
