package iotatest

import (
	"context"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
	testcommon "github.com/iotaledger/wasp/clients/iota-go/test_common"
)

func MakeSignerWithFunds(index int, faucetURL string) iotasigner.Signer {
	// CHeck l1starter -> UseRandomSeed -> switch between testseed/random seed here

	return MakeSignerWithFundsFromSeed(testcommon.TestSeed, index, faucetURL)
}

func MakeSignerWithFundsFromSeed(seed []byte, index int, faucetURL string) iotasigner.Signer {
	keySchemeFlag := iotasigner.KeySchemeFlagIotaEd25519

	// there are only 256 different signers can be generated
	signer := iotasigner.NewSignerByIndex(seed, keySchemeFlag, index)
	err := iotaclient.RequestFundsFromFaucet(context.Background(), signer.Address(), faucetURL)
	if err != nil {
		panic(err)
	}
	return signer
}
