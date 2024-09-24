// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util/bcs"
)

func TestMsgBLSPartialSigSerialization(t *testing.T) {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	require.NoError(t, err)
	msg := &msgBLSPartialSig{
		gpa.BasicMessage{},
		nil,
		b,
	}

	bcs.TestCodec(t, msg)
}
