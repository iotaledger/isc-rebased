// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package solo

import (
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/transaction"
	"github.com/iotaledger/wasp/packages/trie"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	vmerrors "github.com/iotaledger/wasp/packages/vm/core/errors"
	"github.com/iotaledger/wasp/packages/vm/viewcontext"
)

type CallParams struct {
	msg       isc.Message
	ftokens   *isc.Assets // ignored off-ledger
	nft       *isc.NFT
	allowance *isc.Assets
	gasBudget uint64
	nonce     uint64 // ignored for on-ledger
	sender    *cryptolib.Address
}

// NewCallParams creates a structure that wraps in one object call parameters,
// used in PostRequestSync and CallView
func NewCallParams(msg isc.Message) *CallParams {
	return &CallParams{msg: msg}
}

// NewCallParamsEx is a shortcut for NewCallParams
func NewCallParamsEx(c, ep string, params ...any) *CallParams {
	return NewCallParams(isc.NewMessageFromNames(c, ep, codec.DictFromSlice(params)))
}

func (r *CallParams) WithAllowance(allowance *isc.Assets) *CallParams {
	r.allowance = allowance.Clone()
	return r
}

func (r *CallParams) AddAllowance(allowance *isc.Assets) *CallParams {
	if r.allowance == nil {
		r.allowance = allowance.Clone()
	} else {
		r.allowance.Add(allowance)
	}
	return r
}

func (r *CallParams) AddAllowanceBaseTokens(amount uint64) *CallParams {
	return r.AddAllowance(isc.NewAssetsBaseTokens(amount))
}

func (r *CallParams) AddAllowanceNativeTokensVect(nativeTokens ...*iotago.NativeToken) *CallParams {
	if r.allowance == nil {
		r.allowance = isc.NewEmptyAssets()
	}
	r.allowance.Add(&isc.Assets{
		NativeTokens: nativeTokens,
	})
	return r
}

func (r *CallParams) AddAllowanceNativeTokens(nativeTokenID iotago.NativeTokenID, amount *big.Int) *CallParams {
	if r.allowance == nil {
		r.allowance = isc.NewEmptyAssets()
	}
	r.allowance.Add(&isc.Assets{
		NativeTokens: iotago.NativeTokens{&iotago.NativeToken{
			ID:     nativeTokenID,
			Amount: amount,
		}},
	})
	return r
}

func (r *CallParams) AddAllowanceNFTs(nfts ...iotago.NFTID) *CallParams {
	return r.AddAllowance(isc.NewEmptyAssets().AddNFTs(nfts...))
}

func (r *CallParams) WithFungibleTokens(assets *isc.Assets) *CallParams {
	r.ftokens = assets.Clone()
	return r
}

func (r *CallParams) AddFungibleTokens(assets *isc.Assets) *CallParams {
	if r.ftokens == nil {
		r.ftokens = assets.Clone()
	} else {
		r.ftokens.Add(assets)
	}
	return r
}

func (r *CallParams) AddBaseTokens(amount uint64) *CallParams {
	return r.AddFungibleTokens(isc.NewAssets(amount, nil))
}

func (r *CallParams) AddNativeTokensVect(nativeTokens ...*iotago.NativeToken) *CallParams {
	return r.AddFungibleTokens(&isc.Assets{
		NativeTokens: nativeTokens,
	})
}

func (r *CallParams) AddNativeTokens(nativeTokenID iotago.NativeTokenID, amount *big.Int) *CallParams {
	return r.AddFungibleTokens(&isc.Assets{
		NativeTokens: iotago.NativeTokens{&iotago.NativeToken{
			ID:     nativeTokenID,
			Amount: amount,
		}},
	})
}

// Adds an nft to be sent (only applicable when the call is made via on-ledger request)
func (r *CallParams) WithNFT(nft *isc.NFT) *CallParams {
	r.nft = nft
	return r
}

func (r *CallParams) GasBudget() uint64 {
	return r.gasBudget
}

func (r *CallParams) WithGasBudget(gasBudget uint64) *CallParams {
	r.gasBudget = gasBudget
	return r
}

func (r *CallParams) WithMaxAffordableGasBudget() *CallParams {
	r.gasBudget = math.MaxUint64
	return r
}

func (r *CallParams) WithNonce(nonce uint64) *CallParams {
	r.nonce = nonce
	return r
}

