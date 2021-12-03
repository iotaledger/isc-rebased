// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package apilib

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/ed25519"
	"github.com/iotaledger/wasp/client/multiclient"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/registry"
	"github.com/iotaledger/wasp/packages/util"
	"golang.org/x/xerrors"
)

// TODO DeployChain on peering domain, not on committee

type CreateChainParams struct {
	Layer1Client          interface{}
	AllAPIHosts           []string
	AllPeeringHosts       []string
	CommitteeAPIHosts     []string
	CommitteePeeringHosts []string
	N                     uint16
	T                     uint16
	OriginatorPrivateKey  *ed25519.PrivateKey
	Description           string
	Textout               io.Writer
	Prefix                string
}

// DeployChainWithDKG performs all actions needed to deploy the chain
// TODO: [KP] Shouldn't that be in the client packages?
func DeployChainWithDKG(par CreateChainParams) (*iscp.ChainID, iotago.Address, error) {
	if len(par.AllPeeringHosts) > 0 {
		// all committee nodes most also be among allPeers
		if !util.IsSubset(par.CommitteePeeringHosts, par.AllPeeringHosts) {
			return nil, nil, xerrors.Errorf("DeployChainWithDKG: committee nodes must all be among peers")
		}
	}

	dkgInitiatorIndex := uint16(rand.Intn(len(par.CommitteeAPIHosts)))
	stateControllerAddr, err := RunDKG(par.CommitteeAPIHosts, par.CommitteePeeringHosts, par.T, dkgInitiatorIndex)
	if err != nil {
		return nil, nil, err
	}
	chainID, err := DeployChain(par, stateControllerAddr)
	if err != nil {
		return nil, nil, err
	}
	return chainID, stateControllerAddr, nil
}

// DeployChain creates a new chain on specified committee address
// noinspection ALL

func DeployChain(par CreateChainParams, stateControllerAddr iotago.Address) (*iscp.ChainID, error) {
	var err error
	textout := io.Discard
	if par.Textout != nil {
		textout = par.Textout
	}
	originatorAddr := iotago.Ed25519AddressFromPubKey(par.OriginatorPrivateKey.Public().(ed25519.PublicKey))

	fmt.Fprint(textout, par.Prefix)
	fmt.Fprintf(textout, "creating new chain. Owner address: %s. State controller: %s, N = %d, T = %d\n",
		originatorAddr.Bech32(iscp.Bech32Prefix), stateControllerAddr.Bech32(iscp.Bech32Prefix), par.N, par.T)
	fmt.Fprint(textout, par.Prefix)

	chainID, initRequestTx, err := CreateChainOrigin(par.Layer1Client, par.OriginatorPrivateKey, stateControllerAddr, par.Description)
	fmt.Fprint(textout, par.Prefix)
	if err != nil {
		fmt.Fprintf(textout, "creating chain origin and init transaction.. FAILED: %v\n", err)
		return nil, xerrors.Errorf("DeployChain: %w", err)
	}
	txID, err := initRequestTx.ID()
	if err != nil {
		fmt.Fprintf(textout, "creating chain origin and init transaction.. FAILED: %v\n", err)
		return nil, xerrors.Errorf("DeployChain: %w", err)
	}
	fmt.Fprintf(textout, "creating chain origin and init transaction %s.. OK\n", txID)
	fmt.Fprint(textout, "sending committee record to nodes.. OK\n")

	err = ActivateChainOnAccessNodes(par.AllAPIHosts, par.AllPeeringHosts, chainID)
	fmt.Fprint(textout, par.Prefix)
	if err != nil {
		fmt.Fprintf(textout, "activating chain %s.. FAILED: %v\n", chainID.Bech32(iscp.Bech32Prefix), err)
		return nil, xerrors.Errorf("DeployChain: %w", err)
	}
	fmt.Fprintf(textout, "activating chain %s.. OK.\n", chainID.Bech32(iscp.Bech32Prefix))

	peers := multiclient.New(par.CommitteeAPIHosts)

	// ---------- wait until the request is processed at least in all committee nodes
	if err = peers.WaitUntilAllRequestsProcessed(chainID, initRequestTx, 30*time.Second); err != nil {
		fmt.Fprintf(textout, "waiting root init request transaction.. FAILED: %v\n", err)
		return nil, xerrors.Errorf("DeployChain: %w", err)
	}

	fmt.Fprint(textout, par.Prefix)
	fmt.Fprintf(textout, "chain has been created successfully on the Tangle. ChainID: %s, State address: %s, N = %d, T = %d\n",
		chainID.String(), stateControllerAddr.Bech32(iscp.Bech32Prefix), par.N, par.T)

	return chainID, err
}

// CreateChainOrigin creates and confirms origin transaction of the chain and init request transaction to initialize state of it
func CreateChainOrigin(Layer1Client interface{}, originator *ed25519.PrivateKey, stateController iotago.Address, dscr string) (*iscp.ChainID, *iotago.Transaction, error) {
	panic("TODO implement")
	// originatorAddr := iotago.Ed25519AddressFromPubKey(originator.Public().(ed25519.PublicKey))
	// // ----------- request owner address' outputs from the ledger
	// allOuts, err := Layer1Client.GetConfirmedOutputs(originatorAddr)
	// if err != nil {
	// 	return nil, nil, xerrors.Errorf("CreateChainOrigin: %w", err)
	// }

	// // ----------- create origin transaction
	// originTx, chainID, err := transaction.NewChainOriginTransaction(
	// 	originator,
	// 	stateController,
	// 	nil,
	// 	time.Now(),
	// 	allOuts,
	// )
	// if err != nil {
	// 	return nil, nil, xerrors.Errorf("CreateChainOrigin: %w", err)
	// }

	// // ------------- post origin transaction and wait for confirmation
	// err = Layer1Client.PostAndWaitForConfirmation(originTx)
	// if err != nil {
	// 	return nil, nil, xerrors.Errorf("CreateChainOrigin: %w", err)
	// }

	// allOuts, err = Layer1Client.GetConfirmedOutputs(originatorAddr)
	// if err != nil {
	// 	return nil, nil, xerrors.Errorf("CreateChainOrigin: %w", err)
	// }

	// // NOTE: whoever send first init request, is an owner of the chain
	// // create root init transaction
	// reqTx, err := transaction.NewRootInitRequestTransaction(
	// 	originator,
	// 	chainID,
	// 	dscr,
	// 	time.Now(),
	// 	allOuts...,
	// )
	// if err != nil {
	// 	return nil, nil, xerrors.Errorf("CreateChainOrigin: %w", err)
	// }

	// // ---------- post root init request transaction and wait for confirmation
	// err = Layer1Client.PostAndWaitForConfirmation(reqTx)
	// if err != nil {
	// 	return nil, nil, xerrors.Errorf("CreateChainOrigin: %w", err)
	// }

	// return chainID, reqTx, nil
}

// ActivateChainOnAccessNodes puts chain records into nodes and activates its
// TODO needs refactoring and optimization
func ActivateChainOnAccessNodes(apiHosts, peers []string, chainID *iscp.ChainID) error {
	nodes := multiclient.New(apiHosts)
	// ------------ put chain records to hosts
	err := nodes.PutChainRecord(&registry.ChainRecord{
		ChainID: chainID,
		Peers:   peers,
	})
	if err != nil {
		return xerrors.Errorf("ActivateChainOnAccessNodes: %w", err)
	}
	// ------------- activate chain
	err = nodes.ActivateChain(chainID)
	if err != nil {
		return xerrors.Errorf("ActivateChainOnAccessNodes: %w", err)
	}
	return nil
}
