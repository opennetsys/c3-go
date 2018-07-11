package miner

import (
	"crypto/ecdsa"
	"errors"
	"log"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
)

//// NewFromStateBlocks ...
//func NewFromStateBlocks(stateBlocks []*statechain.Block) (*mainchain.Block, error) {
//var (
//list                  []merkletree.Content
//statechainBlockHashes []*string
//)

//for _, statechainBlock := range stateBlocks {
//list = append(list, statechainBlock)
//statechainBlockHashes = append(statechainBlockHashes, statechainBlock.Props().BlockHash)
//}

//t, err := merkletree.NewTree(list)
//if err != nil {
//return nil, err
//}

//// Get the Merkle Root of the tree
//mr := t.MerkleRoot()

//// note: the other missing fields are added by the miner
//return mainchain.New(&mainchain.Props{
//StateBlocksMerkleHash: string(mr),
//StateBlockHashes:      statechainBlockHashes,
//}), nil
//}

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

	difficulty, err := hexutil.DecodeUint64(block.Props().Difficulty)
	if err != nil {
		return false, err
	}

	return CheckHashAgainstDifficulty(*block.Props().BlockHash, difficulty)
}

// CheckHashAgainstDifficulty ...
func CheckHashAgainstDifficulty(hashHex string, difficulty uint64) (bool, error) {
	// TODO: check the difficulty is correct?

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

// VerifyStatechainBlock verifies a statechain block
// TODO: check timestamp?
// TODO: pass all necessary data and remove the p2pSvc
func VerifyStatechainBlock(p2pSvc p2p.Interface, isValid *bool, block *statechain.Block) (bool, error) {
	if block == nil {
		return false, ErrNilBlock
	}
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}

	// 1. block must have a hash
	if block.Props().BlockHash == nil {
		return false, nil
	}

	// TODO: check the block # and StatePrevDiffHash

	// 2. verify the block hash
	tmpHash, err := block.CalculateHash()
	if err != nil {
		return false, err
	}
	// note: checked nil BlockHash, above
	if tmpHash != *block.Props().BlockHash {
		return false, nil
	}

	// 3. verify each tx in the block
	// TODO: do in go funcs
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}
	txCid, err := p2p.GetCIDByHash(block.Props().TxHash)
	if err != nil {
		return false, err
	}

	tx, err := p2pSvc.GetStatechainTransaction(txCid)
	if err != nil {
		return false, err
	}

	ok, err := VerifyTransaction(tx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	// note: just printing to keep the txs var alive
	log.Println(tx)

	// 4. run the tx through the container
	// TODO: step #4

	// 5. verify the statehash and prev diff hash
	// TODO step #5

	return true, nil
}

// VerifyMainchainBlock verifies a mainchain block
// TODO: check block time
// TODO: fetch and check previous block hash
// TODO: pass all necessary data and remove the p2pSvc
func VerifyMainchainBlock(p2pSvc p2p.Interface, isValid *bool, block *mainchain.Block) (bool, error) {
	if block == nil {
		return false, errors.New("block is nil")
	}

	if block.Props().BlockHash == nil {
		return false, errors.New("block hash is nil")
	}

	if isValid == nil {
		return false, errors.New("IsValid is nil")
	}

	if mainchain.ImageHash != block.Props().ImageHash {
		return false, nil
	}

	if block.Props().MinerSig == nil {
		return false, nil
	}

	ok, err := CheckBlockHashAgainstDifficulty(block)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	// hash must verify
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}
	tmpHash, err := block.CalculateHash()
	if err != nil {
		return false, err
	}
	// note: already checked for nil hash
	if *block.Props().BlockHash != tmpHash {
		return false, nil
	}

	// the sig must verify
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}
	pub, err := PubFromBlock(block)
	if err != nil {
		return false, err
	}

	// note: checked for nil sig, above
	sigR, err := hexutil.DecodeBigInt(block.Props().MinerSig.R)
	if err != nil {
		return false, err
	}
	sigS, err := hexutil.DecodeBigInt(block.Props().MinerSig.S)
	if err != nil {
		return false, err
	}

	// note: nil blockhash was checked, above
	ok, err = c3crypto.Verify(pub, []byte(*block.Props().BlockHash), sigR, sigS)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	treeCID, err := p2p.GetCIDByHash(block.Props().StateBlocksMerkleHash)
	if err != nil {
		return false, nil
	}

	tree, err := p2pSvc.GetMerkleTree(treeCID)
	if err != nil {
		return false, err
	}

	// TODO: do in go funcs
	// TODO: check kind?
	for _, stateblockHash := range tree.Props().Hashes {
		if isValid == nil || *isValid == false {
			return false, errors.New("received nil or false isValid")
		}

		stateblockCid, err := p2p.GetCIDByHash(stateblockHash)
		if err != nil {
			return false, err
		}

		stateblock, err := p2pSvc.GetStatechainBlock(stateblockCid)
		if err != nil {
			return false, err
		}

		ok, err := VerifyStatechainBlock(p2pSvc, isValid, stateblock)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}
