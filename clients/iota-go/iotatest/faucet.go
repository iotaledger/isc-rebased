package iotatest

import (
	"context"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaconn"
	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
)

var defaultSeed = []byte{
	50,
	230,
	119,
	9,
	86,
	155,
	106,
	30,
	245,
	81,
	234,
	122,
	116,
	90,
	172,
	148,
	59,
	33,
	88,
	252,
	134,
	42,
	231,
	198,
	208,
	141,
	209,
	116,
	78,
	21,
	216,
	24,
}

func MakeSignerWithFunds(index int, faucetURL string) iotasigner.Signer {
	return MakeSignerWithFundsFromSeed(defaultSeed, index, faucetURL)
}

func MakeSignerWithFundsFromSeed(seed []byte, index int, faucetURL string) iotasigner.Signer {
	keySchemeFlag := iotasigner.KeySchemeFlagEd25519
	if faucetURL == iotaconn.LocalnetFaucetURL {
		keySchemeFlag = iotasigner.KeySchemeFlagIotaEd25519
	}
	// there are only 256 different signers can be generated
	signer := iotasigner.NewSignerByIndex(seed, keySchemeFlag, index)
	err := iotaclient.RequestFundsFromFaucet(context.Background(), signer.Address(), faucetURL)
	if err != nil {
		panic(err)
	}
	return signer
}
