// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

#![allow(dead_code)]

use wasmlib::*;

pub const SC_NAME        : &str = "testcore";
pub const SC_DESCRIPTION : &str = "Wasm equivalent of built-in TestCore contract";
pub const HSC_NAME       : ScHname = ScHname(0x370d33ad);

pub const PARAM_ADDRESS          : &str = "address";
pub const PARAM_AGENT_ID         : &str = "agentID";
pub const PARAM_CALLER           : &str = "caller";
pub const PARAM_CHAIN_ID         : &str = "chainID";
pub const PARAM_CHAIN_OWNER_ID   : &str = "chainOwnerID";
pub const PARAM_CONTRACT_CREATOR : &str = "contractCreator";
pub const PARAM_CONTRACT_ID      : &str = "contractID";
pub const PARAM_COUNTER          : &str = "counter";
pub const PARAM_FAIL             : &str = "initFailParam";
pub const PARAM_GAS_BUDGET       : &str = "gasBudget";
pub const PARAM_HASH             : &str = "Hash";
pub const PARAM_HNAME            : &str = "Hname";
pub const PARAM_HNAME_CONTRACT   : &str = "hnameContract";
pub const PARAM_HNAME_EP         : &str = "hnameEP";
pub const PARAM_HNAME_ZERO       : &str = "Hname-0";
pub const PARAM_INT64            : &str = "int64";
pub const PARAM_INT64_ZERO       : &str = "int64-0";
pub const PARAM_INT_VALUE        : &str = "intParamValue";
pub const PARAM_IOTAS_WITHDRAWAL : &str = "iotasWithdrawal";
pub const PARAM_N                : &str = "n";
pub const PARAM_NAME             : &str = "intParamName";
pub const PARAM_PROG_HASH        : &str = "progHash";
pub const PARAM_STRING           : &str = "string";
pub const PARAM_STRING_ZERO      : &str = "string-0";
pub const PARAM_VAR_NAME         : &str = "varName";

pub const RESULT_CHAIN_OWNER_ID : &str = "chainOwnerID";
pub const RESULT_COUNTER        : &str = "counter";
pub const RESULT_N              : &str = "n";
pub const RESULT_SANDBOX_CALL   : &str = "sandboxCall";
pub const RESULT_VALUES         : &str = "this";
pub const RESULT_VARS           : &str = "this";

pub const STATE_COUNTER : &str = "counter";
pub const STATE_INTS    : &str = "ints";
pub const STATE_STRINGS : &str = "strings";

pub const FUNC_CALL_ON_CHAIN                     : &str = "callOnChain";
pub const FUNC_CHECK_CONTEXT_FROM_FULL_EP        : &str = "checkContextFromFullEP";
pub const FUNC_CLAIM_ALLOWANCE                   : &str = "claimAllowance";
pub const FUNC_DO_NOTHING                        : &str = "doNothing";
pub const FUNC_ESTIMATE_MIN_DUST                 : &str = "estimateMinDust";
pub const FUNC_INC_COUNTER                       : &str = "incCounter";
pub const FUNC_INFINITE_LOOP                     : &str = "infiniteLoop";
pub const FUNC_INIT                              : &str = "init";
pub const FUNC_PASS_TYPES_FULL                   : &str = "passTypesFull";
pub const FUNC_PING_ALLOWANCE_BACK               : &str = "pingAllowanceBack";
pub const FUNC_RUN_RECURSION                     : &str = "runRecursion";
pub const FUNC_SEND_LARGE_REQUEST                : &str = "sendLargeRequest";
pub const FUNC_SEND_NF_TS_BACK                   : &str = "sendNFTsBack";
pub const FUNC_SEND_TO_ADDRESS                   : &str = "sendToAddress";
pub const FUNC_SET_INT                           : &str = "setInt";
pub const FUNC_SPAWN                             : &str = "spawn";
pub const FUNC_SPLIT_FUNDS                       : &str = "splitFunds";
pub const FUNC_SPLIT_FUNDS_NATIVE_TOKENS         : &str = "splitFundsNativeTokens";
pub const FUNC_TEST_BLOCK_CONTEXT1               : &str = "testBlockContext1";
pub const FUNC_TEST_BLOCK_CONTEXT2               : &str = "testBlockContext2";
pub const FUNC_TEST_CALL_PANIC_FULL_EP           : &str = "testCallPanicFullEP";
pub const FUNC_TEST_CALL_PANIC_VIEW_EP_FROM_FULL : &str = "testCallPanicViewEPFromFull";
pub const FUNC_TEST_CHAIN_OWNER_ID_FULL          : &str = "testChainOwnerIDFull";
pub const FUNC_TEST_EVENT_LOG_DEPLOY             : &str = "testEventLogDeploy";
pub const FUNC_TEST_EVENT_LOG_EVENT_DATA         : &str = "testEventLogEventData";
pub const FUNC_TEST_EVENT_LOG_GENERIC_DATA       : &str = "testEventLogGenericData";
pub const FUNC_TEST_PANIC_FULL_EP                : &str = "testPanicFullEP";
pub const FUNC_WITHDRAW_FROM_CHAIN               : &str = "withdrawFromChain";
pub const VIEW_CHECK_CONTEXT_FROM_VIEW_EP        : &str = "checkContextFromViewEP";
pub const VIEW_FIBONACCI                         : &str = "fibonacci";
pub const VIEW_FIBONACCI_INDIRECT                : &str = "fibonacciIndirect";
pub const VIEW_GET_COUNTER                       : &str = "getCounter";
pub const VIEW_GET_INT                           : &str = "getInt";
pub const VIEW_GET_STRING_VALUE                  : &str = "getStringValue";
pub const VIEW_INFINITE_LOOP_VIEW                : &str = "infiniteLoopView";
pub const VIEW_JUST_VIEW                         : &str = "justView";
pub const VIEW_PASS_TYPES_VIEW                   : &str = "passTypesView";
pub const VIEW_TEST_CALL_PANIC_VIEW_EP_FROM_VIEW : &str = "testCallPanicViewEPFromView";
pub const VIEW_TEST_CHAIN_OWNER_ID_VIEW          : &str = "testChainOwnerIDView";
pub const VIEW_TEST_PANIC_VIEW_EP                : &str = "testPanicViewEP";
pub const VIEW_TEST_SANDBOX_CALL                 : &str = "testSandboxCall";

