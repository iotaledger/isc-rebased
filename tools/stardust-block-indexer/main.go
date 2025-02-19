// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	hivedb "github.com/iotaledger/hive.go/kvstore/database"
	"github.com/iotaledger/hive.go/kvstore/rocksdb"
	"github.com/samber/lo"

	old_kvstore "github.com/iotaledger/hive.go/kvstore"
	old_database "github.com/nnikolash/wasp-types-exported/packages/database"
	old_state "github.com/nnikolash/wasp-types-exported/packages/state"
	old_indexedstore "github.com/nnikolash/wasp-types-exported/packages/state/indexedstore"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s <chain-db-dir> <dest-index-file>", os.Args[0])
	}

	targetChainDBDir := os.Args[1]
	destIndexFile := os.Args[2]

	targetChainDBDir = lo.Must(filepath.Abs(targetChainDBDir))
	destIndexFile = lo.Must(filepath.Abs(destIndexFile))

	if strings.HasPrefix(destIndexFile, targetChainDBDir) {
		log.Fatalf("destination file cannot reside inside source database folder")
	}

	if _, err := os.Stat(destIndexFile); !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("destination file already exists: %v", destIndexFile)
	}

	targetKVS := ConnectToDB(targetChainDBDir)
	targetStore := old_indexedstore.New(old_state.NewStoreWithUniqueWriteMutex(targetKVS))
	latestState := lo.Must(targetStore.LatestState())

	fmt.Printf("Last block index: %v\n", latestState.BlockIndex())

	startTime := time.Now()

	totalBlocksCount := latestState.BlockIndex() + 1
	blocksLeft := totalBlocksCount

	var estimateRunTime time.Duration
	var estimateSpeed int
	lastEstimateUpdateTime := time.Now()

	reverseIterateStates(targetStore, func(state old_state.State) bool {
		blocksLeft--

		periodicAction(time.Second, &lastEstimateUpdateTime, func() {
			if state.BlockIndex() != blocksLeft {
				// Just double-checking
				panic(fmt.Errorf("blocks left: state block index %d does not match expected block index %d", state.BlockIndex(), blocksLeft))
			}

			blocksProcessed := totalBlocksCount - blocksLeft
			relProgress := float64(blocksProcessed) / float64(totalBlocksCount)
			estimateRunTime = time.Duration(float64(time.Since(startTime)) / relProgress)
			estimateSpeed = int(float64(blocksProcessed) / time.Since(startTime).Seconds())
		})

		fmt.Printf("\rBlocks left: %v. Avg speed: %v blocks/sec. Estimate time left: %v         ", blocksLeft, estimateSpeed, estimateRunTime)

		return true
	})

	fmt.Println()
	fmt.Printf("Elapsed time: %v\n", time.Since(startTime))
}

func ConnectToDB(dbDir string) old_kvstore.KVStore {
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

	db := old_database.New(
		dbDir,
		rocksdb.New(rocksDatabase),
		hivedb.EngineRocksDB,
		true,
		func() bool { panic("should not be called") },
	)

	kvs := db.KVStore()

	return kvs
}

func reverseIterateStates(s old_indexedstore.IndexedStore, f func(state old_state.State) bool) {
	state := lo.Must(s.LatestState())

	for {
		if !f(state) {
			return
		}

		prevL1Commitment := state.PreviousL1Commitment()
		if prevL1Commitment == nil {
			if state.BlockIndex() != 0 {
				// Just double-checking
				panic(fmt.Errorf("iterating the chain: state block index %d has no previous L1 commitment", state.BlockIndex()))
			}

			// done
			break
		}

		state = lo.Must(s.StateByTrieRoot(prevL1Commitment.TrieRoot()))
	}
}

func periodicAction(period time.Duration, lastActionTime *time.Time, action func()) {
	if lastActionTime == nil || time.Since(*lastActionTime) >= period {
		action()
		*lastActionTime = time.Now()
	}
}
