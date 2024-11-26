package blocklog_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/evm/evmutil"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/isctest"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/util/bcs"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

func TestReceiptCodec(t *testing.T) {
	bcs.TestCodec(t, blocklog.RequestReceipt{
		Request: isc.NewOffLedgerRequest(
			isctest.RandomChainID(),
			isc.NewMessage(isc.Hn("0"), isc.Hn("0")),
			123,
			gas.LimitsDefault.MaxGasPerRequest,
		).Sign(cryptolib.NewKeyPair()),
		Error: &isc.UnresolvedVMError{
			ErrorCode: blocklog.ErrBlockNotFound.Code(),
			Params:    []isc.VMErrorParam{uint8(1), uint8(2), "string"},
		},
	})
}

func TestReceiptCodecEVM(t *testing.T) {
	unsignedTx := types.NewTransaction(
		0,
		common.Address{},
		util.Big0,
		100,
		util.Big0,
		[]byte{1, 2, 3},
	)
	ethKey := lo.Must(crypto.GenerateKey())
	tx := lo.Must(types.SignTx(unsignedTx, evmutil.Signer(big.NewInt(int64(42))), ethKey))
	bcs.TestCodec(t, blocklog.RequestReceipt{
		Request: lo.Must(isc.NewEVMOffLedgerTxRequest(
			isctest.RandomChainID(),
			tx,
		)),
	})
}
