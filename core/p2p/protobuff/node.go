package protobuff

import (
	"bufio"
	"errors"
	"time"

	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
	pb "github.com/c3systems/c3-go/core/p2p/protobuff/pb"
	nodetypes "github.com/c3systems/c3-go/node/types"

	"github.com/gogo/protobuf/proto"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	log "github.com/sirupsen/logrus"
)

// node client version
const clientVersion = "go-p2p-node/0.0.1"

// Props ...
type Props struct {
	Host                   host.Host
	GetHeadBlockFN         func() (mainchain.Block, error)
	BroadcastTransactionFN func(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error)
	AddPendingTxFN         func(tx *statechain.Transaction) error
}

// Node type - a p2p host implementing one or more p2p protocols
type Node struct {
	host.Host           // lib-p2p host
	*Echo               // echo protocol impl
	*HeadBlock          // headblock protocol impl
	*ProcessTransaction // process transaction impl
	// add other protocols here...
}

// NewNode creates a new node with its implemented protocols
func NewNode(props *Props) (*Node, error) {
	if props == nil {
		return nil, errors.New("nil props")
	}

	node := &Node{Host: props.Host}
	node.Echo = NewEcho(node)
	node.HeadBlock = NewHeadBlock(node, props.GetHeadBlockFN)
	node.ProcessTransaction = NewProcessTransaction(node, props.BroadcastTransactionFN, props.AddPendingTxFN)
	return node, nil
}

// Authenticate incoming p2p message
// message: a protobufs go data object
// data: common p2p message data
func (n *Node) authenticateMessage(message proto.Message, data *pb.MessageData) bool {
	// store a temp ref to signature and remove it from message data
	// sign is a string to allow easy reset to zero-value (empty string)
	sign := data.Sign
	data.Sign = ""

	// marshall data without the signature to protobufs3 binary format
	bin, err := proto.Marshal(message)
	if err != nil {
		log.Errorf("[p2p] failed to marshal pb message %s", err)
		return false
	}

	// restore sig in message data (for possible future use)
	data.Sign = sign

	// restore peer id binary format from base58 encoded node id data
	peerID, err := peer.IDB58Decode(data.NodeId)
	if err != nil {
		log.Errorf("[p2p] failed to decode node id from base58 %s", err)
		return false
	}

	// verify the data was authored by the signing peer identified by the public key
	// and signature included in the message
	return n.verifyData(bin, []byte(sign), peerID, data.NodePubKey)
}

// sign an outgoing p2p message payload
func (n *Node) signProtoMessage(message proto.Message) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return n.signData(data)
}

// sign binary data using the local node's private key
func (n *Node) signData(data []byte) ([]byte, error) {
	key := n.Peerstore().PrivKey(n.ID())
	res, err := key.Sign(data)
	return res, err
}

// Verify incoming p2p message data integrity
// data: data to verify
// signature: author signature provided in the message payload
// peerId: author peer id from the message payload
// pubKeyData: author public key from the message payload
func (n *Node) verifyData(data []byte, signature []byte, peerID peer.ID, pubKeyData []byte) bool {
	key, err := crypto.UnmarshalPublicKey(pubKeyData)
	if err != nil {
		log.Errorf("[p2p] failed to extract key from message key data %s", err)
		return false
	}

	// extract node id from the provided public key
	idFromKey, err := peer.IDFromPublicKey(key)

	if err != nil {
		log.Errorf("[p2p] failed to extract peer id from public key %s", err)
		return false
	}

	// verify that message author node id matches the provided node public key
	if idFromKey != peerID {
		log.Error("[p2p] node id and provided public key mismatch")
		return false
	}

	res, err := key.Verify(data, signature)
	if err != nil {
		log.Errorf("[p2p] error authenticating data %s", err)
		return false
	}

	return res
}

// NewMessageData is helper method to generate message data shared between all node's p2p protocols
// messageId: unique for requests, copied from request for responses
func (n *Node) NewMessageData(messageID string, gossip bool) *pb.MessageData {
	// Add protobufs bin data for message author public key
	// this is useful for authenticating  messages forwarded by a node authored by another node
	nodePubKey, err := n.Peerstore().PubKey(n.ID()).Bytes()

	if err != nil {
		panic("Failed to get public key for sender from local peer store.")
	}

	return &pb.MessageData{ClientVersion: clientVersion,
		NodeId:     peer.IDB58Encode(n.ID()),
		NodePubKey: nodePubKey,
		Timestamp:  time.Now().Unix(),
		Id:         messageID,
		Gossip:     gossip}
}

// helper method - writes a protobuf go data object to a network stream
// data: reference of protobuf go data object to send (not the object itself)
// s: network stream to write the data to
func (n *Node) sendProtoMessage(data proto.Message, s inet.Stream) bool {
	writer := bufio.NewWriter(s)
	enc := protobufCodec.Multicodec(nil).Encoder(writer)
	if err := enc.Encode(data); err != nil {
		log.Errorf("[p2p] err encoding\n%v", err)

		return false
	}

	if err := writer.Flush(); err != nil {
		log.Errorf("[p2p] err flushing writer\n%v", err)

		return false
	}

	return true
}
