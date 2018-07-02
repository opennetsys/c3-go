package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "c3",
		Short: "C3 command line interface",
		Long: `The command line interface for C3
For more info visit: https://github.com/c3systems/c3,
		`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy image to registry",
		Long: `Deploys the docker image to the decentralized registry on IPFS
		`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO")
		},
	}

	rootCmd.AddCommand(deployCmd)
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
