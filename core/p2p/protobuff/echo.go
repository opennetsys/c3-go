package protobuff

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	pb "github.com/c3systems/c3-go/core/p2p/protobuff/pb"

	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// pattern: /protocol-name/request-or-response-message/version
const echoRequest = "/echo/echoreq/0.0.1"
const echoResponse = "/echo/echoresp/0.0.1"

type echoRequestWrapper struct {
	resp chan interface{}
	req  *pb.EchoRequest
}

// Echo ...
type Echo struct {
	node     *Node                          // local host
	requests map[string]*echoRequestWrapper // used to access request data from response handlers
}

// NewEcho ...
func NewEcho(node *Node) *Echo {
	e := Echo{node: node, requests: make(map[string]*echoRequestWrapper)}
	node.SetStreamHandler(echoRequest, e.onEchoRequest)
	node.SetStreamHandler(echoResponse, e.onEchoResponse)

	// design note: to implement fire-and-forget style messages you may just skip specifying a response callback.
	// a fire-and-forget message will just include a request and not specify a response object

	return &e
}

// remote peer requests handler
func (e *Echo) onEchoRequest(s inet.Stream) {
	// get request data
	data := &pb.EchoRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	if err != nil {
		log.Errorf("[p2p] %s", err)
		return
	}

	log.Printf("[p2p] %s: Received echo request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)

	valid := e.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Error("Failed to authenticate message")
		return
	}

	log.Printf("[p2p] %s: Sending echo response to %s. Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)

	// send response to the request using the message string he provided

	resp := &pb.EchoResponse{
		MessageData: e.node.NewMessageData(data.MessageData.Id, false),
		Message:     data.Message}

	// sign the data
	signature, err := e.node.signProtoMessage(resp)
	if err != nil {
		log.Error("[p2p] failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = string(signature)

	s, respErr := e.node.NewStream(context.Background(), s.Conn().RemotePeer(), echoResponse)
	if respErr != nil {
		log.Error("[p2p] %s", respErr)
		return
	}

	ok := e.node.sendProtoMessage(resp, s)

	if ok {
		log.Printf("[p2p] %s: Echo response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}
}

// remote echo response handler
func (e *Echo) onEchoResponse(s inet.Stream) {
	data := &pb.EchoResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	if err := decoder.Decode(data); err != nil {
		log.Errorf("[p2p] err decoding echo response\n%v", err)
		return
	}

	// locate request data and remove it if found
	reqW, ok := e.requests[data.MessageData.Id]
	if ok {
		// remove request from map as we have processed it here
		delete(e.requests, data.MessageData.Id)
	} else {
		log.Error("[p2p] failed to locate request data boject for response")
		return
	}

	// authenticate message content
	valid := e.node.authenticateMessage(data, data.MessageData)

	if !valid {
		reqW.resp <- errors.New("Failed to authenticate message")
		return
	}

	reqW.resp <- data
}

// SendEcho ...
func (e *Echo) SendEcho(peerID peer.ID, resp chan interface{}) error {
	// log.Printf("[p2p] %s: Sending echo to: %s....", e.node.ID(), peerID)

	// create message data
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	req := &pb.EchoRequest{
		MessageData: e.node.NewMessageData(id.String(), false),
		Message:     fmt.Sprintf("Echo from %s", e.node.ID()),
	}

	signature, err := e.node.signProtoMessage(req)
	if err != nil {
		log.Error("[p2p] failed to sign message")
		return err
	}

	// add the signature to the message
	req.MessageData.Sign = string(signature)

	s, err := e.node.NewStream(context.Background(), peerID, echoRequest)
	if err != nil {
		log.Errorf("[p2p] %s", err)
		return err
	}

	ok := e.node.sendProtoMessage(req, s)

	if !ok {
		return errors.New("failed to send message")
	}

	// store request so response handler has access to it
	e.requests[req.MessageData.Id] = &echoRequestWrapper{
		resp: resp,
		req:  req,
	}
	// log.Printf("[p2p] %s: Echo to: %s was sent. Message Id: %s, Message: %s", e.node.ID(), peerID, req.MessageData.Id, req.Message)
	return nil
}
