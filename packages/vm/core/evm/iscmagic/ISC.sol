// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

pragma solidity >=0.8.11;

import "./ISCSandbox.sol";
import "./ISCAccounts.sol";
import "./ISCUtil.sol";
import "./ISCPrivileged.sol";
import "./ERC20BaseTokens.sol";
import "./ERC20Coin.sol";
import "./ERC721NFTs.sol";
import "./ERC721NFTCollection.sol";

/**
 * @title ISC Library
 * @dev This library contains various interfaces and functions related to the IOTA Smart Contracts (ISC) system.
 * It provides access to the ISCSandbox, ISCAccounts, ISCUtil, ERC20BaseTokens,
 * ERC20Coin, ERC721NFTs, and ERC721NFTCollection contracts.
 */
library ISC {
    ISCSandbox constant sandbox = __iscSandbox;

    ISCAccounts constant accounts = __iscAccounts;

    ISCUtil constant util = __iscUtil;

    ERC20BaseTokens constant baseTokens = __erc20BaseTokens;

    // Get the ERC20Coin contract for the given foundry serial number
    function erc20Coin(string memory coinType) internal view returns (ERC20Coin) {
        return ERC20Coin(sandbox.ERC20CoinAddress(coinType));
    }

    ERC721NFTs constant nfts = __erc721NFTs;

    // Get the ERC721NFTCollection contract for the given collection
    function erc721NFTCollection(SuiObjectID collectionID) internal view returns (ERC721NFTCollection) {
        return ERC721NFTCollection(sandbox.erc721NFTCollectionAddress(collectionID));
    }

}
