package migrations

import (
	"fmt"
	"strings"

	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/collections"
	old_kv "github.com/nnikolash/wasp-types-exported/packages/kv"
	old_collections "github.com/nnikolash/wasp-types-exported/packages/kv/collections"
)

// Iterate state by prefix.
// Prefix is REMOVED before calling the callback.
// Key is automatically deserialized if it is not old_kv.Key.
// Value is automatically deserialized if it is not []byte.
// WARNING: It is UNSAFE to use in some cases, because the prefix could be subprefix of another prefix.
func IterateByPrefix[OldK, OldV any](oldContractState old_kv.KVStoreReader, prefix string, f func(k OldK, v OldV)) (count uint32) {
	oldContractState.Iterate(old_kv.Key(prefix), func(kBytes old_kv.Key, vBytes []byte) bool {
		keyWithoutPrefix := old_kv.Key(strings.TrimPrefix(string(kBytes), prefix))
		k := DeserializeKey[OldK](keyWithoutPrefix)
		v := DeserializeValue[OldV](vBytes)
		f(k, v)
		count++
		return true
	})

	return count
}

// Get list of entries by prefix.
// Prefix is REMOVED before calling the callback.
// Key is automatically deserialized if it is not old_kv.Key.
// Value is automatically deserialized if it is not []byte.
// WARNING: It is UNSAFE to use in some cases, because the prefix could be subprefix of another prefix.
func ListByPrefix[K comparable, V any](store old_kv.KVStoreReader, prefix string) map[K]V {
	entries := map[K]V{}

	IterateByPrefix(store, prefix, func(k K, v V) {
		entries[k] = v
	})

	return entries
}

// Iterate state by prefix and migrate each entry from old state to new state.
// Old prefix is REMOVED from the key before calling callback and new prefix is ADDED to the new key after the callback.
// If you do not want new prefix to be added automatically, just use empty string.
// Keys and values are automatically serialized/deserialized, if needed (see RecordMigrationFunc).
// WARNING: It is UNSAFE to use in some cases, because the prefix could be subprefix of another prefix.
func MigrateByPrefix[OldK, OldV, NewK, NewV any](oldContractState old_kv.KVStoreReader, newContractState kv.KVStore, oldPrefix, newPrefix string, migrationFunc KVMigrationFunc[OldK, OldV, NewK, NewV]) uint32 {
	return IterateByPrefix(oldContractState, oldPrefix, func(oldKey OldK, oldVal OldV) {
		newKey, newVal := migrationFunc(oldKey, oldVal)

		newKeyBytes := kv.Key(newPrefix) + SerializeKey(newKey)
		newValBytes := SerializeValue(newVal)

		newContractState.Set(newKeyBytes, newValBytes)
	})
}

// Iterate specified keys in state.
// Prefix is REMOVED from the key before calling callback.
// Key is automatically deserialized if it is not old_kv.Key.
// Value is automatically deserialized if it is not []byte.
// This is safe version of IterateByPrefix, because it does not rely on prefix - it uses exact keys.
func IterateByKeys[OldK, OldV any](oldContractState old_kv.KVStoreReader, oldPrefix string, oldKeys []old_kv.Key, f func(k OldK, v OldV)) {
	for _, oldKey := range oldKeys {
		oldKeyBytesWithPrefix := old_kv.Key(oldPrefix + string(oldKey))
		oldValBytes := oldContractState.Get(oldKeyBytesWithPrefix)

		k := DeserializeKey[OldK](oldKey)
		v := DeserializeValue[OldV](oldValBytes)

		f(k, v)
	}
}

// Iterate specified keys in state and migrate each entry from old state to new state.
// Old prefix is REMOVED from the key before calling callback and new prefix is ADDED to the new key after the callback.
// If you do not want to add new prefix automatically, just use empty string.
// Keys and values are automatically serialized/deserialized, if needed (see RecordMigrationFunc).
// This is safe version of MigrateByPrefix, because it does not rely on prefix - it uses exact keys.
func MigrateByKeys[OldK, OldV, NewK, NewV any](oldContractState old_kv.KVStoreReader, newContractState kv.KVStore, oldPrefix, newPrefix string, oldKeys []old_kv.Key, migrationFunc KVMigrationFunc[OldK, OldV, NewK, NewV]) {
	IterateByKeys(oldContractState, oldPrefix, oldKeys, func(oldKey OldK, oldVal OldV) {
		newKey, newVal := migrationFunc(oldKey, oldVal)

		newKeyBytes := kv.Key(newPrefix) + SerializeKey(newKey)
		newValBytes := SerializeValue(newVal)

		newContractState.Set(newKeyBytes, newValBytes)
	})
}

