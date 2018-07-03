package fldsb

import (
	"context"
	"errors"

	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

// New ...
// note: implement RPC? https://github.com/libp2p/go-floodsub/blob/master/floodsub.go#L47
func New(ctx context.Context, h *host.Host) (*floodsub.PubSub, error) {
	if h == nil {
		return nil, errors.New("host is required")
	}

	return floodsub.NewFloodSub(ctx, h)
}