func (r *CallParams) WithSender(sender *cryptolib.Address) *CallParams {
	r.sender = sender
	return r
}

// NewRequestOffLedger creates off-ledger request from parameters
func (r *CallParams) NewRequestOffLedger(ch *Chain, keyPair cryptolib.Signer) isc.OffLedgerRequest {
	if r.nonce == 0 {
		r.nonce = ch.Nonce(isc.NewAgentID(keyPair.Address()))
	}
	ret := isc.NewOffLedgerRequest(ch.ID(), r.msg, r.nonce, r.gasBudget).
		WithAllowance(r.allowance)
	return ret.Sign(keyPair)
}

func (r *CallParams) NewRequestImpersonatedOffLedger(ch *Chain, address *cryptolib.Address) isc.OffLedgerRequest {
	if r.nonce == 0 {
		r.nonce = ch.Nonce(isc.NewAgentID(address))
	}
	ret := isc.NewOffLedgerRequest(ch.ID(), r.msg, r.nonce, r.gasBudget).
		WithAllowance(r.allowance)

	return isc.NewImpersonatedOffLedgerRequest(ret.(*isc.OffLedgerRequestData)).WithSenderAddress(address)
}

func (ch *Chain) createRequestTx(req *CallParams, keyPair cryptolib.Signer) (*iotago.Transaction, error) {
	/*if !cryptolib.IsVariantKeyPairValid(keyPair) {
		keyPair = ch.OriginatorPrivateKey
	}*/ // TODO: is it still needed
	L1BaseTokens := ch.Env.L1BaseTokens(keyPair.Address())
	if L1BaseTokens == 0 {
		return nil, errors.New("PostRequestSync - Signer doesn't own any base tokens on L1")
	}

	tx, err := transaction.NewRequestTransaction(ch.requestTransactionParams(req, keyPair))
	if err != nil {
		return nil, err
	}

	if tx.Essence.Outputs[0].Deposit() == 0 {
		return nil, errors.New("createRequestTx: amount == 0. Consider: solo.InitOptions{AutoAdjustStorageDeposit: true}")
	}
	return tx, err
}

func (ch *Chain) requestTransactionParams(req *CallParams, keyPair cryptolib.Signer) transaction.NewRequestTransactionParams {
	/* if !cryptolib.IsVariantKeyPairValid(keyPair) {
		keyPair = ch.OriginatorPrivateKey
	}*/ // TODO: is it still needed
	sender := req.sender
	if sender == nil {
		sender = keyPair.Address()
	}

	addr := keyPair.Address()
	allOuts, allOutIDs := ch.Env.utxoDB.GetUnspentOutputs(addr)

	return transaction.NewRequestTransactionParams{
		SenderKeyPair:    keyPair,
		SenderAddress:    sender,
		UnspentOutputs:   allOuts,
		UnspentOutputIDs: allOutIDs,
		Request: &isc.RequestParameters{
			TargetAddress: ch.ChainID.AsAddress(),
			Assets:        req.ftokens,
			Metadata: &isc.SendMetadata{
				Message:   req.msg,
				Allowance: req.allowance,
				GasBudget: req.gasBudget,
			},
			Options: isc.SendOptions{},
		},
		NFT:                             req.nft,
		DisableAutoAdjustStorageDeposit: ch.Env.disableAutoAdjustStorageDeposit,
	}
}

// requestFromParams creates an on-ledger request without posting the transaction. It is intended
// mainly for estimating gas.
func (ch *Chain) requestFromParams(req *CallParams, keyPair cryptolib.Signer) (isc.Request, error) {
	ch.Env.ledgerMutex.Lock()
	defer ch.Env.ledgerMutex.Unlock()

	tx, err := ch.createRequestTx(req, keyPair)
	if err != nil {
		return nil, err
	}
	reqs, err := isc.RequestsInTransaction(tx)
	require.NoError(ch.Env.T, err)

	for _, r := range reqs[ch.ChainID] {
		// return the first one
		return r, nil
	}
	panic("unreachable")
}

