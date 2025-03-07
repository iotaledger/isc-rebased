package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/clients/apiextensions"
	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/corecontracts"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/contracts/inccounter"
)

const (
	varCounter    = "counter"
	varNumRepeats = "numRepeats"
	varDelay      = "delay"
)

type contractWithMessageCounterEnv struct {
	*contractEnv
}

func setupContract(env *ChainEnv) *contractWithMessageCounterEnv {
	// deposit funds onto the contract account, so it can post a L1 request
	contractAgentID := isc.NewContractAgentID(env.Chain.ChainID, inccounter.Contract.Hname())
	tx, err := env.NewChainClient().PostRequest(context.Background(), accounts.FuncTransferAllowanceTo.Message(contractAgentID), chainclient.PostRequestParams{
		Transfer:  isc.NewAssets(1_500_000),
		Allowance: isc.NewAssets(1_000_000),
	})
	require.NoError(env.t, err)
	_, err = env.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), env.Chain.ChainID, tx, false, 30*time.Second)
	require.NoError(env.t, err)

	return &contractWithMessageCounterEnv{
		contractEnv: &contractEnv{
			ChainEnv: env,
		},
	}
}

func (e *contractWithMessageCounterEnv) postRequest(contract, entryPoint isc.Hname, tokens int, params map[string]interface{}) {
	transfer := isc.NewAssets(coin.Value(tokens))
	b := isc.NewEmptyAssets()
	if transfer != nil {
		b = transfer
	}
	panic("refactor me: params is nil, used to be codec.MakeDict(params)")
	tx, err := e.NewChainClient().PostRequest(context.Background(), isc.NewMessage(contract, entryPoint, nil), chainclient.PostRequestParams{
		Transfer: b,
	})
	require.NoError(e.t, err)
	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(context.Background(), e.Chain.ChainID, tx, false, 60*time.Second)
	require.NoError(e.t, err)
}

func (e *contractEnv) checkSC(numRequests int) {
	for i := range e.Chain.CommitteeNodes {
		blockIndex, err := e.Chain.BlockIndex(i)
		require.NoError(e.t, err)
		require.Greater(e.t, blockIndex, uint32(numRequests+4))

		cl := e.Chain.Client(nil, i)
		r, err := cl.CallView(context.Background(), governance.ViewGetChainInfo.Message())
		require.NoError(e.t, err)
		info, err := governance.ViewGetChainInfo.DecodeOutput(r)
		require.NoError(e.t, err)

		require.EqualValues(e.t, e.Chain.OriginatorID(), info.ChainOwnerID)

		recs, err := e.Chain.Client(nil, i).CallView(context.Background(), root.ViewGetContractRecords.Message())
		require.NoError(e.t, err)

		contractRegistry, err := root.ViewGetContractRecords.DecodeOutput(recs)
		require.NoError(e.t, err)
		require.EqualValues(e.t, len(corecontracts.All)+1, len(contractRegistry))

		cr := contractRegistry[inccounter.Contract.Hname()]
		panic("refactor me: this equal check")
		require.EqualValues(e.t, inccounter.Contract.Name, cr.B)
	}
}

func (e *ChainEnv) checkContractCounter(expected int64) {
	for i := range e.Chain.CommitteeNodes {
		counterValue, err := e.Chain.GetCounterValue(i)
		require.NoError(e.t, err)
		require.EqualValues(e.t, expected, counterValue)
	}
}

// executed in cluster_test.go
func testIncViewCounter(t *testing.T, env *ChainEnv) {
	e := setupContract(env)
	entryPoint := isc.Hn("increment")
	e.postRequest(inccounter.Contract.Hname(), entryPoint, 0, nil)
	e.checkContractCounter(1)

	ret, err := apiextensions.CallView(
		context.Background(),
		e.Chain.Cluster.WaspClient(0),
		e.Chain.ChainID.String(),
		apiclient.ContractCallViewRequest{
			ContractHName: inccounter.Contract.Hname().String(),
			FunctionName:  "getCounter",
		})
	require.NoError(t, err)
	counter, err := inccounter.ViewGetCounter.DecodeOutput(ret)
	require.NoError(t, err)
	require.EqualValues(t, 1, counter)
}
