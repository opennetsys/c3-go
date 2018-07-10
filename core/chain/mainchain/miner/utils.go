package miner

import (
	"crypto/ecdsa"
	"errors"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"

	"github.com/c3systems/merkletree"
)

// NewFromStateBlocks ...
func NewFromStateBlocks(stateBlocks []*statechain.Block) (*mainchain.Block, error) {
	var (
		list                  []merkletree.Content
		statechainBlockHashes []*string
	)

	for _, statechainBlock := range stateBlocks {
		list = append(list, statechainBlock)
		statechainBlockHashes = append(statechainBlockHashes, statechainBlock.Props().BlockHash)
	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	// Get the Merkle Root of the tree
	mr := t.MerkleRoot()

	// note: the other missing fields are added by the miner
	return mainchain.New(&mainchain.Props{
		ImageHash:             mainchain.ImageHash,
		StateBlocksMerkleHash: string(mr),
		StateBlockHashes:      statechainBlockHashes,
	}), nil
}

// PubFromBlock ...
func PubFromBlock(block *mainchain.Block) (*ecdsa.PublicKey, error) {
	if block == nil {
		return nil, errors.New("block is nil")
	}

	pubStr, err := hexutil.DecodeString(block.Props().MinerAddress)
	if err != nil {
		return nil, err
	}
	pub, err := c3crypto.DeserializePublicKey([]byte(pubStr))
	if err != nil {
		return nil, err
	}
	if pub == nil {
		return nil, errors.New("invalid miner address")
	}

	return pub, nil
}

// CheckBlockHashAgainstDifficulty ...
func CheckBlockHashAgainstDifficulty(block *mainchain.Block) (bool, error) {
	if block == nil {
		return false, ErrNilBlock
	}
	if block.Props().BlockHash == nil {
		return false, nil
	}

	return CheckHashAgainstDifficulty(*block.Props().BlockHash, block.Props().Difficulty)
}

// CheckHashAgainstDifficulty ...
func CheckHashAgainstDifficulty(hashHex, difficultyHex string) (bool, error) {
	difficulty, err := hexutil.DecodeUint64(difficultyHex)
	if err != nil {
		return false, err
	}

	// TODO: check the difficulty is correct

	hashStr, err := hexutil.StripLeader(hashHex)
	if err != nil {
		return false, err
	}

	if len(hashStr) <= int(difficulty) {
		return false, nil
	}

	for i := 0; i < int(difficulty); i++ {
		if hashStr[i:i+1] != "0" {
			return false, nil
		}
	}

	return true, nil
}

// BuildNextState ...
// TODO: everything...
func BuildNextState(imageHash string, transactions []*statechain.Transaction) (*statechain.Block, error) {
	// TODO: add miguel's code, here
	return nil, nil
}

// PubFromTx ...
func PubFromTx(tx *statechain.Transaction) (*ecdsa.PublicKey, error) {
	if tx == nil {
		return nil, ErrNilTx
	}

	pubStr, err := hexutil.DecodeString(tx.Props().From)
	if err != nil {
		return nil, err
	}
	pub, err := c3crypto.DeserializePublicKey([]byte(pubStr))
	if err != nil {
		return nil, err
	}
	if pub == nil {
		return nil, ErrInvalidFromAddress
	}

	return pub, nil
}

// BuildTxsMap ...
func BuildTxsMap(txs []*statechain.Transaction) statechain.TransactionsMap {
	txsMap := make(statechain.TransactionsMap)
	for _, tx := range txs {
		txsMap[tx.Props().ImageHash] = append(txsMap[tx.Props().ImageHash], tx)
	}

	return txsMap
}

// VerifyTransaction ...
func VerifyTransaction(tx *statechain.Transaction) (bool, error) {
	// note: we hash the message and then sign the hash
	// TODO: check the image hash exists?
	// TODO: check for blank inputs?
	if tx == nil {
		return false, ErrNilTx
	}

	// 1. tx must have a hash
	if tx.Props().TxHash == nil {
		return false, nil
	}

	// 2. tx must have a sig
	if tx.Props().Sig == nil {
		return false, nil
	}

	// 3. verify the hash
	tmpHash, err := tx.CalculateHash()
	if err != nil {
		return false, err
	}
	// note: already checked for nil hash
	if *tx.Props().TxHash != tmpHash {
		return false, nil
	}

	// 4. the sig must verify
	pub, err := PubFromTx(tx)
	if err != nil {
		return false, err
	}

	// note: checked for nil sig, above
	r, err := hexutil.DecodeBigInt(tx.Props().Sig.R)
	if err != nil {
		return false, err
	}
	s, err := hexutil.DecodeBigInt(tx.Props().Sig.S)
	if err != nil {
		return false, err
	}

	return c3crypto.Verify(pub, []byte(*tx.Props().TxHash), r, s)
}
