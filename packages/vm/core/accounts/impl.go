package accounts

import (
	"math/big"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/bigint"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
	"github.com/iotaledger/wasp/packages/vm/core/evm"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

func CommonAccount() isc.AgentID {
	return isc.NewAgentID(
		&cryptolib.Address{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
	)
}

var Processor = Contract.Processor(nil,
	// funcs
	FuncDeposit.WithHandler(deposit),
	FuncMintNFT.WithHandler(mintNFT),
	FuncTransferAccountToChain.WithHandler(transferAccountToChain),
	FuncTransferAllowanceTo.WithHandler(transferAllowanceTo),
	FuncWithdraw.WithHandler(withdraw),

	// Kept for compatibility
	FuncFoundryCreateNew.WithHandler(foundryCreateNew),
	//
	FuncNativeTokenCreate.WithHandler(nativeTokenCreate),
	FuncNativeTokenModifySupply.WithHandler(nativeTokenModifySupply),
	FuncNativeTokenDestroy.WithHandler(nativeTokenDestroy),

	// views
	ViewAccountNFTs.WithHandler(viewAccountNFTs),
	ViewAccountNFTAmount.WithHandler(viewAccountNFTAmount),
	ViewAccountNFTsInCollection.WithHandler(viewAccountNFTsInCollection),
	ViewAccountNFTAmountInCollection.WithHandler(viewAccountNFTAmountInCollection),
	ViewNFTIDbyMintID.WithHandler(viewNFTIDbyMintID),
	ViewAccountFoundries.WithHandler(viewAccountFoundries),
	ViewBalance.WithHandler(viewBalance),
	ViewBalanceBaseToken.WithHandler(viewBalanceBaseToken),
	ViewBalanceBaseTokenEVM.WithHandler(viewBalanceBaseTokenEVM),
	ViewBalanceNativeToken.WithHandler(viewBalanceNativeToken),
	ViewNativeToken.WithHandler(viewFoundryOutput),
	ViewGetAccountNonce.WithHandler(viewGetAccountNonce),
	ViewGetNativeTokenIDRegistry.WithHandler(viewGetNativeTokenIDRegistry),
	ViewNFTData.WithHandler(viewNFTData),
	ViewTotalAssets.WithHandler(viewTotalAssets),
)

// this expects the origin amount minus SD
func (s *StateWriter) SetInitialState(baseTokensOnAnchor uint64) {
	// initial load with base tokens from origin anchor output exceeding minimum storage deposit assumption
	s.CreditToAccount(CommonAccount(), isc.NewAssetsBaseTokensU64(baseTokensOnAnchor), isc.ChainID{})
}

// deposit is a function to deposit attached assets to the sender's chain account
// It does nothing because assets are already on the sender's account
// Allowance is ignored
func deposit(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf("accounts.deposit")
	return nil
}

// transferAllowanceTo moves whole allowance from the caller to the specified account on the chain.
// Can be sent as a request (sender is the caller) or can be called
// Params:
// - ParamAgentID. AgentID. Required
func transferAllowanceTo(ctx isc.Sandbox, targetAccount isc.AgentID) dict.Dict {
	allowance := ctx.AllowanceAvailable().Clone()
	ctx.TransferAllowedFunds(targetAccount)

	if targetAccount.Kind() != isc.AgentIDKindEthereumAddress {
		return nil // done
	}
	if !ctx.Caller().Equals(ctx.Request().SenderAccount()) {
		return nil // only issue "custom EVM tx" when this function is called directly by the request sender
	}
	// issue a "custom EVM tx" so the funds appear on the explorer
	ctx.Call(
		isc.NewMessage(
			evm.Contract.Hname(),
			evm.FuncNewL1Deposit.Hname(),
			dict.Dict{
				evm.FieldAddress:                  targetAccount.(*isc.EthereumAddressAgentID).EthAddress().Bytes(),
				evm.FieldAssets:                   allowance.Bytes(),
				evm.FieldAgentIDDepositOriginator: ctx.Caller().Bytes(),
			},
		),
		nil,
	)
	ctx.Log().Debugf("accounts.transferAllowanceTo.success: target: %s\n%s", targetAccount, ctx.AllowanceAvailable())
	return nil
}

var errCallerMustHaveL1Address = coreerrors.Register("caller must have L1 address").Create()

// withdraw sends the allowed funds to the caller's L1 address,
func withdraw(ctx isc.Sandbox) dict.Dict {
	allowance := ctx.AllowanceAvailable()
	ctx.Log().Debugf("accounts.withdraw.begin -- %s", allowance)
	if allowance.IsEmpty() {
		panic(ErrNotEnoughAllowance)
	}
	if len(allowance.NFTs) > 1 {
		panic(ErrTooManyNFTsInAllowance)
	}

	caller := ctx.Caller()
	if _, ok := caller.(*isc.ContractAgentID); ok {
		// cannot withdraw from contract account
		panic(vm.ErrUnauthorized)
	}

	// simple case, caller is not a contract, this is a straightforward withdrawal to L1
	callerAddress, ok := isc.AddressFromAgentID(caller)
	if !ok {
		panic(errCallerMustHaveL1Address)
	}
	remains := ctx.TransferAllowedFunds(ctx.AccountID())
	ctx.Requiref(remains.IsEmpty(), "internal: allowance remains must be empty")
	ctx.Send(isc.RequestParameters{
		TargetAddress: callerAddress,
		Assets:        allowance,
	})
	ctx.Log().Debugf("accounts.withdraw.success. Sent to address %s: %s",
		callerAddress.String(),
		allowance.String(),
	)
	return nil
}

// transferAccountToChain transfers the specified allowance from the sender SC's L2
// account on the target chain to the sender SC's L2 account on the origin chain.
//
// Caller must be a contract, and we will transfer the allowance from its L2 account
// on the target chain to its L2 account on the origin chain. This requires that
// this function takes the allowance into custody and in turn sends the assets as
// allowance to the origin chain, where that chain's accounts.TransferAllowanceTo()
// function then transfers it into the caller's L2 account on that chain.
//
// IMPORTANT CONSIDERATIONS:
// 1. The caller contract needs to provide sufficient base tokens in its
// allowance, to cover the gas fee GAS1 for this request.
// Note that this amount depend on the fee structure of the target chain,
// which can be different from the fee structure of the caller's own chain.
//
// 2. The caller contract also needs to provide sufficient base tokens in
// its allowance, to cover the gas fee GAS2 for the resulting request to
// accounts.TransferAllowanceTo() on the origin chain. The caller needs to
// specify this GAS2 amount through the GasReserve parameter.
//
// 3. The caller contract also needs to provide a storage deposit SD with
// this request, holding enough base tokens *independent* of the GAS1 and
// GAS2 amounts.
// Since this storage deposit is dictated by L1 we can use this amount as
// storage deposit for the resulting accounts.TransferAllowanceTo() request,
// where it will be then returned to the caller as part of the transfer.
//
// 4. This means that the caller contract needs to provide at least
// GAS1 + GAS2 + SD base tokens as assets to this request, and provide an
// allowance to the request that is exactly GAS2 + SD + transfer amount.
// Failure to meet these conditions may result in a failed request and
// worst case the assets sent to accounts.TransferAllowanceTo() could be
// irretrievably locked up in an account on the origin chain that belongs
// to the accounts core contract of the target chain.
//
// 5. The caller contract needs to set the gas budget for this request to
// GAS1 to guard against unanticipated changes in the fee structure that
// raise the gas price, otherwise the request could accidentally cannibalize
// GAS2 or even SD, with potential failure and locked up assets as a result.
func transferAccountToChain(ctx isc.Sandbox, optionalGasReserve *uint64) dict.Dict {
	allowance := ctx.AllowanceAvailable()
	ctx.Log().Debugf("accounts.transferAccountToChain.begin -- %s", allowance)
	if allowance.IsEmpty() {
		panic(ErrNotEnoughAllowance)
	}
	if len(allowance.NFTs) > 1 {
		panic(ErrTooManyNFTsInAllowance)
	}

	caller := ctx.Caller()
	callerContract, ok := caller.(*isc.ContractAgentID)
	if !ok || callerContract.Hname().IsNil() {
		// caller must be contract
		panic(vm.ErrUnauthorized)
	}

	// if the caller contract is on the same chain the transfer would end up
	// in the same L2 account it is taken from, so we do nothing in that case
	if callerContract.ChainID().Equals(ctx.ChainID()) {
		return nil
	}

	// save the assets to send to the transfer request, as specified by the allowance
	assets := allowance.Clone()

	// deduct the gas reserve GAS2 from the allowance, if possible
	gasReserve := coreutil.FromOptional(optionalGasReserve, gas.LimitsDefault.MinGasPerRequest)
	if allowance.BaseTokens < gasReserve {
		panic(ErrNotEnoughAllowance)
	}
	allowance.BaseTokens -= gasReserve

	// Warning: this will transfer all assets into the accounts core contract's L2 account.
	// Be sure everything transfers out again, or assets will be stuck forever.
	ctx.TransferAllowedFunds(ctx.AccountID())

	// Send the specified assets, which should include GAS2 and SD, as part of the
	// accounts.TransferAllowanceTo() request on the origin chain.
	// Note that the assets initially end up in the L2 account of this core accounts
	// contract on the origin chain, from where an allowance of SD plus transfer amount
	// will finally end up in the caller's L2 account on the origin chain.
	ctx.Send(isc.RequestParameters{
		TargetAddress: callerContract.Address(),
		Assets:        assets,
		Metadata: &isc.SendMetadata{
			Message: isc.NewMessage(
				Contract.Hname(),
				FuncTransferAllowanceTo.Hname(),
				dict.Dict{ParamAgentID: callerContract.Bytes()},
			),
			Allowance: allowance,
			GasBudget: gasReserve,
		},
	})
	ctx.Log().Debugf("accounts.transferAccountToChain.success. Sent to contract %s: %s",
		callerContract.String(),
		allowance.String(),
	)
	return nil
}

func nativeTokenCreate(
	ctx isc.Sandbox,
	metadata *isc.IRC30NativeTokenMetadata,
	optionalTokenScheme *iotago.TokenScheme,
) uint32 {
	sn := foundryCreateNewWithMetadata(ctx, optionalTokenScheme, metadata.Bytes())
	// Register native token as an evm ERC20 token
	ctx.Privileged().
		CallOnBehalfOf(ctx.Caller(), evm.FuncRegisterERC20NativeToken.Message(evm.ERC20NativeTokenParams{
			FoundrySN:    sn,
			Name:         metadata.Name,
			TickerSymbol: metadata.Symbol,
			Decimals:     metadata.Decimals,
		}), ctx.AllowanceAvailable())
	return sn
}

func foundryCreateNewWithMetadata(ctx isc.Sandbox, optionalTokenScheme *iotago.TokenScheme, metadata []byte) uint32 {
	ctx.Log().Debugf("accounts.foundryCreateNew")

	tokenScheme := coreutil.FromOptional[iotago.TokenScheme](optionalTokenScheme, &iotago.SimpleTokenScheme{})
	ts := util.MustTokenScheme(tokenScheme)
	ts.MeltedTokens = util.Big0
	ts.MintedTokens = util.Big0

	// create UTXO
	sn, storageDepositConsumed := ctx.Privileged().CreateNewFoundry(tokenScheme, metadata)
	ctx.Requiref(storageDepositConsumed > 0, "storage deposit Consumed > 0: assert failed")
	// storage deposit for the foundry is taken from the allowance and removed from L2 ledger
	debitBaseTokensFromAllowance(ctx, storageDepositConsumed, ctx.ChainID())

	// add to the ownership list of the account
	NewStateWriterFromSandbox(ctx).addFoundryToAccount(ctx.Caller(), sn)

	eventFoundryCreated(ctx, sn)

	return sn
}

// Params:
// - token scheme
// - must be enough allowance for the storage deposit
func foundryCreateNew(ctx isc.Sandbox, optionalTokenScheme *iotago.TokenScheme) dict.Dict {
	sn := foundryCreateNewWithMetadata(ctx, optionalTokenScheme, nil)

	return dict.Dict{
		ParamFoundrySN: codec.Uint32.Encode(sn),
	}
}

var errFoundryWithCirculatingSupply = coreerrors.Register("foundry must have zero circulating supply").Create()

// nativeTokenDestroy destroys foundry if that is possible
func nativeTokenDestroy(ctx isc.Sandbox, sn uint32) dict.Dict {
	ctx.Log().Debugf("accounts.nativeTokenDestroy")
	// check if foundry is controlled by the caller
	state := NewStateWriterFromSandbox(ctx)
	caller := ctx.Caller()
	if !state.hasFoundry(caller, sn) {
		panic(vm.ErrUnauthorized)
	}

	out, _ := state.GetFoundryOutput(sn, ctx.ChainID())
	simpleTokenScheme := util.MustTokenScheme(out.TokenScheme)
	if !bigint.IsZero(big.NewInt(0).Sub(simpleTokenScheme.MintedTokens, simpleTokenScheme.MeltedTokens)) {
		panic(errFoundryWithCirculatingSupply)
	}

	storageDepositReleased := ctx.Privileged().DestroyFoundry(sn)

	state.deleteFoundryFromAccount(caller, sn)
	state.DeleteFoundryOutput(sn)
	// the storage deposit goes to the caller's account
	state.CreditToAccount(
		caller,
		&isc.Assets{BaseTokens: storageDepositReleased},
		ctx.ChainID(),
	)
	eventFoundryDestroyed(ctx, sn)
	return nil
}

// nativeTokenModifySupply inflates (mints) or shrinks supply of token by the foundry, controlled by the caller
func nativeTokenModifySupply(ctx isc.Sandbox, sn uint32, delta *big.Int, destroy bool) {
	if bigint.IsZero(delta) {
		return
	}
	state := NewStateWriterFromSandbox(ctx)
	caller := ctx.Caller()
	// check if foundry is controlled by the caller
	if !state.hasFoundry(caller, sn) {
		panic(vm.ErrUnauthorized)
	}

	out, _ := state.GetFoundryOutput(sn, ctx.ChainID())
	if out == nil {
		panic(errFoundryNotFound)
	}

	nativeTokenID, err := out.NativeTokenID()
	ctx.RequireNoError(err, "internal")

	// accrue change on the caller's account
	// update native tokens on L2 ledger and transit foundry UTXO
	var storageDepositAdjustment int64
	if deltaAssets := isc.NewEmptyAssets().AddNativeTokens(nativeTokenID, delta); destroy {
		// take tokens to destroy from allowance
		accountID := ctx.AccountID()
		ctx.TransferAllowedFunds(accountID,
			isc.NewAssets(0, isc.NativeTokens{
				&isc.NativeToken{
					ID:     nativeTokenID,
					Amount: delta,
				},
			}),
		)
		state.DebitFromAccount(accountID, deltaAssets, ctx.ChainID())
		storageDepositAdjustment = ctx.Privileged().ModifyFoundrySupply(sn, delta.Neg(delta))
	} else {
		state.CreditToAccount(caller, deltaAssets, ctx.ChainID())
		storageDepositAdjustment = ctx.Privileged().ModifyFoundrySupply(sn, delta)
	}

	// adjust base tokens on L2 due to the possible change in storage deposit
	switch {
	case storageDepositAdjustment < 0:
		// storage deposit is taken from the allowance of the caller
		debitBaseTokensFromAllowance(ctx, uint64(-storageDepositAdjustment), ctx.ChainID())
	case storageDepositAdjustment > 0:
		// storage deposit is returned to the caller account
		state.CreditToAccount(caller, isc.NewAssetsBaseTokensU64(uint64(storageDepositAdjustment)), ctx.ChainID())
	}
	eventFoundryModified(ctx, sn)
	return
}
