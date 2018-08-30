package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func peerCmd() *cobra.Command {
	peercmd := &cobra.Command{
		Use:   "peer",
		Short: "Peer command",
		Long:  "Peer command requires a sub command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrSubCommandRequired
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	peerIDCmd := &cobra.Command{
		Use:   "id",
		Short: "Show peer ID",
		Long:  "Shows the IPNS peer ID",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("TODO")

			return nil
		},
	}

	peercmd.AddCommand(peerIDCmd)

	return peercmd
}
