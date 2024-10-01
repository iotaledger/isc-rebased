// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package solo

import (
	"bytes"
	"context"
	"math"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/iotaledger/wasp/packages/util/bcs"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/evm/evmlogger"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/metrics"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/publisher"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/state/indexedstore"
	"github.com/iotaledger/wasp/packages/testutil/l1starter"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/vm/core/coreprocessors"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/migrations"
	"github.com/iotaledger/wasp/packages/vm/core/migrations/allmigrations"
	"github.com/iotaledger/wasp/packages/vm/processors"
	_ "github.com/iotaledger/wasp/packages/vm/sandbox"
	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/suiclient"
	"github.com/iotaledger/wasp/sui-go/suiconn"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
)

const (
	MaxRequestsInBlock = 100
	timeLayout         = "04:05.000000000"
)

// Solo is a structure which contains global parameters of the test: one per test instance
type Solo struct {
	// instance of the test
	T                    Context
	logger               *logger.Logger
	chainsMutex          sync.RWMutex
	chains               map[isc.ChainID]*Chain
	processorConfig      *processors.Config
	enableGasBurnLogging bool
	seed                 cryptolib.Seed
	publisher            *publisher.Publisher
	ctx                  context.Context

	l1Config L1Config
}

// data to be persisted in the snapshot
type chainData struct {
	// Name is the name of the chain
	Name string

	// StateControllerKeyPair signature scheme of the chain address, the one used to control funds owned by the chain.
	// In Solo it is Ed25519 signature scheme (in full Wasp environment is is a BLS address)
	StateControllerKeyPair *cryptolib.KeyPair

	// ChainID is the ID of the chain (in this version alias of the ChainAddress)
	ChainID isc.ChainID

	// OriginatorPrivateKey the key pair used to create the chain (origin transaction).
	// It is a default key pair in many of Solo calls which require private key.
	OriginatorPrivateKey *cryptolib.KeyPair

	// ValidatorFeeTarget is the agent ID to which all fees are accrued. By default, it is equal to OriginatorAgentID
	ValidatorFeeTarget isc.AgentID

	db kvstore.KVStore

	migrationScheme *migrations.MigrationScheme
}

// Chain represents state of individual chain.
// There may be several parallel instances of the chain in the 'solo' test
type Chain struct {
	chainData

	OriginatorAddress *cryptolib.Address
	OriginatorAgentID isc.AgentID

	// Env is a pointer to the global structure of the 'solo' test
	Env *Solo

	// Store is where the chain data (blocks, state) is stored
	store indexedstore.IndexedStore
	// Log is the named logger of the chain
	log *logger.Logger
	// global processor cache
	proc *processors.Cache
	// related to asynchronous backlog processing
	runVMMutex sync.Mutex
	// mempool of the chain is used in Solo to mimic a real node
	mempool Mempool

	RequestsBlock uint32

	metrics *metrics.ChainMetrics

	migrationScheme *migrations.MigrationScheme
}

var _ chain.ChainCore = &Chain{}

type InitOptions struct {
	L1Config          L1Config
	Debug             bool
	PrintStackTrace   bool
	GasBurnLogEnabled bool
	Seed              cryptolib.Seed
	Log               *logger.Logger
}

type L1Config struct {
	SuiRPCURL    string
	SuiFaucetURL string
	ISCPackageID sui.PackageID
}

func DefaultInitOptions() *InitOptions {
	return &InitOptions{
		L1Config: L1Config{
			SuiRPCURL:    suiconn.LocalnetEndpointURL,
			SuiFaucetURL: suiconn.LocalnetFaucetURL,
			ISCPackageID: l1starter.ISCPackageID(),
		},
		Debug:             false,
		PrintStackTrace:   false,
		Seed:              cryptolib.Seed{},
		GasBurnLogEnabled: true, // is ON by default
	}
}

