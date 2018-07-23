package cmd

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/config"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	"github.com/c3systems/c3-go/node"
	nodetypes "github.com/c3systems/c3-go/node/types"
	"github.com/c3systems/c3-go/registry"
	"github.com/spf13/cobra"
)

var (
	// ErrImageIDRequired ...
	ErrImageIDRequired = errors.New("image hash or name is required")
	// ErrOnlyOneArgumentRequired ...
	ErrOnlyOneArgumentRequired = errors.New("only one argument is required")
)

// Build ...
func Build() *cobra.Command {
	var (
		nodeURI    string
		dataDir    string
		peer       string
		pem        string
		password   string
		difficulty int
	)

	dit := registry.NewRegistry(&registry.Config{})

	rootCmd := &cobra.Command{
		Use:   "c3",
		Short: "C3 command line interface",
		Long: `The command line interface for C3
For more info visit: https://github.com/c3systems/c3-go,
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
			if len(args) == 0 {
				return ErrImageIDRequired
			}
			if len(args) != 1 {
				return ErrOnlyOneArgumentRequired
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			hash, err := dit.PushImageByID(args[0])
			if err != nil {
				return err
			}

			log.Printf("[cli] %s", hash)
			return nil
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull image from registry",
		Long: `Pull the docker image from the decentralized registry on IPFS
		`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrImageIDRequired
			}
			if len(args) != 1 {
				return ErrOnlyOneArgumentRequired
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := dit.PullImage(args[0])
			return err
		},
	}

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy image to the blockchain",
		Long:  `Deploys the docker image to the decentralized registry on IPFS and broadcasts a transaction to the blockchain for it to be mined`,
		Args: func(cmd *cobra.Command, args []string) error {
			// c3 deploy {imageID} "foo" "bar"

			if len(args) == 0 {
				return ErrImageIDRequired
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// broadcast transaction as grpc

			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			n := new(node.Service)
			return node.Start(n, &nodetypes.Config{
				URI:     nodeURI,
				Peer:    peer,
				DataDir: dataDir,
				Keys: nodetypes.Keys{
					PEMFile:  pem,
					Password: password,
				},
				BlockDifficulty: difficulty,
			})
		},
	}
	startSubCmd.Flags().StringVarP(&nodeURI, "uri", "u", "/ip4/0.0.0.0/tcp/9000", "The host on which to run the node")
	startSubCmd.Flags().StringVarP(&peer, "peer", "p", "", "A peer to which to connect")
	startSubCmd.Flags().StringVarP(&dataDir, "data-dir", "d", config.DefaultStoreDirectory, "The directory in which to save data")
	startSubCmd.Flags().StringVar(&pem, "pem", "", "A pem file containing an ecdsa private key")
	startSubCmd.Flags().StringVar(&password, "password", "", "A password for the pem file [OPTIONAL]")
	startSubCmd.Flags().IntVar(&difficulty, "difficulty", 6, "The hashing difficulty for mining blocks. (1-15) [OPTIONAL]")

	startSubCmd.MarkFlagRequired("pem")
	// TODO: add more flags for blockstore and nodestore, etc.

	nodeCmd.AddCommand(startSubCmd)
	rootCmd.AddCommand(pushCmd, pullCmd, deployCmd, nodeCmd)

	return rootCmd
}

// Execute ...
func Execute() {
	rootCmd := Build()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
