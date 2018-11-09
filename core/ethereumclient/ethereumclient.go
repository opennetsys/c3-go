package ethereumclient

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	contract "github.com/c3systems/c3-go/core/ethereumclient/contracts"
)

// CheckpointConfig ...
type CheckpointConfig struct {
	URL             string
	PrivateKey      string
	ContractAddress string
	MethodName      string
}

// CheckpointClient ...
type CheckpointClient struct {
	auth     *bind.TransactOpts
	client   *ethclient.Client
	instance *contract.Checkpoint
}

// NewCheckpointClient ...
func NewCheckpointClient(config *CheckpointConfig) *CheckpointClient {
	client, err := ethclient.Dial(config.URL)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(config.ContractAddress)
	instance, err := contract.NewCheckpoint(address, client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	return &CheckpointClient{
		auth:     auth,
		client:   client,
		instance: instance,
	}
}

// CheckpointRoot ...
func (s *CheckpointClient) CheckpointRoot(root string) (string, error) {
	rootBig := new(big.Int)
	rootBig.SetString(root[2:], 16)

	tx, err := s.instance.Checkpoint(s.auth, rootBig)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
