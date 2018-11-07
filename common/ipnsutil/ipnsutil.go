package ipnsutil

import (
	"bufio"
	"bytes"
	"encoding/pem"
	"errors"
	"fmt"
	"log"

	"github.com/c3systems/c3-go/common/c3crypto"
	keystore "github.com/ipfs/go-ipfs-keystore"
	lCrypt "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	b58 "github.com/mr-tron/base58/base58"
	mh "github.com/multiformats/go-multihash"
)

// PEMToIPNS ...
func PEMToIPNS(pemFilepath string, password *string) (string, error) {
	if pemFilepath == "" {
		return "", errors.New("pem filepath is required")
	}

	priv, err := c3crypto.ReadPrivateKeyFromPem(pemFilepath, password)
	if err != nil {
		return "", fmt.Errorf("err reading pem file\n%v", err)
	}

	_, wPub, err := lCrypt.ECDSAKeyPairFromKey(priv)
	if err != nil {
		return "", fmt.Errorf("err generating key pairs\n%v", err)
	}

	pid, err := peer.IDFromPublicKey(wPub)
	if err != nil {
		return "", err
	}

	return pid.Pretty(), nil
}

// KeystorePrivateKeyToIPNS ...
func KeystorePrivateKeyToIPNS(keystorePath string) string {
	ks, err := keystore.NewFSKeystore(keystorePath)
	if err != nil {
		log.Fatal(err)
	}

	priv, err := ks.Get("default")
	if err != nil {
		log.Fatal(err)
	}

	privBytes, err := priv.Bytes()
	if err != nil {
		log.Fatal(err)
	}

	var privPEM bytes.Buffer
	privWriter := bufio.NewWriter(&privPEM)
	pem.Encode(privWriter, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	privWriter.Flush()

	pub := priv.GetPublic()
	var pubPEM bytes.Buffer
	pubBytes, err := pub.Bytes()
	pubWriter := bufio.NewWriter(&pubPEM)
	pem.Encode(pubWriter, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})
	pubWriter.Flush()

	//privPEMString := string(privPEM.Bytes())
	//pubPEMString := string(pubPEM.Bytes())

	var alg uint64 = mh.SHA2_256
	maxInlineKeyLength := 42
	if len(pubBytes) <= maxInlineKeyLength {
		alg = mh.ID
	}
	hash, _ := mh.Sum(pubBytes, alg, -1)
	peerID := b58.Encode(hash)

	return peerID
}
