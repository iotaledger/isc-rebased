package sui

import (
	"context"

	"github.com/iotaledger/wasp/sui-go/models"
	"github.com/iotaledger/wasp/sui-go/sui_types"
)

func (s *ImplSuiAPI) GetAllBalances(ctx context.Context, owner *sui_types.SuiAddress) ([]*models.Balance, error) {
	var resp []*models.Balance
	return resp, s.http.CallContext(ctx, &resp, getAllBalances, owner)
}

// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetAllCoins(ctx context.Context, req *models.GetAllCoinsRequest) (*models.CoinPage, error) {
	var resp models.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getAllCoins, req.Owner, req.Cursor, req.Limit)
}

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (s *ImplSuiAPI) GetBalance(ctx context.Context, req *models.GetBalanceRequest) (*models.Balance, error) {
	resp := models.Balance{}
	if req.CoinType == "" {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, req.Owner)
	} else {
		return &resp, s.http.CallContext(ctx, &resp, getBalance, req.Owner, req.CoinType)
	}
}

func (s *ImplSuiAPI) GetCoinMetadata(ctx context.Context, coinType string) (*models.SuiCoinMetadata, error) {
	var resp models.SuiCoinMetadata
	return &resp, s.http.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (s *ImplSuiAPI) GetCoins(ctx context.Context, req *models.GetCoinsRequest) (*models.CoinPage, error) {
	var resp models.CoinPage
	return &resp, s.http.CallContext(ctx, &resp, getCoins, req.Owner, req.CoinType, req.Cursor, req.Limit)
}

func (s *ImplSuiAPI) GetTotalSupply(ctx context.Context, coinType sui_types.ObjectType) (*models.Supply, error) {
	var resp models.Supply
	return &resp, s.http.CallContext(ctx, &resp, getTotalSupply, coinType)
}
