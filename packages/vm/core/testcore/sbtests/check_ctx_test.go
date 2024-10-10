package sbtests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

func TestMainCallsFromFullEP(t *testing.T) { run2(t, testMainCallsFromFullEP) }
func testMainCallsFromFullEP(t *testing.T) {
	_, chain := setupChain(t, nil)

	user, userAgentID := setupDeployer(t, chain)

	setupTestSandboxSC(t, chain, user)

	req := solo.NewCallParamsEx(ScName, sbtestsc.FuncCheckContextFromFullEP.Name,
		sbtestsc.ParamChainID, chain.ChainID,
		sbtestsc.ParamAgentID, isc.NewContractAgentID(chain.ChainID, HScName),
		sbtestsc.ParamCaller, userAgentID,
		sbtestsc.ParamChainOwnerID, chain.OriginatorAgentID,
	).
		WithGasBudget(10 * gas.LimitsDefault.MinGasPerRequest)
	_, err := chain.PostRequestSync(req, user)
	require.NoError(t, err)
}

func TestMainCallsFromViewEP(t *testing.T) { run2(t, testMainCallsFromViewEP) }
func testMainCallsFromViewEP(t *testing.T) {
	_, chain := setupChain(t, nil)

	user, _ := setupDeployer(t, chain)

	setupTestSandboxSC(t, chain, user)

	_, err := chain.CallViewEx(ScName, sbtestsc.FuncCheckContextFromViewEP.Name,
		sbtestsc.ParamChainID, chain.ChainID,
		sbtestsc.ParamAgentID, isc.NewContractAgentID(chain.ChainID, HScName),
		sbtestsc.ParamChainOwnerID, chain.OriginatorAgentID,
	)

	sbtestsc.FuncCheckContextFromViewEP.
		require.NoError(t, err)
}
