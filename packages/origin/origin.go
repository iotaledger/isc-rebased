package origin

import (
	"fmt"
	"time"

	"github.com/iotaledger/wasp/clients/iota-go/iotago"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/util/bcs"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/errors"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmimpl"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/contracts/inccounter"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/contracts/manyevents"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/contracts/testerrors"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

type InitParams struct {
	ChainOwner          isc.AgentID
	EVMChainID          uint16
	BlockKeepAmount     int32
	DeployTestContracts bool
}

func NewInitParams(
	chainOwner isc.AgentID,
	evmChainID uint16,
	blockKeepAmount int32,
	deployTestContracts bool,
) *InitParams {
	return &InitParams{
		ChainOwner:          chainOwner,
		EVMChainID:          evmChainID,
		BlockKeepAmount:     blockKeepAmount,
		DeployTestContracts: deployTestContracts,
	}
}

func DefaultInitParams(chainOwner isc.AgentID) *InitParams {
	return &InitParams{
		ChainOwner:          chainOwner,
		EVMChainID:          evm.DefaultChainID,
		BlockKeepAmount:     governance.DefaultBlockKeepAmount,
		DeployTestContracts: false,
	}
}

func (p *InitParams) Encode() isc.CallArguments {
	return isc.CallArguments{bcs.MustMarshal(p)}
}

func DecodeInitParams(args isc.CallArguments) (*InitParams, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid init params")
	}
	return bcs.Unmarshal[*InitParams](args[0])
}

// L1Commitment calculates the L1 commitment for the origin state
// originDeposit must exclude the minSD for the AliasOutput
func L1Commitment(
	v isc.SchemaVersion,
	args isc.CallArguments,
	gasCoinObjectID iotago.ObjectID,
	originDeposit coin.Value,
	baseTokenCoinInfo *isc.IotaCoinInfo,
) *state.L1Commitment {
	block, _ := InitChain(
		v,
		state.NewStoreWithUniqueWriteMutex(mapdb.NewMapDB()),
		args,
		gasCoinObjectID,
		originDeposit,
		baseTokenCoinInfo,
	)
	return block.L1Commitment()
}

func L1CommitmentFromAnchorStateMetadata(
	stateMetadataBytes []byte,
	originDeposit coin.Value,
	baseTokenCoinInfo *isc.IotaCoinInfo,
) (*state.L1Commitment, error) {
	stateMetadata, err := transaction.StateMetadataFromBytes(stateMetadataBytes)
	if err != nil {
		return nil, err
	}
	l1c := L1Commitment(
		stateMetadata.SchemaVersion,
		stateMetadata.InitParams,
		*stateMetadata.GasCoinObjectID,
		originDeposit,
		baseTokenCoinInfo,
	)
	return l1c, nil
}

func InitChain(
	v isc.SchemaVersion,
	store state.Store,
	args isc.CallArguments,
	gasCoinObjectID iotago.ObjectID,
	originDeposit coin.Value,
	baseTokenCoinInfo *isc.IotaCoinInfo,
) (state.Block, *transaction.StateMetadata) {
	initParams, err := DecodeInitParams(args)
	if err != nil {
		panic(err)
	}

	d := store.NewOriginStateDraft()
	d.Set(kv.Key(coreutil.StatePrefixBlockIndex), codec.Encode(uint32(0)))
	d.Set(kv.Key(coreutil.StatePrefixTimestamp), codec.Encode(time.Unix(0, 0)))

	contracts := []*coreutil.ContractInfo{
		root.Contract,
		accounts.Contract,
		blocklog.Contract,
		errors.Contract,
		governance.Contract,
		evm.Contract,
	}
	if initParams.DeployTestContracts {
		contracts = append(contracts, inccounter.Contract)
		contracts = append(contracts, manyevents.Contract)
		contracts = append(contracts, testerrors.Contract)
		contracts = append(contracts, sbtestsc.Contract)
	}

	// init the state of each core contract
	root.NewStateWriter(root.Contract.StateSubrealm(d)).SetInitialState(v, contracts)
	accounts.NewStateWriter(v, accounts.Contract.StateSubrealm(d)).SetInitialState(originDeposit, baseTokenCoinInfo)
	blocklog.NewStateWriter(blocklog.Contract.StateSubrealm(d)).SetInitialState()
	errors.NewStateWriter(errors.Contract.StateSubrealm(d)).SetInitialState()
	governance.NewStateWriter(governance.Contract.StateSubrealm(d)).SetInitialState(initParams.ChainOwner, initParams.BlockKeepAmount)
	evmimpl.SetInitialState(evm.Contract.StateSubrealm(d), initParams.EVMChainID)
	if initParams.DeployTestContracts {
		inccounter.SetInitialState(inccounter.Contract.StateSubrealm(d))
	}

	block := store.Commit(d)
	if err := store.SetLatest(block.TrieRoot()); err != nil {
		panic(err)
	}
	return block, transaction.NewStateMetadata(
		v,
		block.L1Commitment(),
		&gasCoinObjectID,
		gas.DefaultFeePolicy(),
		args,
		"",
	)
}

func InitChainByAnchor(
	chainStore state.Store,
	anchor *isc.StateAnchor,
	originDeposit coin.Value,
	baseTokenCoinInfo *isc.IotaCoinInfo,
) (state.Block, error) {
	stateMetadata, err := transaction.StateMetadataFromBytes(anchor.GetStateMetadata())
	if err != nil {
		return nil, err
	}
	originBlock, _ := InitChain(
		stateMetadata.SchemaVersion,
		chainStore,
		stateMetadata.InitParams,
		*stateMetadata.GasCoinObjectID,
		originDeposit,
		baseTokenCoinInfo,
	)
	if !originBlock.L1Commitment().Equals(stateMetadata.L1Commitment) {
		return nil, fmt.Errorf(
			"l1Commitment mismatch between originAO / originBlock: %s / %s",
			stateMetadata.L1Commitment,
			originBlock.L1Commitment(),
		)
	}
	return originBlock, nil
}
