// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package tests

// refactor me: Fix this to actual L1 Rebased

/*
func TestHornetStartup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping privtangle test in short mode")
	}
	l1.StartPrivtangleIfNecessary(t.Logf)

	if l1.Privtangle == nil {
		t.Skip("tests running against live network, skipping pvt tangle tests")
	}
	// pvt tangle is already stated by the cluster l1_init
	ctx := context.Background()

	//
	// Try call the faucet.
	myKeyPair := cryptolib.NewKeyPair()
	myAddress := myKeyPair.GetPublicKey().AsAddress()

	nc := nodeclient.New(l1.Config.APIAddress)
	_, err := nc.Info(ctx)
	require.NoError(t, err)

	log := testlogger.NewSilentLogger(t.Name(), true)
	client := l2connection.NewClient(l1.Config, log)

	initialOutputCount := mustOutputCount(client, myAddress)
	//
	// Check if faucet requests are working.
	client.RequestFunds(myAddress)
	for i := 0; ; i++ {
		t.Log("Waiting for a TX...")
		time.Sleep(100 * time.Millisecond)
		if initialOutputCount != mustOutputCount(client, myAddress) {
			break
		}
	}

	//
	// Check if the TX post works.
	tx, err := l2connection.MakeSimpleValueTX(client, l1.Config.FaucetKey, myAddress, 500_000)
	require.NoError(t, err)
	_, err = client.PostTxAndWaitUntilConfirmation(tx)
	require.NoError(t, err)
	for i := 0; ; i++ {
		t.Log("Waiting for a TX...")
		time.Sleep(100 * time.Millisecond)
		if initialOutputCount != mustOutputCount(client, myAddress) {
			break
		}
	}
}

func mustOutputCount(client l2connection.Client, myAddress *cryptolib.Address) int {
	return len(mustOutputMap(client, myAddress))
}

func mustOutputMap(client l2connection.Client, myAddress *cryptolib.Address) map[iotago.OutputID]iotago.Output {
	outs, err := client.OutputMap(myAddress)
	if err != nil {
		panic(fmt.Errorf("unable to get outputs as a map: %w", err))
	}
	return outs
}
*/
