package main

import (
	"log"

	"github.com/c3systems/c3/common/c3crypto"
)

// Generates ECDA key in PEM format
func main() {
	priv, err := c3crypto.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	fileName := "priv.pem"
	if err := c3crypto.WritePrivateKeyToPemFile(priv, nil, fileName); err != nil {
		log.Fatal(err)
	}

	pub, err := c3crypto.GetPublicKey(priv)
	fileName = "public.pem"
	if err := c3crypto.WritePublicKeyToPemFile(pub, nil, fileName); err != nil {
		log.Fatal(err)
	}
}
