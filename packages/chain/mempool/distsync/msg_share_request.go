// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package distsync

import (
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
)

type msgShareRequest struct {
	gpa.BasicMessage
	ttl     byte        `bcs:""`
	request isc.Request `bcs:""`
}

var _ gpa.Message = new(msgShareRequest)

func newMsgShareRequest(request isc.Request, ttl byte, recipient gpa.NodeID) gpa.Message {
	return &msgShareRequest{
		BasicMessage: gpa.NewBasicMessage(recipient),
		request:      request,
		ttl:          ttl,
	}
}
