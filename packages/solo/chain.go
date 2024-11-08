// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package solo

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/evm/evmutil"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/state/indexedstore"
	"github.com/iotaledger/wasp/packages/util/rwutil"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	vmerrors "github.com/iotaledger/wasp/packages/vm/core/errors"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

// String is string representation for main parameters of the chain
func (ch *Chain) String() string {
	w := new(rwutil.Buffer)
	fmt.Fprintf(w, "Chain ID: %s\n", ch.ChainID)
	fmt.Fprintf(w, "Chain state controller: %s\n", ch.OriginatorAddress)
	block, err := ch.store.LatestBlock()
	require.NoError(ch.Env.T, err)
	fmt.Fprintf(w, "Root commitment: %s\n", block.TrieRoot())
	return string(*w)
}

// DumpAccounts dumps all account balances into the human-readable string
func (ch *Chain) DumpAccounts() string {
	_, chainOwnerID, _ := ch.GetInfo()
	ret := fmt.Sprintf("ChainID: %s\nChain owner: %s\n",
		ch.ChainID.String(),
		chainOwnerID.String(),
	)
	acc := ch.L2Accounts()
	for i := range acc {
		aid := acc[i]
		ret += fmt.Sprintf("  %s:\n", aid.String())
		bals := ch.L2Assets(aid)
		ret += fmt.Sprintf("%s\n", bals.String())
	}
	return ret
}

// FindContract is a view call to the 'root' smart contract on the chain.
// It returns blobCache record of the deployed smart contract with the given name
func (ch *Chain) FindContract(scName string) (*root.ContractRecord, error) {
	ret, err := ch.CallView(root.ViewFindContract.Message(isc.Hn(scName)))
	if err != nil {
		return nil, err
	}
	ok, prec := lo.Must2(root.ViewFindContract.DecodeOutput(ret))
	if !ok {
		return nil, fmt.Errorf("smart contract '%s' not found", scName)
	}
	record := *prec
	if record.Name != scName {
		return nil, fmt.Errorf("smart contract '%s' not found", scName)
	}
	return record, err
}

func (ch *Chain) GetGasFeePolicy() *gas.FeePolicy {
	res, err := ch.CallView(governance.ViewGetFeePolicy.Message())
	require.NoError(ch.Env.T, err)
	return lo.Must(governance.ViewGetFeePolicy.DecodeOutput(res))
}

func (ch *Chain) SetGasFeePolicy(user *cryptolib.KeyPair, fp *gas.FeePolicy) {
	_, err := ch.PostRequestOffLedger(NewCallParams(governance.FuncSetFeePolicy.Message(fp)), user)
	require.NoError(ch.Env.T, err)
}

func (ch *Chain) GetGasLimits() *gas.Limits {
	res, err := ch.CallView(governance.ViewGetGasLimits.Message())
	require.NoError(ch.Env.T, err)
	return lo.Must(governance.ViewGetGasLimits.DecodeOutput(res))
}

func (ch *Chain) SetGasLimits(user *cryptolib.KeyPair, gl *gas.Limits) {
	_, err := ch.PostRequestOffLedger(NewCallParams(governance.FuncSetGasLimits.Message(gl)), user)
	require.NoError(ch.Env.T, err)
}

func EVMCallDataFromArtifacts(t require.TestingT, abiJSON string, bytecode []byte, args ...interface{}) (abi.ABI, []byte) {
	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	require.NoError(t, err)

	constructorArguments, err := contractABI.Pack("", args...)
	require.NoError(t, err)

	data := []byte{}
	data = append(data, bytecode...)
	data = append(data, constructorArguments...)
	return contractABI, data
}

