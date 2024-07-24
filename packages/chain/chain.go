// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"context"
	"fmt"
	"time"

	"github.com/iotaledger/hive.go/logger"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/peering"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/state/indexedstore"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/processors"
)

type ChainNodeConn interface {
	// AttachChain begins processing of the anchor and requests of the given
	// chain, and stops when the ctx is canceled.
	// Upon Attach, the latest anchor is returned FIRST.
	// Upon receiving new Anchor versions from L1, they are returned in the
	// same order.
	// onChainConnect is called immediately.
	// onChainDisconnect is called after the ctx is canceled and the chain is
	// detached.
	AttachChain(
		ctx context.Context,
		chainID isc.ChainID,
		recvRequest RequestHandler,
		recvAnchor AnchorHandler,
		onChainConnect func(),
		onChainDisconnect func(),
	)
	// PublishTX posts the PTB asynchronously and calls the callback when it is
	// confirmed or an error is detected, or the ctx is canceled.
	PublishTX(
		ctx context.Context,
		chainID isc.ChainID,
		tx SignedTx,
		callback TxPostHandler,
	) error
	// RefreshOnLedgerRequests synchronously fetches all owned requests by the
	// previously attached chain, and calls recvRequest for each one.
	RefreshOnLedgerRequests(ctx context.Context, chainID isc.ChainID)
}

type NodeConnection interface {
	ChainNodeConn
	// Run starts the connection to the L1 node, and blocks until the context
	// is canceled.
	Run(ctx context.Context) error
	// WaitUntilInitiallySynced blocks until the connection is established.
	WaitUntilInitiallySynced(context.Context) error
	GetL1Params() *parameters.L1Params
	GetL1ProtocolParams() *iotago.ProtocolParameters
}

type StateFreshness byte

const (
	ActiveOrCommittedState StateFreshness = iota // ActiveState, if exist; Confirmed state otherwise.
	ActiveState                                  // The state the chain build next TX on, can be ahead of ConfirmedState.
	ConfirmedState                               // The state confirmed on L1.
)

func (sf StateFreshness) String() string {
	switch sf {
	case ActiveOrCommittedState:
		return "ActiveOrCommittedState"
	case ActiveState:
		return "ActiveState"
	case ConfirmedState:
		return "ConfirmedState"
	default:
		return fmt.Sprintf("StateFreshness=%v", int(sf))
	}
}

type ChainCore interface {
	ID() isc.ChainID
	// Returns the current latest confirmed alias output and the active one.
	// The active AO can be ahead of the confirmed one by several blocks.
	// Both values can be nil, if the node haven't received an output from
	// L1 yet (after a restart or a chain activation).
	LatestAliasOutput(freshness StateFreshness) (*isc.AliasOutputWithID, error)
	LatestState(freshness StateFreshness) (state.State, error)
	GetCommitteeInfo() *CommitteeInfo // TODO: Review, maybe we can reorganize the CommitteeInfo structure.
	Store() indexedstore.IndexedStore // Use LatestState whenever possible. That will work faster.
	Processors() *processors.Cache
	GetChainNodes() []peering.PeerStatusProvider     // CommitteeNodes + AccessNodes
	GetCandidateNodes() []*governance.AccessNodeInfo // All the current candidates.
	Log() *logger.Logger
}

type ConsensusPipeMetrics interface { // TODO: Review it.
	GetEventStateTransitionMsgPipeSize() int
	GetEventPeerLogIndexMsgPipeSize() int
	GetEventACSMsgPipeSize() int
	GetEventVMResultMsgPipeSize() int
	GetEventTimerMsgPipeSize() int
}

type ConsensusWorkflowStatus interface { // TODO: Review it.
	IsStateReceived() bool
	IsBatchProposalSent() bool
	IsConsensusBatchKnown() bool
	IsVMStarted() bool
	IsVMResultSigned() bool
	IsTransactionFinalized() bool
	IsTransactionPosted() bool
	IsTransactionSeen() bool
	IsInProgress() bool
	GetBatchProposalSentTime() time.Time
	GetConsensusBatchKnownTime() time.Time
	GetVMStartedTime() time.Time
	GetVMResultSignedTime() time.Time
	GetTransactionFinalizedTime() time.Time
	GetTransactionPostedTime() time.Time
	GetTransactionSeenTime() time.Time
	GetCompletedTime() time.Time
	GetCurrentStateIndex() uint32
}
