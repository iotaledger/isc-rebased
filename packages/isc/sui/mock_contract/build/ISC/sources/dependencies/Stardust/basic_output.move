// Copyright (c) 2024 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

module stardust::basic_output {
    // === Imports ===

    // Sui imports.
    use sui::bag::Bag;
    use sui::balance::Balance;
    use sui::sui::SUI;
    use sui::transfer::Receiving;

    // Package imports.
    use stardust::expiration_unlock_condition::ExpirationUnlockCondition;
    use stardust::storage_deposit_return_unlock_condition::StorageDepositReturnUnlockCondition;
    use stardust::timelock_unlock_condition::TimelockUnlockCondition;

    // === Structs ===

    /// A basic output that has unlock conditions/features.
    ///   - basic outputs with expiration unlock condition must be a shared object, since that's the only
    ///     way to handle the two possible addresses that can unlock the output.
    ///   - notice that there is no `store` ability and there is no custom transfer function:
    ///       -  you can call `extract_assets`,
    ///       -  or you can call `receive` in other models to receive a `BasicOutput`.
    public struct BasicOutput has key {
        /// Hash of the `outputId` that was migrated.
        id: UID,

        /// The amount of IOTA coins held by the output.
        iota: Balance<SUI>,

        /// The `Bag` holds native tokens, key-ed by the stringified type of the asset.
        /// Example: key: "0xabcded::soon::SOON", value: Balance<0xabcded::soon::SOON>.
        native_tokens: Bag,

        /// The storage deposit return unlock condition.
        storage_deposit_return_uc: Option<StorageDepositReturnUnlockCondition>,
        /// The timelock unlock condition.
        timelock_uc: Option<TimelockUnlockCondition>,
        /// The expiration unlock condition.
        expiration_uc: Option<ExpirationUnlockCondition>,

        // Possible features, they have no effect and only here to hold data until the object is deleted.

        /// The metadata feature.
        metadata: Option<vector<u8>>,
        /// The tag feature.
        tag: Option<vector<u8>>,
        /// The sender feature.
        sender: Option<address>
    }

    // === Public-Mutative Functions ===

    /// Extract the assets stored inside the output, respecting the unlock conditions.
    ///  - The object will be deleted.
    ///  - The `StorageDepositReturnUnlockCondition` will return the deposit.
    ///  - Remaining assets (IOTA coins and native tokens) will be returned.
    public fun extract_assets(output: BasicOutput, ctx: &mut TxContext) : (Balance<SUI>, Bag) {
        // Unpack the output into its basic part.
        let BasicOutput {
            id,
            iota: mut iota,
            native_tokens,
            storage_deposit_return_uc: mut storage_deposit_return_uc,
            timelock_uc: mut timelock_uc,
            expiration_uc: mut expiration_uc,
            sender: _,
            metadata: _,
            tag: _
        } = output;

        // If the output has a timelock unlock condition, then we need to check if the timelock_uc has expired.
        if (timelock_uc.is_some()) {
            timelock_uc.extract().unlock(ctx);
        };

        // If the output has an expiration unlock condition, then we need to check who can unlock the output.
        if (expiration_uc.is_some()) {
            expiration_uc.extract().unlock(ctx);
        };

        // If the output has an storage deposit return unlock condition, then we need to return the deposit.
        if (storage_deposit_return_uc.is_some()) {
            storage_deposit_return_uc.extract().unlock(&mut iota, ctx);
        };

        // Destroy the unlock conditions.
        option::destroy_none(timelock_uc);
        option::destroy_none(expiration_uc);
        option::destroy_none(storage_deposit_return_uc);

        // Delete the output.
        object::delete(id);

        return (iota, native_tokens)
    }

    // === Public-Package Functions ===

    /// Utility function to receive a basic output in other stardust modules.
    /// Since `BasicOutput` only has `key`, it can not be received via `public_receive`.
    /// The private receiver must be implemented in its defining module (here).
    /// Other modules in the Stardust package can call this function to receive a basic output (alias, NFT).
    public(package) fun receive(parent: &mut UID, output: Receiving<BasicOutput>) : BasicOutput {
        transfer::receive(parent, output)
    }

    // === Test Functions ===

    // test only function to create a basic output
    #[test_only]
    public fun create_for_testing(
        iota: Balance<SUI>,
        native_tokens: Bag,
        storage_deposit_return_uc: Option<StorageDepositReturnUnlockCondition>,
        timelock_uc: Option<TimelockUnlockCondition>,
        expiration_uc: Option<ExpirationUnlockCondition>,
        metadata: Option<vector<u8>>,
        tag: Option<vector<u8>>,
        sender: Option<address>,
        ctx: &mut TxContext
    ): BasicOutput {
        BasicOutput {
            id: object::new(ctx),
            iota,
            native_tokens,
            storage_deposit_return_uc,
            timelock_uc,
            expiration_uc,
            metadata,
            tag,
            sender
        }
    }
}
