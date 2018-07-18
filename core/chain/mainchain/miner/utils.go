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

			if orderedBlocks[0] == nil {
				ch <- errors.New("nil block")

				return
			}

			// TODO: improve this if/else hell
			if orderedBlocks[0].Props().BlockNumber == mainchain.GenesisBlock.Props().BlockNumber {
				// check if there's already a genesis block
				block, err := p2pSvc.FetchMostRecentStateBlock(orderedBlocks[0].Props().ImageHash, minedBlock.NextBlock)
				if err != nil {
					ch <- err

					return
				}

				if block == nil {
					tx, ok := minedBlock.TransactionsMap[orderedBlocks[0].Props().TxHash]
					if !ok {
						ch <- errors.New("tx not included")

						return
					}
					if tx == nil {
						ch <- errors.New("tx is nil")

						return
					}

					genesisBlock, _, err := buildGenesisStateBlock(orderedBlocks[0].Props().ImageHash, tx)
					if err != nil {
						ch <- err

						return
					}

					block = genesisBlock
				}
				if orderedBlocks[0].Props().BlockHash == nil {
					ch <- false

					return
				}
				hash, err := orderedBlocks[0].CalculateHash()
				if err != nil {
					ch <- err

					return
				}
				if hash != *orderedBlocks[0].Props().BlockHash {
					ch <- false

					return
				}

				if block.Props().BlockHash == nil {
					ch <- errors.New("nil blockhash")

					return
				}
				if *block.Props().BlockHash != hash {
					ch <- false

					return
				}

				orderedBlocks = orderedBlocks[1:]
				if len(orderedBlocks) == 0 {
					continue
				}

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
			prevStateHash := hashing.HashToHexString(prevState)
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
	outPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/combined.txt", prevBlock.Props().ImageHash, ts))
	if err != nil {
		return nil, nil, nil, err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err = outPatchFile.Close(); err != nil {
		return nil, nil, nil, err
	}
	prevStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/prevState.txt", prevBlock.Props().ImageHash, ts))
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
		nextStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.txt", prevBlock.Props().ImageHash, ts))
		if err != nil {
			return nil, nil, nil, err
		}
		defer os.Remove(nextStateFile.Name()) // clean up

		if _, err = nextStateFile.Write(nextState); err != nil {
			return nil, nil, nil, err
		}
		if err = nextStateFile.Close(); err != nil {
			return nil, nil, nil, err
		}

		if err = diffing.Diff(prevStateFileName, nextStateFile.Name(), outPatchFile.Name(), false); err != nil {
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
		if err = diffStruct.SetHash(); err != nil {
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

// TODO: improve
func buildGenesisStateBlock(imageHash string, tx *statechain.Transaction) (*statechain.Block, *statechain.Diff, error) {
	log.Printf("[miner] building genesis state block for image hash %s", imageHash)

	ts := time.Now().Unix()

	if tx.Props().TxHash == nil {
		log.Printf("[miner] tx hash is nil for %v", tx.Props())
		return nil, nil, errors.New("nil tx hash")
	}

	if tx.Props().TxHash == nil {
		log.Printf("[miner] tx hash is nil for %v", tx.Props())
		return nil, nil, errors.New("nil tx hash")
	}

	log.Printf("[miner] tx method %s", tx.Props().Method)

	// initial state
	nextState := tx.Props().Payload
	log.Printf("[miner] container initial state: %s", string(nextState))

	nextStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(nextStateFile.Name()) // clean up
	if _, err = nextStateFile.Write(nextState); err != nil {
		return nil, nil, err
	}
	if err = nextStateFile.Close(); err != nil {
		return nil, nil, err
	}

	tmpStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/tmpState.txt", imageHash, ts))
	defer os.Remove(tmpStateFile.Name()) // clean up
	if err != nil {
		return nil, nil, err
	}
	if err = tmpStateFile.Close(); err != nil {
		return nil, nil, err
	}

	outPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/combined.txt", imageHash, ts))
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err = outPatchFile.Close(); err != nil {
		return nil, nil, err
	}

	if err = diffing.Diff(tmpStateFile.Name(), nextStateFile.Name(), outPatchFile.Name(), false); err != nil {
		return nil, nil, err
	}

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
	nextStateStruct := statechain.New(&statechain.BlockProps{
		BlockNumber:       hexutil.EncodeUint64(0),
		BlockTime:         hexutil.EncodeUint64(uint64(ts)),
		ImageHash:         imageHash,
		TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
		PrevBlockHash:     "",
		StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
		StateCurrentHash:  string(nextStateHash),
	})

	if err := nextStateStruct.SetHash(); err != nil {
		return nil, nil, err
	}

	return nextStateStruct, diffStruct, nil
}

func fetchCurrentState(ctx context.Context, p2pSvc p2p.Interface, block *statechain.Block) ([]byte, error) {
	ch := make(chan interface{})

	go func() {
		if block == nil {
			log.Println("[miner] block is nil")
			ch <- errors.New("nil block")

			return
		}
		if block.Props().BlockHash == nil {
			log.Println("[miner] block hash is nil")
			ch <- errors.New("nil block hash")

			return
		}
		if ctx.Err() != nil {
			return
		}

		// gather the diffs
		diffs, err := gatherDiffs(ctx, p2pSvc, block)
		if err != nil {
			log.Printf("[miner] error gathering diffs\n%v", err)
			ch <- err

			return
		}

		// apply the diffs to get the current state
		// TODO: get the genesis state of the block?
		genesisState := []byte("")
		imageHash := block.Props().ImageHash
		state, err := generateStateFromDiffs(ctx, imageHash, genesisState, diffs)
		if err != nil {
			log.Printf("[miner] error reading state file\n%v", err)
			ch <- err

			return
		}
		if ctx.Err() != nil {
			log.Printf("[miner] received context error; %s", err)
			return
		}

		ch <- state
		return
	}()

	select {
	case v := <-ch:
		switch v.(type) {
		case error:
			err, _ := v.(error)
			log.Printf("[miner] channel error; %s", err)
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

func gatherDiffs(ctx context.Context, p2pSvc p2p.Interface, block *statechain.Block) ([]*statechain.Diff, error) {
	var diffs []*statechain.Diff

	// gather the diffs
	diffCID, err := p2p.GetCIDByHash(block.Props().StatePrevDiffHash)
	if err != nil {
		log.Printf("[miner] err getting diff cid by hash\n%v", err)
		return nil, err
	}

	diff, err := p2pSvc.GetStatechainDiff(diffCID)
	if err != nil {
		log.Printf("[miner] err getting diff by cid\n%v", err)
		return nil, err
	}

	log.Printf("[miner] state chain diff; %v", diff)

	// note: prepend
	diffs = append([]*statechain.Diff{diff}, diffs...)

	head := block
	for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
		if ctx.Err() != nil {
			log.Printf("[miner] gather diffs context error; %v", err)
			return nil, ctx.Err()
		}

		prevStateCID, err := p2p.GetCIDByHash(head.Props().PrevBlockHash)
		if err != nil {
			log.Printf("[miner] err getting statechain cid by hash\n%v", err)
			return nil, err
		}

		prevStateBlock, err := p2pSvc.GetStatechainBlock(prevStateCID)
		if err != nil {
			log.Printf("[miner] err getting state chain block by cid\n%v", err)
			return nil, err
		}

		head = prevStateBlock

		diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
		if err != nil {
			log.Printf("[miner] err getting diffCID by hash\n%v", err)
			return nil, err
		}

		diff, err := p2pSvc.GetStatechainDiff(diffCID)
		if err != nil {
			log.Printf("[miner] err gitting diff by cid\n%v", err)
			return nil, err
		}

		// note: prepend
		diffs = append([]*statechain.Diff{diff}, diffs...)
	}

	if diffs == nil {
		log.Println("[miner] error; diffs is nil")
		return nil, errors.New("diffs is nil")
	}

	return diffs, nil
}

func generateStateFromDiffs(ctx context.Context, imageHash string, genesisState []byte, diffs []*statechain.Diff) ([]byte, error) {
	combinedDiff, err := generateCombinedDiffs(ctx, imageHash, genesisState, diffs)
	if err != nil {
		log.Printf("[miner] error generating combined diffs; %s", err)
		return nil, err
	}

	ts := time.Now().Unix()
	var fileNames []string
	defer cleanupFiles(&fileNames)

	tmpStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		log.Printf("[miner] error generating tmp state file; %s", err)
		return nil, err
	}

	fileNames = append(fileNames, tmpStateFile.Name())
	if _, err := tmpStateFile.Write(genesisState); err != nil {
		log.Printf("[miner] error writing to genesis state to tmp state file; %s", err)
		return nil, err
	}
	if err := tmpStateFile.Close(); err != nil {
		log.Printf("[miner] error closing tmp state file; %s", err)
		return nil, err
	}

	combinedPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/combined.patch", imageHash, ts))
	if err != nil {
		log.Printf("[miner] error creating combined patch file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, combinedPatchFile.Name())
	if _, err := combinedPatchFile.Write(combinedDiff); err != nil {
		log.Printf("[miner] error writing combined diff to combined patch file; %s", err)
		return nil, err
	}
	if err := combinedPatchFile.Close(); err != nil {
		log.Printf("[miner] error closing combined patch file; %s", err)
		return nil, err
	}

	// now apply the combined patch file to the state
	if err := diffing.Patch(combinedPatchFile.Name(), false, true); err != nil {
		log.Printf("[miner] error diffing combined patch file; %s", err)
		return nil, err
	}
	state, err := ioutil.ReadFile(tmpStateFile.Name())
	if err != nil {
		log.Printf("[miner] error reading tmp state file; %s", err)
		return nil, err
	}

	return state, nil
}

func generateCombinedDiffs(ctx context.Context, imageHash string, genesisState []byte, diffs []*statechain.Diff) ([]byte, error) {
	ts := time.Now().Unix()
	var fileNames []string
	defer cleanupFiles(&fileNames)

	tmpStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		log.Printf("[miner] error creating tmp state file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, tmpStateFile.Name())
	if _, err := tmpStateFile.Write(genesisState); err != nil {
		log.Printf("[miner] error writing genesis state to tmp state file; %s", err)
		return nil, err
	}
	if err := tmpStateFile.Close(); err != nil {
		log.Printf("[miner] error closing tmp state file; %s", err)
		return nil, err
	}

	combinedPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/combined.patch", imageHash, ts))
	if err != nil {
		log.Printf("[miner] error creating combined patch file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, combinedPatchFile.Name())
	if err := combinedPatchFile.Close(); err != nil {
		log.Printf("[miner] error closing combined patch file; %s", err)
		return nil, err
	}

	tmpPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/tmp.patch", imageHash, ts))
	if err != nil {
		log.Printf("[miner] error creating tmp patch file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, tmpPatchFile.Name())
	if err := tmpPatchFile.Close(); err != nil {
		log.Printf("[miner] error closing tmp patch file; %s", err)
		return nil, err
	}

	for _, diff := range diffs {
		if ctx.Err() != nil {
			log.Printf("[miner] error diffing; %s", err)
			return nil, ctx.Err()
		}

		if err := ioutil.WriteFile(tmpPatchFile.Name(), []byte(diff.Props().Data), os.ModePerm); err != nil {
			log.Printf("[miner] error writing to tmp patch file; %s", err)
			return nil, err
		}

		if err := diffing.CombineDiff(combinedPatchFile.Name(), tmpPatchFile.Name(), combinedPatchFile.Name()); err != nil {
			log.Printf("[miner] error invoking diffing combined diff with combined patch file, tmp patch file and combined patch file; %s", err)
			return nil, err
		}
	}

	patch, err := ioutil.ReadFile(combinedPatchFile.Name())
	if err != nil {
		log.Printf("[miner] error reading combined patch file; %s", err)
		return nil, err
	}

	return patch, nil
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

// note: the transactions array are all for the image hash, passed
func isGenesisTransaction(p2pSvc p2p.Interface, prevBlock *mainchain.Block, imageHash string, transactions []*statechain.Transaction) (bool, *statechain.Transaction, []*statechain.Transaction, error) {
	for idx, tx := range transactions {
		log.Printf("[miner] state block tx method %s", tx.Props().Method)
		if tx.Props().Method == methodTypes.Deploy {
			prevStateBlock, _ := p2pSvc.FetchMostRecentStateBlock(imageHash, prevBlock)
			if prevStateBlock != nil {
				log.Printf("[miner] prev state block exists image hash %s", imageHash)
				return false, nil, nil, errors.New("prev state block exists; can't deploy")
			}

			return true, tx, append(transactions[:idx], transactions[idx+1:]...), nil
		}
	}

	return false, nil, transactions, nil
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

func cleanupFiles(fileNames *[]string) {
	if fileNames == nil {
		return
	}
	for _, fileName := range *fileNames {
		if err := os.Remove(fileName); err != nil {
			log.Printf("[miner] err cleaning up file %s", fileName)
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
