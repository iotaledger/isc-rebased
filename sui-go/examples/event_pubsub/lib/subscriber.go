package serialization

import (
	"context"
	"fmt"
	"log"

	"github.com/iotaledger/wasp/sui-go/sui"
	"github.com/iotaledger/wasp/sui-go/suiclient"
	"github.com/iotaledger/wasp/sui-go/suijsonrpc"
)

type Subscriber struct {
	client *suiclient.WebsocketClient
	// *account.Account
}

func NewSubscriber(client *suiclient.WebsocketClient) *Subscriber {
	return &Subscriber{client: client}
}

func (s *Subscriber) SubscribeEvent(ctx context.Context, packageID *sui.PackageID) {
	resultCh := make(chan suijsonrpc.SuiEvent)
	err := s.client.SubscribeEvent(context.Background(), &suijsonrpc.EventFilter{Package: packageID}, resultCh)
	if err != nil {
		log.Fatal(err)
	}

	for result := range resultCh {
		fmt.Println("result: ", result)
	}
}