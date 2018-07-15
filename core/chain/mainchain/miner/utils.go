package miner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/merkle"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/diffing"
	"github.com/c3systems/c3/core/p2p"
	"github.com/c3systems/c3/core/sandbox"
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
		log.Printf("err calculating block hash\n%v", err)
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
		log.Printf("err getting cid for tx\n%v", err)
		return false, err
	}

	tx, err := p2pSvc.GetStatechainTransaction(txCid)
	if err != nil {
		return false, err
	}

	ok, err := VerifyTransaction(tx)
	if err != nil {
		log.Printf("err verifying tx\n%v", err)
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
		log.Printf("err checking block hash against difficulty\n%v", err)
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
		log.Printf("err calculating tmpHash\n%v", err)
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
		log.Printf("err decoding miner addr\n%v", err)
		return false, err
	}

	// note: checked for nil sig, above
	sigR, err := hexutil.DecodeBigInt(block.Props().MinerSig.R)
	if err != nil {
		log.Printf("err decoding miner sig r\n%v", err)
		return false, err
	}
	sigS, err := hexutil.DecodeBigInt(block.Props().MinerSig.S)
	if err != nil {
		log.Printf("err decoding miner sig s\n%v", err)
		return false, err
	}

	// note: nil blockhash was checked, above
	ok, err = c3crypto.Verify(pub, []byte(*block.Props().BlockHash), sigR, sigS)
	if err != nil {
		log.Printf("err verifying\n%v", err)
		return false, err
	}
	if !ok {
		return false, nil
	}

	treeCID, err := p2p.GetCIDByHash(block.Props().StateBlocksMerkleHash)
	if err != nil {
		log.Printf("[miner] err getting cid by has\n%v", err)
		return false, nil
	}

	tree, err := p2pSvc.GetMerkleTree(treeCID)
	if err != nil {
		log.Printf("[miner] err getting merkle tree\n%v", err)
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
			log.Printf("err getting cid\n%v", err)
			return false, err
		}

		stateblock, err := p2pSvc.GetStatechainBlock(stateblockCid)
		if err != nil {
			log.Printf("err getting statechain block\n%v", err)
			return false, err
		}

		ok, err := VerifyStatechainBlock(p2pSvc, isValid, stateblock)
		if err != nil {
			log.Printf("err verifying statechain block\n%v", err)
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
		log.Println("mined block is nil")
		return false, nil
	}
	if minedBlock.NextBlock == nil {
		log.Println("next block is nil")
		return false, nil
	}
	if minedBlock.PreviousBlock == nil {
		log.Println("mined block is nil")
		return false, nil
	}
	if minedBlock.NextBlock.Props().BlockHash == nil {
		log.Println("next block blockhash is nil")
		return false, nil
	}
	if minedBlock.PreviousBlock.Props().BlockHash == nil {
		log.Println("prev block block hash is nil")
		return false, nil
	}
	// note checked for nil pointer, above
	if *minedBlock.PreviousBlock.Props().BlockHash != minedBlock.NextBlock.Props().PrevBlockHash {
		log.Println("prev block block hash != next block block hash")
		return false, nil
	}
	if isValid == nil {
		return false, errors.New("IsValid is nil")
	}
	if mainchain.ImageHash != minedBlock.NextBlock.Props().ImageHash {
		log.Println("mainchain imagehash != nextblock image hash")
		return false, nil
	}
	if minedBlock.NextBlock.Props().MinerSig == nil {
		log.Println("next block miner sig is nil")
		return false, nil
	}

	ok, err := CheckBlockHashAgainstDifficulty(minedBlock.NextBlock)
	if err != nil {
		log.Printf("err checking block hash against difficulty\n%v", err)
		return false, err
	}
	if !ok {
		log.Println("block hash did not checkout against difficulty")
		return false, nil
	}

	// hash must verify
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}
	tmpHash, err := minedBlock.NextBlock.CalculateHash()
	if err != nil {
		log.Printf("err calculating hash\n%v", err)
		return false, err
	}
	// note: already checked for nil hash
	if *minedBlock.NextBlock.Props().BlockHash != tmpHash {
		log.Printf("next block hash != calced hash\n%s\n%s", *minedBlock.NextBlock.Props().BlockHash, tmpHash)
		return false, nil
	}

	// the sig must verify
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}
	pub, err := c3crypto.DecodeAddress(minedBlock.NextBlock.Props().MinerAddress)
	if err != nil {
		log.Printf("err decoding addr\n%v", err)
		return false, err
	}

	// note: checked for nil sig, above
	sigR, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.R)
	if err != nil {
		log.Printf("err decoding r\n%v", err)
		return false, err
	}
	sigS, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.S)
	if err != nil {
		log.Printf("err decoding s\n%v", err)
		return false, err
	}

	// note: nil blockhash was checked, above
	ok, err = c3crypto.Verify(pub, []byte(*minedBlock.NextBlock.Props().BlockHash), sigR, sigS)
	if err != nil {
		log.Printf("err verifying miner sig\n%v", err)
		return false, err
	}
	if !ok {
		log.Println("block hash did not checkout agains sig")
		return false, nil
	}

	// BlockNumber must be +1 prev block number
	blockNumber, err := hexutil.DecodeUint64(minedBlock.NextBlock.Props().BlockNumber)
	if err != nil {
		log.Printf("err decoding block #\n%v", err)
		return false, err
	}
	prevNumber, err := hexutil.DecodeUint64(minedBlock.PreviousBlock.Props().BlockNumber)
	if err != nil {
		log.Printf("err decoding prev block #\n%v", err)
		return false, err
	}

	if prevNumber+1 != blockNumber {
		log.Println("prevBlockNumber +1 != nextBlockNumber")
		return false, nil
	}

	return VerifyStateBlocksFromMinedBlock(p2pSvc, isValid, minedBlock)
}

