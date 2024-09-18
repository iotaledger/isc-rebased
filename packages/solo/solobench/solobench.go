// package solobench provides tools to benchmark contracts running under solo
package solobench

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
)

type Func func(b *testing.B, chain *solo.Chain, reqs []*solo.CallParams, keyPair *cryptolib.KeyPair)

// RunBenchmarkSync processes requests synchronously, producing 1 block per request
func RunBenchmarkSync(b *testing.B, chain *solo.Chain, reqs []*solo.CallParams, keyPair *cryptolib.KeyPair) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := chain.PostRequestSync(reqs[i], keyPair)
		require.NoError(b, err)
	}
}

// RunBenchmarkAsync processes requests asynchronously, producing 1 block per many requests
func RunBenchmarkAsync(b *testing.B, chain *solo.Chain, reqs []*solo.CallParams, keyPair *cryptolib.KeyPair) {
	_ = keyPair
	txs := make([]isc.RequestID, b.N)
	for i := 0; i < b.N; i++ {
		var err error
		txs[i], err = chain.RequestFromParamsToLedger(reqs[i], nil)
		require.NoError(b, err)
	}

	chain.WaitForRequestsMark()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suiAddress := txs[i].AsSuiAddress()
		request, err := chain.Env.ISCMoveClient().GetRequestFromObjectID(chain.Env.Ctx(), &suiAddress)
		require.NoError(b, err)

		address := cryptolib.NewAddressFromSui(chain.GetLatestAnchor().ObjectID)
		iscRequest, err := isc.OnLedgerFromRequest(request, address)
		require.NoError(b, err)

		requestMap := map[isc.ChainID][]isc.Request{
			chain.ChainID: []isc.Request{
				iscRequest,
			},
		}

		go chain.Env.EnqueueRequests(requestMap)
	}
	require.True(b, chain.WaitForRequestsThrough(b.N, 20*time.Second))
}
