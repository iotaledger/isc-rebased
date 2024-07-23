package suiclient

import (
	"context"

	"github.com/iotaledger/wasp/sui-go/suiconn"
)

type Client struct {
	transport transport
}

type transport interface {
	Call(ctx context.Context, v any, method suiconn.JsonRPCMethod, args ...any) error
	Subscribe(ctx context.Context, v chan<- []byte, method suiconn.JsonRPCMethod, args ...any) error
}
