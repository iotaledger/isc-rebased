// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

#![allow(dead_code)]

use crate::*;
use crate::coreaccounts::*;

pub struct DepositCall<'a> {
    pub func:   ScFunc<'a>,
}

pub struct FoundryCreateNewCall<'a> {
    pub func:    ScFunc<'a>,
    pub params:  MutableFoundryCreateNewParams,
    pub results: ImmutableFoundryCreateNewResults,
}

pub struct FoundryDestroyCall<'a> {
    pub func:   ScFunc<'a>,
    pub params: MutableFoundryDestroyParams,
}

pub struct FoundryModifySupplyCall<'a> {
    pub func:   ScFunc<'a>,
    pub params: MutableFoundryModifySupplyParams,
}

pub struct HarvestCall<'a> {
    pub func:   ScFunc<'a>,
    pub params: MutableHarvestParams,
}

pub struct TransferAccountToChainCall<'a> {
    pub func:   ScFunc<'a>,
    pub params: MutableTransferAccountToChainParams,
}

pub struct TransferAllowanceToCall<'a> {
    pub func:   ScFunc<'a>,
    pub params: MutableTransferAllowanceToParams,
}

pub struct WithdrawCall<'a> {
    pub func:   ScFunc<'a>,
}

pub struct AccountFoundriesCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableAccountFoundriesParams,
    pub results: ImmutableAccountFoundriesResults,
}

pub struct AccountNFTAmountCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableAccountNFTAmountParams,
    pub results: ImmutableAccountNFTAmountResults,
}

pub struct AccountNFTAmountInCollectionCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableAccountNFTAmountInCollectionParams,
    pub results: ImmutableAccountNFTAmountInCollectionResults,
}

pub struct AccountNFTsCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableAccountNFTsParams,
    pub results: ImmutableAccountNFTsResults,
}

pub struct AccountNFTsInCollectionCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableAccountNFTsInCollectionParams,
    pub results: ImmutableAccountNFTsInCollectionResults,
}

pub struct AccountsCall<'a> {
    pub func:    ScView<'a>,
    pub results: ImmutableAccountsResults,
}

pub struct BalanceCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableBalanceParams,
    pub results: ImmutableBalanceResults,
}

pub struct BalanceBaseTokenCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableBalanceBaseTokenParams,
    pub results: ImmutableBalanceBaseTokenResults,
}

pub struct BalanceNativeTokenCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableBalanceNativeTokenParams,
    pub results: ImmutableBalanceNativeTokenResults,
}

pub struct FoundryOutputCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableFoundryOutputParams,
    pub results: ImmutableFoundryOutputResults,
}

pub struct GetAccountNonceCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableGetAccountNonceParams,
    pub results: ImmutableGetAccountNonceResults,
}

pub struct GetNativeTokenIDRegistryCall<'a> {
    pub func:    ScView<'a>,
    pub results: ImmutableGetNativeTokenIDRegistryResults,
}

pub struct NftDataCall<'a> {
    pub func:    ScView<'a>,
    pub params:  MutableNftDataParams,
    pub results: ImmutableNftDataResults,
}

pub struct TotalAssetsCall<'a> {
    pub func:    ScView<'a>,
    pub results: ImmutableTotalAssetsResults,
}

pub struct ScFuncs {
}

impl ScFuncs {
    // A no-op that has the side effect of crediting any transferred tokens to the sender's account.
    pub fn deposit(ctx: &impl ScFuncClientContext) -> DepositCall {
        DepositCall {
            func: ScFunc::new(ctx, HSC_NAME, HFUNC_DEPOSIT),
        }
    }