// New creates an instance of the Solo environment
// If solo is used for unit testing, 't' should be the *testing.T instance;
// otherwise it can be either nil or an instance created with NewTestContext.
func New(t Context, initOptions ...*InitOptions) *Solo {
	opt := DefaultInitOptions()
	if len(initOptions) > 0 {
		opt = initOptions[0]
	}
	if opt.Log == nil {
		opt.Log = testlogger.NewNamedLogger(t.Name(), timeLayout)
		if !opt.Debug {
			opt.Log = testlogger.WithLevel(opt.Log, zapcore.InfoLevel, opt.PrintStackTrace)
		}
	}
	evmlogger.Init(opt.Log)

	ctx, cancelCtx := context.WithCancel(context.Background())
	t.Cleanup(cancelCtx)

	ret := &Solo{
		T:                    t,
		logger:               opt.Log,
		l1Config:             opt.L1Config,
		chains:               make(map[isc.ChainID]*Chain),
		processorConfig:      coreprocessors.NewConfigWithCoreContracts(),
		enableGasBurnLogging: opt.GasBurnLogEnabled,
		seed:                 opt.Seed,
		publisher:            publisher.New(opt.Log.Named("publisher")),
		ctx:                  ctx,
	}
	_ = ret.publisher.Events.Published.Hook(func(ev *publisher.ISCEvent[any]) {
		ret.logger.Infof("solo publisher: %s %s %v", ev.Kind, ev.ChainID, ev.String())
	})

	go ret.publisher.Run(ctx)
	go ret.batchLoop()

	return ret
}

func (env *Solo) batchLoop() {
	for {
		time.Sleep(50 * time.Millisecond)
		chains := func() []*Chain {
			env.chainsMutex.Lock()
			defer env.chainsMutex.Unlock()
			return lo.Values(env.chains)
		}()
		for _, ch := range chains {
			ch.collateAndRunBatch()
		}
	}
}

func (env *Solo) IterateChainTrieDBs(
	f func(chainID *isc.ChainID, k []byte, v []byte),
) {
	env.chainsMutex.Lock()
	defer env.chainsMutex.Unlock()

	chainIDs := lo.Keys(env.chains)
	slices.SortFunc(chainIDs, func(a, b isc.ChainID) int { return bytes.Compare(a.Bytes(), b.Bytes()) })
	for _, chID := range chainIDs {
		chID := chID // prevent loop variable aliasing
		ch := env.chains[chID]
		lo.Must0(ch.db.Iterate(nil, func(k []byte, v []byte) bool {
			f(&chID, k, v)
			return true
		}))
	}
}

func (env *Solo) IterateChainLatestStates(
	prefix kv.Key,
	f func(chainID *isc.ChainID, k []byte, v []byte),
) {
	env.chainsMutex.Lock()
	defer env.chainsMutex.Unlock()

	chainIDs := lo.Keys(env.chains)
	slices.SortFunc(chainIDs, func(a, b isc.ChainID) int { return bytes.Compare(a.Bytes(), b.Bytes()) })
	for _, chID := range chainIDs {
		chID := chID // prevent loop variable aliasing
		ch := env.chains[chID]
		store := indexedstore.New(state.NewStoreWithUniqueWriteMutex(ch.db))
		state, err := store.LatestState()
		require.NoError(env.T, err)
		state.IterateSorted(prefix, func(k kv.Key, v []byte) bool {
			f(&chID, []byte(k), v)
			return true
		})
	}
}

func (env *Solo) SyncLog() {
	_ = env.logger.Sync()
}

func (env *Solo) Publisher() *publisher.Publisher {
	return env.publisher
}

func (env *Solo) GetChainByName(name string) *Chain {
	env.chainsMutex.Lock()
	defer env.chainsMutex.Unlock()
	for _, ch := range env.chains {
		if ch.Name == name {
			return ch
		}
	}
	panic("chain not found")
}

// WithNativeContract registers a native contract so that it may be deployed
func (env *Solo) WithNativeContract(c *coreutil.ContractProcessor) *Solo {
	env.processorConfig.RegisterNativeContract(c)
	return env
}

// WithVMProcessor registers a VM processor for binary contracts
func (env *Solo) WithVMProcessor(vmType string, constructor processors.VMConstructor) *Solo {
	_ = env.processorConfig.RegisterVMType(vmType, constructor)
	return env
}

// NewChain deploys new default chain instance.
func (env *Solo) NewChain(depositFundsForOriginator ...bool) *Chain {
	ret, _ := env.NewChainExt(nil, 0, "chain1", evm.DefaultChainID, governance.DefaultBlockKeepAmount)
	if len(depositFundsForOriginator) == 0 || depositFundsForOriginator[0] {
		// deposit some tokens for the chain originator
		err := ret.DepositAssetsToL2(isc.NewAssets(5*isc.Million), nil)
		require.NoError(env.T, err)
	}
	return ret
}

func (env *Solo) ISCPackageID() sui.PackageID {
	return env.l1Config.ISCPackageID
}

