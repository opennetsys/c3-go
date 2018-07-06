package pubsub

import (
	"context"

	pb "github.com/libp2p/go-floodsub/pb"
)

// Message ...
type Message struct {
	*pb.Message
}

// SubOpt ...
type SubOpt func(sub SubscriptionInterface) error

// SubscriptionInterface ...
type SubscriptionInterface interface {
	Topic() string
	Next(ctx context.Context) (*Message, error)
}
