package pubsub

import (
	floodsub "github.com/libp2p/go-floodsub"
)

// Interface ...
type Interface interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, opts ...SubOpt) (*floodsub.Subscription, error)
}
