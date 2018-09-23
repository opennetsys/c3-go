package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/c3systems/c3-go/config"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/p2p/protobuff"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	log "github.com/sirupsen/logrus"

	pb "github.com/c3systems/c3-go/core/p2p/protobuff/pb"
	csms "github.com/libp2p/go-conn-security-multistream"
	secio "github.com/libp2p/go-libp2p-secio"
	swarm "github.com/libp2p/go-libp2p-swarm"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

// note: https://github.com/libp2p/go-libp2p-swarm/blob/da01184afe4c67bec58c5e73f3350ad80b624c0d/testing/testing.go#L39
func genUpgrader(n *swarm.Swarm) *tptu.Upgrader {
	id := n.LocalPeer()
	pk := n.Peerstore().PrivKey(id)
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(secio.ID, &secio.Transport{
		LocalID:    id,
		PrivateKey: pk,
	})

	stMuxer := msmux.NewBlankTransport()
	stMuxer.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	return &tptu.Upgrader{
		Secure:  secMuxer,
		Muxer:   stMuxer,
		Filters: n.Filters,
	}
}

func fetchHeadBlock(self peer.ID, headBlock *mainchain.Block, peers []peer.ID, pBuff protobuff.Interface) error {
	// TODO: pass contexts to pBuff functions
	ctx1, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
	defer cancel()
	ch := make(chan interface{})

	var peer peer.ID
	for _, peerID := range peers {
		if peerID != self {
			peer = peerID
			break
		}
	}

	if err := pBuff.FetchHeadBlock(peer, ch); err != nil {
		return err
	}

	select {
	case v := <-ch:
		switch v.(type) {
		case error:
			err, _ := v.(error)
			return err

		case *pb.HeadBlockResponse:
			hb, _ := v.(*pb.HeadBlockResponse)
			if hb == nil {
				return errors.New("received nil headblock")
			}

			block := new(mainchain.Block)
			if err := block.Deserialize(hb.HeadBlockBytes); err != nil {
				return err
			}

			headBlock = block
			return nil

		default:
			return fmt.Errorf("received unknown type %T\n%v", v, v)

		}

	case <-ctx1.Done():
		return errors.New("fetching headblock from peer timedout")

	}
}

func sendEcho(self peer.ID, peers []peer.ID, pBuff protobuff.Interface) error {
	ctx1, cancel := context.WithTimeout(context.Background(), config.IPFSTimeout)
	defer cancel()
	ch := make(chan interface{})

	var peer peer.ID
	for _, peerID := range peers {
		if peerID != self {
			peer = peerID
			break
		}
	}

	if err := pBuff.SendEcho(peer, ch); err != nil {
		return err
	}

	select {
	case v := <-ch:
		switch v.(type) {
		case error:
			err, _ := v.(error)
			return err

		case *pb.EchoResponse:
			eb, _ := v.(*pb.EchoResponse)
			log.Printf("[node] received echo response\n%v", eb)

			return nil

		default:
			return fmt.Errorf("received unknown type %T\n%v", v, v)

		}

	case <-ctx1.Done():
		return errors.New("echo timedout")

	}
}

func onConn(network net.Network, conn net.Conn) {
	log.Printf("[node] peer did connect\nid %v peerAddr %v", conn.RemotePeer().Pretty(), conn.RemoteMultiaddr())

	addAddr(conn)
}

func addAddr(conn net.Conn) {
	for _, peer := range h.Peerstore().Peers() {
		if conn.RemotePeer() == peer {
			// note: we already have info on this peer
			log.Println("[node] already have peer in our peerstore")
			return
		}
	}

	// note: we don't have this peer's info
	h.Peerstore().AddAddr(conn.RemotePeer(), conn.RemoteMultiaddr(), peerstore.PermanentAddrTTL)
	log.Printf("[node] added %s to our peerstore", conn.RemoteMultiaddr())

	if _, err := h.Network().DialPeer(context.Background(), conn.RemotePeer()); err != nil {
		log.Errorf("[node] err connecting to a peer\n%v", err)

		return
	}

	log.Printf("[node] connected to %s", conn.RemoteMultiaddr())
}
