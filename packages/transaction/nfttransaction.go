package transaction

import (
	`github.com/iotaledger/wasp/clients/iota-go/iotago`
	"github.com/iotaledger/wasp/packages/cryptolib"
)

// TODO: Keeping it to give context for further refactoring
type MintNFTsTransactionParams struct {
	IssuerKeyPair      cryptolib.Signer
	CollectionOutputID iotago.ObjectID
	Target             *cryptolib.Address
	ImmutableMetadata  [][]byte
}
