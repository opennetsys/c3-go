package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/core/chain/statechain"
	methodTypes "github.com/c3systems/c3-go/core/types/methods"
	"github.com/spf13/cobra"
)

func signCmd() *cobra.Command {
	var (
		image   string
		payload string
		privPEM string
		genesis bool
	)

	signcmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign a transaction",
		Long:  "Sign a transaction with a private key",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			priv, err := c3crypto.ReadPrivateKeyFromPem(privPEM, nil)
			if err != nil {
				return err
			}

			pub, err := c3crypto.GetPublicKey(priv)
			if err != nil {
				return err
			}

			encodedPub, err := c3crypto.EncodeAddress(pub)
			if err != nil {
				return err
			}

			method := methodTypes.InvokeMethod
			if genesis {
				method = methodTypes.Deploy
			}

			tx := statechain.NewTransaction(&statechain.TransactionProps{
				ImageHash: image,
				Method:    method,
				Payload:   []byte(payload),
				From:      encodedPub,
			})

			err = tx.SetHash()
			if err != nil {
				return err
			}

			err = tx.Sign(priv)
			if err != nil {
				return err
			}

			txJSON, err := json.Marshal(tx)
			if err != nil {
				return err
			}

			fmt.Println(string(txJSON))

			return nil
		},
	}

	signcmd.Flags().StringVarP(&privPEM, "priv", "k", "", "The private key to sign the transaction with")
	signcmd.Flags().StringVarP(&image, "image", "i", "", "The image hash for the transaction")
	signcmd.Flags().StringVarP(&payload, "payload", "p", "", "The transaction payload to sign")
	signcmd.Flags().BoolVarP(&genesis, "genesis", "g", false, "Set to true if this is a genesis transaction")

	return signcmd
}