func (env *Solo) deployChain(
	chainOriginator *cryptolib.KeyPair,
	initBaseTokens coin.Value,
	name string,
	evmChainID uint16,
	blockKeepAmount int32,
) (chainData, *iscmove.AnchorWithRef) {
	env.logger.Debugf("deploying new chain '%s'", name)

	if chainOriginator == nil {
		chainOriginator = env.NewKeyPairFromIndex(-1000 + len(env.chains)) // making new originator for each new chain
		originatorAddr := chainOriginator.GetPublicKey().AsAddress()
		env.GetFundsFromFaucet(originatorAddr)
	}

	initParams := origin.EncodeInitParams(
		isc.NewAddressAgentID(chainOriginator.Address()),
		evmChainID,
		blockKeepAmount,
	)

	originatorAddr := chainOriginator.GetPublicKey().AsAddress()
	originatorAgentID := isc.NewAddressAgentID(originatorAddr)

	baseTokenCoinInfo := env.L1CoinInfo(coin.BaseTokenType)

	schemaVersion := allmigrations.DefaultScheme.LatestSchemaVersion()
	db := mapdb.NewMapDB()
	store := indexedstore.New(state.NewStoreWithUniqueWriteMutex(db))
	block, stateMetadata := origin.InitChain(
		schemaVersion,
		store,
		initParams,
		initBaseTokens,
		baseTokenCoinInfo,
	)

	panic("refactor me: validate StartNewChain call (initCoinRef, gasPayments)")
	anchorRef, err := env.ISCMoveClient().StartNewChain(
		env.ctx,
		chainOriginator,
		env.ISCPackageID(),
		stateMetadata.Bytes(),
		nil,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
		false,
	)
	chainID := isc.ChainIDFromObjectID(anchorRef.Object.ID)

	env.logger.Infof(
		"deployed chain '%s' - ID: %s - state controller address: %s - origin trie root: %s",
		name,
		chainID,
		originatorAddr,
		block.TrieRoot(),
	)

	return chainData{
		Name:                 name,
		ChainID:              chainID,
		OriginatorPrivateKey: chainOriginator,
		ValidatorFeeTarget:   originatorAgentID,
		db:                   db,
	}, anchorRef
}

// NewChainExt returns also origin and init transactions. Used for core testing
//
// If 'chainOriginator' is nil, new one is generated and utxodb.FundsFromFaucetAmount (many) base tokens are loaded from the UTXODB faucet.
// ValidatorFeeTarget will be set to OriginatorAgentID, and can be changed after initialization.
// To deploy a chain instance the following steps are performed:
//   - chain signature scheme (private key), chain address and chain ID are created
//   - empty virtual state is initialized
//   - origin transaction is created by the originator and added to the UTXODB
//   - 'init' request transaction to the 'root' contract is created and added to UTXODB
//   - backlog processing threads (goroutines) are started
//   - VM processor cache is initialized
//   - 'init' request is run by the VM. The 'root' contracts deploys the rest of the core contracts:
//
// Upon return, the chain is fully functional to process requests
func (env *Solo) NewChainExt(
	chainOriginator *cryptolib.KeyPair,
	initBaseTokens coin.Value,
	name string,
	evmChainID uint16,
	blockKeepAmount int32,
) (*Chain, *iscmove.AnchorWithRef) {
	chData, anchorRef := env.deployChain(chainOriginator, initBaseTokens, name, evmChainID, blockKeepAmount)

	env.chainsMutex.Lock()
	defer env.chainsMutex.Unlock()
	ch := env.addChain(chData)

	ch.log.Infof("chain '%s' deployed. Chain ID: %s", ch.Name, ch.ChainID.String())
	return ch, anchorRef
}

func (env *Solo) addChain(chData chainData) *Chain {
	ch := &Chain{
		chainData:         chData,
		OriginatorAddress: chData.OriginatorPrivateKey.GetPublicKey().AsAddress(),
		OriginatorAgentID: isc.NewAddressAgentID(chData.OriginatorPrivateKey.GetPublicKey().AsAddress()),
		Env:               env,
		store:             indexedstore.New(state.NewStoreWithUniqueWriteMutex(chData.db)),
		proc:              processors.MustNew(env.processorConfig),
		log:               env.logger.Named(chData.Name),
		metrics:           metrics.NewChainMetricsProvider().GetChainMetrics(chData.ChainID),
		mempool:           newMempool(env.GlobalTime, chData.ChainID),
		migrationScheme:   chData.migrationScheme,
	}
	env.chains[chData.ChainID] = ch
	return ch
}

func (env *Solo) Ctx() context.Context {
	return env.ctx
}

// AddRequestsToMempool adds all the requests to the chain mempool,
func (env *Solo) AddRequestsToMempool(ch *Chain, reqs []isc.Request) {
	ch.mempool.ReceiveRequests(reqs...)
}

