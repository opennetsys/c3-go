package cmd

import (
	"log"
	"os"

	"github.com/c3systems/c3/ditto"
	"github.com/c3systems/c3/node"
	nodetypes "github.com/c3systems/c3/node/types"
	"github.com/go-openapi/errors"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// Build ...
// note: don't want to use init bc it will init even for tests, benchmarks, etc!
func Build() {
	var (
		nodeURI string
		dataDir string
	)

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

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push image to registry",
		Long: `Push the docker image to the decentralized registry on IPFS
		`,
		Args: func(cmd *cobra.Command, args []string) error {
			require(len(args) != 0, "image hash or name is required")
			require(len(args) == 1, "only one argument is required")
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			must(dittoSvc.PushImageByID(args[0]))
			log.Println("success")
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull image from registry",
		Long: `Pull the docker image from the decentralized registry on IPFS
		`,
		Args: func(cmd *cobra.Command, args []string) error {
			require(len(args) != 0, "image hash or name is required")
			require(len(args) == 1, "only one argument is required")
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			must(dittoSvc.PullImage(args[0], "", ""))
			log.Println("success")
		},
	}

	nodeCmd := &cobra.Command{
		Use:              "node [OPTIONS] [COMMANDS] [OPTIONS]",
		Short:            "c3 node commands",
		TraverseChildren: true,
	}

	startSubCmd := &cobra.Command{
		Use:   "start [OPTIONS]",
		Short: "Start a c3 node",
		Long:  "By starting a c3 node, you will participate in the c3 network: mining and storing blocks and responding to RPC requests. Thank you, you are making the c3 network stronger by participating.",
		Run: func(cmd *cobra.Command, args []string) {
			must(node.Start(&nodetypes.CFG{
				URI:     nodeURI,
				DataDir: dataDir,
			}))
		},
	}
	startSubCmd.Flags().StringVarP(&nodeURI, "uri", "u", "/ip4/127.0.0.1/tcp/9000/ipfs/QmdRa9h1mAxthj4ACrHULZC5yQmuiHzXDV56rWvnQaMA9o", "The host on which to run the node")
	//startSubCmd.MarkFlagRequired("uri")
	startSubCmd.Flags().StringVarP(&dataDir, "data-dir", "d", "~/c3-data/", "The directory in which to save data")
	//startSubCmd.MarkFlagRequired("data-dir")
	// TODO: add more flags for blockstore and nodestore, etc.

	nodeCmd.AddCommand(startSubCmd)
	rootCmd.AddCommand(pushCmd, pullCmd, nodeCmd)
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
		log.Println(v)
	case string:
		log.Println(v)
	//case *errors.errorString:
	//log.Println(v)
	default:
		log.Printf("%T\n%v", v, ierr)
	}
	os.Exit(1)
}
