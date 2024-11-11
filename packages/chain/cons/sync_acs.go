// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"fmt"
	"strings"
	"time"

	"github.com/iotaledger/wasp/clients/iota-go/iotago"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/gpa/acs"
	"github.com/iotaledger/wasp/packages/isc"
)

type SyncACS interface {
	StateProposalReceived(proposedBaseAliasOutput *isc.StateAnchor) gpa.OutMessages
	MempoolRequestsReceived(requestRefs []*isc.RequestRef) gpa.OutMessages
	DSSIndexProposalReceived(dssIndexProposal []int) gpa.OutMessages
	TimeDataReceived(timeData time.Time) gpa.OutMessages
	GasInfoReceived(gasCoins []*iotago.ObjectRef, gasPrice uint64) gpa.OutMessages
	ACSOutputReceived(output gpa.Output) gpa.OutMessages
	String() string
}

// > UPON Reception of responses from Mempool, StateMgr and DSS NonceIndexes:
// >     Produce a batch proposal.
// >     Start the ACS.
type syncACSImpl struct {
	BaseAliasOutput  *isc.StateAnchor
	RequestRefs      []*isc.RequestRef
	DSSIndexProposal []int
	TimeData         time.Time
	gasCoins         []*iotago.ObjectRef
	gasPrice         uint64

	inputsReady   bool
	inputsReadyCB func(
		baseAliasOutput *isc.StateAnchor,
		requestRefs []*isc.RequestRef,
		dssIndexProposal []int,
		timeData time.Time,
		gasCoins []*iotago.ObjectRef,
		gasPrice uint64,
	) gpa.OutMessages

	outputReady   bool
	outputReadyCB func(output map[gpa.NodeID][]byte) gpa.OutMessages

	terminated   bool
	terminatedCB func()
}

func NewSyncACS(
	inputsReadyCB func(
		baseAliasOutput *isc.StateAnchor,
		requestRefs []*isc.RequestRef,
		dssIndexProposal []int,
		timeData time.Time,
		gasCoins []*iotago.ObjectRef,
		gasPrice uint64,
	) gpa.OutMessages,
	outputReadyCB func(output map[gpa.NodeID][]byte) gpa.OutMessages,
	terminatedCB func(),
) SyncACS {
	return &syncACSImpl{
		inputsReadyCB: inputsReadyCB,
		outputReadyCB: outputReadyCB,
		terminatedCB:  terminatedCB,
	}
}

func (sub *syncACSImpl) StateProposalReceived(proposedBaseAliasOutput *isc.StateAnchor) gpa.OutMessages {
	if sub.BaseAliasOutput != nil {
		return nil
	}
	sub.BaseAliasOutput = proposedBaseAliasOutput
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) MempoolRequestsReceived(requestRefs []*isc.RequestRef) gpa.OutMessages {
	if sub.RequestRefs != nil {
		return nil
	}
	sub.RequestRefs = requestRefs
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) DSSIndexProposalReceived(dssIndexProposal []int) gpa.OutMessages {
	if sub.DSSIndexProposal != nil {
		return nil
	}
	sub.DSSIndexProposal = dssIndexProposal
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) TimeDataReceived(timeData time.Time) gpa.OutMessages {
	if timeData.After(sub.TimeData) {
		sub.TimeData = timeData
		return sub.tryCompleteInput()
	}
	return nil
}

func (sub *syncACSImpl) GasInfoReceived(gasCoins []*iotago.ObjectRef, gasPrice uint64) gpa.OutMessages {
	if gasCoins != nil {
		sub.gasCoins = gasCoins
	}
	if gasPrice != 0 {
		sub.gasPrice = gasPrice
	}
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) tryCompleteInput() gpa.OutMessages {
	if sub.inputsReady || sub.BaseAliasOutput == nil || sub.RequestRefs == nil || sub.DSSIndexProposal == nil || sub.TimeData.IsZero() || sub.gasCoins == nil || sub.gasPrice == 0 {
		return nil
	}
	sub.inputsReady = true
	return sub.inputsReadyCB(sub.BaseAliasOutput, sub.RequestRefs, sub.DSSIndexProposal, sub.TimeData, sub.gasCoins, sub.gasPrice)
}

func (sub *syncACSImpl) ACSOutputReceived(output gpa.Output) gpa.OutMessages {
	if output == nil {
		return nil
	}
	acsOutput, ok := output.(*acs.Output)
	if !ok {
		panic(fmt.Errorf("acs returned unexpected output: %v", output))
	}
	if !sub.terminated && acsOutput.Terminated {
		sub.terminated = true
		sub.terminatedCB()
	}
	if sub.outputReady {
		return nil
	}
	sub.outputReady = true
	return sub.outputReadyCB(acsOutput.Values)
}

// Try to provide useful human-readable compact status.
func (sub *syncACSImpl) String() string {
	str := "ACS"
	if sub.outputReady {
		str += statusStrOK
	} else if sub.inputsReady {
		str += "/WAIT[ACS to complete]"
	} else {
		wait := []string{}
		if sub.BaseAliasOutput == nil {
			wait = append(wait, "BaseAliasOutput")
		}
		if sub.RequestRefs == nil {
			wait = append(wait, "RequestRefs")
		}
		if sub.DSSIndexProposal == nil {
			wait = append(wait, "DSSIndexProposal")
		}
		if sub.TimeData.IsZero() {
			wait = append(wait, "TimeData")
		}
		if sub.gasCoins == nil {
			wait = append(wait, "GasCoins")
		}
		if sub.gasPrice == 0 {
			wait = append(wait, "GasPrice")
		}
		str += fmt.Sprintf("/WAIT[%v]", strings.Join(wait, ","))
	}
	return str
}
