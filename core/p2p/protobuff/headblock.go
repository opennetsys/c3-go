package protobuff

import (
	"bufio"
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"

	"github.com/c3systems/c3-go/core/chain/mainchain"
	pb "github.com/c3systems/c3-go/core/p2p/protobuff/pb"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const headBlockRequest = "/headblock/headblockreq/0.0.1"
const headBlockResponse = "/headblock/headblockresp/0.0.1"

type headBlockRequestWrapper struct {
	resp chan interface{}
	req  *pb.HeadBlockRequest
}

// HeadBlock ...
type HeadBlock struct {
	node           *Node                               // local host
	requests       map[string]*headBlockRequestWrapper // used to access request data from response handlers
	getHeadBlockFN func() (mainchain.Block, error)
}

// NewHeadBlock ...
func NewHeadBlock(node *Node, getHeadBlockFN func() (mainchain.Block, error)) *HeadBlock {
	h := HeadBlock{
		node:           node,
		requests:       make(map[string]*headBlockRequestWrapper),
		getHeadBlockFN: getHeadBlockFN,
	}
	node.SetStreamHandler(headBlockRequest, h.onHeadBlockRequest)
	node.SetStreamHandler(headBlockResponse, h.onHeadBlockResponse)

	// design note: to implement fire-and-forget style messages you may just skip specifying a response callback.
	// a fire-and-forget message will just include a request and not specify a response object

	return &h
}

// remote peer requests handler
func (h *HeadBlock) onHeadBlockRequest(s inet.Stream) {
	// get request data
	data := &pb.HeadBlockRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	if err != nil {
		log.Errorf("[p2p] %s", err)
		return
	}

	// log.Printf("[p2p] %s: Received headblock request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)
	valid := h.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Error("[p2p] failed to authenticate message")
		return
	}

	// log.Printf("[p2p] %s: Sending headblock response to %s. Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)
	// fetch our head block
	headBlock, err := h.getHeadBlockFN()
	if err != nil {
		log.Errorf("[p2p] err getting headblock\n%v", err)
		return
	}
	bytes, err := headBlock.Serialize()
	if err != nil {
		log.Errorf("[p2p] err serializing headblock\n%v", err)
		return
	}

	// send response to the request using the message string he provided

	resp := &pb.HeadBlockResponse{
		MessageData:    h.node.NewMessageData(data.MessageData.Id, false),
		HeadBlockBytes: bytes,
	}

	// sign the data
	signature, err := h.node.signProtoMessage(resp)
	if err != nil {
		log.Errorf("[p2p] failed to sign response\n%v", err)
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = string(signature)

	s, respErr := h.node.NewStream(context.Background(), s.Conn().RemotePeer(), headBlockResponse)
	if respErr != nil {
		log.Errorf("[p2p] %s", respErr)
		return
	}

	ok := h.node.sendProtoMessage(resp, s)

	if ok {
		log.Printf("[p2p] %s: headblock response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}
}

// remote peer response handler
func (h *HeadBlock) onHeadBlockResponse(s inet.Stream) {
	data := &pb.HeadBlockResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	if err := decoder.Decode(data); err != nil {
		log.Errorf("[p2p] err decoding headblock response\n%v", err)
		return
	}

	// locate request data and remove it if found
	reqW, ok := h.requests[data.MessageData.Id]
	if ok {
		// remove request from map as we have processed it here
		delete(h.requests, data.MessageData.Id)
	} else {
		log.Error("[p2p] ailed to locate request data boject for response")
		return
	}

	// authenticate message content
	valid := h.node.authenticateMessage(data, data.MessageData)

	if !valid {
		reqW.resp <- errors.New("Failed to authenticate message")
		return
	}

	reqW.resp <- data
}

// FetchHeadBlock ...
func (h *HeadBlock) FetchHeadBlock(peerID peer.ID, resp chan interface{}) error {
	// log.Printf("[p2p] %s: Sending headblock to: %s....", e.node.ID(), peerID)

	// create message data
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	req := &pb.HeadBlockRequest{
		MessageData: h.node.NewMessageData(id.String(), false),
	}

	signature, err := h.node.signProtoMessage(req)
	if err != nil {
		log.Error("[p2p] failed to sign message")
		return err
	}

	// add the signature to the message
	req.MessageData.Sign = string(signature)

	s, err := h.node.NewStream(context.Background(), peerID, headBlockRequest)
	if err != nil {
		log.Errorf("[p2p] %s", err)
		return err
	}

	ok := h.node.sendProtoMessage(req, s)

	if !ok {
		return errors.New("failed to send message")
	}

	// store request so response handler has access to it
	h.requests[req.MessageData.Id] = &headBlockRequestWrapper{
		resp: resp,
		req:  req,
	}
	// log.Printf("[p2p] %s: Headblock to: %s was sent. Message Id: %s, Message: %s", e.node.ID(), peerID, req.MessageData.Id, req.Message)
	return nil
}
