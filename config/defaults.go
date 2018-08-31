package config

import (
	"fmt"
	"time"
)

// ServerHost ...
const ServerHost = "0.0.0.0"

// DefaultServerPort ...
const DefaultServerPort = 3333

// DefaultConfigDirectory is the default directory where the config settings will live
var DefaultConfigDirectory = "~/.c3"

// DefaultConfigFilename is the default filename for the config settings
var DefaultConfigFilename = "config.toml"

// DefaultStoreDirectory is the default directory where the file system store will live.
var DefaultStoreDirectory = "~/.c3/chaindata"

// DefaultPrivateKeyFilename is the filename for the account's private key
var DefaultPrivateKeyFilename = "priv.pem"

// TempContainerStatePath ...
var TempContainerStatePath = "/tmp"

// TempContainerStateFileName ...
var TempContainerStateFileName = "state.json"

// TempContainerStateFilePath ...
var TempContainerStateFilePath = fmt.Sprintf("%s/%s", TempContainerStatePath, TempContainerStateFileName)

// DockerRegistryPort ...
const DockerRegistryPort = 5000

// DefaultBlockDifficulty ...
const DefaultBlockDifficulty = 6

// MinedBlockVerificationTimeout ...
const MinedBlockVerificationTimeout = 10 * time.Minute

// IPFSTimeout ...
const IPFSTimeout = 20 * time.Second