// DeployEVMContract deploys an evm contract on the chain
func (ch *Chain) DeployEVMContract(creator *ecdsa.PrivateKey, abiJSON string, bytecode []byte, value *big.Int, args ...interface{}) (common.Address, abi.ABI) {
	creatorAddress := crypto.PubkeyToAddress(creator.PublicKey)

	nonce := ch.Nonce(isc.NewEthereumAddressAgentID(ch.ChainID, creatorAddress))

	contractABI, data := EVMCallDataFromArtifacts(ch.Env.T, abiJSON, bytecode, args...)

	gasLimit, err := ch.EVM().EstimateGas(ethereum.CallMsg{
		From:  creatorAddress,
		Value: value,
		Data:  data,
	}, nil)
	require.NoError(ch.Env.T, err)

	tx, err := types.SignTx(
		types.NewContractCreation(nonce, value, gasLimit, ch.EVM().GasPrice(), data),
		evmutil.Signer(big.NewInt(int64(ch.EVM().ChainID()))),
		creator,
	)
	require.NoError(ch.Env.T, err)

	err = ch.EVM().SendTransaction(tx)
	require.NoError(ch.Env.T, err)
	return crypto.CreateAddress(creatorAddress, nonce), contractABI
}

// GetInfo returns information about the chain:
//   - chainID
//   - agentID of the chain owner
//   - list of contracts deployed on the chain
func (ch *Chain) GetInfo() (isc.ChainID, isc.AgentID, map[isc.Hname]*root.ContractRecord) {
	res, err := ch.CallView(governance.ViewGetChainOwner.Message())
	require.NoError(ch.Env.T, err)

	chainOwnerID, err := governance.ViewGetChainOwner.DecodeOutput(res)
	require.NoError(ch.Env.T, err)

	res, err = ch.CallView(root.ViewGetContractRecords.Message())
	require.NoError(ch.Env.T, err)

	contracts, err := root.ViewGetContractRecords.DecodeOutput(res)
	require.NoError(ch.Env.T, err)
	return ch.ChainID, chainOwnerID, lo.Associate(contracts, func(item lo.Tuple2[*isc.Hname, *root.ContractRecord]) (isc.Hname, *root.ContractRecord) {
		return *item.A, item.B
	})
}

// GetEventsForRequest calls the view in the 'blocklog' core smart contract to retrieve events for a given request.
func (ch *Chain) GetEventsForRequest(reqID isc.RequestID) ([]*isc.Event, error) {
	viewResult, err := ch.CallView(blocklog.ViewGetEventsForRequest.Message(reqID))
	if err != nil {
		return nil, err
	}
	return blocklog.ViewGetEventsForRequest.DecodeOutput(viewResult)
}

// GetEventsForBlock calls the view in the 'blocklog' core smart contract to retrieve events for a given block.
func (ch *Chain) GetEventsForBlock(blockIndex uint32) ([]*isc.Event, error) {
	viewResult, err := ch.CallView(blocklog.ViewGetEventsForBlock.Message(&blockIndex))
	if err != nil {
		return nil, err
	}
	_, events := lo.Must2(blocklog.ViewGetEventsForBlock.DecodeOutput(viewResult))
	return events, nil
}

// GetLatestBlockInfo return BlockInfo for the latest block in the chain
func (ch *Chain) GetLatestBlockInfo() *blocklog.BlockInfo {
	ret, err := ch.CallView(blocklog.ViewGetBlockInfo.Message(nil))
	require.NoError(ch.Env.T, err)
	_, bi := lo.Must2(blocklog.ViewGetBlockInfo.DecodeOutput(ret))
	return bi
}

func (ch *Chain) GetErrorMessageFormat(code isc.VMErrorCode) (string, error) {
	ret, err := ch.CallView(vmerrors.ViewGetErrorMessageFormat.Message(code))
	if err != nil {
		return "", err
	}
	return vmerrors.ViewGetErrorMessageFormat.DecodeOutput(ret)
}

// GetBlockInfo return BlockInfo for the particular block index in the chain
func (ch *Chain) GetBlockInfo(blockIndex ...uint32) (*blocklog.BlockInfo, error) {
	ret, err := ch.CallView(blocklog.ViewGetBlockInfo.Message(coreutil.Optional(blockIndex...)))
	if err != nil {
		return nil, err
	}
	_, bi := lo.Must2(blocklog.ViewGetBlockInfo.DecodeOutput(ret))
	return bi, nil
}

