// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
)

type inputConsensusOutputSkip struct {
	logIndex LogIndex
}

// This message is internal one, but should be sent by other components (e.g. consensus or the chain).
func NewInputConsensusOutputSkip(
	logIndex LogIndex,
) gpa.Input {
	return &inputConsensusOutputSkip{
		logIndex: logIndex,
	}
}

func (inp *inputConsensusOutputSkip) String() string {
	return fmt.Sprintf(
		"{cmtLog.inputConsensusOutputSkip, logIndex=%v}",
		inp.logIndex,
	)
}
