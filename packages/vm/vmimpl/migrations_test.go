package vmimpl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/migrations"
	"github.com/iotaledger/wasp/packages/vm/core/root"
)

type migrationsTestEnv struct {
	t     *testing.T
	db    kvstore.KVStore
	cs    state.Store
	vmctx *vmContext

	counter    int
	incCounter migrations.Migration
	panic      migrations.Migration
}

func (e *migrationsTestEnv) getSchemaVersion() (ret isc.SchemaVersion) {
	e.vmctx.withStateUpdate(func(chainState kv.KVStore) {
		ret = root.NewStateReaderFromChainState(chainState).GetSchemaVersion()
	})
	return
}

func (e *migrationsTestEnv) setSchemaVersion(v isc.SchemaVersion) {
	e.vmctx.withStateUpdate(func(chainState kv.KVStore) {
		root.NewStateWriter(root.Contract.StateSubrealm(chainState)).SetSchemaVersion(v)
	})
}

func newMigrationsTest(t *testing.T, stateIndex uint32) *migrationsTestEnv {
	db := mapdb.NewMapDB()
	cs := state.NewStoreWithUniqueWriteMutex(db)
	origin.InitChain(0, cs, nil, 0)
	latest, err := cs.LatestBlock()
	require.NoError(t, err)
	stateDraft, err := cs.NewStateDraft(time.Now(), latest.L1Commitment())
	require.NoError(t, err)
	task := &vm.VMTask{
		Anchor: &isc.StateAnchor{
			Ref: &iscmove.RefWithObject[iscmove.Anchor]{
				Object: &iscmove.Anchor{
					StateIndex: stateIndex,
				},
			},
		},
	}
	vmctx := &vmContext{
		task:       task,
		stateDraft: stateDraft,
	}
	vmctx.loadChainConfig()

	env := &migrationsTestEnv{
		t:     t,
		db:    db,
		cs:    cs,
		vmctx: vmctx,
	}

	env.incCounter = migrations.Migration{
		Contract: governance.Contract,
		Apply: func(state kv.KVStore, log *zap.SugaredLogger) error {
			env.counter++
			return nil
		},
	}

	env.panic = migrations.Migration{
		Contract: governance.Contract,
		Apply: func(state kv.KVStore, log *zap.SugaredLogger) error {
			panic("should not be called")
		},
	}

	return env
}

func TestMigrationsStateIndex1(t *testing.T) {
	env := newMigrationsTest(t, 1)

	require.EqualValues(t, 0, env.getSchemaVersion())

	env.vmctx.withStateUpdate(func(chainState kv.KVStore) {
		env.vmctx.runMigrations(chainState, &migrations.MigrationScheme{
			BaseSchemaVersion: 0,
			Migrations:        []migrations.Migration{env.incCounter, env.incCounter, env.incCounter},
		})
	})

	require.EqualValues(t, 3, env.counter)
	require.EqualValues(t, 3, env.getSchemaVersion())
}

func TestMigrationsStateIndex1Current1(t *testing.T) {
	env := newMigrationsTest(t, 1)

	env.setSchemaVersion(1)

	env.vmctx.withStateUpdate(func(chainState kv.KVStore) {
		env.vmctx.runMigrations(chainState, &migrations.MigrationScheme{
			BaseSchemaVersion: 0,
			Migrations:        []migrations.Migration{env.panic, env.incCounter, env.incCounter},
		})
	})

	require.EqualValues(t, 2, env.counter)
	require.EqualValues(t, 3, env.getSchemaVersion())
}

func TestMigrationsStateIndex1Current2Base1(t *testing.T) {
	env := newMigrationsTest(t, 1)

	env.setSchemaVersion(2)

	env.vmctx.withStateUpdate(func(chainState kv.KVStore) {
		env.vmctx.runMigrations(chainState, &migrations.MigrationScheme{
			BaseSchemaVersion: 1,
			Migrations:        []migrations.Migration{env.panic, env.incCounter, env.incCounter},
		})
	})

	require.EqualValues(t, 2, env.counter)
	require.EqualValues(t, 4, env.getSchemaVersion())
}
