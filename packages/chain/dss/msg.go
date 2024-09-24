// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dss

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
)

const (
	msgTypePartialSig gpa.MessageType = iota
	msgTypeWrapped
)

func (d *dssImpl) msgWrapperFunc(subsystem byte, index int) (gpa.GPA, error) {
	if subsystem == subsystemDKG {
		if index != 0 {
			return nil, fmt.Errorf("unexpected DKG index: %v", index)
		}
		return d.dkg, nil
	}
	return nil, fmt.Errorf("unexpected subsystem: %v", subsystem)
}

func (d *dssImpl) MarshalMessage(msg gpa.Message) ([]byte, error) {
	switch msg := msg.(type) {
	case *msgPartialSig:
		return gpa.MarshalMessage(msgTypePartialSig, msg)
	default:
		return gpa.MarshalWrappedMessage(msgTypeWrapped, msg, d.msgWrapper)
	}
}

func (d *dssImpl) UnmarshalMessage(data []byte) (gpa.Message, error) {
	return gpa.UnmarshalMessage(data, gpa.Mapper{
		msgTypePartialSig: func() gpa.Message { return &msgPartialSig{suite: d.suite} },
	}, gpa.Fallback{
		msgTypeWrapped: d.msgWrapper.UnmarshalMessage,
	})
}
