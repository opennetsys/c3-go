package cmd

import (
	"errors"
	"time"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/config"
	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/node"
	nodetypes "github.com/c3systems/c3-go/node/types"
)

func broadcastTx(txType, image, payloadStr, peer, privPEM string) (string, error) {
	cnf := config.New()

	nodeURI := "/ip4/0.0.0.0/tcp/9911"
	dataDir := "~/.c3-2"
	n, err := node.NewFullNode(&nodetypes.Config{
		URI:     nodeURI,
		Peer:    peer,
		DataDir: dataDir,
		Keys: nodetypes.Keys{
			PEMFile:  privPEM,
			Password: "",
		},
		BlockDifficulty: cnf.BlockDifficulty(),
	})

	if err != nil {
		return "", err
	}

	priv, err := c3crypto.ReadPrivateKeyFromPem(privPEM, nil)
	if err != nil {
		return "", err
	}

	pub, err := c3crypto.GetPublicKey(priv)
	if err != nil {
		return "", err
	}

	encodedPub, err := c3crypto.EncodeAddress(pub)
	if err != nil {
		return "", err
	}

	payload := []byte(payloadStr)

	tx := statechain.NewTransaction(&statechain.TransactionProps{
		ImageHash: image,
		Method:    txType,
		Payload:   payload,
		From:      encodedPub,
	})

	err = tx.SetHash()
	if err != nil {
		return "", err
	}

	err = tx.SetSig(priv)
	if err != nil {
		return "", err
	}

	resp, err := n.BroadcastTransaction(tx)
	if err != nil {
		return "", err
	}

	if resp.TxHash == nil {
		return "", errors.New("expected hash")
	}

	time.Sleep(3 * time.Second)
	return *resp.TxHash, err
}
