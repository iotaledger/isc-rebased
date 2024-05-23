package serialization

import (
	"context"
	"fmt"
	"log"

	"github.com/iotaledger/isc-private/sui-go/models"
	"github.com/iotaledger/isc-private/sui-go/sui"
	"github.com/iotaledger/isc-private/sui-go/sui_types"
)

type Subscriber struct {
	client *sui.ImplSuiAPI
	// *account.Account
}

func NewSubscriber(client *sui.ImplSuiAPI) *Subscriber {
	return &Subscriber{client: client}
}

func (s *Subscriber) SubscribeEvent(ctx context.Context, packageID *sui_types.PackageID) {
	resultCh := make(chan models.SuiEvent)
	err := s.client.SubscribeEvent(context.Background(), &models.EventFilter{Package: packageID}, resultCh)
	if err != nil {
		log.Fatal(err)
	}

	for result := range resultCh {
		fmt.Println("result: ", result)
	}
}
