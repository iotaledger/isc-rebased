// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"path"

	"go.uber.org/dig"

	"github.com/iotaledger/hive.go/app"
	hivep2p "github.com/iotaledger/hive.go/crypto/p2p"
	"github.com/iotaledger/hive.go/runtime/ioutils"
	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/chain/cmt_log"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/registry"
)

func init() {
	Component = &app.Component{
		Name:    "Registry",
		Params:  params,
		Provide: provide,
	}
}

var Component *app.Component

func provide(c *dig.Container) error {
	if err := c.Provide(func() registry.NodeIdentityProvider {
		return nodeIdentityRegistry()
	}); err != nil {
		Component.LogPanic(err)
	}

	if err := c.Provide(func() registry.ChainRecordRegistryProvider {
		chainRecordRegistryProvider, err := registry.NewChainRecordRegistryImpl(ParamsRegistries.Chains.FilePath)
		if err != nil {
			Component.LogPanic(err)
		}
		return chainRecordRegistryProvider
	}); err != nil {
		Component.LogPanic(err)
	}

	type consensusRegistryDeps struct {
		dig.In

		NodeConnection chain.NodeConnection
	}

	if err := c.Provide(func(deps consensusRegistryDeps) cmt_log.ConsensusStateRegistry {
		consensusStateRegistry, err := registry.NewConsensusStateRegistry(ParamsRegistries.ConsensusState.Path, "IOTA")
		if err != nil {
			Component.LogPanic(err)
		}
		return consensusStateRegistry
	}); err != nil {
		Component.LogPanic(err)
	}

	type dkSharesRegistryDeps struct {
		dig.In

		NodeIdentityProvider registry.NodeIdentityProvider
		NodeConnection       chain.NodeConnection
	}

	if err := c.Provide(func(deps dkSharesRegistryDeps) registry.DKShareRegistryProvider {
		dkSharesRegistry, err := registry.NewDKSharesRegistry(ParamsRegistries.DKShares.Path, deps.NodeIdentityProvider.NodeIdentity().GetPrivateKey(), "IOTA")
		if err != nil {
			Component.LogPanic(err)
		}
		return dkSharesRegistry
	}); err != nil {
		Component.LogPanic(err)
	}

	if err := c.Provide(func() registry.TrustedPeersRegistryProvider {
		trustedPeersRegistryProvider, err := registry.NewTrustedPeersRegistryImpl(ParamsRegistries.TrustedPeers.FilePath)
		if err != nil {
			Component.LogPanic(err)
		}
		return trustedPeersRegistryProvider
	}); err != nil {
		Component.LogPanic(err)
	}

	return nil
}

func nodeIdentityRegistry() *registry.NodeIdentity {
	if err := ioutils.CreateDirectory(ParamsP2P.Database.Path, 0o700); err != nil {
		Component.LogPanicf("could not create peer store database dir '%s': %w", ParamsP2P.Database.Path, err)
	}

	// make sure nobody copies around the peer store since it contains the private key of the node
	Component.LogInfof(`WARNING: never share your "%s" or "%s" folder as both contain your node's private key!`, ParamsP2P.Database.Path, path.Dir(ParamsP2P.Identity.FilePath))

	// load up the previously generated identity or create a new one
	privKey, newlyCreated, err := hivep2p.LoadOrCreateIdentityPrivateKey(ParamsP2P.Identity.FilePath, ParamsP2P.Identity.PrivateKey)
	if err != nil {
		Component.LogPanic(err)
	}

	if newlyCreated {
		Component.LogInfof(`stored new private key for peer identity under "%s"`, ParamsP2P.Identity.FilePath)
	} else {
		Component.LogInfof(`loaded existing private key for peer identity from "%s"`, ParamsP2P.Identity.FilePath)
	}

	privKeyBytes, err := privKey.Raw()
	if err != nil {
		Component.LogPanicf("unable to convert private key for peer identity: %s", err)
	}

	waspPrivKey, err := cryptolib.PrivateKeyFromBytes(privKeyBytes)
	if err != nil {
		Component.LogPanicf("unable to convert private key for peer identity: %s", err)
	}

	waspKeyPair := cryptolib.KeyPairFromPrivateKey(waspPrivKey)
	Component.LogInfof("this node identity: %v", waspKeyPair.GetPublicKey())
	return registry.NewNodeIdentity(waspKeyPair)
}
