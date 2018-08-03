package protobuff

import (
	"bufio"
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/core/miner"
	pb "github.com/c3systems/c3-go/core/p2p/protobuff/pb"
	nodetypes "github.com/c3systems/c3-go/node/types"

	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const processTransactionRequest = "/processtransaction/processtransactionreq/0.0.1"
const processTransactionResponse = "/processtransaction/processtransactionresp/0.0.1"

type processTransactionRequestWrapper struct {
	resp chan interface{}
	req  *pb.ProcessTransactionRequest
}

// ProcessTransaction ...
type ProcessTransaction struct {
	node                   *Node                                        // local host
	requests               map[string]*processTransactionRequestWrapper // used to access request data from response handlers
	broadcastTransactionFN func(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error)
	addPendingTxFN         func(tx *statechain.Transaction) error
}

// NewProcessTransaction ...
func NewProcessTransaction(node *Node, broadcastTransactionFN func(tx *statechain.Transaction) (*nodetypes.SendTxResponse, error), addPendingTxFN func(tx *statechain.Transaction) error) *ProcessTransaction {
	p := ProcessTransaction{
		node:                   node,
		requests:               make(map[string]*processTransactionRequestWrapper),
		broadcastTransactionFN: broadcastTransactionFN,
		addPendingTxFN:         addPendingTxFN,
	}
	node.SetStreamHandler(processTransactionRequest, p.onProcessTransactionRequest)
	node.SetStreamHandler(processTransactionResponse, p.onProcessTransactionResponse)

	// design note: to implement fire-and-forget style messages you may just skip specifying a response callback.
	// a fire-and-forget message will just include a request and not specify a response object

	return &p
}

// remote peer requests handler
func (p *ProcessTransaction) onProcessTransactionRequest(s inet.Stream) {
	// get request data
	data := &pb.ProcessTransactionRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	if err != nil {
		log.Errorf("[p2p] %s", err)
		return
	}

	// log.Printf("[p2p] %s: Received process transaction request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)
	valid := p.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Error("[p2p] failed to authenticate message")
		return
	}

	resp := &pb.ProcessTransactionResponse{
		MessageData: p.node.NewMessageData(data.MessageData.Id, false),
	}
	// log.Printf("[p2p] %s: Sending process transaction response to %s. Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)
	// interpret the tx
	tx := new(statechain.Transaction)
	if err := tx.Deserialize(data.TxBytes); err != nil {
		resp.Success = false
		resp.Message = fmt.Sprintf("err deserializing tx: %v", err)

		p.sendResp(resp, s)
		return
	}

	ok, err := miner.VerifyTransaction(tx)
	if err != nil {
		resp.Success = false
		resp.Message = fmt.Sprintf("err verifying tx: %v", err)

		p.sendResp(resp, s)
		return
	}
	if !ok {
		resp.Success = false
		resp.Message = "invalid transaction"

		p.sendResp(resp, s)
		return
	}

	// should already be checked in the VerifyTransaction fn, but just to be safe...
	hash := tx.Props().TxHash
	if hash == nil {
		resp.Success = false
		resp.Message = "nil tx hash"

		p.sendResp(resp, s)
		return

	}

	if _, err := p.broadcastTransactionFN(tx); err != nil {
		resp.Success = false
		resp.Message = fmt.Sprintf("err verifying tx: %v", err)

		p.sendResp(resp, s)
		return
	}

	if err := p.addPendingTxFN(tx); err != nil {
		log.Errorf("err adding pending tx\n%v", err)
	}

	resp.Success = true
	resp.Message = "tx sent"
	resp.Hash = *hash // note: checked for nil hash, above
	p.sendResp(resp, s)
	return
}

func (p *ProcessTransaction) sendResp(resp *pb.ProcessTransactionResponse, s inet.Stream) {
	// sign the data
	signature, err := p.node.signProtoMessage(resp)
	if err != nil {
		log.Errorf("[p2p] failed to sign response\n%v", err)
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = string(signature)

	s, respErr := p.node.NewStream(context.Background(), s.Conn().RemotePeer(), processTransactionResponse)
	if respErr != nil {
		log.Errorf("[p2p] %s", respErr)
		return
	}

	ok := p.node.sendProtoMessage(resp, s)

	if ok {
		log.Printf("[p2p] %s: process transaction response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}
}

// remote peer response handler
func (p *ProcessTransaction) onProcessTransactionResponse(s inet.Stream) {
	data := &pb.ProcessTransactionResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	if err := decoder.Decode(data); err != nil {
		log.Errorf("[p2p] err decoding process transaction response\n%v", err)
		return
	}

	// locate request data and remove it if found
	reqW, ok := p.requests[data.MessageData.Id]
	if ok {
		// remove request from map as we have processed it here
		delete(p.requests, data.MessageData.Id)
	} else {
		log.Error("[p2p] failed to locate request data object for response")
		return
	}

	// authenticate message content
	valid := p.node.authenticateMessage(data, data.MessageData)

	if !valid {
		reqW.resp <- errors.New("Failed to authenticate message")
		return
	}

	reqW.resp <- data
}

// SendTransaction ...
func (p *ProcessTransaction) SendTransaction(peerID peer.ID, txBytes []byte, resp chan interface{}) error {
	// log.Printf("[p2p] %s: Sending process transaction to: %s....", e.node.ID(), peerID)

	// create message data
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	req := &pb.ProcessTransactionRequest{
		MessageData: p.node.NewMessageData(id.String(), false),
		TxBytes:     txBytes,
	}

	signature, err := p.node.signProtoMessage(req)
	if err != nil {
		log.Error("[p2p] failed to sign message")
		return err
	}

	// add the signature to the message
	req.MessageData.Sign = string(signature)

	s, err := p.node.NewStream(context.Background(), peerID, processTransactionRequest)
	if err != nil {
		log.Errorf("[p2p] %s", err)
		return err
	}

	ok := p.node.sendProtoMessage(req, s)

	if !ok {
		return errors.New("failed to send message")
	}

	// store request so response handler has access to it
	p.requests[req.MessageData.Id] = &processTransactionRequestWrapper{
		resp: resp,
		req:  req,
	}
	// log.Printf("[p2p] %s: process transaction to: %s was sent. Message Id: %s, Message: %s", e.node.ID(), peerID, req.MessageData.Id, req.Message)
	return nil
}
