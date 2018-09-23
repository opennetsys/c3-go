package node

import (
	"context"

	log "github.com/sirupsen/logrus"

	host "github.com/libp2p/go-libp2p-host"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

var h host.Host

// DiscoveryNotifee ...
type DiscoveryNotifee struct {
	h host.Host
}

// HandlePeerFound ...
func (n *DiscoveryNotifee) HandlePeerFound(pi peerstore.PeerInfo) {
	n.h.Peerstore().AddAddrs(pi.ID, pi.Addrs, peerstore.PermanentAddrTTL)
	if err := n.h.Connect(context.Background(), pi); err != nil {
		log.Printf("[node] found peer %s\nerr connecting %v", pi.Addrs, err)

		return
	}

	log.Printf("[node] found peer %s\nadded to peerstore and connected", pi.Addrs)
}
