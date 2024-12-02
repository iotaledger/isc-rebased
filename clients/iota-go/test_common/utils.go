package testcommon

import (
	"context"
	"errors"

	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaconn"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
)

func GetValidatorAddress(ctx context.Context) (iotago.Address, error) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	apy, err := client.GetValidatorsApy(ctx)
	if err != nil {
		return iotago.Address{}, err
	}
	validator1 := apy.Apys[0].Address
	address, err := iotago.AddressFromHex(validator1)
	if err != nil {
		return iotago.Address{}, err
	}

	return *address, nil
}

func GetValidatorAddressWithCoins(ctx context.Context) (iotago.Address, error) {
	client := iotaclient.NewHTTP(iotaconn.AlphanetEndpointURL)
	apy, err := client.GetValidatorsApy(ctx)
	if err != nil {
		return iotago.Address{}, err
	}

	for _, apy := range apy.Apys {
		coins, err := client.GetCoins(ctx, iotaclient.GetCoinsRequest{
			Owner: iotago.MustAddressFromHex(apy.Address),
			Limit: 10,
		})
		if err != nil {
			return iotago.Address{}, err
		}
		if len(coins.Data) > 0 {
			return *iotago.MustAddressFromHex(apy.Address), nil
		}
	}

	return iotago.Address{}, errors.New("validator with coins not found")
}
