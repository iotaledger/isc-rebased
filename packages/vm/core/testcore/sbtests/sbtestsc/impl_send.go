package sbtestsc

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/coin"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

// testSplitFunds calls Send in a loop by sending 200 base tokens back to the caller
func testSplitFunds(ctx isc.Sandbox) isc.CallArguments {
	addr, ok := isc.AddressFromAgentID(ctx.Caller())
	ctx.Requiref(ok, "caller must have L1 address")
	// claim 1Mi base tokens from allowance at a time
	var baseTokensToTransfer coin.Value = 1 * isc.Million
	for !ctx.AllowanceAvailable().IsEmpty() && ctx.AllowanceAvailable().BaseTokens() >= baseTokensToTransfer {
		// send back to caller's address
		// depending on the amount of base tokens, it will exceed number of outputs or not
		ctx.TransferAllowedFunds(ctx.AccountID(), isc.NewAssets(baseTokensToTransfer))
		ctx.Send(
			isc.RequestParameters{
				TargetAddress: addr,
				Assets:        isc.NewAssets(baseTokensToTransfer),
			},
		)
	}
	return nil
}

// testSplitFundsNativeTokens calls Send for each Native token
func testSplitFundsNativeTokens(ctx isc.Sandbox) isc.CallArguments {
	addr, ok := isc.AddressFromAgentID(ctx.Caller())
	ctx.Requiref(ok, "caller must have L1 address")
	// claims all base tokens from allowance
	accountID := ctx.AccountID()
	ctx.TransferAllowedFunds(accountID, isc.NewAssets(ctx.AllowanceAvailable().BaseTokens()))
	for coinType, coinValue := range ctx.AllowanceAvailable().Coins {
		for coinValue > 0 {
			// claim 1 token from allowance at a time
			// send back to caller's address
			// depending on the amount of tokens, it will exceed number of outputs or not
			assets := isc.NewEmptyAssets().AddCoin(coinType, 1)
			rem := ctx.TransferAllowedFunds(accountID, assets)
			fmt.Printf("%s\n", rem)
			ctx.Send(
				isc.RequestParameters{
					TargetAddress: addr,
					Assets:        assets,
				},
			)
		}
	}
	return nil
}

func pingAllowanceBack(ctx isc.Sandbox) isc.CallArguments {
	caller := ctx.Caller()
	addr, ok := isc.AddressFromAgentID(caller)
	// assert caller is L1 address, not a SC
	ctx.Requiref(ok && !ctx.ChainID().IsSameChain(caller),
		"pingAllowanceBack: caller expected to be a L1 address")
	// save allowance budget because after transfer it will be modified
	toSend := ctx.AllowanceAvailable()
	if toSend.IsEmpty() {
		// nothing to send back, NOP
		return nil
	}
	// claim all transfer to the current account
	left := ctx.TransferAllowedFunds(ctx.AccountID())
	// assert what has left is empty. Only for testing
	ctx.Requiref(left.IsEmpty(), "pingAllowanceBack: inconsistency")

	// send the funds to the caller L1 address on-ledger
	ctx.Send(
		isc.RequestParameters{
			TargetAddress: addr,
			Assets:        toSend,
		},
	)
	return nil
}

// testEstimateMinimumStorageDeposit returns true if the provided allowance is enough to pay for a L1 request, panics otherwise
func testEstimateMinimumStorageDeposit(ctx isc.Sandbox) isc.CallArguments {
	addr, ok := isc.AddressFromAgentID(ctx.Caller())
	ctx.Requiref(ok, "caller must have L1 address")

	provided := ctx.AllowanceAvailable().BaseTokens

	requestParams := isc.RequestParameters{
		TargetAddress: addr,
		Metadata: &isc.SendMetadata{
			Message: isc.NewMessage(isc.Hn("foo"), isc.Hn("bar")),
		},
	}

	required := ctx.EstimateRequiredStorageDeposit(requestParams)
	ctx.Requiref(provided >= required, "not enough funds")
	return nil
}

// tries to sendback whaever NFTs are specified in allowance
func sendNFTsBack(ctx isc.Sandbox) isc.CallArguments {
	addr, ok := isc.AddressFromAgentID(ctx.Caller())
	ctx.Requiref(ok, "caller must have L1 address")

	allowance := ctx.AllowanceAvailable()
	ctx.TransferAllowedFunds(ctx.AccountID())
	for nftID := range allowance.Objects {
		ctx.Send(isc.RequestParameters{
			TargetAddress: addr,
			Assets:        isc.NewEmptyAssets().AddObject(nftID),

			Metadata: &isc.SendMetadata{},
			Options:  isc.SendOptions{},
		})
	}
	return nil
}

// just claims everything from allowance and does nothing with it
// tests the "getData" sandbox call for every NFT sent in allowance
func claimAllowance(ctx isc.Sandbox) isc.CallArguments {
	initialNFTset := ctx.OwnedObjects()
	allowance := ctx.AllowanceAvailable()
	ctx.TransferAllowedFunds(ctx.AccountID())
	ctx.Requiref(len(ctx.OwnedObjects())-len(initialNFTset) == len(allowance.Objects), "must get all NFTs from allowance")
	for _, id := range allowance.Objects {
		nftData := ctx.GetNFTData(id)
		ctx.Requiref(!nftData.ID.Empty(), "must have NFTID")
		ctx.Requiref(len(nftData.Metadata) > 0, "must have metadata")
		ctx.Requiref(nftData.Issuer != nil, "must have issuer")
	}

	return nil
}

func sendLargeRequest(ctx isc.Sandbox) isc.CallArguments {
	req := isc.RequestParameters{
		TargetAddress: cryptolib.NewRandomAddress(),
		Metadata: &isc.SendMetadata{
			Message: isc.NewMessage(
				isc.Hn("foo"),
				isc.Hn("bar"),
				dict.Dict{"x": make([]byte, ctx.Params().MustGetInt32(ParamSize))},
			),
		},

		Assets: ctx.AllowanceAvailable(),
	}
	storageDeposit := ctx.EstimateRequiredStorageDeposit(req)
	provided := ctx.AllowanceAvailable().BaseTokens
	if provided < storageDeposit {
		panic("not enough funds for storage deposit")
	}
	ctx.TransferAllowedFunds(ctx.AccountID(), isc.NewAssets(storageDeposit))
	req.Assets.Coins[coin.BaseTokenType] = storageDeposit
	ctx.Send(req)
	return nil
}
