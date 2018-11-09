package cmd

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/config"
	"github.com/c3systems/c3-go/core/eosclient"
	"github.com/c3systems/c3-go/core/ethereumclient"
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
	// ErrSubCommandRequired ...
	ErrSubCommandRequired = errors.New("sub command is required")
)

// Build ...
func Build() *cobra.Command {
	var (
		configPath              string
		nodeURI                 string
		dataDir                 string
		peer                    string
		pem                     string
		password                string
		outputPath              string
		dockerLocalRegistryHost string
		mempoolType             string
		rpcHost                 string
		ipfsHost                string
		blockDifficulty         int

		eosURL         string
		eosWifPrivKey  string
		eosAccountName string
		eosActionName  string
		eosPermissions string

		ethereumURL             string
		ethereumPrivKey         string
		ethereumContractAddress string
		ethereumMethodName      string
	)

	cnf := config.New()

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
		Long:  "Push the docker image to the decentralized registry on IPFS",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errw(ErrImageIDRequired)
			}
			if len(args) != 1 {
				return errw(ErrOnlyOneArgumentRequired)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := registry.NewRegistry(&registry.Config{
				DockerLocalRegistryHost: dockerLocalRegistryHost,
				IPFSHost:                ipfsHost,
			})

			hash, err := reg.PushImageByID(args[0])
			if err != nil {
				return errw(err)
			}

			log.Printf("[cli] %s", hash)
			return nil
		},
	}

	pushCmd.Flags().StringVarP(&ipfsHost, "ipfs-host", "", "", "A remote IPFS API host to push the image to. Example: 127.0.0.1:5001")

	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull image from registry",
		Long:  "Pull the docker image from the decentralized registry on IPFS",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errw(ErrImageIDRequired)
			}
			if len(args) != 1 {
				return errw(ErrOnlyOneArgumentRequired)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := registry.NewRegistry(&registry.Config{
				DockerLocalRegistryHost: dockerLocalRegistryHost,
			})
			_, err := reg.PullImage(args[0])
			return errw(err)
		},
	}

	pullCmd.Flags().StringVarP(&dockerLocalRegistryHost, "host", "", "", "Docker local registry host")

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
			// use settings from config file if config filepath is passed as argument to CLI
			if configPath != "" {
				cnf := config.NewFromFile(configPath)
				nodeURI = cnf.NodeURI()
				dataDir = cnf.DataDir()
				pem = cnf.PrivateKeyPath()
				peer = cnf.Peer()
				blockDifficulty = cnf.BlockDifficulty()
			}

			if _, err := os.Stat(pem); os.IsNotExist(err) {
				return errw(fmt.Errorf("%s does not exist", pem))
			}

			var eosClient *eosclient.CheckpointClient
			if eosWifPrivKey != "" {
				eosClient = eosclient.NewCheckpointClient(&eosclient.CheckpointConfig{
					URL:           eosURL,
					WifPrivateKey: eosWifPrivKey,
					AccountName:   eosAccountName,
					ActionName:    eosActionName,
					Permissions:   eosPermissions,
					Debug:         false,
				})
			}

			var ethClient *ethereumclient.CheckpointClient
			if ethereumPrivKey != "" {
				ethClient = ethereumclient.NewCheckpointClient(&ethereumclient.CheckpointConfig{
					URL:             ethereumURL,
					PrivateKey:      ethereumPrivKey,
					ContractAddress: ethereumContractAddress,
					MethodName:      ethereumMethodName,
				})
			}

			n, err := node.NewFullNode(&nodetypes.Config{
				URI:     nodeURI,
				Peer:    peer,
				DataDir: dataDir,
				Keys: nodetypes.Keys{
					PEMFile:  pem,
					Password: password,
				},
				BlockDifficulty: blockDifficulty,
				MempoolType:     mempoolType,
				RPCHost:         rpcHost,
				EOSClient:       eosClient,
				EthereumClient:  ethClient,
			})
			if err != nil {
				return errw(err)
			}

			if err = n.Start(); err != nil {
				return errw(err)
			}

			return nil
		},
	}

	startSubCmd.Flags().StringVarP(&configPath, "config", "c", "", "filepath of the config file to use [OPTIONAL]")
	startSubCmd.Flags().StringVarP(&nodeURI, "uri", "u", "/ip4/0.0.0.0/tcp/9000", "The host on which to run the node")
	startSubCmd.Flags().StringVarP(&peer, "peer", "p", cnf.Peer(), "A peer to which to connect")
	startSubCmd.Flags().StringVarP(&dataDir, "data-dir", "d", cnf.DataDir(), "The directory in which to save data")
	startSubCmd.Flags().StringVar(&pem, "pem", cnf.PrivateKeyPath(), "A pem file containing an ecdsa private key")
	startSubCmd.Flags().StringVar(&password, "password", "", "A password for the pem file [OPTIONAL]")
	startSubCmd.Flags().StringVar(&mempoolType, "mempool-type", "memory", "The mempool type to use (memory, redis) [OPTIONAL]")
	startSubCmd.Flags().StringVarP(&rpcHost, "rpc", "", "0.0.0.0:5005", "The port to run rpc on")
	startSubCmd.Flags().IntVar(&blockDifficulty, "difficulty", cnf.BlockDifficulty(), "The hashing difficulty for mining blocks. (1-15) [OPTIONAL]. This feature will be deprecated when C3 soon moves to Delegated Proof-of-Stake.")

	startSubCmd.Flags().StringVarP(&eosURL, "checkpoint-eos-url", "", "", "EOS block producer URL for checkpointing")
	startSubCmd.Flags().StringVarP(&eosWifPrivKey, "checkpoint-eos-wif-private-key", "", "", "EOS private key for EOS account that will be used for checkpointing")
	startSubCmd.Flags().StringVarP(&eosAccountName, "checkpoint-eos-account-name", "", "", "EOS account name that will be used for checkpointing")
	startSubCmd.Flags().StringVarP(&eosActionName, "checkpoint-eos-action-name", "", "", "EOS action name for checkpointing")
	startSubCmd.Flags().StringVarP(&eosPermissions, "checkpoint-eos-action-permissions", "", "", `EOS action permissions for checkpointing. eg. "myaccount@active"`)

	startSubCmd.Flags().StringVarP(&ethereumURL, "checkpoint-ethereum-url", "", "", "Ethereum host node URL for checkpointing")
	startSubCmd.Flags().StringVarP(&ethereumPrivKey, "checkpoint-ethereum-private-key", "", "", "Ethereum private key of account for checkpointing")
	startSubCmd.Flags().StringVarP(&ethereumContractAddress, "checkpoint-ethereum-contract-address", "", "", "Ethereum smart contract address for checkpointing")
	startSubCmd.Flags().StringVarP(&ethereumMethodName, "checkpoint-ethereum-method-name", "", "", "Ethereum smart contract method for checkpointing")

	// TODO: add more flags for blockstore and nodestore, etc.
	nodeCmd.AddCommand(startSubCmd)

	generateCmd := &cobra.Command{
		Use:   "generate [OPTIONS] [COMMANDS]",
		Short: "Generate command",
		Long:  "Generate command requires a sub command",
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

	generateKeyCmd := &cobra.Command{
		Use:   "key",
		Short: "Generate new private key",
		Long:  "Generates a new private key in PEM format",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			priv, err := c3crypto.NewPrivateKey()
			if err != nil {
				return errw(err)
			}

			var pwd *string
			if password == "" {
				log.Println("creating pem file...")
			} else {
				pwd = &password
				log.Println("creating pem file with password...")
			}

			if err := c3crypto.WritePrivateKeyToPemFile(priv, pwd, outputPath); err != nil {
				return errw(err)
			}

			log.Printf("private key saved in %s", outputPath)
			return nil
		},
	}

	generateKeyCmd.Flags().StringVarP(&outputPath, "output", "o", "priv.pem", "Output file path")
	generateKeyCmd.Flags().StringVarP(&password, "password", "p", "", "Password for private key")
	startSubCmd.MarkFlagRequired("output")
	generateCmd.AddCommand(generateKeyCmd)

	rootCmd.AddCommand(pushCmd, pullCmd, nodeCmd, generateCmd, deployCmd(), invokeMethodCmd(), encodeCmd(), peerCmd(), snapshotCmd())

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
