// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/tracers"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/util/bcs"
	"github.com/iotaledger/wasp/packages/vm/gas"
	"github.com/iotaledger/wasp/sui-go/sui"
)

// SandboxBase is the common interface of Sandbox and SandboxView
type SandboxBase interface {
	Helpers
	Balance
	// Params returns the parameters of the current call
	Params() CallArguments
	// ChainID returns the chain ID
	ChainID() ChainID
	// ChainOwnerID returns the AgentID of the current owner of the chain
	ChainOwnerID() AgentID
	// ChainInfo returns information and configuration parameters of the chain
	ChainInfo() *ChainInfo
	// Contract returns the Hname of the current contract in the context
	Contract() Hname
	// AccountID returns the agentID of the current contract (i.e. chainID + contract hname)
	AccountID() AgentID
	// Caller is the agentID of the caller.
	Caller() AgentID
	// Timestamp returns the Unix timestamp of the current state in seconds
	Timestamp() time.Time
	// Log returns a logger that outputs on the local machine. It includes Panicf method
	Log() LogInterface
	// Utils provides access to common necessary functionality
	Utils() Utils
	// Gas returns sub-interface for gas related functions. It is stateful but does not modify chain's state
	Gas() Gas
	// GetObjectBCS returns the BCS-encoded contents of an object known by the chain
	GetObjectBCS(id sui.ObjectID) ([]byte, bool)
	// GetCoinInfo returns information about a coin known by the chain
	GetCoinInfo(coinType coin.Type) (*SuiCoinInfo, bool)
	// CallView calls another contract. Only calls view entry points
	CallView(Message) CallArguments
	// StateR returns the immutable k/v store of the current call (in the context of the smart contract)
	StateR() kv.KVStoreReader
	// SchemaVersion returns the schema version of the current state
	SchemaVersion() SchemaVersion
}

type SchemaVersion uint32

type Helpers interface {
	Requiref(cond bool, format string, args ...interface{})
	RequireNoError(err error, str ...string)
}

type Authorize interface {
	RequireCaller(agentID AgentID)
	RequireCallerAnyOf(agentID []AgentID)
	RequireCallerIsChainOwner()
}

type Balance interface {
	// BalanceBaseTokens returns number of base tokens in the balance of the smart contract
	BaseTokensBalance() (bts coin.Value, remainder *big.Int)
	// CoinBalance returns the balance of the given coin
	CoinBalance(p coin.Type) coin.Value
	// CoinBalances returns the balance of all coins owned by the smart contract
	CoinBalances() CoinBalances
	// OwnedObjects returns the ids of objects owned by the smart contract
	OwnedObjects() []sui.ObjectID
	// returns whether a given user owns a given amount of tokens
	HasInAccount(AgentID, *Assets) bool
}

// Sandbox is an interface given to the processor to access the VMContext
// and virtual state, transaction builder and request parameters through it.
type Sandbox interface {
	SandboxBase
	Authorize

	// State k/v store of the current call (in the context of the smart contract)
	State() kv.KVStore
	// Request return the request in the context of which the smart contract is called
	Request() Calldata

	// Call calls the entry point of the contract with parameters and allowance.
	// If the entry point is full entry point, allowance tokens are available to be moved from the caller's
	// accounts (if enough). If the entry point is view, 'allowance' has no effect
	Call(msg Message, allowance *Assets) CallArguments
	// DeployContract deploys contract on the same chain. 'initParams' are passed to the 'init' entry point
	DeployContract(programHash hashing.HashValue, name string, initParams CallArguments)
	// Event emits an event
	Event(topic string, payload []byte)
	// RegisterError registers an error
	RegisterError(messageFormat string) *VMErrorTemplate
	// GetEntropy 32 random bytes based on the hash of the current state transaction
	GetEntropy() hashing.HashValue
	// AllowanceAvailable specifies max remaining (after transfers) budget of assets the smart contract can take
	// from the caller with TransferAllowedFunds. Nil means no allowance left (zero budget)
	AllowanceAvailable() *Assets
	// TransferAllowedFunds moves assets from the caller's account to specified account within the budget set by Allowance.
	// Skipping 'assets' means transfer all Allowance().
	// The TransferAllowedFunds call mutates AllowanceAvailable
	// Returns remaining budget
	TransferAllowedFunds(target AgentID, transfer ...*Assets) *Assets
	// Send sends an on-ledger request (or a regular transaction to any L1 Address)
	Send(metadata RequestParameters)
	// StateAnchor properties of the anchor request
	StateAnchor() *StateAnchor

	RequestIndex() uint16

	// EVMTracer returns a non-nil tracer if an EVM tx is being traced
	// (e.g. with the debug_traceTransaction JSONRPC method).
	EVMTracer() *EVMTracer

	// TakeStateSnapshot takes a snapshot of the state. This is useful to implement the try/catch
	// behavior in Solidity, where the state is reverted after a low level call fails.
	TakeStateSnapshot() int
	RevertToStateSnapshot(int)

	// Privileged is a sub-interface of the sandbox which should not be called by VM plugins
	Privileged() Privileged
}

// Privileged is a sub-interface for core contracts. Should not be called by VM plugins
type Privileged interface {
	TryLoadContract(programHash hashing.HashValue) error
	GasBurnEnable(enable bool)
	GasBurnEnabled() bool
	OnWriteReceipt(CoreCallbackFunc)
	CallOnBehalfOf(caller AgentID, msg Message, allowance *Assets) CallArguments
	SendOnBehalfOf(caller ContractIdentity, metadata RequestParameters)

	// only called from EVM
	MustMoveBetweenAccounts(fromAgentID, toAgentID AgentID, assets *Assets)
	DebitFromAccount(AgentID, *big.Int)
	CreditToAccount(AgentID, *big.Int)
}

