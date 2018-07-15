package miner

import (
	"errors"
	"log"
	"sort"

	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/merkle"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/p2p"
)

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
	pub, err := c3crypto.DecodeAddress(tx.Props().From)
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
	pub, err := c3crypto.DecodeAddress(block.Props().MinerAddress)
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

// VerifyMinedBlock ...
func VerifyMinedBlock(p2pSvc p2p.Interface, isValid *bool, minedBlock *MinedBlock) (bool, error) {
	if minedBlock == nil {
		return false, nil
	}
	if minedBlock.NextBlock == nil {
		return false, nil
	}
	if minedBlock.PreviousBlock == nil {
		return false, nil
	}
	if minedBlock.NextBlock.Props().BlockHash == nil {
		return false, nil
	}
	if minedBlock.PreviousBlock.Props().BlockHash == nil {
		return false, nil
	}
	// note checked for nil pointer, above
	if *minedBlock.PreviousBlock.Props().BlockHash != minedBlock.NextBlock.Props().PrevBlockHash {
		return false, nil
	}
	if isValid == nil {
		return false, errors.New("IsValid is nil")
	}
	if mainchain.ImageHash != minedBlock.NextBlock.Props().ImageHash {
		return false, nil
	}
	if minedBlock.NextBlock.Props().MinerSig == nil {
		return false, nil
	}

	ok, err := CheckBlockHashAgainstDifficulty(minedBlock.NextBlock)
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
	tmpHash, err := minedBlock.NextBlock.CalculateHash()
	if err != nil {
		return false, err
	}
	// note: already checked for nil hash
	if *minedBlock.NextBlock.Props().BlockHash != tmpHash {
		return false, nil
	}

	// the sig must verify
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}
	pub, err := c3crypto.DecodeAddress(minedBlock.NextBlock.Props().MinerAddress)
	if err != nil {
		return false, err
	}

	// note: checked for nil sig, above
	sigR, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.R)
	if err != nil {
		return false, err
	}
	sigS, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.S)
	if err != nil {
		return false, err
	}

	// note: nil blockhash was checked, above
	ok, err = c3crypto.Verify(pub, []byte(*minedBlock.NextBlock.Props().BlockHash), sigR, sigS)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	return VerifyStateBlocksFromMinedBlock(p2pSvc, isValid, minedBlock)
}

// VerifyStateBlocksFromMinedBlock ...
// note: this function also checks the merkle tree. That check is not required to be performed, separately.
func VerifyStateBlocksFromMinedBlock(p2pSvc p2p.Interface, isValid *bool, minedBlock *MinedBlock) (bool, error) {
	if minedBlock.NextBlock == nil {
		return false, nil
	}
	if minedBlock.StatechainBlocksMap == nil {
		return false, nil
	}
	if minedBlock.MerkleTreesMap == nil {
		return false, nil
	}
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}

	// 1. Verify state blocks merkle hash
	if ok, err := VerifyMerkleTreeFromMinedBlock(isValid, minedBlock); !ok || err != nil {
		return false, nil
	}

	// 2. Verify each state block
	// first, order them by block number
	orderedBlocks, err := orderStatechainBlocks(minedBlock.StatechainBlocksMap)
	if err != nil {
		return false, err
	}
	if len(orderedBlocks) == 0 {
		if len(minedBlock.TransactionsMap) == 0 {
			return true, nil
		}

		return false, nil
	}

	prevBlockHash := orderedBlocks[0].Props().PrevBlockHash
	prevBlockCID, err := p2p.GetCIDByHash(prevBlockHash)
	if err != nil {
		return false, err
	}
	// TODO: check that this is the actual prev block on the blockchain
	prevBlock, err := p2pSvc.GetStatechainBlock(prevBlockCID)
	if err != nil {
		return false, err
	}
	if prevBlock == nil {
		return false, errors.New("got nil prev block")
	}
	if prevBlock.Props().BlockHash == nil {
		return false, errors.New("got nil prev block hash")
	}
	// note: checked for nil pointer, above
	if *prevBlock.Props().BlockHash != prevBlockHash {
		return false, nil
	}

	// TODO: walk the blockchain to generate current state
	prevState := ""
	prevStateHash := hashing.HashToHexString([]byte(prevState))
	if prevStateHash != prevBlock.Props().StateCurrentHash {
		return false, nil
	}

	for _, block := range orderedBlocks {
		if isValid == nil || *isValid == false {
			return false, errors.New("received nil or false isValid")
		}

		// 2a. block must have a hash
		if block == nil || block.Props().BlockHash == nil {
			return false, nil
		}

		// 2b. Block #'s must be sequential
		prevBlockNumber, err := hexutil.DecodeUint64(prevBlock.Props().BlockNumber)
		if err != nil {
			return false, err
		}
		blockNumber, err := hexutil.DecodeUint64(block.Props().BlockNumber)
		if err != nil {
			return false, err
		}
		if prevBlockNumber+1 != blockNumber {
			return false, nil
		}

		// 2c. verify the block hash
		tmpHash, err := block.CalculateHash()
		if err != nil {
			return false, err
		}
		// note: checked nil BlockHash, above
		if tmpHash != *block.Props().BlockHash {
			return false, nil
		}

		// 2d. verify the block tx
		// note: can't have a state block without transactions?
		tx, ok := minedBlock.TransactionsMap[block.Props().TxHash]
		if !ok || tx == nil {
			return false, err
		}

		ok, err = VerifyTransaction(tx)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}

		// note: just printing to keep the txs var alive
		log.Println(tx)

		// TODO: run the tx through the container
		// 2e. verify current state hash
		currentState := ""
		currentStateHash := hashing.HashToHexString([]byte(currentState))
		if currentStateHash != block.Props().StateCurrentHash {
			return false, nil
		}

		// 2f. verify prevDiff
		prevDiff := ""
		prevDiffHash := hashing.HashToHexString([]byte(prevDiff))
		if prevDiffHash != block.Props().StatePrevDiffHash {
			return false, nil
		}

		// set prev to current for next loop
		prevState = currentState
		prevBlock = block
	}

	return true, nil

}

