package vmimpl

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"
	"github.com/iotaledger/wasp/packages/kv/subrealm"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/execution"
	"github.com/iotaledger/wasp/packages/vm/sandbox"
)

// Call implements sandbox logic of the call between contracts on-chain
func (reqctx *requestContext) Call(msg isc.Message, allowance *isc.Assets) dict.Dict {
	reqctx.Debugf("Call: targetContract: %s entry point: %s", msg.Target.Contract, msg.Target.EntryPoint)
	return reqctx.callProgram(msg, allowance, reqctx.CurrentContractAgentID())
}

func (reqctx *requestContext) withoutGasBurn(f func()) {
	prev := reqctx.gas.burnEnabled
	reqctx.GasBurnEnable(false)
	f()
	reqctx.GasBurnEnable(prev)
}

func (reqctx *requestContext) callProgram(msg isc.Message, allowance *isc.Assets, caller isc.AgentID) dict.Dict {
	// don't charge gas for finding the contract (otherwise EVM requests may not produce EVM receipt)
	var ep isc.VMProcessorEntryPoint
	reqctx.withoutGasBurn(func() {
		contractRecord := reqctx.GetContractRecord(msg.Target.Contract)
		ep = execution.GetEntryPointByProgHash(reqctx, msg.Target.Contract, msg.Target.EntryPoint, contractRecord.ProgramHash)
	})

	reqctx.pushCallContext(msg.Target.Contract, msg.Params, allowance, caller)
	defer reqctx.popCallContext()

	// distinguishing between two types of entry points. Passing different types of sandboxes
	if ep.IsView() {
		return ep.Call(sandbox.NewSandboxView(reqctx))
	}
	// prevent calling 'init' not from root contract
	if msg.Target.EntryPoint == isc.EntryPointInit && !caller.Equals(isc.NewContractAgentID(reqctx.vm.ChainID(), root.Contract.Hname())) {
		panic(fmt.Errorf("%v: target=(%s, %s)", vm.ErrRepeatingInitCall, msg.Target.Contract, msg.Target.EntryPoint))
	}
	return ep.Call(NewSandbox(reqctx))
}

const traceStack = false

func (reqctx *requestContext) pushCallContext(contract isc.Hname, params dict.Dict, allowance *isc.Assets, caller isc.AgentID) {
	ctx := &callContext{
		caller:   caller,
		contract: contract,
		params: isc.Params{
			Dict:      params,
			KVDecoder: kvdecoder.New(params, reqctx.vm.task.Log),
		},
		allowanceAvailable: allowance.Clone(), // we have to clone it because it will be mutated by TransferAllowedFunds
	}
	if traceStack {
		reqctx.Debugf("+++++++++++ PUSH %d, stack depth = %d caller = %s", contract, len(reqctx.callStack), ctx.caller)
	}
	reqctx.callStack = append(reqctx.callStack, ctx)
}

func (reqctx *requestContext) popCallContext() {
	if traceStack {
		reqctx.Debugf("+++++++++++ POP @ depth %d", len(reqctx.callStack))
	}
	reqctx.callStack[len(reqctx.callStack)-1] = nil // for GC
	reqctx.callStack = reqctx.callStack[:len(reqctx.callStack)-1]
}

func (reqctx *requestContext) getCallContext() *callContext {
	if len(reqctx.callStack) == 0 {
		panic("getCallContext: stack is empty")
	}
	return reqctx.callStack[len(reqctx.callStack)-1]
}

func withContractState(chainState kv.KVStore, c *coreutil.ContractInfo, f func(s kv.KVStore)) {
	f(subrealm.New(chainState, kv.Key(c.Hname().Bytes())))
}

func (vm *vmContext) withAccountsState(chainState kv.KVStore, f func(*accounts.StateWriter)) {
	withContractState(chainState, accounts.Contract, func(contractState kv.KVStore) {
		f(accounts.NewStateWriter(vm.schemaVersion, contractState))
	})
}

func (reqctx *requestContext) callAccounts(f func(*accounts.StateWriter)) {
	reqctx.callCore(accounts.Contract, func(contractState kv.KVStore) {
		f(accounts.NewStateWriter(reqctx.vm.schemaVersion, contractState))
	})
}

func (reqctx *requestContext) callCore(c *coreutil.ContractInfo, f func(s kv.KVStore)) {
	var caller isc.AgentID
	if len(reqctx.callStack) > 0 {
		caller = reqctx.CurrentContractAgentID()
	} else {
		caller = reqctx.req.SenderAccount()
	}
	reqctx.pushCallContext(c.Hname(), nil, nil, caller)
	defer reqctx.popCallContext()

	f(reqctx.contractStateWithGasBurn())
}
