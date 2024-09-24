package inccounter_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/solo/solobench"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/vm/core/inccounter"
)

func checkCounter(e *solo.Chain, expected int64) {
	ret, err := e.CallView(inccounter.ViewGetCounter.Message())
	require.NoError(e.Env.T, err)

	output := lo.Must(inccounter.ViewGetCounter.DecodeOutput(ret))

	require.EqualValues(e.Env.T, expected, output)
}

func initSolo(t *testing.T) *solo.Solo {
	return solo.New(t, &solo.InitOptions{
		Debug:           true,
		PrintStackTrace: true,
	}).WithNativeContract(inccounter.Processor)
}

func TestDeployIncInitParams(t *testing.T) {
	env := initSolo(t)
	chain := env.NewChain()

	// err := chain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash, inccounter.InitParams(17))
	// require.NoError(t, err)
	checkCounter(chain, 17)
	chain.CheckAccountLedger()
}

func TestIncDefaultParam(t *testing.T) {
	env := initSolo(t)
	chain := env.NewChain()

	// err := chain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash, inccounter.InitParams(17))
	// require.NoError(t, err)
	checkCounter(chain, 17)

	req := solo.NewCallParams(inccounter.FuncIncCounter.Message(nil)).
		AddBaseTokens(1).
		WithMaxAffordableGasBudget()
	_, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)
	checkCounter(chain, 18)
	chain.CheckAccountLedger()
}

func TestIncParam(t *testing.T) {
	env := initSolo(t)
	chain := env.NewChain()

	// err := chain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash, inccounter.InitParams(17))
	// require.NoError(t, err)
	checkCounter(chain, 17)

	n := int64(3)
	req := solo.NewCallParams(inccounter.FuncIncCounter.Message(&n)).
		AddBaseTokens(1).
		WithMaxAffordableGasBudget()
	_, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)
	checkCounter(chain, 20)

	chain.CheckAccountLedger()
}

// func TestIncWith1Post(t *testing.T) {
// 	env := initSolo(t)
// 	chain := env.NewChain()

// 	// err := chain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash, inccounter.InitParams(17))
// 	// require.NoError(t, err)
// 	checkCounter(chain, 17)

// 	chain.WaitForRequestsMark()

// 	req := solo.NewCallParams(inccounter.FuncIncAndRepeatOnceAfter2s.Message()).
// 		AddBaseTokens(1_000_000).
// 		WithAllowance(isc.NewAssets(1_000_000)).
// 		WithMaxAffordableGasBudget()
// 	_, err := chain.PostRequestSync(req, nil)
// 	require.NoError(t, err)

// 	// advance logical clock to unlock that timelocked request
// 	env.AdvanceClockBy(6 * time.Second)
// 	require.True(t, chain.WaitForRequestsThrough(2))

// 	checkCounter(chain, 19)
// 	chain.CheckAccountLedger()
// }

func initBenchmark(b *testing.B) (*solo.Chain, []*solo.CallParams) {
	// setup: deploy the inccounter contract
	log := testlogger.NewSilentLogger(b.Name(), true)
	opts := solo.DefaultInitOptions()
	opts.Log = log
	env := solo.New(b, opts).WithNativeContract(inccounter.Processor)
	chain := env.NewChain()

	// err := chain.DeployContract(nil, inccounter.Contract.Name, inccounter.Contract.ProgramHash, inccounter.InitParams(0))
	// require.NoError(b, err)

	// setup: prepare N requests that call FuncIncCounter
	reqs := make([]*solo.CallParams, b.N)
	for i := 0; i < b.N; i++ {
		reqs[i] = solo.NewCallParams(inccounter.FuncIncCounter.Message(nil)).AddBaseTokens(1)
	}

	return chain, reqs
}

// BenchmarkIncSync is a benchmark for the inccounter native contract running under solo,
// processing requests synchronously, and producing 1 block per request.
// run with: go test -benchmem -cpu=1 -run=' ' -bench='Bench.*'
func BenchmarkIncSync(b *testing.B) {
	chain, reqs := initBenchmark(b)
	solobench.RunBenchmarkSync(b, chain, reqs, nil)
}

// BenchmarkIncAsync is a benchmark for the inccounter native contract running under solo,
// processing requests synchronously, and producing 1 block per many requests.
// run with: go test -benchmem -cpu=1 -run=' ' -bench='Bench.*'
func BenchmarkIncAsync(b *testing.B) {
	chain, reqs := initBenchmark(b)
	solobench.RunBenchmarkAsync(b, chain, reqs, nil)
}