    // Creates a new foundry with the specified token scheme, and assigns the foundry to the sender.
    pub fn foundry_create_new(ctx: &impl ScFuncClientContext) -> FoundryCreateNewCall {
        let mut f = FoundryCreateNewCall {
            func:    ScFunc::new(ctx, HSC_NAME, HFUNC_FOUNDRY_CREATE_NEW),
            params:  MutableFoundryCreateNewParams { proxy: Proxy::nil() },
            results: ImmutableFoundryCreateNewResults { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        ScFunc::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Destroys a given foundry output on L1, reimbursing the storage deposit to the caller.
    // The foundry must be owned by the caller.
    pub fn foundry_destroy(ctx: &impl ScFuncClientContext) -> FoundryDestroyCall {
        let mut f = FoundryDestroyCall {
            func:    ScFunc::new(ctx, HSC_NAME, HFUNC_FOUNDRY_DESTROY),
            params:  MutableFoundryDestroyParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    // Mints or destroys tokens for the given foundry, which must be owned by the caller.
    pub fn foundry_modify_supply(ctx: &impl ScFuncClientContext) -> FoundryModifySupplyCall {
        let mut f = FoundryModifySupplyCall {
            func:    ScFunc::new(ctx, HSC_NAME, HFUNC_FOUNDRY_MODIFY_SUPPLY),
            params:  MutableFoundryModifySupplyParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    // Moves all tokens from the chain common account to the sender's L2 account.
    // The chain owner is the only one who can call this entry point.
    pub fn harvest(ctx: &impl ScFuncClientContext) -> HarvestCall {
        let mut f = HarvestCall {
            func:    ScFunc::new(ctx, HSC_NAME, HFUNC_HARVEST),
            params:  MutableHarvestParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    // Transfers the specified allowance from the sender SC's L2 account on
    // the target chain to the sender SC's L2 account on the origin chain.
    pub fn transfer_account_to_chain(ctx: &impl ScFuncClientContext) -> TransferAccountToChainCall {
        let mut f = TransferAccountToChainCall {
            func:    ScFunc::new(ctx, HSC_NAME, HFUNC_TRANSFER_ACCOUNT_TO_CHAIN),
            params:  MutableTransferAccountToChainParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    // Transfers the specified allowance from the sender's L2 account
    // to the given L2 account on the chain.
    pub fn transfer_allowance_to(ctx: &impl ScFuncClientContext) -> TransferAllowanceToCall {
        let mut f = TransferAllowanceToCall {
            func:    ScFunc::new(ctx, HSC_NAME, HFUNC_TRANSFER_ALLOWANCE_TO),
            params:  MutableTransferAllowanceToParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    // Moves tokens from the caller's on-chain account to the caller's L1 address.
    // The number of tokens to be withdrawn must be specified via the allowance of the request.
    pub fn withdraw(ctx: &impl ScFuncClientContext) -> WithdrawCall {
        WithdrawCall {
            func: ScFunc::new(ctx, HSC_NAME, HFUNC_WITHDRAW),
        }
    }

    // Returns a set of all foundries owned by the given account.
    pub fn account_foundries(ctx: &impl ScViewClientContext) -> AccountFoundriesCall {
        let mut f = AccountFoundriesCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_ACCOUNT_FOUNDRIES),
            params:  MutableAccountFoundriesParams { proxy: Proxy::nil() },
            results: ImmutableAccountFoundriesResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the amount of NFTs owned by the given account.
    pub fn account_nft_amount(ctx: &impl ScViewClientContext) -> AccountNFTAmountCall {
        let mut f = AccountNFTAmountCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_ACCOUNT_NFT_AMOUNT),
            params:  MutableAccountNFTAmountParams { proxy: Proxy::nil() },
            results: ImmutableAccountNFTAmountResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the amount of NFTs in the specified collection owned by the given account.
    pub fn account_nft_amount_in_collection(ctx: &impl ScViewClientContext) -> AccountNFTAmountInCollectionCall {
        let mut f = AccountNFTAmountInCollectionCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_ACCOUNT_NFT_AMOUNT_IN_COLLECTION),
            params:  MutableAccountNFTAmountInCollectionParams { proxy: Proxy::nil() },
            results: ImmutableAccountNFTAmountInCollectionResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the NFT IDs for all NFTs owned by the given account.
    pub fn account_nf_ts(ctx: &impl ScViewClientContext) -> AccountNFTsCall {
        let mut f = AccountNFTsCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_ACCOUNT_NF_TS),
            params:  MutableAccountNFTsParams { proxy: Proxy::nil() },
            results: ImmutableAccountNFTsResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the NFT IDs for all NFTs in the specified collection owned by the given account.
    pub fn account_nf_ts_in_collection(ctx: &impl ScViewClientContext) -> AccountNFTsInCollectionCall {
        let mut f = AccountNFTsInCollectionCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_ACCOUNT_NF_TS_IN_COLLECTION),
            params:  MutableAccountNFTsInCollectionParams { proxy: Proxy::nil() },
            results: ImmutableAccountNFTsInCollectionResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns a set of all agent IDs that own assets on the chain.
    pub fn accounts(ctx: &impl ScViewClientContext) -> AccountsCall {
        let mut f = AccountsCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_ACCOUNTS),
            results: ImmutableAccountsResults { proxy: Proxy::nil() },
        };
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the fungible tokens owned by the given Agent ID on the chain.
    pub fn balance(ctx: &impl ScViewClientContext) -> BalanceCall {
        let mut f = BalanceCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_BALANCE),
            params:  MutableBalanceParams { proxy: Proxy::nil() },
            results: ImmutableBalanceResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the amount of base tokens owned by an agent on the chain
    pub fn balance_base_token(ctx: &impl ScViewClientContext) -> BalanceBaseTokenCall {
        let mut f = BalanceBaseTokenCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_BALANCE_BASE_TOKEN),
            params:  MutableBalanceBaseTokenParams { proxy: Proxy::nil() },
            results: ImmutableBalanceBaseTokenResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the amount of specific native tokens owned by an agent on the chain
    pub fn balance_native_token(ctx: &impl ScViewClientContext) -> BalanceNativeTokenCall {
        let mut f = BalanceNativeTokenCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_BALANCE_NATIVE_TOKEN),
            params:  MutableBalanceNativeTokenParams { proxy: Proxy::nil() },
            results: ImmutableBalanceNativeTokenResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns specified foundry output in serialized form.
    pub fn foundry_output(ctx: &impl ScViewClientContext) -> FoundryOutputCall {
        let mut f = FoundryOutputCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_FOUNDRY_OUTPUT),
            params:  MutableFoundryOutputParams { proxy: Proxy::nil() },
            results: ImmutableFoundryOutputResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the current account nonce for an Agent.
    // The account nonce is used to issue unique off-ledger requests.
    pub fn get_account_nonce(ctx: &impl ScViewClientContext) -> GetAccountNonceCall {
        let mut f = GetAccountNonceCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_GET_ACCOUNT_NONCE),
            params:  MutableGetAccountNonceParams { proxy: Proxy::nil() },
            results: ImmutableGetAccountNonceResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns a set of all native tokenIDs that are owned by the chain.
    pub fn get_native_token_id_registry(ctx: &impl ScViewClientContext) -> GetNativeTokenIDRegistryCall {
        let mut f = GetNativeTokenIDRegistryCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_GET_NATIVE_TOKEN_ID_REGISTRY),
            results: ImmutableGetNativeTokenIDRegistryResults { proxy: Proxy::nil() },
        };
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the data for a given NFT that is on the chain.
    pub fn nft_data(ctx: &impl ScViewClientContext) -> NftDataCall {
        let mut f = NftDataCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_NFT_DATA),
            params:  MutableNftDataParams { proxy: Proxy::nil() },
            results: ImmutableNftDataResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }

    // Returns the balances of all fungible tokens controlled by the chain.
    pub fn total_assets(ctx: &impl ScViewClientContext) -> TotalAssetsCall {
        let mut f = TotalAssetsCall {
            func:    ScView::new(ctx, HSC_NAME, HVIEW_TOTAL_ASSETS),
            results: ImmutableTotalAssetsResults { proxy: Proxy::nil() },
        };
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }
}
