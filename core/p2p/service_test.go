package p2p

import (
	"context"
	"crypto/rand"
	"testing"

	ds "github.com/ipfs/go-datastore"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	blankhost "github.com/libp2p/go-libp2p-blankhost"
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
)

func TestNew(t *testing.T) {
	priv, pub, err := ci.GenerateKeyPairWithReader(ci.RSA, 2048, rand.Reader)
	if err != nil {
		t.Error(err)
	}

	pid, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		t.Error(err)
	}

	ps := peerstore.NewPeerstore()
	if err := ps.AddPubKey(pid, pub); err != nil {
		t.Error(err)
	}
	if err := ps.AddPrivKey(pid, priv); err != nil {
		t.Error(err)
	}

	swrm := swarm.NewSwarm(context.Background(), pid, ps, nil)
	host := blankhost.NewBlankHost(swrm)
	svc, err := New(&Props{
		Host:       host,
		BlockStore: bstore.NewBlockstore(ds.NewMapDatastore()),
	})
	if err != nil {
		t.Error(err)
	}

	if svc == nil {
		t.FailNow()
	}
}
