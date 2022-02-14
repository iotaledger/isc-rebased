package viewcontext

// import (
// 	"fmt"
// 	"runtime/debug"

// 	"github.com/iotaledger/hive.go/logger"
// 	"github.com/iotaledger/wasp/packages/chain"
// 	"github.com/iotaledger/wasp/packages/hashing"
// 	"github.com/iotaledger/wasp/packages/iscp"
// 	"github.com/iotaledger/wasp/packages/kv"
// 	"github.com/iotaledger/wasp/packages/kv/dict"
// 	"github.com/iotaledger/wasp/packages/kv/subrealm"
// 	"github.com/iotaledger/wasp/packages/state"
// 	"github.com/iotaledger/wasp/packages/vm/core/blob"
// 	"github.com/iotaledger/wasp/packages/vm/core/root"
// 	"github.com/iotaledger/wasp/packages/vm/processors"
// 	"golang.org/x/xerrors"
// )

// type Viewcontext struct {
// 	processors  *processors.Cache
// 	stateReader state.OptimisticStateReader
// 	chainID     *iscp.ChainID
// 	log         *logger.Logger
// }

// func NewFromChain(ch chain.ChainCore) *Viewcontext {
// 	// optimisticStateAccess := state.WrapMustOptimisticVirtualStateAccess(task.VirtualStateAccess, task.SolidStateBaseline)

// 	// ret := &VMContext{
// 	// 	task:                 task,
// 	// 	virtualState:         optimisticStateAccess,
// 	// 	blockContext:         make(map[iscp.Hname]*blockContext),
// 	// 	blockContextCloseSeq: make([]iscp.Hname, 0),
// 	// 	entropy:              task.Entropy,
// 	// 	callStack:            make([]*callContext, 0),
// 	// }

// 	// VMContext{}
// 	task := &vm.VMTask{
// 		Processors: ch.Processors(),
// 		SolidStateBaseline: ch.GlobalSync()
// 	}

// 	return New(ch.ID(), ch.GetStateReader(), ch.Processors(), ch.Log().Named("view"))
// }

// func New(chainID *iscp.ChainID, stateReader state.OptimisticStateReader, proc *processors.Cache, log *logger.Logger) *Viewcontext {
// 	return &Viewcontext{
// 		processors:  proc,
// 		stateReader: stateReader,
// 		chainID:     chainID,
// 		log:         log,
// 	}
// }

// // CallView in viewcontext implements own panic catcher.
// func (v *Viewcontext) CallView(contractHname, epCode iscp.Hname, params dict.Dict) (dict.Dict, error) {
// 	var ret dict.Dict
// 	var err error
// 	func() {
// 		defer func() {
// 			r := recover()
// 			if r == nil {
// 				return
// 			}
// 			ret = nil
// 			switch err1 := r.(type) {
// 			case *kv.DBError:
// 				v.log.Panicf("DB error: %v", err1)
// 			case error:
// 				err = err1
// 			default:
// 				err = xerrors.Errorf("viewcontext: panic in VM: %v", err1)
// 			}
// 			v.log.Debugf("CallView: %v", err)
// 			v.log.Debugf(string(debug.Stack()))
// 		}()
// 		ret = v.callView(contractHname, epCode, params)
// 	}()
// 	return ret, err
// }

// func (v *Viewcontext) callView(contractHname, epCode iscp.Hname, params dict.Dict) dict.Dict {
// 	var err error
// 	contractRecord := root.FindContract(contractStateSubpartition(v.stateReader.KVStoreReader(), root.Contract.Hname()), contractHname)
// 	if contractRecord == nil {
// 		panic(fmt.Sprintf("contract not found %s", contractHname))
// 	}
// 	proc, err := v.processors.GetOrCreateProcessor(contractRecord, func(programHash hashing.HashValue) (string, []byte, error) {
// 		if vmtype, ok := v.processors.Config.GetNativeProcessorType(programHash); ok {
// 			return vmtype, nil, nil
// 		}
// 		return blob.LocateProgram(contractStateSubpartition(v.stateReader.KVStoreReader(), blob.Contract.Hname()), programHash)
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	ep, ok := proc.GetEntryPoint(epCode)
// 	if !ok {
// 		panic(fmt.Sprintf("trying to call contract '%s': can't find entry point '%s'", proc.GetDescription(), epCode.String()))
// 	}

// 	if !ep.IsView() {
// 		panic("only view entry point can be called in this context")
// 	}
// 	return ep.Call(newSandboxView(v, contractHname, params))
// }

// func contractStateSubpartition(stateKvReader kv.KVStoreReader, contractHname iscp.Hname) kv.KVStoreReader {
// 	return subrealm.NewReadOnly(stateKvReader, kv.Key(contractHname.Bytes()))
// }

// func (v *Viewcontext) Infof(format string, params ...interface{}) {
// 	v.log.Infof(format, params...)
// }

// func (v *Viewcontext) Debugf(format string, params ...interface{}) {
// 	v.log.Debugf(format, params...)
// }

// func (v *Viewcontext) Panicf(format string, params ...interface{}) {
// 	v.log.Panicf(format, params...)
}
