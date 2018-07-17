package miner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/c3systems/c3/common/c3crypto"
	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/merkle"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/diffing"
	"github.com/c3systems/c3/core/p2p"
	"github.com/c3systems/c3/core/sandbox"
	methodTypes "github.com/c3systems/c3/core/types/methods"

	log "github.com/sirupsen/logrus"
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
/*
func VerifyStatechainBlock(ctx context.Context, p2pSvc p2p.Interface, block *statechain.Block) (bool, error) {
	if block == nil {
		return false, ErrNilBlock
	}

	// 1. block must have a hash
	if block.Props().BlockHash == nil {
		return false, nil
	}

	// TODO: check the block # and StatePrevDiffHash

	// 2. verify the block hash
	tmpHash, err := block.CalculateHash()
	if err != nil {
		log.Printf("[miner] err calculating block hash\n%v", err)
		return false, err
	}
	// note: checked nil BlockHash, above
	if tmpHash != *block.Props().BlockHash {
		return false, nil
	}

	// 3. verify each tx in the block
	// TODO: do in go funcs
	txCid, err := p2p.GetCIDByHash(block.Props().TxHash)
	if err != nil {
		log.Printf("[miner] err getting cid for tx\n%v", err)
		return false, err
	}

	tx, err := p2pSvc.GetStatechainTransaction(txCid)
	if err != nil {
		return false, err
	}

	ok, err := VerifyTransaction(tx)
	if err != nil {
		log.Printf("[miner] err verifying tx\n%v", err)
		return false, err
	}
	if !ok {
		return false, nil
	}

	// note: just printing to keep the txs var alive
	spew.Dump(tx)

	// 4. run the tx through the container
	// TODO: step #4

	// 5. verify the statehash and prev diff hash
	// TODO step #5

	return true, nil
}
*/

// VerifyMainchainBlock verifies a mainchain block
// TODO: check block time
// TODO: fetch and check previous block hash
// TODO: pass all necessary data and remove the p2pSvc
/*
func VerifyMainchainBlock(ctx context.Context, p2pSvc p2p.Interface, block *mainchain.Block) (bool, error) {
	if block == nil {
		return false, errors.New("block is nil")
	}
	if block.Props().BlockHash == nil {
		return false, errors.New("block hash is nil")
	}
	if mainchain.ImageHash != block.Props().ImageHash {
		return false, nil
	}
	if block.Props().MinerSig == nil {
		return false, nil
	}

	ok, err := CheckBlockHashAgainstDifficulty(block)
	if err != nil {
		log.Printf("[miner] err checking block hash against difficulty\n%v", err)
		return false, err
	}
	if !ok {
		return false, nil
	}

	// hash must verify
	tmpHash, err := block.CalculateHash()
	if err != nil {
		log.Printf("[miner] err calculating tmpHash\n%v", err)
		return false, err
	}
	// note: already checked for nil hash
	if *block.Props().BlockHash != tmpHash {
		return false, nil
	}

	// the sig must verify
	pub, err := c3crypto.DecodeAddress(block.Props().MinerAddress)
	if err != nil {
		log.Printf("[miner] err decoding miner addr\n%v", err)
		return false, err
	}

	// note: checked for nil sig, above
	sigR, err := hexutil.DecodeBigInt(block.Props().MinerSig.R)
	if err != nil {
		log.Printf("[miner] err decoding miner sig r\n%v", err)
		return false, err
	}
	sigS, err := hexutil.DecodeBigInt(block.Props().MinerSig.S)
	if err != nil {
		log.Printf("[miner] err decoding miner sig s\n%v", err)
		return false, err
	}

	// note: nil blockhash was checked, above
	ok, err = c3crypto.Verify(pub, []byte(*block.Props().BlockHash), sigR, sigS)
	if err != nil {
		log.Printf("[miner] err verifying\n%v", err)
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
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		stateblockCid, err := p2p.GetCIDByHash(stateblockHash)
		if err != nil {
			log.Printf("[miner] err getting cid\n%v", err)
			return false, err
		}

		stateblock, err := p2pSvc.GetStatechainBlock(stateblockCid)
		if err != nil {
			log.Printf("[miner] err getting statechain block\n%v", err)
			return false, err
		}

		ok, err := VerifyStatechainBlock(ctx, p2pSvc, stateblock)
		if err != nil {
			log.Printf("[miner] err verifying statechain block\n%v", err)
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}
*/

