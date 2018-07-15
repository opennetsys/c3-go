package node

import (
	"encoding/hex"
	"log"
	"os"
	"testing"
	"time"

	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/statechain"
	nodetypes "github.com/c3systems/c3/node/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
)

// TODO: finish

func TestBroadcast(t *testing.T) {
	privPEM := "./test_data/priv.pem"
	nodeURI := "/ip4/0.0.0.0/tcp/9006"
	peer := os.Getenv("PEER")
	dataDir := "~/.c3"
	n := new(Service)
	ready := make(chan bool)
	go func() {
		go func() {
			err := Start(n, &nodetypes.Config{
				URI:     nodeURI,
				Peer:    peer,
				DataDir: dataDir,
				Keys: nodetypes.Keys{
					PEMFile:  privPEM,
					Password: "",
				},
			})

			if err != nil {
				log.Fatal(err)
			}
		}()

		time.Sleep(10 * time.Second)
		ready <- true
	}()

	<-ready

	priv, err := c3crypto.ReadPrivateKeyFromPem(privPEM, nil)
	if err != nil {
		t.Error(err)
	}

	pub, err := c3crypto.GetPublicKey(priv)
	pubBytes := crypto.FromECDSAPub(pub)

	imageHash := "hello-world"
	tx := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    "c3_transaction",
		Payload:   []byte(`["foo", "bar"]`),
		From:      hex.EncodeToString(pubBytes),
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
