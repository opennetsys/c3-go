package node

import (
	"testing"

	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/statechain"
	nodetypes "github.com/c3systems/c3/node/types"
	"github.com/davecgh/go-spew/spew"
)

func TestBroadcast(t *testing.T) {
	pem := "./test_data/key.pem"
	password := ""
	p := "/ip4/192.168.0.9/tcp/9000/ipfs/QmcAfcm2kJZaiunoRwh8H7ihWbHez3w1hu8EAXKeMZD7pj"
	nodeURI := "/ip4/0.0.0.0/tcp/9000"
	dataDir := "~/.c3"
	n, err := Start(&nodetypes.Config{
		URI:     nodeURI,
		Peer:    p,
		DataDir: dataDir,
		Keys: nodetypes.Keys{
			PEMFile:  pem,
			Password: password,
		},
	})
	if err != nil {
		t.Error(err)
	}

	priv, err := c3crypto.ReadPrivateKeyFromPem(pem, &password)
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

	//hash := "foobar"

	tx.SetSig(priv)

	resp, err := n.BroadcastTransaction(tx)
	if err != nil {
		t.Error(err)
	}

	spew.Dump(resp)
}
