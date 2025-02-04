package migrations

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmimpl"
	old_isc "github.com/nnikolash/wasp-types-exported/packages/isc"
	old_kv "github.com/nnikolash/wasp-types-exported/packages/kv"
	old_evm "github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	old_evmimpl "github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/evmimpl"
)

func MigrateEVMContract(oldChainState old_kv.KVStoreReader, newChainState state.StateDraft) {
	log.Print("Migrating evm contract...\n")

	oldContractState := old_evm.ContractPartitionR(oldChainState)
	newContractState := evm.ContractPartition(newChainState)

	migrateEVMEmulator(oldContractState, newContractState)

	oldMagicState := old_evm.ISCMagicSubrealmR(oldContractState)
	newMagicState := evm.ISCMagicSubrealm(newContractState)

	migrateISCMagicPrivileged(oldMagicState, newMagicState)
	migrateISCMagicAllowance(oldMagicState, newMagicState)
	//migrateISCMagicERC20ExternalNativeTokens(oldMagicState, newMagicState)

	log.Print("Migrated evm contract\n")
}

func migrateEVMEmulator(oldContractState old_kv.KVStoreReader, newContractState kv.KVStore) {
	log.Print("Migrating evm/emulator...\n")

	oldEmulatorState := old_evm.EmulatorStateSubrealmR(oldContractState)
	newEmulatorState := evm.EmulatorStateSubrealm(newContractState)

	progress := NewProgressPrinter(500)

	// We simply copy all bytes
	oldEmulatorState.Iterate("", func(key old_kv.Key, value []byte) bool {
		newEmulatorState.Set(kv.Key(key), value)
		progress.Print()
		return true
	})

	log.Print("Migrated %v keys for evm/emulator\n", progress.Count)
}

func migrateISCMagicPrivileged(oldMagicState old_kv.KVStoreReader, newMagicState kv.KVStore) {
	log.Print("Migrating iscmagic/privileged...\n")

	count := 0

	// Simply copy all bytes
	oldMagicState.Iterate(old_evmimpl.PrefixPrivileged, func(k old_kv.Key, v []byte) bool {
		newMagicState.Set(kv.Key(k), v)
		count++
		return true
	})

	log.Print("Migrated %v keys for iscmagic/privileged\n", count)
}

func migrateISCMagicAllowance(oldMagicState old_kv.KVStoreReader, newMagicState kv.KVStore) {
	log.Print("Migrating iscmagic/allowance...\n")

	progress := NewProgressPrinter()

	oldMagicState.Iterate(old_evmimpl.PrefixAllowance, func(k old_kv.Key, v []byte) bool {
		k = MustRemovePrefix(k, old_evmimpl.PrefixAllowance)
		if len(k) != 2*common.AddressLength {
			log.Panicf("unexpected key length: %v", len(k))
		}

		oldFromBytes := []byte(k[:common.AddressLength])
		oldToBytes := []byte(k[common.AddressLength:])
		oldAllowance := old_isc.MustAssetsFromBytes(v)

		from := common.BytesToAddress(oldFromBytes)
		to := common.BytesToAddress(oldToBytes)
		newAllowance := OldAssetsToNewAssets(oldAllowance)

		newMagicState.Set(evmimpl.KeyAllowance(from, to), newAllowance.Bytes())

		progress.Print()
		return true
	})

	log.Print("Migrated %v keys for iscmagic/allowance\n", progress.Count)
}

func migrateISCMagicERC20ExternalNativeTokens(oldMagicState old_kv.KVStoreReader, newMagicState kv.KVStore) {
	log.Print("Migrating iscmagic/erc20_external_native_tokens...\n")

	count := 0

	panic("TODO: implement")

	log.Print("Migrated %v keys for iscmagic/erc20_external_native_tokens\n", count)
}