// EnqueueRequests adds requests contained in the transaction to mempools of respective target chains
func (env *Solo) EnqueueRequests(requests map[isc.ChainID][]isc.Request) {
	env.chainsMutex.RLock()
	defer env.chainsMutex.RUnlock()

	for chainID, reqs := range requests {
		ch, ok := env.chains[chainID]
		if !ok {
			env.logger.Infof("dispatching requests. Unknown chain: %s", chainID.String())
			continue
		}
		ch.mempool.ReceiveRequests(reqs...)
	}
}

func (ch *Chain) GetLatestAnchor() *iscmove.AnchorWithRef {
	anchor, err := ch.Env.ISCMoveClient().GetAnchorFromObjectID(ch.Env.ctx, ch.ChainID.AsAddress().AsSuiAddress())
	require.NoError(ch.Env.T, err)
	return anchor
}

// collateBatch selects requests which are not time locked
// returns batch and and 'remains unprocessed' flag
func (ch *Chain) collateBatch() []isc.Request {
	// emulating variable sized blocks
	maxBatch := MaxRequestsInBlock - rand.Intn(MaxRequestsInBlock/3)
	requests := ch.mempool.RequestBatchProposal()
	batchSize := len(requests)

	if batchSize > maxBatch {
		batchSize = maxBatch
	}

	return requests[:batchSize]
}

func (ch *Chain) collateAndRunBatch() {
	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()
	if ch.Env.ctx.Err() != nil {
		return
	}
	batch := ch.collateBatch()
	if len(batch) > 0 {
		results := ch.runRequestsNolock(batch, "batchLoop")
		for _, res := range results {
			if res.Receipt.Error != nil {
				ch.log.Errorf("runRequestsSync: %v", res.Receipt.Error)
			}
		}
	}
}

func (ch *Chain) AddMigration(m migrations.Migration) {
	ch.migrationScheme.Migrations = append(ch.migrationScheme.Migrations, m)
}

func (ch *Chain) GetCandidateNodes() []*governance.AccessNodeInfo {
	panic("unimplemented")
}

func (ch *Chain) GetChainNodes() []peering.PeerStatusProvider {
	panic("unimplemented")
}

func (ch *Chain) GetCommitteeInfo() *chain.CommitteeInfo {
	panic("unimplemented")
}

func (ch *Chain) ID() isc.ChainID {
	return ch.ChainID
}

func (ch *Chain) Log() *logger.Logger {
	return ch.log
}

func (ch *Chain) Processors() *processors.Cache {
	return ch.proc
}

// ---------------------------------------------

func (env *Solo) L1CoinInfo(coinType coin.Type) *isc.SuiCoinInfo {
	md, err := env.SuiClient().GetCoinMetadata(env.ctx, string(coinType))
	require.NoError(env.T, err)
	ts, err := env.SuiClient().GetTotalSupply(env.ctx, string(coinType))
	require.NoError(env.T, err)
	return isc.SuiCoinInfoFromL1Metadata(coinType, md, coin.Value(ts.Value.Uint64()))
}

func (env *Solo) L1BaseTokenCoins(addr *cryptolib.Address) []*suijsonrpc.Coin {
	return env.L1Coins(addr, coin.BaseTokenType)
}

func (env *Solo) L1AllCoins(addr *cryptolib.Address) []*suijsonrpc.Coin {
	r, err := env.SuiClient().GetCoins(env.ctx, suiclient.GetCoinsRequest{
		Owner: addr.AsSuiAddress(),
		Limit: math.MaxUint,
	})
	require.NoError(env.T, err)
	return r.Data
}

func (env *Solo) L1Coins(addr *cryptolib.Address, coinType coin.Type) []*suijsonrpc.Coin {
	r, err := env.SuiClient().GetCoins(env.ctx, suiclient.GetCoinsRequest{
		Owner:    addr.AsSuiAddress(),
		CoinType: (*string)(&coinType),
		Limit:    math.MaxUint,
	})
	require.NoError(env.T, err)
	return r.Data
}

func (env *Solo) L1BaseTokens(addr *cryptolib.Address) coin.Value {
	return env.L1CoinBalance(addr, coin.BaseTokenType)
}

func (env *Solo) L1CoinBalance(addr *cryptolib.Address, coinType coin.Type) coin.Value {
	r, err := env.SuiClient().GetBalance(env.ctx, suiclient.GetBalanceRequest{
		Owner:    addr.AsSuiAddress(),
		CoinType: string(coinType),
	})
	require.NoError(env.T, err)
	return coin.Value(r.TotalBalance)
}

func (env *Solo) L1NFTs(addr *cryptolib.Address) []sui.ObjectID {
	panic("TODO")
}

