// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainmanager

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc/sui"
)

type inputAliasOutputConfirmed struct {
	aliasOutput *sui.Anchor
}

func NewInputAliasOutputConfirmed(aliasOutput *sui.Anchor) gpa.Input {
	return &inputAliasOutputConfirmed{
		aliasOutput: aliasOutput,
	}
}

func (inp *inputAliasOutputConfirmed) String() string {
	return fmt.Sprintf("{chainMgr.inputAliasOutputConfirmed, %v}", inp.aliasOutput)
}
