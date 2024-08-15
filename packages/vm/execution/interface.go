package execution

import (
	"math/big"
	"time"

	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/packages/vm/processors"
	"github.com/iotaledger/wasp/sui-go/sui"
)

// The following interfaces define the common functionality for SC execution (VM/external view calls)

type WaspContext interface {
	LocateProgram(programHash hashing.HashValue) (vmtype string, binary []byte, err error)
	GetContractRecord(contractHname isc.Hname) (ret *root.ContractRecord)
	Processors() *processors.Cache
}

type GasContext interface {
	GasBurnEnabled() bool
	GasBurnEnable(enable bool)
	GasBurn(burnCode gas.BurnCode, par ...uint64)
	GasEstimateMode() bool
}

type WaspCallContext interface {
	WaspContext
	GasContext
	isc.LogInterface
	Timestamp() time.Time
	CurrentContractAccountID() isc.AgentID
	Caller() isc.AgentID
	GetCoinBalances(agentID isc.AgentID) isc.CoinBalances
	GetBaseTokensBalance(agentID isc.AgentID) (coin.Value, *big.Int)
	GetCoinBalance(agentID isc.AgentID, coinType coin.Type) coin.Value
	Call(msg isc.Message, allowance *isc.Assets) isc.CallArguments
	ChainID() isc.ChainID
	ChainOwnerID() isc.AgentID
	ChainInfo() *isc.ChainInfo
	CurrentContractHname() isc.Hname
	Params() isc.CallArguments
	ContractStateReaderWithGasBurn() kv.KVStoreReader
	SchemaVersion() isc.SchemaVersion
	GasBurned() uint64
	GasBudgetLeft() uint64
	GetAccountNFTs(agentID isc.AgentID) []sui.ObjectID
	GetObjectBCS(id sui.ObjectID) ([]byte, bool)
	GetCoinInfo(coinType coin.Type) (*isc.SuiCoinInfo, bool)
}
