package disrec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	hivep2p "github.com/iotaledger/hive.go/crypto/p2p"
	"github.com/iotaledger/hive.go/logger"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/registry"
	"github.com/iotaledger/wasp/packages/testutil/testpeers"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

func initDisrecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disrec <command>",
		Short: "Disaster recovery utils.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log.Check(cmd.Help())
		},
	}
}

func initSignAndPostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign_post <tx_bytes_file> <committee_address> <committee_keys_dir> <iota_ws_url>",
		Short: "Read unsigned TX byted from the file, sign it using the committee partial keys and send to the L1 network.",
		Long: `
			We assume contents of <tx_bytes_file> contains bytes corresponding to the\n
			serialized *iotago.TransactionData.\n
			\n
			The <committee_keys_dir> should contain directories from distinct nodes,\n
			each containing the waspdb contents, namely:\n
			  - <committee_keys_dir>/<any>/identity/identity.key\n
			  - <committee_keys_dir>/<any>/dkshares/0x...hex....json\n
			`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			//
			// Read the serialized TX Data.
			txBytesFile := args[0]
			txBytes := lo.Must(os.ReadFile(txBytesFile))
			//
			// Parse the committee address.
			committeeAddressStr := args[1]
			committeeAddress := lo.Must(cryptolib.AddressFromHex(committeeAddressStr))
			//
			// Read the node keys and construct the DK Registries and the signer.
			committeeKeysDir := args[2]
			if !lo.Must(os.Stat(committeeKeysDir)).IsDir() {
				panic("must be dir")
			}
			nodeIDs := []gpa.NodeID{}
			peerIdentities := []*cryptolib.KeyPair{}
			dkRegistries := []registry.DKShareRegistryProvider{}
			for _, entry := range lo.Must(os.ReadDir(committeeKeysDir)) {
				if !entry.IsDir() {
					continue
				}
				identityPath := filepath.Join(committeeKeysDir, entry.Name(), "identity", "identity.key")
				if lo.Must(os.Stat(identityPath)).IsDir() {
					continue
				}
				dkSharesDir := filepath.Join(committeeKeysDir, entry.Name(), "dkshares")
				dkSharePath := filepath.Join(dkSharesDir, committeeAddressStr+".key")
				if lo.Must(os.Stat(dkSharePath)).IsDir() {
					continue
				}
				privKeyRaw, newlyCreated, err := hivep2p.LoadOrCreateIdentityPrivateKey(identityPath, "")
				if err != nil || newlyCreated {
					continue
				}
				privKey := lo.Must(cryptolib.PrivateKeyFromBytes(lo.Must(privKeyRaw.Raw())))
				keyPair := cryptolib.KeyPairFromPrivateKey(privKey)
				nodeID := gpa.NodeIDFromPublicKey(keyPair.GetPublicKey())
				nodeIDs = append(nodeIDs, nodeID)
				peerIdentities = append(peerIdentities, keyPair)
				dkRegistries = append(dkRegistries, lo.Must(registry.NewDKSharesRegistry(dkSharesDir, privKey)))
			}
			log := logger.NewLogger("disrec")
			signer := testpeers.NewTestDSSSigner(committeeAddress, dkRegistries, nodeIDs, peerIdentities, log)
			//
			// Sign and Post the TX to the L1.
			iotaWsUrl := args[3]
			ctx := context.Background()
			wsClient := lo.Must(iscmoveclient.NewWebsocketClient(ctx, iotaWsUrl, "", iotaclient.WaitForEffectsDisabled, log))
			res, err := wsClient.SignAndExecuteTransaction(ctx, &iotaclient.SignAndExecuteTransactionRequest{
				TxDataBytes: txBytes,
				Signer:      cryptolib.SignerToIotaSigner(signer),
				Options: &iotajsonrpc.IotaTransactionBlockResponseOptions{
					ShowEffects: true,
				},
			})
			if err != nil {
				panic(fmt.Errorf("error executing tx: %s Digest: %s", err, res.Digest))
			}
			if !res.Effects.Data.IsSuccess() {
				panic(fmt.Errorf("error executing tx: %s Digest: %s", res.Effects.Data.V1.Status.Error, res.Digest))
			}
		},
	}
	return cmd
}

func Init(rootCmd *cobra.Command) {
	disrecCmd := initDisrecCmd()
	disrecCmd.AddCommand(initSignAndPostCmd())
	rootCmd.AddCommand(disrecCmd)
}
