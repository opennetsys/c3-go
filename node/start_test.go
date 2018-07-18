package node

import (
	"os"
	"testing"
	"time"

	"github.com/c3systems/c3/common/c3crypto"
	"github.com/c3systems/c3/core/chain/statechain"
	//docker "github.com/c3systems/c3/core/docker"
	methodTypes "github.com/c3systems/c3/core/types/methods"
	nodetypes "github.com/c3systems/c3/node/types"
	"github.com/davecgh/go-spew/spew"
)

func TestBroadcast(t *testing.T) {
	// imageHash := "Qmf9XFxbFDGv4yssc7YvAisxxUBU89BFbimAAYgT33ZTAf"
	imageHash := "e8758b300c09"
	/*
			dockerclient := docker.NewClient()
			err := dockerclient.LoadImageByFilepath("./test_data/go_example_image.tar")
			if err != nil {
				t.Error(err)
			}

		registry := registry.NewRegistry(&registry.Config{})
		imageHash, err := registry.PushImageByID("goexample")
		if err != nil {
			t.Error(err)
		}
	*/

	privPEM := "./test_data/priv.pem"
	nodeURI := "/ip4/0.0.0.0/tcp/9004"
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
				t.Error(err)
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
	if err != nil {
		t.Error(err)
	}

	encodedPub, err := c3crypto.EncodeAddress(pub)
	if err != nil {
		t.Error(err)
	}

	tx1 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    methodTypes.Deploy,
		Payload:   []byte(`{"hello": "world"}`),
		From:      encodedPub,
	})

	tx2 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    methodTypes.InvokeMethod,
		Payload:   []byte(`[""setItem", "foo", "bar"]`),
		From:      encodedPub,
	})

	tx := tx1
	_ = tx1
	_ = tx2

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