// IsRequestProcessed checks if the request is booked on the chain as processed
func (ch *Chain) IsRequestProcessed(reqID isc.RequestID) bool {
	ret, err := ch.CallView(blocklog.ViewIsRequestProcessed.Message(reqID))
	require.NoError(ch.Env.T, err)
	return lo.Must(blocklog.ViewIsRequestProcessed.DecodeOutput(ret))
}

// GetRequestReceipt gets the log records for a particular request, the block index and request index in the block
func (ch *Chain) GetRequestReceipt(reqID isc.RequestID) (*blocklog.RequestReceipt, bool) {
	ret, err := ch.CallView(blocklog.ViewGetRequestReceipt.Message(reqID))
	require.NoError(ch.Env.T, err)
	rec, err := blocklog.ViewGetRequestReceipt.DecodeOutput(ret)
	require.NoError(ch.Env.T, err)
	return rec, rec != nil
}

// GetRequestReceiptsForBlock returns all request log records for a particular block
func (ch *Chain) GetRequestReceiptsForBlock(blockIndex ...uint32) []*blocklog.RequestReceipt {
	res, err := ch.CallView(blocklog.ViewGetRequestReceiptsForBlock.Message(coreutil.Optional(blockIndex...)))
	if err != nil {
		return nil
	}
	recs, err := blocklog.ViewGetRequestReceiptsForBlock.DecodeOutput(res)
	if err != nil {
		ch.Log().Warn(err)
		return nil
	}
	return recs.Receipts
}

// GetRequestIDsForBlock returns the list of requestIDs settled in a particular block
func (ch *Chain) GetRequestIDsForBlock(blockIndex uint32) []isc.RequestID {
	res, err := ch.CallView(blocklog.ViewGetRequestIDsForBlock.Message(&blockIndex))
	require.NoError(ch.Env.T, err)
	_, ids := lo.Must2(blocklog.ViewGetRequestIDsForBlock.DecodeOutput(res))
	return ids
}

// GetRequestReceiptsForBlockRange returns all request log records for range of blocks, inclusively.
// Upper bound is 'latest block' is set to 0
func (ch *Chain) GetRequestReceiptsForBlockRange(fromBlockIndex, toBlockIndex uint32) []*blocklog.RequestReceipt {
	if toBlockIndex == 0 {
		toBlockIndex = ch.GetLatestBlockInfo().BlockIndex
	}
	if fromBlockIndex > toBlockIndex {
		return nil
	}
	ret := make([]*blocklog.RequestReceipt, 0)
	for i := fromBlockIndex; i <= toBlockIndex; i++ {
		recs := ch.GetRequestReceiptsForBlock(i)
		require.True(ch.Env.T, i == 0 || len(recs) != 0)
		ret = append(ret, recs...)
	}
	return ret
}

func (ch *Chain) GetRequestReceiptsForBlockRangeAsStrings(fromBlockIndex, toBlockIndex uint32) []string {
	recs := ch.GetRequestReceiptsForBlockRange(fromBlockIndex, toBlockIndex)
	ret := make([]string, len(recs))
	for i := range ret {
		ret[i] = recs[i].String()
	}
	return ret
}

// AddAllowedStateController adds the address to the allowed state controlled address list
func (ch *Chain) AddAllowedStateController(addr *cryptolib.Address, keyPair *cryptolib.KeyPair) error {
	req := NewCallParams(governance.FuncAddAllowedStateControllerAddress.Message(addr)).
		WithMaxAffordableGasBudget()
	_, err := ch.PostRequestSync(req, keyPair)
	return err
}

// AddAllowedStateController adds the address to the allowed state controlled address list
func (ch *Chain) RemoveAllowedStateController(addr *cryptolib.Address, keyPair *cryptolib.KeyPair) error {
	req := NewCallParams(governance.FuncRemoveAllowedStateControllerAddress.Message(addr)).
		WithMaxAffordableGasBudget()
	_, err := ch.PostRequestSync(req, keyPair)
	return err
}

// AddAllowedStateController adds the address to the allowed state controlled address list
func (ch *Chain) GetAllowedStateControllerAddresses() []*cryptolib.Address {
	res, err := ch.CallView(governance.ViewGetAllowedStateControllerAddresses.Message())
	require.NoError(ch.Env.T, err)
	return lo.Must(governance.ViewGetAllowedStateControllerAddresses.DecodeOutput(res))
}

