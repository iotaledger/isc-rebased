package tests

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/evm/evmtest"
	"github.com/iotaledger/wasp/packages/evm/evmutil"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/tools/cluster"
)

// TODO remove this?
func setupWithNoChain(t *testing.T, opt ...waspClusterOpts) *ChainEnv {
	clu := newCluster(t, opt...)
	return &ChainEnv{t: t, Clu: clu}
}

type ChainEnv struct {
	t               *testing.T
	Clu             *cluster.Cluster
	Chain           *cluster.Chain
	testContractEnv *TestContractEnv
}

func newChainEnv(t *testing.T, clu *cluster.Cluster, chain *cluster.Chain) *ChainEnv {
	env := &ChainEnv{
		t:     t,
		Clu:   clu,
		Chain: chain,
	}
	env.testContractEnv = env.NewTestContractEnv(t)
	return env
}

type contractEnv struct {
	*ChainEnv
}

func (e *ChainEnv) createNewClient() *chainclient.Client {
	keyPair, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)
	return e.Chain.Client(keyPair)
}

func SetupWithChain(t *testing.T, opt ...waspClusterOpts) *ChainEnv {
	e := setupWithNoChain(t, opt...)
	chain, err := e.Clu.DeployDefaultChain()
	require.NoError(t, err)

	env := newChainEnv(e.t, e.Clu, chain)

	return env
}

func (e *ChainEnv) NewChainClient() *chainclient.Client {
	wallet, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)
	return chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), e.Chain.ChainID, e.Clu.Config.ISCPackageID(), wallet)
}

func (e *ChainEnv) DepositFunds(amount coin.Value, keyPair *cryptolib.KeyPair) {
	client := e.Chain.Client(keyPair)
	params := chainclient.PostRequestParams{
		Transfer:  isc.NewAssets(amount),
		Allowance: isc.NewAssets(amount - iotaclient.DefaultGasBudget),
	}
	tx, err := client.PostRequest(context.Background(), accounts.FuncDeposit.Message(), *params.WithGasBudget(iotaclient.DefaultGasBudget))
	require.NoError(e.t, err)
	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, false, 30*time.Second)
	require.NoError(e.t, err, "Error while WaitUntilAllRequestsProcessedSuccessfully for tx.ID=%v", tx.Digest)
}

func (e *ChainEnv) TransferFundsTo(assets *isc.Assets, nft *isc.NFT, keyPair *cryptolib.KeyPair, targetAccount isc.AgentID) {
	transferAssets := assets.Clone()
	transferAssets.AddBaseTokens(1 * isc.Million) // to pay for the fees
	tx, err := e.Chain.Client(keyPair).PostRequest(context.Background(), accounts.FuncTransferAllowanceTo.Message(targetAccount), chainclient.PostRequestParams{
		Transfer:  transferAssets,
		NFT:       nft,
		Allowance: assets,
	})
	require.NoError(e.t, err)
	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, false, 30*time.Second)
	require.NoError(e.t, err, "Error while WaitUntilAllRequestsProcessedSuccessfully for tx.ID=%v", tx.Digest)
}

// DeploySolidityContract deploys a given solidity contract with a given private key, returns the create contract address
// it will send the EVM request to node #0, using the default EVM chainID, that can be changed if needed
func (e *ChainEnv) DeploySolidityContract(creator *ecdsa.PrivateKey, abiJSON string, bytecode []byte, args ...interface{}) (common.Address, abi.ABI) {
	creatorAddress := crypto.PubkeyToAddress(creator.PublicKey)

	nonce := e.GetNonceEVM(creatorAddress)

	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	require.NoError(e.t, err)
	constructorArguments, err := contractABI.Pack("", args...)
	require.NoError(e.t, err)

	data := []byte{}
	data = append(data, bytecode...)
	data = append(data, constructorArguments...)

	value := big.NewInt(0)

	jsonRPCClient := e.EVMJSONRPClient(0) // send request to node #0
	gasLimit, err := jsonRPCClient.EstimateGas(context.Background(),
		ethereum.CallMsg{
			From:  creatorAddress,
			Value: value,
			Data:  data,
		})
	require.NoError(e.t, err)

	tx, err := types.SignTx(
		types.NewContractCreation(nonce, value, gasLimit, e.GetGasPriceEVM(), data),
		EVMSigner(),
		creator,
	)
	require.NoError(e.t, err)

	err = jsonRPCClient.SendTransaction(context.Background(), tx)
	require.NoError(e.t, err)

	// await tx confirmed
	fmt.Println("!!!!e.Chain.ChainID: ", e.Chain.ChainID)
	_, err = e.Clu.MultiClient().WaitUntilEVMRequestProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx.Hash(), false, 5*time.Second)
	require.NoError(e.t, err)

	return crypto.CreateAddress(creatorAddress, nonce), contractABI
}