// VerifyStateBlocksFromMinedBlock ...
// note: this function also checks the merkle tree. That check is not required to be performed, separately.
func VerifyStateBlocksFromMinedBlock(p2pSvc p2p.Interface, isValid *bool, minedBlock *MinedBlock) (bool, error) {
	if minedBlock.NextBlock == nil {
		log.Println("nil next block")
		return false, nil
	}
	if len(minedBlock.StatechainBlocksMap) != len(minedBlock.TransactionsMap) {
		log.Println("len state blocks map != len tx map")
		return false, nil
	}
	// note: ok to have nil map?
	//if minedBlock.StatechainBlocksMap == nil {
	//log.Println("nil state blocks map")
	//return false, nil
	//}
	if minedBlock.MerkleTreesMap == nil {
		log.Println("nil merkle trees map")
		return false, nil
	}
	if isValid == nil || *isValid == false {
		return false, errors.New("received nil or false isValid")
	}

	// 1. Verify state blocks merkle hash
	ok, err := VerifyMerkleTreeFromMinedBlock(isValid, minedBlock)
	if err != nil {
		log.Printf("[miner] err verifying merkle tree\n%v", err)
		return false, err
	}

	if !ok {
		log.Println("[miner] merkle tree didn't verify")
		return false, nil
	}

	// 2. Verify each state block
	// first, group them by image hash
	groupedBlocks, err := groupStateBlocksByImageHash(minedBlock.StatechainBlocksMap)
	if err != nil {
		log.Printf("err grouping state blocks\n%v", err)
		return false, err
	}

	for _, blocks := range groupedBlocks {
		if isValid == nil || *isValid == false {
			return false, errors.New("received nil or false isValid")
		}

		// then, order them by block number
		orderedBlocks, err := orderStatechainBlocks(blocks)
		if err != nil {
			log.Printf("err ordering state blocks\n%v", err)
			return false, err
		}

		if len(orderedBlocks) == 0 {
			continue
		}
		prevBlockHash := orderedBlocks[0].Props().PrevBlockHash
		prevBlockCID, err := p2p.GetCIDByHash(prevBlockHash)
		if err != nil {
			log.Printf("err getting cid by has\n%v", err)
			return false, err
		}
		// TODO: check that this is the actual prev block on the blockchain
		prevBlock, err := p2pSvc.GetStatechainBlock(prevBlockCID)
		if err != nil {
			log.Printf("err getting state block\n%v", err)
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

		prevState, err := fetchCurrentState(p2pSvc, prevBlock)
		if err != nil {
			log.Printf("err fetching current state\n%v", err)
			return false, err
		}
		prevStateHash := hashing.HashToHexString([]byte(prevState))
		if prevStateHash != prevBlock.Props().StateCurrentHash {
			return false, nil
		}

		for _, block := range orderedBlocks {
			// 2a. block must have a hash
			if block == nil || block.Props().BlockHash == nil {
				return false, nil
			}

			// 2b. Block #'s must be sequential
			prevBlockNumber, err := hexutil.DecodeUint64(prevBlock.Props().BlockNumber)
			if err != nil {
				log.Printf("err decoding prev block # 2b\n%v", err)
				return false, err
			}
			blockNumber, err := hexutil.DecodeUint64(block.Props().BlockNumber)
			if err != nil {
				log.Printf("err decoding block # 2b\n %v", err)
				return false, err
			}
			if prevBlockNumber+1 != blockNumber {
				return false, nil
			}

			// 2c. verify the block hash
			tmpHash, err := block.CalculateHash()
			if err != nil {
				log.Printf("err calculating block hash 2c\n%v", err)
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
				log.Printf("err verifying tx 2d\n %v", err)
				return false, err
			}
			if !ok {
				return false, nil
			}

			// TODO: run the tx through the container
			nextStateBlock, nextDiff, nextState, err := buildNextStateFromPrevState(p2pSvc, prevState, block, tx)
			if err != nil {
				log.Printf("err building next state from prev state\n %v", err)
				return false, err
			}
			if nextStateBlock == nil {
				return false, errors.New("nil state block")
			}
			if nextDiff == nil {
				return false, errors.New("nil diff")
			}
			if nextState == nil {
				return false, errors.New("nil next state")
			}
			if nextStateBlock.Props().BlockHash == nil {
				return false, errors.New("nil block hash")
			}
			if nextDiff.Props().DiffHash == nil {
				return false, errors.New("nil diff hash")
			}

			// 2e. verify current state hash
			if nextStateBlock.Props().StateCurrentHash != block.Props().StateCurrentHash {
				return false, nil
			}

			// 2f. verify prevDiff
			if nextStateBlock.Props().StatePrevDiffHash != block.Props().StatePrevDiffHash {
				return false, nil
			}

			// set prev to current for next loop
			prevState = nextState
			prevBlock = block
		}
	}

	return true, nil
}

func buildNextStateFromPrevState(p2pSvc p2p.Interface, prevState []byte, prevBlock *statechain.Block, tx *statechain.Transaction) (*statechain.Block, *statechain.Diff, []byte, error) {
	if prevState == nil {
		return nil, nil, nil, errors.New("nil state")
	}
	if prevBlock == nil {
		return nil, nil, nil, errors.New("nil prev block")
	}
	if tx == nil {
		return nil, nil, nil, errors.New("nil tx")
	}

	ts := time.Now().Unix()
	outPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/combined.txt", prevBlock.Props().ImageHash, ts))
	if err != nil {
		return nil, nil, nil, err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err := outPatchFile.Close(); err != nil {
		return nil, nil, nil, err
	}
	prevStateFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/prevState.txt", prevBlock.Props().ImageHash, ts))
	if err != nil {
		return nil, nil, nil, err
	}
	defer os.Remove(prevStateFile.Name()) // clean up
	if _, err := prevStateFile.Write(prevState); err != nil {
		return nil, nil, nil, err
	}
	if err := prevStateFile.Close(); err != nil {
		return nil, nil, nil, err
	}
	prevStateFileName := prevStateFile.Name()

	var nextState []byte
	if tx.Props().Method == "c3_invokeMethod" {
		payload, ok := tx.Props().Payload.([]byte)
		if !ok {
			return nil, nil, nil, errors.New("could not parse payload")
		}

		var parsed []string
		if err := json.Unmarshal(payload, &parsed); err != nil {
			return nil, nil, nil, err
		}

		inputsJSON, err := json.Marshal(struct {
			Method string   `json:"method"`
			Params []string `json:"params"`
		}{
			Method: parsed[0],
			Params: parsed[1:],
		})
		if err != nil {
			return nil, nil, nil, err
		}

		// run container, passing the tx inputs
		sb := sandbox.NewSandbox(&sandbox.Config{})
		nextState, err = sb.Play(&sandbox.PlayConfig{
			ImageID:      tx.Props().ImageHash,
			Payload:      inputsJSON,
			InitialState: prevState,
		})

		if err != nil {
			return nil, nil, nil, err
		}

		log.Printf("container new state: %s", string(nextState))
		nextStateFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/state.txt", prevBlock.Props().ImageHash, ts))
		if err != nil {
			return nil, nil, nil, err
		}
		defer os.Remove(nextStateFile.Name()) // clean up

		if _, err := nextStateFile.Write(nextState); err != nil {
			return nil, nil, nil, err
		}
		if err := nextStateFile.Close(); err != nil {
			return nil, nil, nil, err
		}

		if err := diffing.Diff(prevStateFileName, nextStateFile.Name(), outPatchFile.Name(), false); err != nil {
			return nil, nil, nil, err
		}

		// build the diff struct
		diffData, err := ioutil.ReadFile(outPatchFile.Name())
		if err != nil {
			return nil, nil, nil, err
		}

		diffStruct := statechain.NewDiff(&statechain.DiffProps{
			Data: string(diffData),
		})
		if err := diffStruct.SetHash(); err != nil {
			return nil, nil, nil, err
		}

		prevBlockNumber, err := hexutil.DecodeUint64(prevBlock.Props().BlockNumber)
		if err != nil {
			return nil, nil, nil, err
		}
		prevBlockNumber++

		nextStateHash := hexutil.EncodeBytes(nextState)
		nextStateStruct := statechain.New(&statechain.BlockProps{
			BlockNumber:       hexutil.EncodeUint64(prevBlockNumber),
			BlockTime:         hexutil.EncodeUint64(uint64(ts)),
			ImageHash:         prevBlock.Props().ImageHash,
			TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
			PrevBlockHash:     *prevBlock.Props().BlockHash,
			StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
			StateCurrentHash:  string(nextStateHash),
		})
		if err := nextStateStruct.SetHash(); err != nil {
			return nil, nil, nil, err
		}

		return nextStateStruct, diffStruct, nextState, nil
	}

	// TODO: is this what we want?
	return nil, nil, nil, errors.New("tx doesn't affect state")
}

func fetchCurrentState(p2pSvc p2p.Interface, block *statechain.Block) ([]byte, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}
	if block.Props().BlockHash == nil {
		return nil, errors.New("nil block hash")
	}

	var (
		diffs []*statechain.Diff
	)

	// gather the diffs
	diffCID, err := p2p.GetCIDByHash(block.Props().StatePrevDiffHash)
	if err != nil {
		return nil, err
	}
	diff, err := p2pSvc.GetStatechainDiff(diffCID)
	if err != nil {
		return nil, err
	}
	// note: prepend
	diffs = append([]*statechain.Diff{diff}, diffs...)

	head := block
	imageHash := block.Props().ImageHash
	for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
		prevStateCID, err := p2p.GetCIDByHash(head.Props().PrevBlockHash)
		if err != nil {
			return nil, err
		}

		prevStateBlock, err := p2pSvc.GetStatechainBlock(prevStateCID)
		if err != nil {
			return nil, err
		}
		head = prevStateBlock

		diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
		if err != nil {
			return nil, err
		}
		diff, err := p2pSvc.GetStatechainDiff(diffCID)
		if err != nil {
			return nil, err
		}
		// note: prepend
		diffs = append([]*statechain.Diff{diff}, diffs...)
	}

	// apply the diffs to get the current state
	// TODO: get the genesis state of the block
	genesisState := ""
	ts := time.Now().Unix()
	tmpStateFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpStateFile.Name()) // clean up

	if _, err := tmpStateFile.Write([]byte(genesisState)); err != nil {
		return nil, err
	}
	if err := tmpStateFile.Close(); err != nil {
		return nil, err
	}

	outPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/combined.txt", imageHash, ts))
	if err != nil {
		return nil, err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err := outPatchFile.Close(); err != nil {
		return nil, err
	}

	for i, diff := range diffs {
		tmpPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/patch.%d.txt", imageHash, ts, i))
		if err != nil {
			return nil, err
		}
		defer os.Remove(tmpPatchFile.Name()) // clean up

		if _, err := tmpPatchFile.Write([]byte(diff.Props().Data)); err != nil {
			return nil, err
		}
		if err := tmpPatchFile.Close(); err != nil {
			return nil, err
		}

		if err := diffing.CombineDiff(outPatchFile.Name(), tmpPatchFile.Name(), outPatchFile.Name()); err != nil {
			return nil, err
		}
	}

	// now apply the combined patch file to the state
	if err := diffing.Patch(outPatchFile.Name(), false, true); err != nil {
		return nil, err
	}
	state, err := ioutil.ReadFile(tmpStateFile.Name())
	if err != nil {
		return nil, err
	}

	return []byte(state), nil
}

