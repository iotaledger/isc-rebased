// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc

import "github.com/iotaledger/wasp/packages/cryptolib"

type ControlAddresses struct {
	StateAddress     *cryptolib.Address
	GoverningAddress *cryptolib.Address
	SinceBlockIndex  uint32
}
