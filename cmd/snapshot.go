package cmd

import (
	"fmt"

	colorlog "github.com/c3systems/c3-go/log/color"
	"github.com/c3systems/c3-go/node"
	nodetypes "github.com/c3systems/c3-go/node/types"
	"github.com/c3systems/c3-go/snapshot"
	"github.com/spf13/cobra"
)

func snapshotCmd() *cobra.Command {
	var (
		image            string
		privPEM          string
		stateBlockNumber int
	)

	snapshotcmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Snaphost image with state to local registry",
		Long:  "Snapshot image with state to local registry",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeURI := "/ip4/0.0.0.0/tcp/3330"
			dataDir := ".tmp"
			peer := ""

			n, err := node.NewFullNode(&nodetypes.Config{
				URI:     nodeURI,
				Peer:    peer,
				DataDir: dataDir,
				Keys: nodetypes.Keys{
					PEMFile:  privPEM,
					Password: "",
				},
				//BlockDifficulty: 5,
				MempoolType: "memory",
				RPCHost:     ":5005",
			})
			if err != nil {
				return errw(err)
			}

			svc := snapshot.New(&snapshot.Config{
				P2P:     n.Props().P2P,
				Mempool: n.Props().Store,
			})

			snapshotImageID, err := svc.Snapshot(image, stateBlockNumber)
			if err != nil {
				return err
			}

			fmt.Printf(colorlog.Green("snapshot image ID: %s", snapshotImageID))

			return nil
		},
	}

	snapshotcmd.Flags().StringVarP(&privPEM, "priv", "k", "", "The private key of the host")
	snapshotcmd.Flags().StringVarP(&image, "image", "i", "", "The image to snapshot")
	snapshotcmd.Flags().IntVarP(&stateBlockNumber, "stateblock", "b", 0, "The state block number to snapshot image at")

	return snapshotcmd
}