// Iterate named map.
// Key is automatically deserialized if it is not old_kv.Key.
// Value is automatically deserialized if it is not []byte.
func IterateMapByName[OldK, OldV any](oldContractState old_kv.KVStoreReader, oldMapName string, f func(k OldK, v OldV)) (count uint32) {
	oldMap := old_collections.NewMapReadOnly(oldContractState, oldMapName)
	IterateMap(oldMap, f)
	return oldMap.Len()
}

// Iterate map.
// Key is automatically deserialized if it is not old_kv.Key.
// Value is automatically deserialized if it is not []byte.
func IterateMap[OldK, OldV any](oldRecords *old_collections.ImmutableMap, f func(k OldK, v OldV)) {
	oldRecords.Iterate(func(kBytes []byte, vBytes []byte) bool {
		k := DeserializeKey[OldK](old_kv.Key(kBytes))
		v := DeserializeValue[OldV](vBytes)
		f(k, v)
		return true
	})
}

// Migrate records from named map into another named map.
// Keys and values are automatically serialized/deserialized, if needed (see RecordMigrationFunc).
func MigrateMapByName[OldK, OldV, NewK, NewV any](oldContractState old_kv.KVStoreReader, newContractState kv.KVStore, oldMapName, newMapName string, migrationFunc KVMigrationFunc[OldK, OldV, NewK, NewV]) uint32 {
	if oldMapName == "" {
		panic("oldMapName is empty")
	}
	if newMapName == "" {
		panic("newMapName is empty")
	}

	oldRecords := old_collections.NewMapReadOnly(oldContractState, oldMapName)
	newRecords := collections.NewMap(newContractState, newMapName)

	MigrateMap(oldRecords, newRecords, migrationFunc)

	return oldRecords.Len()
}

// Migrate records from state map into another state map.
// Keys and values are automatically serialized/deserialized, if needed (see RecordMigrationFunc).
func MigrateMap[OldK, OldV, NewK, NewV any](oldRecords *old_collections.ImmutableMap, newRecords *collections.Map, migrationFunc KVMigrationFunc[OldK, OldV, NewK, NewV]) {
	IterateMap(oldRecords, func(k OldK, v OldV) {
		newKey, newVal := migrationFunc(k, v)

		newKeyBytes := SerializeKey(newKey)
		newValBytes := SerializeValue(newVal)

		newRecords.SetAt([]byte(newKeyBytes), newValBytes)
	})
}

// Iterate named array.
// Value is automatically deserialized if it is not []byte.
func IterateArrayByName[OldV any](oldContractState old_kv.KVStoreReader, oldArrName string, f func(k uint32, v OldV)) (count uint32) {
	oldRecords := old_collections.NewArrayReadOnly(oldContractState, oldArrName)
	IterateArray(oldRecords, f)
	return oldRecords.Len()
}

// Iterate map.
// Value is automatically deserialized if it is not []byte.
func IterateArray[OldV any](oldArr *old_collections.ArrayReadOnly, f func(k uint32, v OldV)) {
	for i := uint32(0); i < oldArr.Len(); i++ {
		oldBytes := oldArr.GetAt(i)
		v := DeserializeValue[OldV](oldBytes)
		f(i, v)
	}
}

// Migrate records from the named array into another named array.
// Values are automatically serialized/deserialized, if needed (see RecordMigrationFunc).
func MigrateArrayByName[OldV, NewV any](oldContractState old_kv.KVStoreReader, newContractState kv.KVStore, oldArrName, newArrName string, migrationFunc ArrayMigrationFunc[OldV, NewV]) uint32 {
	if oldArrName == "" {
		panic("oldArrName is empty")
	}
	if newArrName == "" {
		panic("newArrName is empty")
	}

	oldRecords := old_collections.NewArrayReadOnly(oldContractState, oldArrName)
	newRecords := collections.NewArray(newContractState, newArrName)

	MigrateArray(oldRecords, newRecords, migrationFunc)

	return oldRecords.Len()
}

