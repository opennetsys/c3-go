package pubsub

// Interface ...
type Interface interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, opts ...SubOpt) (*Subscription, error)
}
