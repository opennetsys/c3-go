package cmd

import (
	"fmt"

	"github.com/c3systems/c3-go/common/netutil"
	"github.com/c3systems/c3-go/config"
	"github.com/spf13/cobra"
)

func peerCmd() *cobra.Command {
	peercmd := &cobra.Command{
		Use:   "peer",
		Short: "Peer command",
		Long:  "Peer command requires a sub command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errw(ErrSubCommandRequired)
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
			cnf := config.New()
			ip, err := netutil.LocalIP()
			if err != nil {
				return errw(err)
			}

			fmt.Printf("Your Peer ID:\n/ip4/%s/tcp/%v/%s", ip.String(), cnf.Port(), cnf.PrivateKeyIPNS())

			return nil
		},
	}

	peercmd.AddCommand(peerIDCmd)

	return peercmd
}
