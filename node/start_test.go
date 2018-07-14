package node

import (
	"testing"
	"time"

	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/statechain"
	nodetypes "github.com/c3systems/c3/node/types"
	"github.com/davecgh/go-spew/spew"
)

func TestBroadcast(t *testing.T) {
	pem := "./test_data/key.pem"
	nodeURI := "/ip4/0.0.0.0/tcp/9000"
	dataDir := "~/.c3"
	n := new(Service)
	ready := make(chan bool)
	go func() {
		go Start(n, &nodetypes.Config{
			URI:     nodeURI,
			Peer:    "/ip4/192.168.84.20/tcp/9005/ipfs/QmTwGKMGvVL9Txti4AhB14bRZ6rhQNeAX4ZAhn6r1xfUxK",
			DataDir: dataDir,
			Keys: nodetypes.Keys{
				PEMFile:  pem,
				Password: "",
			},
		})

		time.Sleep(60 * time.Second)
		ready <- true
	}()

	<-ready

	priv, err := c3crypto.ReadPrivateKeyFromPem(pem, nil)
	if err != nil {
		t.Error(err)
	}

	imageHash := "hello-world"
	tx := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    "c3_deploy",
		Payload:   []byte(`["foo", "bar"]`),
		From:      "abc",
	})

	err = tx.SetHash()
	if err != nil {
		t.Error(err)
	}

	err = tx.SetSig(priv)
	if err != nil {
		t.Error(err)
	}

	resp, err := n.BroadcastTransaction(tx)
	if err != nil {
		t.Error(err)
	}

	if resp.TxHash == nil {
		t.Error("expected txhash")
	}

	spew.Dump(resp)
}
