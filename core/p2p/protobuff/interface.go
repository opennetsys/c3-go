package protobuff

import (
	pb "github.com/c3systems/c3-go/core/p2p/protobuff/pb"
	peer "github.com/libp2p/go-libp2p-peer"
)

// Interface ...
type Interface interface {
	NewMessageData(messageID string, gossip bool) *pb.MessageData
	SendEcho(peerID peer.ID, resp chan interface{}) error
	FetchHeadBlock(peerID peer.ID, resp chan interface{}) error
	SendTransaction(peerID peer.ID, txBytes []byte, resp chan interface{}) error
}