// Migrate records from state array into another state array.
// Values are automatically serialized/deserialized, if needed (see RecordMigrationFunc).
func MigrateArray[OldV, NewV any](oldRecords *old_collections.ArrayReadOnly, newRecords *collections.Array, migrationFunc ArrayMigrationFunc[OldV, NewV]) {
	IterateArray(oldRecords, func(k uint32, v OldV) {
		newVal := migrationFunc(k, v)
		newValBytes := SerializeValue(newVal)
		newRecords.Push(newValBytes)
	})
}

// Migrate simple high-level variable from old state to new state.
func MigrateVariable[OldV, NewV any](oldContractState old_kv.KVStoreReader, newContractState kv.KVStore, oldKey old_kv.Key, newKey kv.Key, migrationFunc VariableMigrationFunc[OldV, NewV]) {
	oldValueBytes := oldContractState.Get(oldKey)
	oldValue := DeserializeValue[OldV](oldValueBytes)
	newVal := migrationFunc(oldValue)
	newValBytes := SerializeValue(newVal)
	newContractState.Set(newKey, newValBytes)
}

// Can be used with MigrateVariable to migrate variable as by re-encoding.
// Also can be used as an argument of ConvertKV.
func AsIs[Value any](value Value) Value {
	return value
}

// Creates KVMigrationFunc from specified separate converters for key and value.
func ConvertKV[OldK, OldV, NewK, NewV any](
	convertKey func(OldK) NewK,
	convertValue func(OldV) NewV,
) KVMigrationFunc[OldK, OldV, NewK, NewV] {
	return func(oldKey OldK, oldVal OldV) (NewK, NewV) {
		return convertKey(oldKey), convertValue(oldVal)
	}
}

// This is a migration function for a single KV pair.
// Used for maps and for KV storages itself.
// Keys and values are automatically deserialized/serialized if they are not of following types:
// * OldKey - old_kv.Key
// * OldValue - []byte
// * NewKey - kv.Key
// * NewValue - []byte
type KVMigrationFunc[OldKey, OldValue, NewKey, NewValue any] func(OldKey, OldValue) (NewKey, NewValue)

// This is a migration function for a single KV pair.
// Used for maps and for KV storages itself.
// Values are automatically deserialized/serialized if they are not of type []byte.
type ArrayMigrationFunc[OldValue, NewValue any] func(index uint32, oldVal OldValue) NewValue

// This is a migration function for a single high-level variable migration.
// Value is automatically deserialized/serialized if it is not of type []byte.
type VariableMigrationFunc[OldValue, NewValue any] func(oldVal OldValue) NewValue

// Split map key into map name and element key
func SplitMapKey(storeKey old_kv.Key, prefixToRemove ...string) (mapName, elemKey old_kv.Key) {
	if len(prefixToRemove) > 0 {
		storeKey = MustRemovePrefix(storeKey, prefixToRemove[0])
	}

	const elemSep = "."
	pos := strings.Index(string(storeKey), elemSep)

	sepFound := pos >= 0
	sepIsNotLastChar := pos < len(storeKey)-1
	isMapElement := sepFound && sepIsNotLastChar

	if isMapElement {
		return storeKey[:pos], storeKey[pos+1:]
	}

	// Not a map element - maybe map itself or just something else
	return storeKey, ""
}

func MustRemovePrefix[T ~string](v T, prefix string) T {
	r, found := strings.CutPrefix(string(v), prefix)
	if !found {
		panic(fmt.Sprintf("Prefix '%v' not found: %v", prefix, v))
	}

	return T(r)
}

func DecodeKey[Value any](k old_kv.Key, prefixToRemove string, decodeFunc func([]byte) (Value, error)) Value {
	valueBytes, hadPrefix := strings.CutPrefix(string(k), prefixToRemove)
	if !hadPrefix {
		panic(fmt.Sprintf("Key does not have prefix '%v': %v", prefixToRemove, k))
	}

	v, err := decodeFunc([]byte(valueBytes))
	if err != nil {
		panic(err)
	}

	return v
}

// Wraps a function and adds to it printing of number of times it was called
func p[OldK, OldV, NewK, NewV any](f KVMigrationFunc[OldK, OldV, NewK, NewV]) KVMigrationFunc[OldK, OldV, NewK, NewV] {
	callCount := 0

	return func(oldKey OldK, oldVal OldV) (NewK, NewV) {
		callCount++
		if callCount%100 == 0 {
			fmt.Printf("\rProcessed: %v         ", callCount)
		}

		return f(oldKey, oldVal)
	}
}
