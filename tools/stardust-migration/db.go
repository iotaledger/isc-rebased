package main

import (
	"log"
	"runtime"

	"github.com/iotaledger/hive.go/kvstore"
	hivedb "github.com/iotaledger/hive.go/kvstore/database"
	"github.com/iotaledger/hive.go/kvstore/rocksdb"
	"github.com/nnikolash/wasp-types-exported/packages/database"
	"github.com/samber/lo"
)

func createDB(dbDir string) kvstore.KVStore {
	log.Printf("Creating DB in %v\n", dbDir)

	// TODO: BlockCacheSize - what value should be there?
	rocksDB := lo.Must(database.NewRocksDB(dbDir, database.CacheSizeDefault))

	db := database.New(
		dbDir,
		rocksdb.New(rocksDB),
		hivedb.EngineRocksDB,
		true,
		func() bool { panic("should not be called") },
	)

	kvs := db.KVStore()

	return kvs
}

func connectDB(dbDir string) kvstore.KVStore {
	log.Printf("Connecting to DB in %v\n", dbDir)

	rocksDatabase := lo.Must(rocksdb.OpenDBReadOnly(dbDir,
		rocksdb.IncreaseParallelism(runtime.NumCPU()-1),
		rocksdb.Custom([]string{
			"periodic_compaction_seconds=43200",
			"level_compaction_dynamic_level_bytes=true",
			"keep_log_file_num=2",
			"max_log_file_size=50000000", // 50MB per log file
		}),
	))

	db := database.New(
		dbDir,
		rocksdb.New(rocksDatabase),
		hivedb.EngineRocksDB,
		true,
		func() bool { panic("should not be called") },
	)

	kvs := db.KVStore()

	return kvs
}
