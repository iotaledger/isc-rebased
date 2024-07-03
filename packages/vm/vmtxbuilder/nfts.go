package vmtxbuilder

import (
	"bytes"
	"slices"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/vm/vmexceptions"
)

type nftIncluded struct {
	ID                iotago.NFTID
	accountingInputID iotago.OutputID // only available when the input is already accounted for (NFT was deposited in a previous block)
	accountingInput   *iotago.NFTOutput
	resultingOutput   *iotago.NFTOutput // this is not the same as in the `nativeTokenBalance` struct, this can be the accounting output, or the output leaving the chain. // TODO should refactor to follow the same logic so its easier to grok
	sentOutside       bool
}

// 3 cases of handling NFTs in txbuilder
// - NFT comes in
// - NFT goes out
// - NFT comes in and goes out in the same block
// all cases need 1 input and 1 output, but in the last case we don't need to keep the "accounting" for the NFT

func (n *nftIncluded) Clone() *nftIncluded {
	nftID := iotago.NFTID{}
	copy(nftID[:], n.ID[:])

	outputID := iotago.OutputID{}
	copy(outputID[:], n.accountingInputID[:])

	return &nftIncluded{
		ID:                nftID,
		accountingInputID: outputID,
		accountingInput:   cloneInternalNFTOutputOrNil(n.accountingInput),
		resultingOutput:   cloneInternalNFTOutputOrNil(n.resultingOutput),
	}
}

func cloneInternalNFTOutputOrNil(o *iotago.NFTOutput) *iotago.NFTOutput {
	if o == nil {
		return nil
	}
	return o.Clone().(*iotago.NFTOutput)
}

func (txb *AnchorTransactionBuilder) nftsSorted() []*nftIncluded {
	ret := make([]*nftIncluded, 0, len(txb.nftsIncluded))
	for _, nft := range txb.nftsIncluded {
		ret = append(ret, nft)
	}
	slices.SortFunc(ret, func(a, b *nftIncluded) int {
		return bytes.Compare(a.ID[:], b.ID[:])
	})
	return ret
}

func (txb *AnchorTransactionBuilder) NFTOutputs() []*iotago.NFTOutput {
	outs := make([]*iotago.NFTOutput, 0)
	for _, nft := range txb.nftsSorted() {
		if !nft.sentOutside {
			// outputs sent outside are already added to txb.postedOutputs
			outs = append(outs, nft.resultingOutput)
		}
	}
	return outs
}

func (txb *AnchorTransactionBuilder) NFTOutputsToBeUpdated() (toBeAdded, toBeRemoved []*iotago.NFTOutput, minted []iotago.Output) {
	toBeAdded = make([]*iotago.NFTOutput, 0, len(txb.nftsIncluded))
	toBeRemoved = make([]*iotago.NFTOutput, 0, len(txb.nftsIncluded))
	for _, nft := range txb.nftsSorted() {
		if nft.accountingInput != nil && nft.sentOutside {
			// to remove if input is not nil (nft exists in accounting), and it's sent to outside the chain
			toBeRemoved = append(toBeRemoved, nft.resultingOutput)
			continue
		}
		if nft.sentOutside {
			// do nothing if input is nil (doesn't exist in accounting) and it's sent outside (comes in and leaves on the same block)
			continue
		}
		// to add if input is nil (doesn't exist in accounting), and it's not sent outside the chain
		toBeAdded = append(toBeAdded, nft.resultingOutput)
	}
	return toBeAdded, toBeRemoved, txb.nftsMinted
}

func (txb *AnchorTransactionBuilder) sendNFT(o *iotago.NFTOutput) int64 {
	if txb.outputsAreFull() {
		panic(vmexceptions.ErrOutputLimitExceeded)
	}

	if txb.nftsIncluded[o.NFTID] != nil {
		// NFT comes in and out in the same block
		txb.nftsIncluded[o.NFTID].sentOutside = true
		sd := txb.nftsIncluded[o.NFTID].resultingOutput.Amount // reimburse the SD cost
		txb.nftsIncluded[o.NFTID].resultingOutput = o
		return int64(sd)
	}
	if txb.InputsAreFull() {
		panic(vmexceptions.ErrInputLimitExceeded)
	}

	// using NFT already owned by the chain
	in, outputID := txb.accountsView.NFTOutput(o.NFTID)
	toInclude := &nftIncluded{
		ID:                o.NFTID,
		accountingInput:   in,
		accountingInputID: outputID,
		resultingOutput:   o,
		sentOutside:       true,
	}

	txb.nftsIncluded[o.NFTID] = toInclude

	return int64(in.Deposit())
}
