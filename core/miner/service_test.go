package miner

import (
	"testing"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/core/chain/statechain"
	methodTypes "github.com/c3systems/c3-go/core/types/methods"
)

func TestBuildGenesisStateBlock(t *testing.T) {
	imageHash := "QmQpXfKvirguQaMG7khqvLrqWcxEzh2qVApfC1Ts7QyFK7"

	privPEM := "./test_data/priv.pem"
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

	tx := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: imageHash,
		Method:    methodTypes.InvokeMethod,
		Payload:   []byte(`["setItem", "foo", "bar"]`),
		From:      encodedPub,
	})
	err = tx.SetHash()
	if err != nil {
		t.Error(err)
	}

	txs := []*statechain.Transaction{tx}
	// ctx := context.TODO()
	// mnr, err := New(&Props{
	// 	Context:             ctx,
	// 	PreviousBlock:       nil,
	// 	Difficulty:          uint64(5),
	// 	Channel:             make(chan interface{}),
	// 	Async:               true,
	// 	EncodedMinerAddress: "",
	// 	PendingTransactions: txs,
	// })

	// if err != nil {
	// 	t.Error(err)
	// }

	_, _, err = buildGenesisStateBlock(imageHash, txs[0])
	if err != nil {
		t.Error(err)
	}
}
