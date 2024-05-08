package sbtests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/vm/core/corecontracts"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
)

func TestInitSuccess(t *testing.T) {
	_, chain := setupChain(t, nil)
	err := chain.DeployContract(nil, ScName, sbtestsc.Contract.ProgramHash)
	require.NoError(t, err)
}

func TestInitFail(t *testing.T) {
	_, chain := setupChain(t, nil)
	err := chain.DeployContract(nil, ScName, sbtestsc.Contract.ProgramHash,
		codec.MakeDict(map[string]any{sbtestsc.ParamFail: 1}))
	require.Error(t, err)
}

func TestInitFailRepeat(t *testing.T) {
	_, chain := setupChain(t, nil)
	err := chain.DeployContract(nil, ScName, sbtestsc.Contract.ProgramHash,
		codec.MakeDict(map[string]any{sbtestsc.ParamFail: 1}))
	require.Error(t, err)
	_, _, rec := chain.GetInfo()
	require.EqualValues(t, len(corecontracts.All), len(rec))

	// repeat must succeed
	err = chain.DeployContract(nil, ScName, sbtestsc.Contract.ProgramHash)
	require.NoError(t, err)
	_, _, rec = chain.GetInfo()
	require.EqualValues(t, len(corecontracts.All)+1, len(rec))
}

func TestInitSuccess2(t *testing.T) {
	_, chain := setupChain(t, nil)
	err := chain.DeployContract(nil, ScName, sbtestsc.Contract.ProgramHash)
	require.NoError(t, err)
}
