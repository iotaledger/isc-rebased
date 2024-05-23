#[test_only]
module isc::assets_test {
    use isc::assets;
    use sui::coin;
    use sui::sui::SUI;
    use sui::test_utils;

    public struct IOTA has drop {}

    #[test]
    fun test_assets() {
        let mut context = tx_context::dummy();
        let ctx = &mut context;
        let sui = coin::mint_for_testing<SUI>(5, ctx);
        let iota = coin::mint_for_testing<IOTA>(10, ctx);

        let mut assets = assets::new(ctx);
        assert!(!assets.has_coin<SUI>(), 0);
        assert!(!assets.has_coin<IOTA>(), 0);

        assets.add_balance(sui.into_balance());
        assert!(assets.has_coin<SUI>(), 0);

        assets.add_balance(iota.into_balance());
        assert!(assets.has_coin<IOTA>(), 0);

        let suis = assets.take_coin<SUI>(5);
        test_utils::destroy(suis);
        assert!(!assets.has_coin<SUI>(), 0);

        let iotas = assets.take_coin<IOTA>(4);
        test_utils::destroy(iotas);
        assert!(assets.has_coin<IOTA>(), 0);

        let iotas2 = assets.take_coin<IOTA>(6);
        test_utils::destroy(iotas2);
        assert!(!assets.has_coin<IOTA>(), 0);

        assets.destroy_empty();
   }
}