func (e *ChainEnv) GetNonceEVM(addr common.Address) uint64 {
	nonce, err := e.EVMJSONRPClient(0).NonceAt(context.Background(), addr, nil)
	require.NoError(e.t, err)
	return nonce
}

func (e *ChainEnv) GetGasPriceEVM() *big.Int {
	res, err := e.EVMJSONRPClient(0).SuggestGasPrice(context.Background())
	require.NoError(e.t, err)
	return res
}

func (e *ChainEnv) EVMJSONRPClient(nodeIndex int) *ethclient.Client {
	return NewEVMJSONRPClient(e.t, e.Chain.ChainID.String(), e.Clu, nodeIndex)
}

func NewEVMJSONRPClient(t *testing.T, chainID string, clu *cluster.Cluster, nodeIndex int) *ethclient.Client {
	evmJSONRPCPath := fmt.Sprintf("/v1/chains/%v/evm", chainID)
	jsonRPCEndpoint := clu.Config.APIHost(nodeIndex) + evmJSONRPCPath
	rawClient, err := rpc.DialHTTP(jsonRPCEndpoint)
	require.NoError(t, err)
	jsonRPCClient := ethclient.NewClient(rawClient)
	t.Cleanup(jsonRPCClient.Close)
	return jsonRPCClient
}

func EVMSigner() types.Signer {
	return evmutil.Signer(big.NewInt(int64(evm.DefaultChainID))) // use default evm chainID
}

type TestContractEnv struct {
	t                   *testing.T
	EvmTesterAddr       common.Address
	EvmTestContractAddr common.Address
	EvmTestContractABI  abi.ABI
	EvmTesterPrivateKey *ecdsa.PrivateKey
}

func (e *ChainEnv) NewTestContractEnv(t *testing.T) *TestContractEnv {
	keyPair, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	evmPvtKey, evmAddr := solo.NewEthereumAccount()
	evmAgentID := isc.NewEthereumAddressAgentID(e.Chain.ChainID, evmAddr)
	e.TransferFundsTo(isc.NewAssets(1*isc.Million), nil, keyPair, evmAgentID)
	contractAddr, contractABI := e.DeploySolidityContract(evmPvtKey, evmtest.StorageContractABI, evmtest.StorageContractBytecode, uint32(42))
	return &TestContractEnv{
		t:                   t,
		EvmTestContractAddr: contractAddr,
		EvmTestContractABI:  contractABI,
		EvmTesterAddr:       evmAddr,
		EvmTesterPrivateKey: evmPvtKey,
	}
}

func (e *ChainEnv) CallStore(archiveClient, lightClient *ethclient.Client, input uint64) *types.Transaction {
	if archiveClient == nil {
		archiveClient = e.EVMJSONRPClient(0)
	}
	if lightClient == nil {
		lightClient = e.EVMJSONRPClient(1)
	}

	callArguments, err := e.testContractEnv.EvmTestContractABI.Pack("store", uint32(input))
	require.NoError(e.testContractEnv.t, err)
	nonce := e.GetNonceEVM(e.testContractEnv.EvmTestContractAddr)
	tx, err := types.SignTx(
		types.NewTransaction(nonce, e.testContractEnv.EvmTestContractAddr, big.NewInt(0), 100000, e.GetGasPriceEVM(), callArguments),
		EVMSigner(),
		e.testContractEnv.EvmTesterPrivateKey,
	)
	require.NoError(e.testContractEnv.t, err)
	err = archiveClient.SendTransaction(context.Background(), tx)
	require.NoError(e.testContractEnv.t, err)
	// await tx confirmed
	for i := 0; i < 3; i++ {
		fmt.Println("i: ", i)
		_, err = e.Clu.MultiClient().WaitUntilEVMRequestProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx.Hash(), false, 30*time.Second)
		if err == nil {
			break
		}
		time.Sleep(15 * time.Second)
	}
	require.NoError(e.testContractEnv.t, err)
	return tx
}

func (e *ChainEnv) CallRetrieve(archiveClient *ethclient.Client) uint32 {
	if archiveClient == nil {
		archiveClient = e.EVMJSONRPClient(0)
	}
	callArgs, err := e.testContractEnv.EvmTestContractABI.Pack("retrieve")
	require.NoError(e.t, err)
	callMsg := ethereum.CallMsg{
		To:   &e.testContractEnv.EvmTestContractAddr,
		Data: callArgs,
	}
	ret, err := archiveClient.CallContract(context.Background(), callMsg, nil)
	require.NoError(e.t, err)
	val, err := e.testContractEnv.EvmTestContractABI.Unpack("retrieve", ret)
	require.NoError(e.t, err)
	return val[0].(uint32)
}
