package snapshot

import (
	"testing"

	"github.com/c3systems/c3-go/node"
	nodetypes "github.com/c3systems/c3-go/node/types"
	log "github.com/sirupsen/logrus"
	//ds "github.com/ipfs/go-datastore"
)

func TestSnapshot(t *testing.T) {
	nodeURI := "/ip4/0.0.0.0/tcp/3330"
	dataDir := "../.tmp"
	//privPEM := "../node/test_data/priv1.pem"
	privPEM := "../priv.pem"
	//peer := "/ip4/127.0.0.1/tcp/3330/ipfs/QmZPNaCnnR59Dtw5nUuxv33pNXxRqKurnZTHLNJ6LaqEnx"
	peer := ""
	n, err := node.NewFullNode(&nodetypes.Config{
		URI:     nodeURI,
		Peer:    peer,
		DataDir: dataDir,
		Keys: nodetypes.Keys{
			PEMFile:  privPEM,
			Password: "",
		},
		BlockDifficulty: 5,
		MempoolType:     "memory",
		RPCHost:         ":5005",
	})
	if err != nil {
		t.Error(err)
	}

	svc := New(&Config{
		P2P:     n.Props().P2P,
		Mempool: n.Props().Store,
	})

	imageHash := "d50ada614c01"
	stateBlockNumber := 2
	snapshotImageID, err := svc.Snapshot(imageHash, stateBlockNumber)
	if err != nil {
		log.Fatal(err)
	}

	if snapshotImageID == "" {
		log.Fatal("expected image ID")
	}

	t.Log(snapshotImageID)
}
