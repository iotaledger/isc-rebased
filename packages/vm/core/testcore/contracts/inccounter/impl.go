package inccounter

import (
	"fmt"
	"math"
	"time"

	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

var Processor = Contract.Processor(nil,
	FuncIncCounter.WithHandler(incCounter),
	FuncIncAndRepeatOnceAfter2s.WithHandler(incCounterAndRepeatOnce),
	FuncIncAndRepeatMany.WithHandler(incCounterAndRepeatMany),
	ViewGetCounter.WithHandler(getCounter),
)

func InitParams(initialValue int64) dict.Dict {
	return dict.Dict{VarCounter: codec.Encode(initialValue)}
}

const (
	incCounterKey = "incCounter"

	VarNumRepeats = "numRepeats"
	VarCounter    = "counter"
	VarName       = "name"
)

func SetInitialState(contractPartition kv.KVStore) {
	contractPartition.Set(VarCounter, codec.Encode[int64](0))
}

func incCounter(ctx isc.Sandbox, incOpt *int64) {
	inc := coreutil.FromOptional(incOpt, 1)
	ctx.Log().Debugf("inccounter.incCounter in %s", ctx.Contract().String())

	val := codec.MustDecode[int64](ctx.State().Get(VarCounter))
	ctx.Log().Infof("incCounter: increasing counter value %d by %d, anchor version: #%d",
		val, inc, ctx.StateAnchor().GetObjectRef().Version)
	tra := "(empty)"
	if ctx.AllowanceAvailable() != nil {
		tra = ctx.AllowanceAvailable().String()
	}
	ctx.Log().Infof("incCounter: allowance available: %s", tra)
	ctx.State().Set(VarCounter, codec.Encode(val+inc))
	eventCounter(ctx, val+inc)
}

func incCounterAndRepeatOnce(ctx isc.Sandbox) {
	ctx.Log().Debugf("inccounter.incCounterAndRepeatOnce")
	state := ctx.State()
	val := codec.MustDecode[int64](state.Get(VarCounter), 0)

	ctx.Log().Debugf(fmt.Sprintf("incCounterAndRepeatOnce: increasing counter value: %d", val))
	state.Set(VarCounter, codec.Encode(val+1))
	eventCounter(ctx, val+1)
	allowance := ctx.AllowanceAvailable()
	ctx.TransferAllowedFunds(ctx.AccountID())
	ctx.Send(isc.RequestParameters{
		TargetAddress: ctx.ChainID().AsAddress(),
		Assets:        isc.NewAssets(allowance.BaseTokens()),
		Metadata: &isc.SendMetadata{
			Message:   isc.NewMessage(ctx.Contract(), FuncIncCounter.Hname()),
			GasBudget: math.MaxUint64,
		},
		Options: isc.SendOptions{
			Timelock: ctx.Timestamp().Add(2 * time.Second),
		},
	})
	ctx.Log().Debugf("incCounterAndRepeatOnce: PostRequestToSelfWithDelay RequestInc 2 sec")
}

func incCounterAndRepeatMany(ctx isc.Sandbox, valOpt, numRepeatsOpt *int64) {
	val := coreutil.FromOptional(valOpt, 0)
	numRepeats := coreutil.FromOptional(valOpt, lo.Must(codec.Decode[int64](ctx.State().Get(VarNumRepeats), 0)))
	ctx.Log().Debugf("inccounter.incCounterAndRepeatMany")

	state := ctx.State()

	state.Set(VarCounter, codec.Encode(val+1))
	eventCounter(ctx, val+1)
	ctx.Log().Debugf("inccounter.incCounterAndRepeatMany: increasing counter value: %d", val)

	if numRepeats == 0 {
		ctx.Log().Debugf("inccounter.incCounterAndRepeatMany: finished chain of requests. counter value: %d", val)
		return
	}

	ctx.Log().Debugf("chain of %d requests ahead", numRepeats)

	state.Set(VarNumRepeats, codec.Encode(numRepeats-1))
	ctx.TransferAllowedFunds(ctx.AccountID())
	ctx.Send(isc.RequestParameters{
		TargetAddress: ctx.ChainID().AsAddress(),
		Assets:        isc.NewAssets(1000),
		Metadata: &isc.SendMetadata{
			Message:   isc.NewMessage(ctx.Contract(), FuncIncAndRepeatMany.Hname()),
			GasBudget: math.MaxUint64,
			Allowance: isc.NewAssets(1000),
		},
		Options: isc.SendOptions{
			Timelock: ctx.Timestamp().Add(2 * time.Second),
		},
	})

	ctx.Log().Debugf("incCounterAndRepeatMany. remaining repeats = %d", numRepeats-1)
}

func getCounter(ctx isc.SandboxView) int64 {
	return lo.Must(codec.Decode[int64](ctx.StateR().Get(VarCounter), 0))
}
