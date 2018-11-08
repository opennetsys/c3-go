package cmd

import (
	methodTypes "github.com/c3systems/c3-go/core/types/methods"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func deployCmd() *cobra.Command {
	var (
		peer    string
		image   string
		genesis string
		privPEM string
	)

	deploycmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy image to the blockchain",
		Long:  "Deploys the docker image to the decentralized registry on IPFS and broadcasts a transaction to the blockchain for it to be mined",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("pem: %s", privPEM)
			log.Printf("image: %s", image)
			log.Printf("peer: %s", peer)
			log.Printf("genesis: %s", genesis)

			txHash, err := broadcastTx(methodTypes.Deploy, image, genesis, peer, privPEM)
			if err != nil {
				return errw(err)
			}

			log.Printf("Broadcasted tx hash: %s", txHash)

			return nil
		},
	}

	deploycmd.Flags().StringVarP(&peer, "peer", "p", "/ip4/0.0.0.0/tcp/3330", "The host to broadcast to transaction to")
	deploycmd.Flags().StringVarP(&genesis, "genesis", "g", "", "The genesis data for the dApp")
	deploycmd.Flags().StringVarP(&privPEM, "priv", "k", "", "The private key to sign the transaction with")
	deploycmd.Flags().StringVarP(&image, "image", "i", "", "The image hash to deploy")

	return deploycmd
}
