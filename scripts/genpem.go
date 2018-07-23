package main

import (
	"log"
	"os"

	"github.com/c3systems/c3-go/common/c3crypto"
)

// Generates ECDA key in PEM format
func main() {
	priv, err := c3crypto.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	var password *string
	if len(os.Args) > 1 {
		log.Println("creating pem files with password")
		tmpPassword := os.Args[1]
		password = &tmpPassword
	}

	privFileName := "priv.pem"
	if err := c3crypto.WritePrivateKeyToPemFile(priv, password, privFileName); err != nil {
		log.Fatal(err)
	}

	pub, err := c3crypto.GetPublicKey(priv)
	pubFileName := "public.pem"
	if err := c3crypto.WritePublicKeyToPemFile(pub, password, pubFileName); err != nil {
		log.Fatal(err)
	}

	log.Printf("created %s and %s", privFileName, pubFileName)
}