// VerifyMerkleTreeFromMinedBlock ...
func VerifyMerkleTreeFromMinedBlock(isValid *bool, minedBlock *MinedBlock) (bool, error) {
	if minedBlock.NextBlock == nil {
		return false, nil
	}
	if minedBlock.StatechainBlocksMap == nil {
		return false, nil
	}
	if minedBlock.MerkleTreesMap == nil {
		return false, nil
	}
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}

	tree, ok := minedBlock.MerkleTreesMap[minedBlock.NextBlock.Props().StateBlocksMerkleHash]
	if !ok || tree == nil {
		return false, nil
	}
	if tree.Props().MerkleTreeRootHash == nil {
		return false, nil
	}
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}

	tmpTree, err := merkle.New(&merkle.TreeProps{
		Hashes: tree.Props().Hashes,
		Kind:   merkle.StatechainBlocksKindStr,
	})
	if err != nil {
		return false, err
	}
	if err := tmpTree.SetHash(); err != nil {
		return false, err
	}

	if *tmpTree.Props().MerkleTreeRootHash != *tree.Props().MerkleTreeRootHash {
		return false, nil
	}

	if len(tmpTree.Props().Hashes) != len(minedBlock.StatechainBlocksMap) {
		return false, nil
	}

	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}
	for _, hash := range tmpTree.Props().Hashes {
		statechainBlock, ok := minedBlock.StatechainBlocksMap[hash]
		if !ok || statechainBlock == nil {
			return false, nil
		}

		tmpHash, err := statechainBlock.CalculateHash()
		if err != nil {
			return false, err
		}
		if hash != tmpHash {
			return false, nil
		}
	}

	return true, nil
}

// note: private types for the sort function, below
type blockWrapper struct {
	Block       *statechain.Block
	BlockNumber uint64
}
type byBlockNumber []blockWrapper

func (b byBlockNumber) Len() int           { return len(b) }
func (b byBlockNumber) Less(i, j int) bool { return b[i].BlockNumber < b[j].BlockNumber }
func (b byBlockNumber) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func orderStatechainBlocks(stateBlocksMap map[string]*statechain.Block) ([]*statechain.Block, error) {
	var (
		stateBlocks   []*statechain.Block
		blockWrappers byBlockNumber
	)

	for _, stateBlock := range stateBlocksMap {
		if stateBlock == nil {
			return nil, errors.New("nil sate block")
		}

		blockNumber, err := hexutil.DecodeUint64(stateBlock.Props().BlockNumber)
		if err != nil {
			return nil, err
		}

		blockWrappers = append(blockWrappers, blockWrapper{
			Block:       stateBlock,
			BlockNumber: blockNumber,
		})
	}

	sort.Sort(blockWrappers)

	for _, blockWrapper := range blockWrappers {
		stateBlocks = append(stateBlocks, blockWrapper.Block)
	}

	return stateBlocks, nil
}
