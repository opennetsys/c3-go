package cmd

import (
	"github.com/spf13/cobra"
)

func snapshotCmd() *cobra.Command {
	// TODO
	snapshotcmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Snaphost image with state to local registry",
		Long:  "Snapshot image with state to local registry",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	return snapshotcmd
}
