package fldsb

import (
	"context"
	"errors"

	host "github.com/libp2p/go-libp2p-host"
	// host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"

	floodsub "github.com/libp2p/go-floodsub"
)

// New ...
// note: implement RPC? https://github.com/libp2p/go-floodsub/blob/master/floodsub.go#L47
func New(ctx context.Context, h *host.Host) (*floodsub.PubSub, error) {
	if h == nil {
		return nil, errors.New("host is required")
	}

	return floodsub.NewFloodSub(ctx, *h)
}