pub const HFUNC_CALL_ON_CHAIN                     : ScHname = ScHname(0x95a3d123);
pub const HFUNC_CHECK_CONTEXT_FROM_FULL_EP        : ScHname = ScHname(0xa56c24ba);
pub const HFUNC_CLAIM_ALLOWANCE                   : ScHname = ScHname(0x40bec0e6);
pub const HFUNC_DO_NOTHING                        : ScHname = ScHname(0xdda4a6de);
pub const HFUNC_ESTIMATE_MIN_DUST                 : ScHname = ScHname(0xe700e7db);
pub const HFUNC_INC_COUNTER                       : ScHname = ScHname(0x7b287419);
pub const HFUNC_INFINITE_LOOP                     : ScHname = ScHname(0xf571430a);
pub const HFUNC_INIT                              : ScHname = ScHname(0x1f44d644);
pub const HFUNC_PASS_TYPES_FULL                   : ScHname = ScHname(0x733ea0ea);
pub const HFUNC_PING_ALLOWANCE_BACK               : ScHname = ScHname(0x66f43c0b);
pub const HFUNC_RUN_RECURSION                     : ScHname = ScHname(0x833425fd);
pub const HFUNC_SEND_LARGE_REQUEST                : ScHname = ScHname(0xfdaaca3c);
pub const HFUNC_SEND_NF_TS_BACK                   : ScHname = ScHname(0x8f6ef428);
pub const HFUNC_SEND_TO_ADDRESS                   : ScHname = ScHname(0x63ce4634);
pub const HFUNC_SET_INT                           : ScHname = ScHname(0x62056f74);
pub const HFUNC_SPAWN                             : ScHname = ScHname(0xec929d12);
pub const HFUNC_SPLIT_FUNDS                       : ScHname = ScHname(0xc7ea86c9);
pub const HFUNC_SPLIT_FUNDS_NATIVE_TOKENS         : ScHname = ScHname(0x16532a28);
pub const HFUNC_TEST_BLOCK_CONTEXT1               : ScHname = ScHname(0x796d4136);
pub const HFUNC_TEST_BLOCK_CONTEXT2               : ScHname = ScHname(0x758b0452);
pub const HFUNC_TEST_CALL_PANIC_FULL_EP           : ScHname = ScHname(0x4c878834);
pub const HFUNC_TEST_CALL_PANIC_VIEW_EP_FROM_FULL : ScHname = ScHname(0xfd7e8c1d);
pub const HFUNC_TEST_CHAIN_OWNER_ID_FULL          : ScHname = ScHname(0x2aff1167);
pub const HFUNC_TEST_EVENT_LOG_DEPLOY             : ScHname = ScHname(0x96ff760a);
pub const HFUNC_TEST_EVENT_LOG_EVENT_DATA         : ScHname = ScHname(0x0efcf939);
pub const HFUNC_TEST_EVENT_LOG_GENERIC_DATA       : ScHname = ScHname(0x6a16629d);
pub const HFUNC_TEST_PANIC_FULL_EP                : ScHname = ScHname(0x24fdef07);
pub const HFUNC_WITHDRAW_FROM_CHAIN               : ScHname = ScHname(0x405c0b0a);
pub const HVIEW_CHECK_CONTEXT_FROM_VIEW_EP        : ScHname = ScHname(0x88ff0167);
pub const HVIEW_FIBONACCI                         : ScHname = ScHname(0x7940873c);
pub const HVIEW_FIBONACCI_INDIRECT                : ScHname = ScHname(0x6dd98513);
pub const HVIEW_GET_COUNTER                       : ScHname = ScHname(0xb423e607);
pub const HVIEW_GET_INT                           : ScHname = ScHname(0x1887e5ef);
pub const HVIEW_GET_STRING_VALUE                  : ScHname = ScHname(0xcf0a4d32);
pub const HVIEW_INFINITE_LOOP_VIEW                : ScHname = ScHname(0x1a383295);
pub const HVIEW_JUST_VIEW                         : ScHname = ScHname(0x33b8972e);
pub const HVIEW_PASS_TYPES_VIEW                   : ScHname = ScHname(0x1a5b87ea);
pub const HVIEW_TEST_CALL_PANIC_VIEW_EP_FROM_VIEW : ScHname = ScHname(0x91b10c99);
pub const HVIEW_TEST_CHAIN_OWNER_ID_VIEW          : ScHname = ScHname(0x26586c33);
pub const HVIEW_TEST_PANIC_VIEW_EP                : ScHname = ScHname(0x22bc4d72);
pub const HVIEW_TEST_SANDBOX_CALL                 : ScHname = ScHname(0x42d72b63);
