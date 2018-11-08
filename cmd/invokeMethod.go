package cmd

import (
	methodTypes "github.com/c3systems/c3-go/core/types/methods"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func invokeMethodCmd() *cobra.Command {
	var (
		peer    string
		image   string
		payload string
		privPEM string
	)

	invokemethodcmd := &cobra.Command{
		Use:   "invokeMethod",
		Short: "Invoke a method on a dApp",
		Long:  "Broadcasts a transation on the network to invoke a method on a dApp given the IPFS image ID",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("pem: %s", privPEM)
			log.Printf("image: %s", image)
			log.Printf("peer: %s", peer)
			log.Printf("payload: %s", payload)

			txHash, err := broadcastTx(methodTypes.InvokeMethod, image, payload, peer, privPEM)
			if err != nil {
				return errw(err)
			}

			log.Printf("Broadcasted tx hash: %s", txHash)

			return nil
		},
	}

	invokemethodcmd.Flags().StringVarP(&peer, "peer", "p", "/ip4/0.0.0.0/tcp/3330", "The host to broadcast to transaction to")
	invokemethodcmd.Flags().StringVarP(&payload, "payload", "d", "", "The payload to send")
	invokemethodcmd.Flags().StringVarP(&privPEM, "priv", "k", "", "The private key to sign the transaction with")
	invokemethodcmd.Flags().StringVarP(&image, "image", "i", "", "The image hash to deploy")

	return invokemethodcmd
}