// BuildNextState ...
func BuildNextState(p2pSvc p2p.Interface, block *statechain.Block, tx *statechain.Transaction) (*statechain.Block, *statechain.Diff, error) {
	if block == nil {
		return nil, nil, errors.New("nil block")
	}
	if block.Props().BlockHash == nil {
		return nil, nil, errors.New("nil block hash")
	}
	if tx == nil {
		return nil, nil, errors.New("nil tx")
	}
	if tx.Props().TxHash == nil {
		return nil, nil, errors.New("nil tx hash")
	}

	var (
		diffs []*statechain.Diff
	)

	// gather the diffs
	diffCID, err := p2p.GetCIDByHash(block.Props().StatePrevDiffHash)
	if err != nil {
		return nil, nil, err
	}
	diff, err := p2pSvc.GetStatechainDiff(diffCID)
	if err != nil {
		return nil, nil, err
	}
	// note: prepend
	diffs = append([]*statechain.Diff{diff}, diffs...)

	head := block
	imageHash := block.Props().ImageHash
	for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
		prevStateCID, err := p2p.GetCIDByHash(head.Props().PrevBlockHash)
		if err != nil {
			return nil, nil, err
		}

		prevStateBlock, err := p2pSvc.GetStatechainBlock(prevStateCID)
		if err != nil {
			return nil, nil, err
		}
		head = prevStateBlock

		diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
		if err != nil {
			return nil, nil, err
		}
		diff, err := p2pSvc.GetStatechainDiff(diffCID)
		if err != nil {
			return nil, nil, err
		}
		// note: prepend
		diffs = append([]*statechain.Diff{diff}, diffs...)
	}

	// apply the diffs to get the current state
	// TODO: get the genesis state of the block
	genesisState := ""
	ts := time.Now().Unix()
	tmpStateFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(tmpStateFile.Name()) // clean up

	if _, err := tmpStateFile.Write([]byte(genesisState)); err != nil {
		return nil, nil, err
	}
	if err := tmpStateFile.Close(); err != nil {
		return nil, nil, err
	}

	outPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/combined.txt", imageHash, ts))
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err := outPatchFile.Close(); err != nil {
		return nil, nil, err
	}

	for i, diff := range diffs {
		tmpPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/patch.%d.txt", imageHash, ts, i))
		if err != nil {
			return nil, nil, err
		}
		defer os.Remove(tmpPatchFile.Name()) // clean up

		if _, err := tmpPatchFile.Write([]byte(diff.Props().Data)); err != nil {
			return nil, nil, err
		}
		if err := tmpPatchFile.Close(); err != nil {
			return nil, nil, err
		}

		if err := diffing.CombineDiff(outPatchFile.Name(), tmpPatchFile.Name(), outPatchFile.Name()); err != nil {
			return nil, nil, err
		}
	}

	// now apply the combined patch file to the state
	if err := diffing.Patch(outPatchFile.Name(), false, true); err != nil {
		return nil, nil, err
	}
	state, err := ioutil.ReadFile(tmpStateFile.Name())
	if err != nil {
		return nil, nil, err
	}

	log.Printf("state\n%s", string(state))
	headStateFileName := tmpStateFile.Name()
	runningBlockNumber, err := hexutil.DecodeUint64(block.Props().BlockNumber)
	if err != nil {
		return nil, nil, err
	}
	runningBlockHash := *block.Props().BlockHash // note: already checked nil pointer, above

	// apply state to container and run tx
	var nextState []byte

	if tx.Props().Method == "c3_invokeMethod" {
		payload, ok := tx.Props().Payload.([]byte)
		if !ok {
			return nil, nil, errors.New("could not parse payload")
		}

		var parsed []string
		if err := json.Unmarshal(payload, &parsed); err != nil {
			return nil, nil, err
		}

		inputsJSON, err := json.Marshal(struct {
			Method string   `json:"method"`
			Params []string `json:"params"`
		}{
			Method: parsed[0],
			Params: parsed[1:],
		})
		if err != nil {
			return nil, nil, err
		}

		// run container, passing the tx inputs
		sb := sandbox.NewSandbox(&sandbox.Config{})
		nextState, err = sb.Play(&sandbox.PlayConfig{
			ImageID:      tx.Props().ImageHash,
			Payload:      inputsJSON,
			InitialState: state,
		})

		if err != nil {
			return nil, nil, err
		}

		log.Printf("container new state: %s", string(nextState))

		nextStateFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
		if err != nil {
			return nil, nil, err
		}
		defer os.Remove(nextStateFile.Name()) // clean up

		if _, err := nextStateFile.Write(nextState); err != nil {
			return nil, nil, err
		}
		if err := nextStateFile.Close(); err != nil {
			return nil, nil, err
		}

		if err := diffing.Diff(headStateFileName, nextStateFile.Name(), outPatchFile.Name(), false); err != nil {
			return nil, nil, err
		}
		headStateFileName = nextStateFile.Name()

		// build the diff struct
		diffData, err := ioutil.ReadFile(outPatchFile.Name())
		if err != nil {
			return nil, nil, err
		}

		diffStruct := statechain.NewDiff(&statechain.DiffProps{
			Data: string(diffData),
		})
		if err := diffStruct.SetHash(); err != nil {
			return nil, nil, err
		}

		nextStateHash := hexutil.EncodeBytes(nextState)
		runningBlockNumber++
		nextStateStruct := statechain.New(&statechain.BlockProps{
			BlockNumber:       hexutil.EncodeUint64(runningBlockNumber),
			BlockTime:         hexutil.EncodeUint64(uint64(ts)),
			ImageHash:         block.Props().ImageHash,
			TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
			PrevBlockHash:     runningBlockHash,
			StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
			StateCurrentHash:  string(nextStateHash),
		})
		if err := nextStateStruct.SetHash(); err != nil {
			return nil, nil, err
		}

		return nextStateStruct, diffStruct, nil
	}

	// TODO: what to do when invoke method not called?
	return nil, nil, errors.New("no statechange")
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
func orderStatechainBlocks(blocks []*statechain.Block) ([]*statechain.Block, error) {
	if blocks == nil {
		return nil, errors.New("nil blocks")
	}

	var (
		stateBlocks   []*statechain.Block
		blockWrappers byBlockNumber
	)

	for _, stateBlock := range blocks {
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

// the return is a map with keys on the image hash
func groupStateBlocksByImageHash(stateBlocksMap map[string]*statechain.Block) (map[string][]*statechain.Block, error) {
	if stateBlocksMap == nil {
		return nil, errors.New("nil stateblocks map")
	}

	ret := make(map[string][]*statechain.Block)

	for _, block := range stateBlocksMap {
		if block == nil {
			return nil, errors.New("nil block")
		}

		ret[block.Props().ImageHash] = append(ret[block.Props().ImageHash], block)
	}

	return ret, nil
}
