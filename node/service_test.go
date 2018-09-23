package node

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/common/txparamcoder"
	"github.com/c3systems/c3-go/core/chain/statechain"
	log "github.com/sirupsen/logrus"
	//docker "github.com/c3systems/c3-go/core/docker"
	methodTypes "github.com/c3systems/c3-go/core/types/methods"
	nodetypes "github.com/c3systems/c3-go/node/types"
	"github.com/davecgh/go-spew/spew"
)

func TestBroadcast(t *testing.T) {
	imageHash := os.Getenv("IMAGEID")
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
	nodeURI := "/ip4/0.0.0.0/tcp/3332"
	peer := os.Getenv("PEERID")
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
				BlockDifficulty: 5,
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
		//Payload:   []byte(`{"hello": "world"}`),
		Payload: []byte(``),
		From:    encodedPub,
	})

	fileBytes, err := ioutil.ReadFile("../../example-go-image-recognition/images/cat/cat.jpg")
	if err != nil {
		log.Error(err)
	}

	payload := txparamcoder.ToJSONArray(
		txparamcoder.EncodeMethodName("processImage"),
		txparamcoder.EncodeParam(hex.EncodeToString(fileBytes)),
		txparamcoder.EncodeParam("jpg"),
	)

	tx2 := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    methodTypes.InvokeMethod,
		Payload:   payload,
		From:      encodedPub,
	})

	tx := tx2
	if os.Getenv("METHOD") == "deploy" {
		tx = tx1
	}

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
	time.Sleep(5 * time.Second) // needed
}