// L1Assets returns all ftokens of the address contained in the UTXODB ledger
func (env *Solo) L1CoinBalances(addr *cryptolib.Address) isc.CoinBalances {
	r, err := env.SuiClient().GetAllBalances(env.ctx, addr.AsSuiAddress())
	require.NoError(env.T, err)
	cb := isc.NewCoinBalances()
	for _, b := range r {
		cb.Add(coin.Type(b.CoinType), coin.Value(b.TotalBalance))
	}
	return cb
}

// MintNFTL1 mints a single NFT with the `issuer` account and sends it to a `target` account.
// Base tokens in the NFT output are sent to the minimum storage deposit and are taken from the issuer account.
func (env *Solo) MintNFTL1(issuer *cryptolib.KeyPair, target *cryptolib.Address, immutableMetadata []byte) (*isc.NFT, error) {
	nfts, err := env.MintNFTsL1(issuer, target, nil, [][]byte{immutableMetadata})
	if err != nil {
		return nil, err
	}
	return nfts[0], nil
}

// MintNFTsL1 mints len(metadata) NFTs with the `issuer` account and sends them
// to a `target` account.
//
// If collectionID is not nil, it must be the ID of an NFT owned by the issuer.
// All minted NFTs will belong to the given collection.
// See: https://github.com/iotaledger/tips/blob/main/tips/TIP-0027/tip-0027.md
//
// Base tokens in the NFT outputs are sent to the minimum storage deposit and are taken from the issuer account.
func (env *Solo) MintNFTsL1(
	issuer *cryptolib.KeyPair,
	target *cryptolib.Address,
	collectionID *sui.ObjectID,
	metadata [][]byte,
) ([]*isc.NFT, error) {
	panic("TODO")
	/*
		err := errors.New("refactor me: MintNFTsL1")
		if err != nil {
			return nil, nil, err
		}
		err = env.AddToLedger(tx)
		if err != nil {
			return nil, nil, err
		}

		outSet, err := tx.OutputsSet()
		if err != nil {
			return nil, nil, err
		}

		var nfts []*isc.NFT
		var infos []*NFTMintedInfo
		for id, out := range outSet {
			if out, ok := out.(*iotago.NFTOutput); ok { //nolint:gocritic // false positive
				nftID := util.NFTIDFromNFTOutput(out, id)
				info := &NFTMintedInfo{
					OutputID: id,
					Output:   out,
					NFTID:    nftID,
				}
				nft := &isc.NFT{
					ID:       info.NFTID,
					Issuer:   cryptolib.NewAddressFromIotago(out.ImmutableFeatureSet().IssuerFeature().Address),
					Metadata: out.ImmutableFeatureSet().MetadataFeature().Data,
				}
				nfts = append(nfts, nft)
				infos = append(infos, info)
			}
		}
		return nfts, infos, nil
	*/
}

func (env *Solo) executePTB(ptb sui.ProgrammableTransaction, wallet *cryptolib.KeyPair) *suijsonrpc.SuiTransactionBlockResponse {
	tx := sui.NewProgrammable(
		wallet.Address().AsSuiAddress(),
		ptb,
		nil,
		suiclient.DefaultGasPrice,
		suiclient.DefaultGasBudget,
	)

	txnBytes, err := bcs.Marshal(&tx)
	require.NoError(env.T, err)

	execRes, err := env.SuiClient().SignAndExecuteTransaction(
		env.ctx,
		cryptolib.SignerToSuiSigner(wallet),
		txnBytes,
		&suijsonrpc.SuiTransactionBlockResponseOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		},
	)
	require.NoError(env.T, err)
	require.True(env.T, execRes.Effects.Data.IsSuccess())
	return execRes
}

// SendL1 sends coins to another L1 address
func (env *Solo) SendL1(targetAddress *cryptolib.Address, coins isc.CoinBalances, wallet *cryptolib.KeyPair) {
	ptb := sui.NewProgrammableTransactionBuilder()
	allCoins := env.L1AllCoins(wallet.Address())
	coins.IterateSorted(func(coinType coin.Type, amount coin.Value) bool {
		ptb.Pay(
			lo.Map(
				lo.Filter(allCoins, func(c *suijsonrpc.Coin, _ int) bool {
					return c.CoinType == string(coinType)
				}),
				func(c *suijsonrpc.Coin, _ int) *sui.ObjectRef {
					return c.Ref()
				},
			),
			[]*sui.Address{targetAddress.AsSuiAddress()},
			[]uint64{uint64(amount)},
		)
		return true
	})

	env.executePTB(ptb.Finish(), wallet)
}
