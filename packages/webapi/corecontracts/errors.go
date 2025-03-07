package corecontracts

import (
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/errors"
	"github.com/iotaledger/wasp/packages/webapi/common"
)

func ErrorMessageFormat(ch chain.Chain, contractID isc.Hname, errorID uint16, blockIndexOrTrieRoot string) (string, error) {
	errorCode := isc.NewVMErrorCode(contractID, errorID)
	ret, err := common.CallView(ch, errors.ViewGetErrorMessageFormat.Message(errorCode), blockIndexOrTrieRoot)
	if err != nil {
		return "", err
	}
	return errors.ViewGetErrorMessageFormat.DecodeOutput(ret)
}
