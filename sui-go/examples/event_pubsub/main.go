package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/sui/conn"
	"github.com/iotaledger/wasp/sui-go/sui_signer"
	"github.com/iotaledger/wasp/sui-go/sui_types"

	serialization "github.com/iotaledger/wasp/sui-go/examples/event_pubsub/lib"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	sender, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC, sui_signer.KeySchemeFlagDefault)
	if err != nil {
		log.Panic(err)
	}
	err = sui.RequestFundFromFaucet(sender.Address, conn.TestnetFaucetUrl)
	if err != nil {
		log.Panic(err)
	}

	packageID, err := sui_types.PackageIDFromHex("")
	if err != nil {
		log.Panic(err)
	}

	log.Println("sender: ", sender.Address)
	publisher := serialization.NewPublisher(api, sender)
	subscriber := serialization.NewSubscriber(api)

	go func() {
		for {
			publisher.PublishEvents(context.Background(), packageID)
		}
	}()

	go func() {
		for {
			subscriber.SubscribeEvent(context.Background(), packageID)
		}
	}()

	<-done
}