// RequestFromParamsToLedger creates transaction with one request based on parameters and sigScheme
// Then it adds it to the ledger, atomically.
// Locking on the mutex is needed to prevent mess when several goroutines work on the same address
func (ch *Chain) RequestFromParamsToLedger(req *CallParams, keyPair cryptolib.Signer) (*iotago.Transaction, isc.RequestID, error) {
	ch.Env.ledgerMutex.Lock()
	defer ch.Env.ledgerMutex.Unlock()

	tx, err := ch.createRequestTx(req, keyPair)
	if err != nil {
		return nil, isc.RequestID{}, err
	}
	err = ch.Env.AddToLedger(tx)
	// once we created transaction successfully, it should be added to the ledger smoothly
	require.NoError(ch.Env.T, err)
	txid, err := tx.ID()
	require.NoError(ch.Env.T, err)

	return tx, isc.NewRequestID(txid, 0), nil
}

// PostRequestSync posts a request synchronously sent by the test program to the smart contract on the same or another chain:
//   - creates a request transaction with the request block on it. The sigScheme is used to
//     sign the inputs of the transaction or OriginatorKeyPair is used if parameter is nil
//   - adds request transaction to UTXODB
//   - runs the request in the VM. It results in new updated virtual state and a new transaction
//     which anchors the state.
//   - adds the resulting transaction to UTXODB
//   - posts requests, contained in the resulting transaction to backlog queues of respective chains
//   - returns the result of the call to the smart contract's entry point
//
// Note that in real network of Wasp nodes (the committee) posting the transaction is completely
// asynchronous, i.e. result of the call is not available to the originator of the post.
//
// Unlike the real Wasp environment, the 'solo' environment makes PostRequestSync a synchronous call.
// It makes it possible step-by-step debug of the smart contract logic.
// The call should be used only from the main thread (goroutine)
func (ch *Chain) PostRequestSync(req *CallParams, keyPair cryptolib.Signer) (dict.Dict, error) {
	_, ret, err := ch.PostRequestSyncTx(req, keyPair)
	return ret, err
}

func (ch *Chain) PostRequestOffLedger(req *CallParams, keyPair cryptolib.Signer) (dict.Dict, error) {
	/* if !cryptolib.IsVariantKeyPairValid(keyPair) {
		keyPair = ch.OriginatorPrivateKey
	}*/ // TODO: is it still needed
	r := req.NewRequestOffLedger(ch, keyPair)
	return ch.RunOffLedgerRequest(r)
}

func (ch *Chain) PostRequestSyncTx(req *CallParams, keyPair cryptolib.Signer) (*iotago.Transaction, dict.Dict, error) {
	tx, receipt, res, err := ch.PostRequestSyncExt(req, keyPair)
	if err != nil {
		return tx, res, err
	}
	return tx, res, ch.ResolveVMError(receipt.Error).AsGoError()
}

// LastReceipt returns the receipt for the latest request processed by the chain, will return nil if the last block is empty
func (ch *Chain) LastReceipt() *isc.Receipt {
	lastBlockReceipts := ch.GetRequestReceiptsForBlock()
	if len(lastBlockReceipts) == 0 {
		return nil
	}
	blocklogReceipt := lastBlockReceipts[len(lastBlockReceipts)-1]
	return blocklogReceipt.ToISCReceipt(ch.ResolveVMError(blocklogReceipt.Error))
}

func (ch *Chain) PostRequestSyncExt(req *CallParams, keyPair cryptolib.Signer) (*iotago.Transaction, *blocklog.RequestReceipt, dict.Dict, error) {
	defer ch.logRequestLastBlock()

	tx, _, err := ch.RequestFromParamsToLedger(req, keyPair)
	require.NoError(ch.Env.T, err)
	reqs, err := ch.Env.RequestsForChain(tx, ch.ChainID)
	require.NoError(ch.Env.T, err)
	results := ch.RunRequestsSync(reqs, "post")
	if len(results) == 0 {
		return tx, nil, nil, errors.New("request has been skipped")
	}
	res := results[0]
	return tx, res.Receipt, res.Return, nil
}

// EstimateGasOnLedger executes the given on-ledger request without committing
// any changes in the ledger. It returns the amount of gas consumed.
// WARNING: Gas estimation is just an "estimate", there is no guarantees that the real call will bear the same cost, due to the turing-completeness of smart contracts
// TODO only a senderAddr, not a keyPair should be necessary to estimate (it definitely shouldn't fallback to the chain originator)
func (ch *Chain) EstimateGasOnLedger(req *CallParams, keyPair cryptolib.Signer) (dict.Dict, *blocklog.RequestReceipt, error) {
	reqCopy := *req
	r, err := ch.requestFromParams(&reqCopy, keyPair)
	if err != nil {
		return nil, nil, err
	}

	res := ch.estimateGas(r)

	return res.Return, res.Receipt, ch.ResolveVMError(res.Receipt.Error).AsGoError()
}

