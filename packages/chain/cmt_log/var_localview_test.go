// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/clients/iscmove"
	"github.com/iotaledger/wasp/packages/chain/cmt_log"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/sui-go/sui"
)

func TestVarLocalView(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()
	j := cmt_log.NewVarLocalView(-1, func(ao *iscmove.Anchor) {}, log)
	require.Nil(t, j.Value())
	tipAO, ok, _ := j.AliasOutputConfirmed(
		isc.NewAliasOutputWithID(
			&iotago.AliasOutput{
				StateMetadata: []byte{},
			},
			sui.ObjectID{},
		),
	)
	require.True(t, ok)
	require.NotNil(t, tipAO)
	require.NotNil(t, j.Value())
	require.Equal(t, tipAO, j.Value())
}
