package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/c3systems/c3/ditto"
	"github.com/go-openapi/errors"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	dittoSvc := ditto.New(&ditto.Config{})

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
		Args: func(cmd *cobra.Command, args []string) error {
			require(len(args) != 0, "image hash or name is required")
			require(len(args) == 1, "only one argument is required")
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			must(dittoSvc.PushImageByID(args[0]))
			fmt.Println("success")
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

func require(cond bool, err string) {
	if !cond {
		logFatal(err)
	}
}

func must(err error) {
	if err != nil {
		logFatal(err)
	}
}

func logFatal(ierr interface{}) {
	switch v := ierr.(type) {
	case errors.Error:
		fmt.Println(v)
	case string:
		fmt.Println(v)
	}
	os.Exit(1)
}