// EstimateGasOffLedger executes the given on-ledger request without committing
// any changes in the ledger. It returns the amount of gas consumed.
// WARNING: Gas estimation is just an "estimate", there is no guarantees that the real call will bear the same cost, due to the turing-completeness of smart contracts
func (ch *Chain) EstimateGasOffLedger(req *CallParams, keyPair cryptolib.Signer) (dict.Dict, *blocklog.RequestReceipt, error) {
	reqCopy := *req
	/*if !cryptolib.IsVariantKeyPairValid(keyPair) {
		keyPair = ch.OriginatorPrivateKey
	}*/ // TODO: is it still needed
	r := reqCopy.NewRequestImpersonatedOffLedger(ch, keyPair.Address())
	res := ch.estimateGas(r)
	return res.Return, res.Receipt, ch.ResolveVMError(res.Receipt.Error).AsGoError()
}

// EstimateNeededStorageDeposit estimates the amount of base tokens that will be
// needed to add to the request (if any) in order to cover for the storage
// deposit.
func (ch *Chain) EstimateNeededStorageDeposit(req *CallParams, keyPair cryptolib.Signer) uint64 {
	out := transaction.MakeRequestTransactionOutput(ch.requestTransactionParams(req, keyPair))
	storageDeposit := parameters.L1().Protocol.RentStructure.MinRent(out)

	reqDeposit := uint64(0)
	if req.ftokens != nil {
		reqDeposit = req.ftokens.BaseTokens
	}

	if reqDeposit >= storageDeposit {
		return 0
	}
	return storageDeposit - reqDeposit
}

func (ch *Chain) ResolveVMError(e *isc.UnresolvedVMError) *isc.VMError {
	resolved, err := vmerrors.Resolve(e, ch.CallView)
	require.NoError(ch.Env.T, err)
	return resolved
}

// CallView calls a view entry point of a smart contract.
func (ch *Chain) CallView(msg isc.Message) (dict.Dict, error) {
	latestState, err := ch.LatestState(chain.ActiveOrCommittedState)
	if err != nil {
		return nil, err
	}
	return ch.CallViewAtState(latestState, msg)
}

// CallViewEx is a shortcut for CallView
func (ch *Chain) CallViewEx(c, ep string, params ...any) (dict.Dict, error) {
	return ch.CallView(isc.NewMessageFromNames(c, ep, codec.DictFromSlice(params)))
}

func (ch *Chain) CallViewAtState(chainState state.State, msg isc.Message) (dict.Dict, error) {
	return ch.callViewByHnameAtState(chainState, msg)
}

func (ch *Chain) CallViewByHname(msg isc.Message) (dict.Dict, error) {
	latestState, err := ch.store.LatestState()
	require.NoError(ch.Env.T, err)
	return ch.callViewByHnameAtState(latestState, msg)
}

func (ch *Chain) callViewByHnameAtState(chainState state.State, msg isc.Message) (dict.Dict, error) {
	ch.Log().Debugf("callView: %s::%s", msg.Target.Contract, msg.Target.EntryPoint)

	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()

	vmctx, err := viewcontext.New(ch, chainState, false)
	if err != nil {
		return nil, err
	}
	return vmctx.CallViewExternal(msg)
}

// GetMerkleProofRaw returns Merkle proof of the key in the state
func (ch *Chain) GetMerkleProofRaw(key []byte) *trie.MerkleProof {
	ch.Log().Debugf("GetMerkleProof")

	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()

	latestState, err := ch.LatestState(chain.ActiveOrCommittedState)
	require.NoError(ch.Env.T, err)
	vmctx, err := viewcontext.New(ch, latestState, false)
	require.NoError(ch.Env.T, err)
	ret, err := vmctx.GetMerkleProof(key)
	require.NoError(ch.Env.T, err)
	return ret
}