type CallArguments [][]byte

func NewCallArguments(args ...[]byte) CallArguments {
	callArguments := make(CallArguments, len(args))
	for i, v := range args {
		callArguments[i] = make([]byte, len(v))
		copy(callArguments[i], v)
	}
	return callArguments
}

func (c CallArguments) Length() int {
	return len(c)
}

func (c CallArguments) Clone() CallArguments {
	clone := make(CallArguments, len(c))
	for i, v := range c {
		clone[i] = make([]byte, len(v))
		copy(clone[i], v)
	}
	return clone
}

func (c CallArguments) At(index int) ([]byte, error) {
	if (index < 0) || (index >= len(c)) {
		return nil, fmt.Errorf("index out of range")
	}

	return (c)[index], nil
}

func (c CallArguments) MustAt(index int) []byte {
	ret, err := c.At(index)
	if err != nil {
		panic(err)
	}
	return ret
}

func (c CallArguments) OrNil(index int) []byte {
	if (index < 0) || (index >= len(c)) {
		return nil
	}
	return c[index]
}

func (c CallArguments) String() string {
	return hexutil.Encode(c.Bytes())
}

func (c CallArguments) Bytes() []byte {
	return bcs.MustMarshal(&c)
}

func CallArgumentsFromBytes(b []byte) (CallArguments, error) {
	return bcs.Unmarshal[CallArguments](b)
}

func (c CallArguments) MarshalJSON() ([]byte, error) {
	d := make([]string, len(c))

	for i, arg := range c {
		d[i] = hexutil.Encode(arg)
	}

	return json.Marshal(d)
}

func (c *CallArguments) UnmarshalJSON(data []byte) error {
	var args []string
	err := json.Unmarshal(data, &args)
	if err != nil {
		return err
	}

	cTemp := make([][]byte, len(args))

	for i, v := range args {
		(cTemp)[i], err = hexutil.Decode(v)
		if err != nil {
			return err
		}
	}

	*c = cTemp

	return nil
}

type Message struct {
	Target CallTarget    `json:"target"`
	Params CallArguments `json:"params"`
}

func NewMessage(contract Hname, ep Hname, params ...CallArguments) Message {
	msg := Message{
		Target: CallTarget{Contract: contract, EntryPoint: ep},
	}
	if len(params) > 0 {
		msg.Params = params[0]
	}
	return msg
}

func (m Message) String() string {
	return fmt.Sprintf("Message(%s, %s, %s)", m.Target.Contract, m.Target.EntryPoint, m.Params)
}

func NewMessageFromNames(contract string, ep string, params ...CallArguments) Message {
	return NewMessage(Hn(contract), Hn(ep), params...)
}

func (m Message) Clone() Message {
	return Message{
		Target: m.Target,
		Params: m.Params.Clone(),
	}
}

type CoreCallbackFunc func(contractPartition kv.KVStore, gasBurned uint64)

// RequestParameters represents parameters of the on-ledger request. The request is build from these parameters
type RequestParameters struct {
	// TargetAddress is the target address. It may represent another chain or L1 address
	TargetAddress *cryptolib.Address
	// Assets attached to the request, always taken from the caller's account.
	// It expected to contain base tokens at least the amount required for storage deposit
	// It depends on the context how it is handled when base tokens are not enough for storage deposit
	Assets *Assets
	// AdjustToMinimumStorageDeposit if true base tokens in attached fungible tokens will be added to meet minimum storage deposit requirements
	AdjustToMinimumStorageDeposit bool
	// Metadata is a request metadata. It may be nil if the request is just sending assets to L1 address
	Metadata *SendMetadata
	// SendOptions includes options of the request, such as time lock or expiry parameters
	Options SendOptions
}

type Gas interface {
	Burn(burnCode gas.BurnCode, par ...uint64)
	Budget() uint64
	Burned() uint64
	EstimateGasMode() bool
}

// StateAnchor contains properties of the anchor request/transaction in the current context
type StateAnchor struct {
	Ref        *iscmove.RefWithObject[iscmove.Anchor]
	Owner      *cryptolib.Address
	ISCPackage sui.Address
}

type SendOptions struct {
	Timelock   time.Time
	Expiration *Expiration
}

type Expiration struct {
	Time          time.Time
	ReturnAddress *cryptolib.Address
}

// SendMetadata represents content of the data payload of the request
type SendMetadata struct {
	Message   Message
	Allowance *Assets
	GasBudget uint64
}

// Utils provides various utilities that are faster on host side than on VM
// interpreter side.
type Utils interface {
	Hashing() Hashing
	ED25519() ED25519
	BLS() BLS
}

type Hashing interface {
	Blake2b(data []byte) hashing.HashValue
	Hname(name string) Hname
	Keccak(data []byte) hashing.HashValue
	Sha3(data []byte) hashing.HashValue
}

type ED25519 interface {
	// ValidSignature(data []byte, pubKey []byte, signature []byte) bool	// TODO: is it needed?
	AddressFromPublicKey(pubKey []byte) (*cryptolib.Address, error)
}

type BLS interface {
	ValidSignature(data []byte, pubKey []byte, signature []byte) bool
	AddressFromPublicKey(pubKey []byte) (iotago.Address, error)
	AggregateBLSSignatures(pubKeysBin [][]byte, sigsBin [][]byte) ([]byte, []byte, error)
}

type EVMTracer struct {
	Tracer  *tracers.Tracer
	TxIndex uint64
}

type SuiCoinInfo struct {
	CoinType    coin.Type
	Decimals    uint8
	Name        string
	Symbol      string
	Description string
	IconURL     string
	TotalSupply coin.Value
}
