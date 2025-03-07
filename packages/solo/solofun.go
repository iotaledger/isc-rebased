package solo

import (
	"math"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient"
	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/testutil/l1starter"
	"github.com/iotaledger/wasp/packages/testutil/testkey"
)

func (env *Solo) IotaClient() *iotaclient.Client {
	return iotaclient.NewHTTP(env.l1Config.IotaRPCURL, l1starter.WaitUntilEffectsVisible)
}

func (env *Solo) ISCMoveClient() *iscmoveclient.Client {
	return iscmoveclient.NewHTTPClient(
		env.l1Config.IotaRPCURL,
		env.l1Config.IotaFaucetURL,
		l1starter.WaitUntilEffectsVisible,
	)
}

func (env *Solo) NewKeyPairFromIndex(index int) *cryptolib.KeyPair {
	var seed cryptolib.Seed
	if testkey.UseRandomSeed() {
		seed = cryptolib.NewSeed()
	} else {
		seed = *env.NewSeedFromIndex(index)
	}
	return cryptolib.KeyPairFromSeed(seed)
}

func (env *Solo) NewSeedFromIndex(index int) *cryptolib.Seed {
	if index < 0 {
		// SubSeed takes a "uint31"
		index += math.MaxUint32 / 2
	}
	seed := cryptolib.SubSeed(env.seed[:], uint32(index))
	return &seed
}

// NewKeyPairWithFunds generates new ed25519 signature scheme
// and requests some tokens from the UTXODB faucet.
// The amount of tokens is equal to utxodb.FundsFromFaucetAmount (=1000Mi) base tokens
// Returns signature scheme interface and public key in binary form
func (env *Solo) NewKeyPairWithFunds(seed ...*cryptolib.Seed) (*cryptolib.KeyPair, *cryptolib.Address) {
	keyPair, addr := env.NewKeyPair(seed...)
	env.GetFundsFromFaucet(addr)
	return keyPair, addr
}

func (env *Solo) GetFundsFromFaucet(target *cryptolib.Address) {
	err := iotaclient.RequestFundsFromFaucet(env.ctx, target.AsIotaAddress(), env.l1Config.IotaFaucetURL)
	require.NoError(env.T, err)
	require.GreaterOrEqual(env.T, env.L1BaseTokens(target), coin.Value(iotaclient.FundsFromFaucetAmount))
}

func (env *Solo) NewKeyPair(seedOpt ...*cryptolib.Seed) (*cryptolib.KeyPair, *cryptolib.Address) {
	return testkey.GenKeyAddr(seedOpt...)
}