// GetBlockProof returns Merkle proof of the key in the state
func (ch *Chain) GetBlockProof(blockIndex uint32) (*blocklog.BlockInfo, *trie.MerkleProof, error) {
	ch.Log().Debugf("GetBlockProof")

	ch.runVMMutex.Lock()
	defer ch.runVMMutex.Unlock()

	latestState, err := ch.LatestState(chain.ActiveOrCommittedState)
	require.NoError(ch.Env.T, err)
	vmctx, err := viewcontext.New(ch, latestState, false)
	if err != nil {
		return nil, nil, err
	}
	retBlockInfo, retProof, err := vmctx.GetBlockProof(blockIndex)
	if err != nil {
		return nil, nil, err
	}
	return retBlockInfo, retProof, nil
}

// GetMerkleProof return the merkle proof of the key in the smart contract. Assumes Merkle model is used
func (ch *Chain) GetMerkleProof(scHname isc.Hname, key []byte) *trie.MerkleProof {
	return ch.GetMerkleProofRaw(append(scHname.Bytes(), key...))
}

// GetL1Commitment returns state commitment taken from the anchor output
func (ch *Chain) GetL1Commitment() *state.L1Commitment {
	anchorOutput := ch.GetAnchorOutputFromL1()
	ret, err := transaction.L1CommitmentFromAliasOutput(anchorOutput.GetAliasOutput())
	require.NoError(ch.Env.T, err)
	return ret
}

// GetRootCommitment returns the root commitment of the latest state index
func (ch *Chain) GetRootCommitment() trie.Hash {
	block, err := ch.store.LatestBlock()
	require.NoError(ch.Env.T, err)
	return block.TrieRoot()
}

// GetContractStateCommitment returns commitment to the state of the specific contract, if possible
func (ch *Chain) GetContractStateCommitment(hn isc.Hname) ([]byte, error) {
	latestState, err := ch.LatestState(chain.ActiveOrCommittedState)
	require.NoError(ch.Env.T, err)
	vmctx, err := viewcontext.New(ch, latestState, false)
	if err != nil {
		return nil, err
	}
	return vmctx.GetContractStateCommitment(hn)
}

// WaitUntil waits until the condition specified by the given predicate yields true
func (ch *Chain) WaitUntil(p func() bool, maxWait ...time.Duration) bool {
	ch.Env.T.Helper()
	maxw := 10 * time.Second
	var deadline time.Time
	if len(maxWait) > 0 {
		maxw = maxWait[0]
	}
	deadline = time.Now().Add(maxw)
	for {
		if p() {
			return true
		}
		if time.Now().After(deadline) {
			ch.Env.T.Logf("WaitUntil failed waiting max %v", maxw)
			return false
		}
		time.Sleep(10 * time.Millisecond)
	}
}

const waitUntilMempoolIsEmptyDefaultTimeout = 5 * time.Second

func (ch *Chain) WaitUntilMempoolIsEmpty(timeout ...time.Duration) bool {
	realTimeout := waitUntilMempoolIsEmptyDefaultTimeout
	if len(timeout) > 0 {
		realTimeout = timeout[0]
	}

	deadline := time.Now().Add(realTimeout)
	for {
		if ch.mempool.Info().TotalPool == 0 {
			return true
		}
		time.Sleep(10 * time.Millisecond)
		if time.Now().After(deadline) {
			return false
		}
	}
}

// WaitForRequestsMark marks the amount of requests processed until now
// This allows the WaitForRequestsThrough() function to wait for the
// specified of number of requests after the mark point.
func (ch *Chain) WaitForRequestsMark() {
	ch.RequestsBlock = ch.LatestBlockIndex()
}

// WaitForRequestsThrough waits until the specified number of requests
// have been processed since the last call to WaitForRequestsMark()
func (ch *Chain) WaitForRequestsThrough(numReq int, maxWait ...time.Duration) bool {
	ch.Env.T.Helper()
	ch.Env.T.Logf("WaitForRequestsThrough: start -- block #%d -- numReq = %d", ch.RequestsBlock, numReq)
	return ch.WaitUntil(func() bool {
		ch.Env.T.Helper()
		latest := ch.LatestBlockIndex()
		for ; ch.RequestsBlock < latest; ch.RequestsBlock++ {
			receipts := ch.GetRequestReceiptsForBlock(ch.RequestsBlock + 1)
			numReq -= len(receipts)
			ch.Env.T.Logf("WaitForRequestsThrough: new block #%d with %d requests -- numReq = %d", ch.RequestsBlock, len(receipts), numReq)
		}
		return numReq <= 0
	}, maxWait...)
}
