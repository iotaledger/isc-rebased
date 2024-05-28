package sui_test

import (
	"testing"

	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/sui/conn"
	"github.com/iotaledger/wasp/sui-go/sui_signer"

	"github.com/stretchr/testify/require"
)

func TestRequestFundFromFaucet_Devnet(t *testing.T) {
	err := sui.RequestFundFromFaucet(sui_signer.TEST_ADDRESS, conn.DevnetFaucetUrl)
	require.NoError(t, err)
}

func TestRequestFundFromFaucet_Testnet(t *testing.T) {
	err := sui.RequestFundFromFaucet(sui_signer.TEST_ADDRESS, conn.TestnetFaucetUrl)
	require.NoError(t, err)
}

func TestRequestFundFromFaucet_Localnet(t *testing.T) {
	t.Skip("only run with local node is set up")
	err := sui.RequestFundFromFaucet(sui_signer.TEST_ADDRESS, conn.LocalnetFaucetUrl)
	require.NoError(t, err)
}
