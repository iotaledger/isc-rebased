package clients

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"github.com/iotaledger/wasp/clients/iota-go/contracts"
	"github.com/iotaledger/wasp/clients/iota-go/iotaclient"
	"github.com/iotaledger/wasp/clients/iota-go/iotaconn"
	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/clients/iota-go/iotajsonrpc"
	"github.com/iotaledger/wasp/clients/iota-go/iotasigner"
	"github.com/iotaledger/wasp/clients/iscmove/iscmoveclient"
	"github.com/iotaledger/wasp/packages/cryptolib"
)

type L1Config struct {
	APIURL    string
	FaucetURL string
	GraphURL  string
}

type L1Client interface {
	GetDynamicFieldObject(
		ctx context.Context,
		req iotaclient.GetDynamicFieldObjectRequest,
	) (*iotajsonrpc.IotaObjectResponse, error)
	GetDynamicFields(
		ctx context.Context,
		req iotaclient.GetDynamicFieldsRequest,
	) (*iotajsonrpc.DynamicFieldPage, error)
	GetOwnedObjects(
		ctx context.Context,
		req iotaclient.GetOwnedObjectsRequest,
	) (*iotajsonrpc.ObjectsPage, error)
	QueryEvents(
		ctx context.Context,
		req iotaclient.QueryEventsRequest,
	) (*iotajsonrpc.EventPage, error)
	QueryTransactionBlocks(
		ctx context.Context,
		req iotaclient.QueryTransactionBlocksRequest,
	) (*iotajsonrpc.TransactionBlocksPage, error)
	ResolveNameServiceAddress(ctx context.Context, iotaName string) (*iotago.Address, error)
	ResolveNameServiceNames(
		ctx context.Context,
		req iotaclient.ResolveNameServiceNamesRequest,
	) (*iotajsonrpc.IotaNamePage, error)
	DevInspectTransactionBlock(
		ctx context.Context,
		req iotaclient.DevInspectTransactionBlockRequest,
	) (*iotajsonrpc.DevInspectResults, error)
	DryRunTransaction(
		ctx context.Context,
		txDataBytes iotago.Base64Data,
	) (*iotajsonrpc.DryRunTransactionBlockResponse, error)
	ExecuteTransactionBlock(
		ctx context.Context,
		req iotaclient.ExecuteTransactionBlockRequest,
	) (*iotajsonrpc.IotaTransactionBlockResponse, error)
	GetCommitteeInfo(
		ctx context.Context,
		epoch *iotajsonrpc.BigInt, // optional
	) (*iotajsonrpc.CommitteeInfo, error)
	GetLatestIotaSystemState(ctx context.Context) (*iotajsonrpc.IotaSystemStateSummary, error)
	GetReferenceGasPrice(ctx context.Context) (*iotajsonrpc.BigInt, error)
	GetStakes(ctx context.Context, owner *iotago.Address) ([]*iotajsonrpc.DelegatedStake, error)
	GetStakesByIds(ctx context.Context, stakedIotaIds []iotago.ObjectID) ([]*iotajsonrpc.DelegatedStake, error)
	GetValidatorsApy(ctx context.Context) (*iotajsonrpc.ValidatorsApy, error)
	BatchTransaction(
		ctx context.Context,
		req iotaclient.BatchTransactionRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	MergeCoins(
		ctx context.Context,
		req iotaclient.MergeCoinsRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	MoveCall(
		ctx context.Context,
		req iotaclient.MoveCallRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	Pay(
		ctx context.Context,
		req iotaclient.PayRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	PayAllIota(
		ctx context.Context,
		req iotaclient.PayAllIotaRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	PayIota(
		ctx context.Context,
		req iotaclient.PayIotaRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	Publish(
		ctx context.Context,
		req iotaclient.PublishRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	RequestAddStake(
		ctx context.Context,
		req iotaclient.RequestAddStakeRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	RequestWithdrawStake(
		ctx context.Context,
		req iotaclient.RequestWithdrawStakeRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	SplitCoin(
		ctx context.Context,
		req iotaclient.SplitCoinRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	SplitCoinEqual(
		ctx context.Context,
		req iotaclient.SplitCoinEqualRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	TransferObject(
		ctx context.Context,
		req iotaclient.TransferObjectRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	TransferIota(
		ctx context.Context,
		req iotaclient.TransferIotaRequest,
	) (*iotajsonrpc.TransactionBytes, error)
	GetCoinObjsForTargetAmount(
		ctx context.Context,
		address *iotago.Address,
		targetAmount uint64,
	) (iotajsonrpc.Coins, error)
	SignAndExecuteTransaction(
		ctx context.Context,
		req *iotaclient.SignAndExecuteTransactionRequest,
	) (*iotajsonrpc.IotaTransactionBlockResponse, error)
	PublishContract(
		ctx context.Context,
		signer iotasigner.Signer,
		modules []*iotago.Base64Data,
		dependencies []*iotago.Address,
		gasBudget uint64,
		options *iotajsonrpc.IotaTransactionBlockResponseOptions,
	) (*iotajsonrpc.IotaTransactionBlockResponse, *iotago.PackageID, error)
	UpdateObjectRef(
		ctx context.Context,
		ref *iotago.ObjectRef,
	) (*iotago.ObjectRef, error)
	MintToken(
		ctx context.Context,
		signer iotasigner.Signer,
		packageID *iotago.PackageID,
		tokenName string,
		treasuryCap *iotago.ObjectID,
		mintAmount uint64,
		options *iotajsonrpc.IotaTransactionBlockResponseOptions,
	) (*iotajsonrpc.IotaTransactionBlockResponse, error)
	GetIotaCoinsOwnedByAddress(ctx context.Context, address *iotago.Address) (iotajsonrpc.Coins, error)
	BatchGetObjectsOwnedByAddress(
		ctx context.Context,
		address *iotago.Address,
		options *iotajsonrpc.IotaObjectDataOptions,
		filterType string,
	) ([]iotajsonrpc.IotaObjectResponse, error)
	BatchGetFilteredObjectsOwnedByAddress(
		ctx context.Context,
		address *iotago.Address,
		options *iotajsonrpc.IotaObjectDataOptions,
		filter func(*iotajsonrpc.IotaObjectData) bool,
	) ([]iotajsonrpc.IotaObjectResponse, error)
	GetAllBalances(ctx context.Context, owner *iotago.Address) ([]*iotajsonrpc.Balance, error)
	GetAllCoins(ctx context.Context, req iotaclient.GetAllCoinsRequest) (*iotajsonrpc.CoinPage, error)
	GetBalance(ctx context.Context, req iotaclient.GetBalanceRequest) (*iotajsonrpc.Balance, error)
	GetCoinMetadata(ctx context.Context, coinType string) (*iotajsonrpc.IotaCoinMetadata, error)
	GetCoins(ctx context.Context, req iotaclient.GetCoinsRequest) (*iotajsonrpc.CoinPage, error)
	GetTotalSupply(ctx context.Context, coinType iotago.ObjectType) (*iotajsonrpc.Supply, error)
	GetChainIdentifier(ctx context.Context) (string, error)
	GetCheckpoint(ctx context.Context, checkpointId *iotajsonrpc.BigInt) (*iotajsonrpc.Checkpoint, error)
	GetCheckpoints(ctx context.Context, req iotaclient.GetCheckpointsRequest) (*iotajsonrpc.CheckpointPage, error)
	GetEvents(ctx context.Context, digest *iotago.TransactionDigest) ([]*iotajsonrpc.IotaEvent, error)
	GetLatestCheckpointSequenceNumber(ctx context.Context) (string, error)
	GetObject(ctx context.Context, req iotaclient.GetObjectRequest) (*iotajsonrpc.IotaObjectResponse, error)
	GetProtocolConfig(
		ctx context.Context,
		version *iotajsonrpc.BigInt, // optional
	) (*iotajsonrpc.ProtocolConfig, error)
	GetTotalTransactionBlocks(ctx context.Context) (string, error)
	GetTransactionBlock(ctx context.Context, req iotaclient.GetTransactionBlockRequest) (*iotajsonrpc.IotaTransactionBlockResponse, error)
	MultiGetObjects(ctx context.Context, req iotaclient.MultiGetObjectsRequest) ([]iotajsonrpc.IotaObjectResponse, error)
	MultiGetTransactionBlocks(
		ctx context.Context,
		req iotaclient.MultiGetTransactionBlocksRequest,
	) ([]*iotajsonrpc.IotaTransactionBlockResponse, error)
	TryGetPastObject(
		ctx context.Context,
		req iotaclient.TryGetPastObjectRequest,
	) (*iotajsonrpc.IotaPastObjectResponse, error)
	TryMultiGetPastObjects(
		ctx context.Context,
		req iotaclient.TryMultiGetPastObjectsRequest,
	) ([]*iotajsonrpc.IotaPastObjectResponse, error)
	RequestFunds(ctx context.Context, address cryptolib.Address) error
	Health(ctx context.Context) error
	L2() L2Client
	DeployISCContracts(ctx context.Context, signer iotasigner.Signer) (iotago.PackageID, error)
}

var _ L1Client = &l1Client{}

type l1Client struct {
	*iotaclient.Client

	Config L1Config
}

func (c *l1Client) RequestFunds(ctx context.Context, address cryptolib.Address) error {
	faucetURL := c.Config.FaucetURL
	if faucetURL == "" {
		faucetURL = iotaconn.FaucetURL(c.Config.APIURL)
	}
	return iotaclient.RequestFundsFromFaucet(ctx, address.AsIotaAddress(), faucetURL)
}

func (c *l1Client) Health(ctx context.Context) error {
	_, err := c.Client.GetLatestIotaSystemState(ctx)
	return err
}

func (c *l1Client) DeployISCContracts(ctx context.Context, signer iotasigner.Signer) (iotago.PackageID, error) {
	iscBytecode := contracts.ISC()
	txnBytes := lo.Must(c.Publish(ctx, iotaclient.PublishRequest{
		Sender:          signer.Address(),
		CompiledModules: iscBytecode.Modules,
		Dependencies:    iscBytecode.Dependencies,
		GasBudget:       iotajsonrpc.NewBigInt(iotaclient.DefaultGasBudget * 10),
	}))
	txnResponse := lo.Must(c.SignAndExecuteTransaction(
		ctx,
		&iotaclient.SignAndExecuteTransactionRequest{
			TxDataBytes: txnBytes.TxBytes,
			Signer:      signer,
			Options: &iotajsonrpc.IotaTransactionBlockResponseOptions{
				ShowEffects:       true,
				ShowObjectChanges: true,
			},
		},
	))
	if !txnResponse.Effects.Data.IsSuccess() {
		return iotago.PackageID{}, errors.New("publish ISC contracts failed")
	}
	packageID := lo.Must(txnResponse.GetPublishedPackageID())
	return *packageID, nil
}

func (c *l1Client) L2() L2Client {
	return iscmoveclient.NewClient(c.Client, c.Config.FaucetURL)
}

func NewL1Client(l1Config L1Config) L1Client {
	return &l1Client{
		iotaclient.NewHTTP(l1Config.APIURL),
		l1Config,
	}
}

func NewLocalnetClient() L1Client {
	return NewL1Client(L1Config{
		APIURL:    iotaconn.LocalnetEndpointURL,
		FaucetURL: iotaconn.LocalnetFaucetURL,
	})
}
