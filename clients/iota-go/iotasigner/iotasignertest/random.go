package iotasignertest

import (
	"crypto/rand"

	"github.com/iotaledger/wasp/clients/iota-go/iotago/iotatest"
	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
	"github.com/iotaledger/wasp/packages/util/bcs"
)

func RandomSigner() iotasigner.Signer {
	b := make([]byte, 32)
	rand.Read(b)
	return iotasigner.NewSigner(b, iotasigner.KeySchemeFlagIotaEd25519)
}

func RandomSignedTransaction(signers ...iotasigner.Signer) iotasigner.SignedTransaction {
	tx := iotatest.RandomTransactionData()
	txBytes, err := bcs.Marshal(&tx.V1.Kind)
	if err != nil {
		panic(err)
	}
	var signer iotasigner.Signer
	if len(signers) == 0 {
		signer = RandomSigner()
	}
	signature, err := signer.SignTransactionBlock(txBytes, iotasigner.DefaultIntent())
	if err != nil {
		panic(err)
	}
	return *iotasigner.NewSignedTransaction(tx, signature)
}
