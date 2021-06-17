// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
//////// DO NOT CHANGE THIS FILE! ////////
// Change the json schema instead

package test

import "github.com/iotaledger/wasp/packages/coretypes"

const (
	ScName  = "inccounter"
	HScName = coretypes.Hname(0xaf2438e9)
)

const (
	ParamCounter    = "counter"
	ParamNumRepeats = "numRepeats"
)

const ResultCounter = "counter"

const (
	VarCounter    = "counter"
	VarNumRepeats = "numRepeats"
)

const (
	FuncCallIncrement          = "callIncrement"
	FuncCallIncrementRecurse5x = "callIncrementRecurse5x"
	FuncIncrement              = "increment"
	FuncInit                   = "init"
	FuncLocalStateInternalCall = "localStateInternalCall"
	FuncLocalStatePost         = "localStatePost"
	FuncLocalStateSandboxCall  = "localStateSandboxCall"
	FuncLoop                   = "loop"
	FuncPostIncrement          = "postIncrement"
	FuncRepeatMany             = "repeatMany"
	FuncTestLeb128             = "testLeb128"
	FuncWhenMustIncrement      = "whenMustIncrement"
	ViewGetCounter             = "getCounter"
)

const (
	HFuncCallIncrement          = coretypes.Hname(0xeb5dcacd)
	HFuncCallIncrementRecurse5x = coretypes.Hname(0x8749fbff)
	HFuncIncrement              = coretypes.Hname(0xd351bd12)
	HFuncInit                   = coretypes.Hname(0x1f44d644)
	HFuncLocalStateInternalCall = coretypes.Hname(0xecfc5d33)
	HFuncLocalStatePost         = coretypes.Hname(0x3fd54d13)
	HFuncLocalStateSandboxCall  = coretypes.Hname(0x7bd22c53)
	HFuncLoop                   = coretypes.Hname(0xa9a20fa9)
	HFuncPostIncrement          = coretypes.Hname(0x81c772f5)
	HFuncRepeatMany             = coretypes.Hname(0x4ff450d3)
	HFuncTestLeb128             = coretypes.Hname(0xd8364cb9)
	HFuncWhenMustIncrement      = coretypes.Hname(0xb4c3e7a6)
	HViewGetCounter             = coretypes.Hname(0xb423e607)
)