// RotateStateController rotates the chain to the new controller address.
// We assume self-governed chain here.
// Mostly use for the testing of committee rotation logic, otherwise not much needed for smart contract testing
func (ch *Chain) RotateStateController(newStateAddr *cryptolib.Address, newStateKeyPair, ownerKeyPair *cryptolib.KeyPair) error {
	panic("TODO")
	// req := NewCallParams(governance.FuncRotateStateController.Message(newStateAddr)).
	// 	WithMaxAffordableGasBudget()
	// result := ch.postRequestSyncTxSpecial(req, ownerKeyPair)
	// if result.Receipt.Error == nil {
	// 	ch.StateControllerKeyPair = newStateKeyPair
	// }
	// return ch.ResolveVMError(result.Receipt.Error).AsGoError()
}

func (ch *Chain) postRequestSyncTxSpecial(req *CallParams, keyPair *cryptolib.KeyPair) *vm.RequestResult {
	panic("TODO")
	/*
		reqID, err := ch.RequestFromParamsToLedger(req, keyPair)
		require.NoError(ch.Env.T, err)
		results := ch.RunRequestsSync(reqs, "postSpecial")
		return results[0]
	*/
}

type L1L2CoinBalances struct {
	Address *cryptolib.Address
	L1      isc.CoinBalances
	L2      isc.CoinBalances
}

func (a *L1L2CoinBalances) String() string {
	return fmt.Sprintf("Address: %s\nL1 ftokens:\n  %s\nL2 ftokens:\n  %s", a.Address, a.L1, a.L2)
}

func (ch *Chain) L1L2Funds(addr *cryptolib.Address) *L1L2CoinBalances {
	return &L1L2CoinBalances{
		Address: addr,
		L1:      ch.Env.L1CoinBalances(addr),
		L2:      ch.L2Assets(isc.NewAddressAgentID(addr)).Coins,
	}
}

func (ch *Chain) GetL2FundsFromFaucet(agentID isc.AgentID, baseTokens ...coin.Value) {
	// find a deterministic L1 address that has 0 balance
	walletKey, walletAddr := func() (*cryptolib.KeyPair, *cryptolib.Address) {
		masterSeed := []byte("GetL2FundsFromFaucet")
		i := uint32(0)
		for {
			ss := cryptolib.SubSeed(masterSeed, i)
			key, addr := ch.Env.NewKeyPair(&ss)
			ch.Env.GetFundsFromFaucet(addr)
			if ch.L2BaseTokens(isc.NewAddressAgentID(addr)) == 0 {
				return key, addr
			}
			i++
		}
	}()

	var amount coin.Value
	if len(baseTokens) > 0 {
		amount = baseTokens[0]
	} else {
		amount = ch.Env.L1BaseTokens(walletAddr) - TransferAllowanceToGasBudgetBaseTokens
	}
	err := ch.TransferAllowanceTo(
		isc.NewAssets(amount),
		agentID,
		walletKey,
	)
	require.NoError(ch.Env.T, err)
}

func (ch *Chain) Store() indexedstore.IndexedStore {
	return ch.store
}

func (ch *Chain) LatestState() (state.State, error) {
	return ch.store.LatestState()
}

func (ch *Chain) LatestBlock() state.Block {
	b, err := ch.store.LatestBlock()
	require.NoError(ch.Env.T, err)
	return b
}

func (ch *Chain) Nonce(agentID isc.AgentID) uint64 {
	if evmAgentID, ok := agentID.(*isc.EthereumAddressAgentID); ok {
		nonce, err := ch.EVM().TransactionCount(evmAgentID.EthAddress(), nil)
		require.NoError(ch.Env.T, err)
		return nonce
	}
	res, err := ch.CallView(accounts.ViewGetAccountNonce.Message(&agentID))
	require.NoError(ch.Env.T, err)
	return lo.Must(accounts.ViewGetAccountNonce.DecodeOutput(res))
}

func (ch *Chain) LatestBlockIndex() uint32 {
	return ch.GetLatestBlockInfo().BlockIndex
}