// VerifyMinedBlock ...
func VerifyMinedBlock(ctx context.Context, p2pSvc p2p.Interface, minedBlock *MinedBlock) (bool, error) {
	ch := make(chan interface{})

	go func() {
		if minedBlock == nil {
			log.Println("[miner] mined block is nil")
			ch <- false

			return
		}
		if minedBlock.NextBlock == nil {
			log.Println("[miner] next block is nil")
			ch <- false

			return
		}
		if minedBlock.PreviousBlock == nil {
			log.Println("[miner] mined block is nil")
			ch <- false

			return
		}
		if minedBlock.NextBlock.Props().BlockHash == nil {
			log.Println("[miner] next block blockhash is nil")
			ch <- false

			return
		}
		if minedBlock.PreviousBlock.Props().BlockHash == nil {
			log.Println("[miner] prev block block hash is nil")
			ch <- false

			return
		}
		// note checked for nil pointer, above
		if *minedBlock.PreviousBlock.Props().BlockHash != minedBlock.NextBlock.Props().PrevBlockHash {
			log.Println("[miner] prev block block hash != next block block hash")
			ch <- false

			return
		}
		if mainchain.ImageHash != minedBlock.NextBlock.Props().ImageHash {
			log.Println("[miner] mainchain imagehash != nextblock image hash")
			ch <- false

			return
		}
		if minedBlock.NextBlock.Props().MinerSig == nil {
			log.Println("[miner] next block miner sig is nil")
			ch <- false

			return
		}

		ok, err := CheckBlockHashAgainstDifficulty(minedBlock.NextBlock)
		if err != nil {
			log.Printf("[miner] err checking block hash against difficulty\n%v", err)
			ch <- err

			return
		}
		if !ok {
			log.Println("[miner] block hash did not checkout against difficulty")
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		tmpHash, err := minedBlock.NextBlock.CalculateHash()
		if err != nil {
			log.Printf("[miner] err calculating hash\n%v", err)
			ch <- err

			return
		}
		// note: already checked for nil hash
		if *minedBlock.NextBlock.Props().BlockHash != tmpHash {
			log.Printf("[miner] next block hash != calced hash\n%s\n%s", *minedBlock.NextBlock.Props().BlockHash, tmpHash)
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		pub, err := c3crypto.DecodeAddress(minedBlock.NextBlock.Props().MinerAddress)
		if err != nil {
			log.Printf("[miner] err decoding addr\n%v", err)
			ch <- err

			return
		}

		// note: checked for nil sig, above
		sigR, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.R)
		if err != nil {
			log.Printf("[miner] err decoding r\n%v", err)
			ch <- err

			return
		}
		sigS, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.S)
		if err != nil {
			log.Printf("[miner] err decoding s\n%v", err)
			ch <- err

			return
		}

		// note: nil blockhash was checked, above
		ok, err = c3crypto.Verify(pub, []byte(*minedBlock.NextBlock.Props().BlockHash), sigR, sigS)
		if err != nil {
			log.Printf("[miner] err verifying miner sig\n%v", err)
			ch <- err

			return
		}
		if !ok {
			log.Println("[miner] block hash did not checkout agains sig")
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		// BlockNumber must be +1 prev block number
		blockNumber, err := hexutil.DecodeUint64(minedBlock.NextBlock.Props().BlockNumber)
		if err != nil {
			log.Printf("[miner] err decoding block #\n%v", err)
			ch <- err

			return
		}
		prevNumber, err := hexutil.DecodeUint64(minedBlock.PreviousBlock.Props().BlockNumber)
		if err != nil {
			log.Printf("[miner] err decoding prev block #\n%v", err)
			ch <- err

			return
		}
		if prevNumber+1 != blockNumber {
			log.Println("[miner] prevBlockNumber +1 != nextBlockNumber")
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		ch <- true
	}()

	select {
	case v := <-ch:
		switch v.(type) {
		case error:
			err, _ := v.(error)
			return false, err

		case bool:
			ok, _ := v.(bool)
			if !ok {
				return false, nil
			}

			return VerifyStateBlocksFromMinedBlock(ctx, p2pSvc, minedBlock)

		default:
			log.Printf("[miner] received unknown message of type %T\n%v", v, v)
			return false, errors.New("received message of unknown type")

		}

	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// VerifyStateBlocksFromMinedBlock ...
// note: this function also checks the merkle tree. That check is not required to be performed, separately.
func VerifyStateBlocksFromMinedBlock(ctx context.Context, p2pSvc p2p.Interface, minedBlock *MinedBlock) (bool, error) {
	ch := make(chan interface{})

	go func() {
		if minedBlock.NextBlock == nil {
			log.Println("[miner] nil next block")
			ch <- false

			return
		}
		if len(minedBlock.StatechainBlocksMap) != len(minedBlock.TransactionsMap) {
			log.Println("[miner] len state blocks map != len tx map")
			ch <- false

			return
		}
		// note: ok to have nil map? e.g. in the case that no transactions were included in the main block
		//if minedBlock.StatechainBlocksMap == nil {
		//log.Println("nil state blocks map")
		//return false, nil
		//}
		if minedBlock.MerkleTreesMap == nil {
			log.Println("[miner] nil merkle trees map")
			ch <- false

			return
		}

		// 1. Verify state blocks merkle hash
		ok, err := VerifyMerkleTreeFromMinedBlock(ctx, minedBlock)
		if err != nil {
			log.Printf("[miner] err verifying merkle tree\n%v", err)
			ch <- err

			return
		}
		if !ok {
			log.Println("[miner] merkle tree didn't verify")
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		// 2. Verify each state block
		// first, group them by image hash
		groupedBlocks, err := groupStateBlocksByImageHash(minedBlock.StatechainBlocksMap)
		if err != nil {
			log.Printf("[miner] err grouping state blocks\n%v", err)
			ch <- err

			return
		}

		for _, blocks := range groupedBlocks {
			if ctx.Err() != nil {
				return
			}

			// order by block number
			orderedBlocks, err := orderStatechainBlocks(blocks)
			if err != nil {
				log.Printf("[miner] err ordering state blocks\n%v", err)
				ch <- err

				return
			}

			if orderedBlocks == nil || len(orderedBlocks) == 0 {
				continue
			}
			prevBlockHash := orderedBlocks[0].Props().PrevBlockHash
			prevBlockCID, err := p2p.GetCIDByHash(prevBlockHash)
			if err != nil {
				log.Printf("[miner] err getting cid by has\n%v", err)
				ch <- err

				return
			}
			// TODO: check that this is the actual prev block on the blockchain
			prevBlock, err := p2pSvc.GetStatechainBlock(prevBlockCID)
			if err != nil {
				log.Printf("[miner] err getting state block\n%v", err)
				ch <- err

				return
			}
			if prevBlock == nil {
				ch <- errors.New("got nil prev block")

				return
			}
			if prevBlock.Props().BlockHash == nil {
				ch <- errors.New("got nil prev block hash")

				return
			}
			// note: checked for nil pointer, above
			if *prevBlock.Props().BlockHash != prevBlockHash {
				ch <- false

				return
			}

			prevState, err := fetchCurrentState(ctx, p2pSvc, prevBlock)
			if err != nil {
				log.Printf("[miner] err fetching current state\n%v", err)
				ch <- err

				return
			}
			prevStateHash := hashing.HashToHexString([]byte(prevState))
			if prevStateHash != prevBlock.Props().StateCurrentHash {
				ch <- false

				return
			}

			for _, block := range orderedBlocks {
				if ctx.Err() != nil {
					return
				}

				// 2a. block must have a hash
				if block == nil || block.Props().BlockHash == nil {
					ch <- false

					return
				}

				// 2b. Block #'s must be sequential
				prevBlockNumber, err := hexutil.DecodeUint64(prevBlock.Props().BlockNumber)
				if err != nil {
					log.Printf("[miner] err decoding prev block # 2b\n%v", err)
					ch <- err

					return
				}
				blockNumber, err := hexutil.DecodeUint64(block.Props().BlockNumber)
				if err != nil {
					log.Printf("[miner] err decoding block # 2b\n %v", err)
					ch <- err

					return
				}
				if prevBlockNumber+1 != blockNumber {
					ch <- false

					return
				}

				// 2c. verify the block hash
				tmpHash, err := block.CalculateHash()
				if err != nil {
					log.Printf("[miner] err calculating block hash 2c\n%v", err)
					ch <- err

					return
				}
				// note: checked nil BlockHash, above
				if tmpHash != *block.Props().BlockHash {
					ch <- false

					return
				}

				// 2d. verify the block tx
				// note: can't have a state block without transactions?
				tx, ok := minedBlock.TransactionsMap[block.Props().TxHash]
				if !ok || tx == nil {
					txCID, err := p2p.GetCIDByHash(block.Props().TxHash)
					if err != nil {
						ch <- err

						return
					}

					tx, err = p2pSvc.GetStatechainTransaction(txCID)
					if err != nil {
						ch <- err

						return
					}
					if tx == nil {
						ch <- errors.New("nil tx")

						return
					}
				}

				ok, err = VerifyTransaction(tx)
				if err != nil {
					log.Printf("[miner] err verifying tx 2d\n %v", err)
					ch <- err

					return
				}
				if !ok {
					ch <- false

					return
				}
				if ctx.Err() != nil {
					return
				}

				nextStateBlock, nextDiff, nextState, err := buildNextStateFromPrevState(p2pSvc, prevState, block, tx)
				if err != nil {
					log.Printf("[miner] err building next state from prev state\n %v", err)
					ch <- err

					return
				}
				if nextStateBlock == nil {
					ch <- errors.New("nil state block")

					return
				}
				if nextDiff == nil {
					ch <- errors.New("nil diff")

					return
				}
				if nextState == nil {
					ch <- errors.New("nil next state")

					return
				}
				if nextStateBlock.Props().BlockHash == nil {
					ch <- errors.New("nil block hash")

					return
				}
				if nextDiff.Props().DiffHash == nil {
					ch <- errors.New("nil diff hash")

					return
				}
				if ctx.Err() != nil {
					return
				}

				// 2e. verify current state hash
				if nextStateBlock.Props().StateCurrentHash != block.Props().StateCurrentHash {
					ch <- false

					return
				}

				// 2f. verify prevDiff
				if nextStateBlock.Props().StatePrevDiffHash != block.Props().StatePrevDiffHash {
					ch <- false

					return
				}

				// set prev to current for next loop
				prevState = nextState
				prevBlock = block
			}
		}
		if ctx.Err() != nil {
			return
		}

		ch <- true
		return
	}()

	select {
	case v := <-ch:
		switch v.(type) {
		case error:
			err, _ := v.(error)
			return false, err

		case bool:
			ok, _ := v.(bool)
			if !ok {
				return false, nil
			}

			return true, nil

		default:
			log.Printf("[miner] received unknown message of type %T\n%v", v, v)
			return false, errors.New("received message of unknown type")

		}

	case <-ctx.Done():
		return false, ctx.Err()
	}
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
	if tx.Props().Method == methodTypes.InvokeMethod {
		payload := tx.Props().Payload

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

		//log.Printf("[miner] container new state: %s", string(nextState))
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

func fetchCurrentState(ctx context.Context, p2pSvc p2p.Interface, block *statechain.Block) ([]byte, error) {
	ch := make(chan interface{})

	go func() {
		if block == nil {
			ch <- errors.New("nil block")

			return
		}
		if block.Props().BlockHash == nil {
			ch <- errors.New("nil block hash")

			return
		}
		if ctx.Err() != nil {
			return
		}

		var (
			diffs []*statechain.Diff
		)

		// gather the diffs
		diffCID, err := p2p.GetCIDByHash(block.Props().StatePrevDiffHash)
		if err != nil {
			log.Printf("[miner] err getting diff cid by hash\n%v", err)
			ch <- err

			return
		}
		diff, err := p2pSvc.GetStatechainDiff(diffCID)
		if err != nil {
			log.Printf("[miner] err getting diff by cid\n%v", err)
			ch <- err

			return
		}
		// note: prepend
		diffs = append([]*statechain.Diff{diff}, diffs...)

		head := block
		imageHash := block.Props().ImageHash
		for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
			if ctx.Err() != nil {
				return
			}

			prevStateCID, err := p2p.GetCIDByHash(head.Props().PrevBlockHash)
			if err != nil {
				log.Printf("[miner] err getting statechain cid by hash\n%v", err)
				ch <- err

				return
			}

			prevStateBlock, err := p2pSvc.GetStatechainBlock(prevStateCID)
			if err != nil {
				log.Printf("[miner] err getting state chain block by cid\n%v", err)
				ch <- err

				return
			}
			head = prevStateBlock

			diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
			if err != nil {
				log.Printf("[miner] err getting diffCID by hash\n%v", err)
				ch <- err

				return
			}
			diff, err := p2pSvc.GetStatechainDiff(diffCID)
			if err != nil {
				log.Printf("[miner] err gitting diff by cid\n%v", err)
				ch <- err

				return
			}
			// note: prepend
			diffs = append([]*statechain.Diff{diff}, diffs...)
		}

		// apply the diffs to get the current state
		// TODO: get the genesis state of the block?
		genesisState := ""
		ts := time.Now().Unix()
		tmpStateFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
		if err != nil {
			log.Printf("[miner] err creating state file\n%v", err)
			ch <- err

			return
		}
		defer os.Remove(tmpStateFile.Name()) // clean up

		if _, err := tmpStateFile.Write([]byte(genesisState)); err != nil {
			log.Printf("[miner] err writing to state file\n%v", err)
			ch <- err

			return
		}
		if err := tmpStateFile.Close(); err != nil {
			log.Printf("[miner] err closing state file\n%v", err)
			ch <- err

			return
		}

		outPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/combined.txt", imageHash, ts))
		if err != nil {
			log.Printf("[miner] err creating combined patch file\n%v", err)
			ch <- err

			return
		}
		defer os.Remove(outPatchFile.Name()) // clean up
		if err := outPatchFile.Close(); err != nil {
			log.Printf("[miner] err closing patch file\n%v", err)
			ch <- err

			return
		}

		for i, diff := range diffs {
			if ctx.Err() != nil {
				return
			}

			tmpPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/patch.%d.txt", imageHash, ts, i))
			if err != nil {
				log.Printf("[miner] err writing tmp patch file\n%v", err)
				ch <- err

				return
			}
			defer os.Remove(tmpPatchFile.Name()) // clean up

			if _, err := tmpPatchFile.Write([]byte(diff.Props().Data)); err != nil {
				log.Printf("[miner] err writing diff to patch file\n%v", err)
				ch <- err

				return
			}
			if err := tmpPatchFile.Close(); err != nil {
				log.Printf("[miner] err closing patch file\n%v", err)
				ch <- err

				return
			}

			if err := diffing.CombineDiff(outPatchFile.Name(), tmpPatchFile.Name(), outPatchFile.Name()); err != nil {
				log.Printf("[miner] err combining diff files\n%v", err)
				ch <- err

				return
			}
		}

		// now apply the combined patch file to the state
		if err := diffing.Patch(outPatchFile.Name(), false, true); err != nil {
			log.Printf("[miner] err diffing patch file\n%v", err)
			ch <- err

			return
		}
		state, err := ioutil.ReadFile(tmpStateFile.Name())
		if err != nil {
			log.Printf("[miner] err reading state file\n%v", err)
			ch <- err

			return
		}
		if ctx.Err() != nil {
			return
		}

		ch <- []byte(state)
		return
	}()

	select {
	case v := <-ch:
		switch v.(type) {
		case error:
			err, _ := v.(error)
			return nil, err

		case []byte:
			bytes, _ := v.([]byte)

			return bytes, nil

		default:
			log.Printf("[miner] received unknown message of type %T\n%v", v, v)
			return nil, errors.New("received message of unknown type")

		}

	case <-ctx.Done():
		return nil, ctx.Err()
	}
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

	headStateFileName := tmpStateFile.Name()
	runningBlockNumber, err := hexutil.DecodeUint64(block.Props().BlockNumber)
	if err != nil {
		return nil, nil, err
	}
	runningBlockHash := *block.Props().BlockHash // note: already checked nil pointer, above

	// apply state to container and run tx
	var nextState []byte

	if tx.Props().Method == methodTypes.InvokeMethod {
		payload := tx.Props().Payload

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

		//log.Printf("[miner] container new state: %s", string(nextState))

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
func VerifyMerkleTreeFromMinedBlock(ctx context.Context, minedBlock *MinedBlock) (bool, error) {
	if minedBlock.NextBlock == nil {
		log.Printf("[miner] verify mined block error - block data is nil")
		return false, nil
	}
	/*
		// note: mined block with no tx will have nil state chain blocks map
		if minedBlock.StatechainBlocksMap == nil {
			log.Printf("[miner] verify mined block error - state chain blocks is nil")
			return false, nil
		}
	*/
	if minedBlock.MerkleTreesMap == nil {
		log.Printf("[miner] verify mined block error - merkle tree map is nil")
		return false, nil
	}

	tree, ok := minedBlock.MerkleTreesMap[minedBlock.NextBlock.Props().StateBlocksMerkleHash]
	if !ok || tree == nil {
		log.Printf("[miner] verify mined block error - merkle trees map %s is not ok or nil", minedBlock.NextBlock.Props().StateBlocksMerkleHash)
		return false, nil
	}
	if tree.Props().MerkleTreeRootHash == nil {
		log.Println("[miner] verify mined block error - merkle tree root hash is nil")
		return false, nil
	}

	tmpTree, err := merkle.New(&merkle.TreeProps{
		Hashes: tree.Props().Hashes,
		Kind:   merkle.StatechainBlocksKindStr,
	})
	if err != nil {
		log.Printf("[miner] verify mined block error - new merkle tree error: %s", err)
		return false, err
	}
	if err := tmpTree.SetHash(); err != nil {
		log.Printf("[miner] verify mined block error - set hash error: %s", err)
		return false, err
	}

	if *tmpTree.Props().MerkleTreeRootHash != *tree.Props().MerkleTreeRootHash {
		log.Printf("[miner] verify mined block error - merkle tree root hash doesn't match; %s", *tmpTree.Props().MerkleTreeRootHash)
		return false, nil
	}

	if len(tmpTree.Props().Hashes) != len(minedBlock.StatechainBlocksMap) {
		log.Printf("[miner] verify mined block error - tree hashes length doesn't match; %v", len(tmpTree.Props().Hashes))
		return false, nil
	}

	for _, hash := range tmpTree.Props().Hashes {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		statechainBlock, ok := minedBlock.StatechainBlocksMap[hash]
		if !ok || statechainBlock == nil {
			log.Println("[miner] verify mined block error - state chain block from map is nil or not ok")
			return false, nil
		}

		tmpHash, err := statechainBlock.CalculateHash()
		if err != nil {
			log.Printf("[miner] verify mined block error - state chain calculate hash error: %s", err)
			return false, err
		}
		if hash != tmpHash {
			log.Printf("[miner] verify mined block error - hash does not match %s", tmpHash)
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
	/*
		// note: mined block with no transactions will have nil state block map
		if stateBlocksMap == nil {
			log.Printf("[miner] state blocks map is nil")
			return nil, errors.New("nil stateblocks map")
		}
	*/

	ret := make(map[string][]*statechain.Block)

	if stateBlocksMap != nil {
		for _, block := range stateBlocksMap {
			if block == nil {
				return nil, errors.New("nil block")
			}

			ret[block.Props().ImageHash] = append(ret[block.Props().ImageHash], block)
		}
	}

	return ret, nil
}

func cleanupFiles(fileNames []string) {
	for idx := range fileNames {
		if err := os.Remove(fileNames[idx]); err != nil {
			log.Printf("[miner] err cleaning up file %s", fileNames[idx])
		}
	}
}

func makeTempFile(filename string) (*os.File, error) {
	paths := strings.Split(filename, "/")
	prefix := strings.Join(paths[:len(paths)-1], "_") // does not like slashes for some reason
	filename = strings.Join(paths[len(paths)-1:len(paths)], "")

	tmpdir, err := ioutil.TempDir("/tmp", prefix)
	if err != nil {
		return nil, err
	}

	filepath := fmt.Sprintf("%s/%s", tmpdir, filename)

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}
