package miner

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/common/fileutil"
	"github.com/c3systems/c3-go/common/hashutil"
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/merkle"
	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/core/diffing"
	"github.com/c3systems/c3-go/core/p2p"
	"github.com/c3systems/c3-go/core/sandbox"
	methodTypes "github.com/c3systems/c3-go/core/types/methods"
	colorlog "github.com/c3systems/c3-go/log/color"

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

	hashStr, err := hexutil.RemovePrefix(hashHex)
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
		return false, ErrInvalidTx
	}

	// 2. tx must have a sig
	if tx.Props().Sig == nil {
		return false, ErrInvalidTx
	}

	// 3. verify the hash
	tmpHash, err := tx.CalculateHash()
	if err != nil {
		return false, err
	}

	// note: already checked for nil hash
	if *tx.Props().TxHash != tmpHash {
		return false, ErrInvalidTx
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

	ok, err := c3crypto.Verify(pub, []byte(*tx.Props().TxHash), r, s)
	if !ok || err != nil {
		return false, ErrInvalidTx
	}

	return true, nil
}

// VerifyMinedBlock ...
func VerifyMinedBlock(ctx context.Context, p2pSvc p2p.Interface, sbSvc sandbox.Interface, minedBlock *MinedBlock) (bool, error) {
	ch := make(chan interface{})

	go func() {
		if minedBlock == nil {
			log.Error("[miner] mined block is nil")
			ch <- false

			return
		}
		if minedBlock.NextBlock == nil {
			log.Error("[miner] next block is nil")
			ch <- false

			return
		}
		if minedBlock.PreviousBlock == nil {
			log.Error("[miner] mined block is nil")
			ch <- false

			return
		}
		if minedBlock.NextBlock.Props().BlockHash == nil {
			log.Error("[miner] next block blockhash is nil")
			ch <- false

			return
		}
		if minedBlock.PreviousBlock.Props().BlockHash == nil {
			log.Error("[miner] prev block block hash is nil")
			ch <- false

			return
		}
		// note checked for nil pointer, above
		if *minedBlock.PreviousBlock.Props().BlockHash != minedBlock.NextBlock.Props().PrevBlockHash {
			log.Error("[miner] prev block block hash != next block block hash")
			ch <- false

			return
		}
		if mainchain.ImageHash != minedBlock.NextBlock.Props().ImageHash {
			log.Error("[miner] mainchain imagehash != nextblock image hash")
			ch <- false

			return
		}
		if minedBlock.NextBlock.Props().MinerSig == nil {
			log.Error("[miner] next block miner sig is nil")
			ch <- false

			return
		}

		ok, err := CheckBlockHashAgainstDifficulty(minedBlock.NextBlock)
		if err != nil {
			log.Errorf("[miner] err checking block hash against difficulty\n%v", err)
			ch <- err

			return
		}
		if !ok {
			log.Error("[miner] block hash did not checkout against difficulty")
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		tmpHash, err := minedBlock.NextBlock.CalculateHash()
		if err != nil {
			log.Errorf("[miner] err calculating hash\n%v", err)
			ch <- err

			return
		}
		// note: already checked for nil hash
		if *minedBlock.NextBlock.Props().BlockHash != tmpHash {
			log.Errorf("[miner] next block hash != calced hash\n%s\n%s", *minedBlock.NextBlock.Props().BlockHash, tmpHash)
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		pub, err := c3crypto.DecodeAddress(minedBlock.NextBlock.Props().MinerAddress)
		if err != nil {
			log.Errorf("[miner] err decoding addr\n%v", err)
			ch <- err

			return
		}

		// note: checked for nil sig, above
		sigR, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.R)
		if err != nil {
			log.Errorf("[miner] err decoding r\n%v", err)
			ch <- err

			return
		}
		sigS, err := hexutil.DecodeBigInt(minedBlock.NextBlock.Props().MinerSig.S)
		if err != nil {
			log.Errorf("[miner] err decoding s\n%v", err)
			ch <- err

			return
		}

		// note: nil blockhash was checked, above
		ok, err = c3crypto.Verify(pub, []byte(*minedBlock.NextBlock.Props().BlockHash), sigR, sigS)
		if err != nil {
			log.Errorf("[miner] err verifying miner sig\n%v", err)
			ch <- err

			return
		}
		if !ok {
			log.Error("[miner] block hash did not checkout agains sig")
			ch <- false

			return
		}
		if ctx.Err() != nil {
			return
		}

		// BlockNumber must be +1 prev block number
		blockNumber, err := hexutil.DecodeUint64(minedBlock.NextBlock.Props().BlockNumber)
		if err != nil {
			log.Errorf("[miner] err decoding block #\n%v", err)
			ch <- err

			return
		}
		prevNumber, err := hexutil.DecodeUint64(minedBlock.PreviousBlock.Props().BlockNumber)
		if err != nil {
			log.Errorf("[miner] err decoding prev block #\n%v", err)
			ch <- err

			return
		}
		if prevNumber+1 != blockNumber {
			log.Error("[miner] prevBlockNumber +1 != nextBlockNumber")
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

			return VerifyStateBlocksFromMinedBlock(ctx, p2pSvc, sbSvc, minedBlock)

		default:
			log.Errorf("[miner] received unknown message of type %T\n%v", v, v)
			return false, errors.New("received message of unknown type")

		}

	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// VerifyStateBlocksFromMinedBlock ...
// note: this function also checks the merkle tree. That check is not required to be performed, separately.
func VerifyStateBlocksFromMinedBlock(ctx context.Context, p2pSvc p2p.Interface, sbSvc sandbox.Interface, minedBlock *MinedBlock) (bool, error) {
	ch := make(chan interface{})

	go func() {
		if minedBlock.NextBlock == nil {
			log.Error("[miner] nil next block")
			ch <- false

			return
		}
		if len(minedBlock.StatechainBlocksMap) != len(minedBlock.TransactionsMap) {
			log.Error("[miner] len state blocks map != len tx map")
			ch <- false

			return
		}
		// note: ok to have nil map? e.g. in the case that no transactions were included in the main block
		//if minedBlock.StatechainBlocksMap == nil {
		//log.Println("nil state blocks map")
		//return false, nil
		//}
		if minedBlock.MerkleTreesMap == nil {
			log.Error("[miner] nil merkle trees map")
			ch <- false

			return
		}

		// 1. Verify state blocks merkle hash
		ok, err := VerifyMerkleTreeFromMinedBlock(ctx, minedBlock)
		if err != nil {
			log.Errorf("[miner] err verifying merkle tree\n%v", err)
			ch <- err

			return
		}
		if !ok {
			log.Error("[miner] merkle tree didn't verify")
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
			log.Errorf("[miner] err grouping state blocks\n%v", err)
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
				log.Errorf("[miner] err ordering state blocks\n%v", err)
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
				log.Errorf("[miner] err getting cid by has\n%v", err)
				ch <- err

				return
			}
			// TODO: check that this is the actual prev block on the blockchain
			prevBlock, err := p2pSvc.GetStatechainBlock(prevBlockCID)
			if err != nil {
				log.Errorf("[miner] err getting state block\n%v", err)
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
				log.Errorf("[miner] err fetching current state\n%v", err)
				ch <- err

				return
			}
			prevStateHash := hashutil.HashToHexString(prevState)
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
					log.Errorf("[miner] err decoding prev block # 2b\n%v", err)
					ch <- err

					return
				}
				blockNumber, err := hexutil.DecodeUint64(block.Props().BlockNumber)
				if err != nil {
					log.Errorf("[miner] err decoding block # 2b\n %v", err)
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
					log.Errorf("[miner] err calculating block hash 2c\n%v", err)
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
					log.Errorf("[miner] err verifying tx 2d\n %v", err)
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

				nextStateBlock, nextDiff, nextState, err := buildNextStateFromPrevState(p2pSvc, sbSvc, prevState, block, tx)
				if err != nil {
					log.Errorf("[miner] err building next state from prev state\n %v", err)
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
			log.Errorf("[miner] received unknown message of type %T\n%v", v, v)
			return false, errors.New("received message of unknown type")

		}

	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func buildNextStateFromPrevState(p2pSvc p2p.Interface, sbSvc sandbox.Interface, prevState []byte, prevBlock *statechain.Block, tx *statechain.Transaction) (*statechain.Block, *statechain.Diff, []byte, error) {
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
	outPatchFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/diff.patch", prevBlock.Props().ImageHash, ts))
	if err != nil {
		return nil, nil, nil, err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err = outPatchFile.Close(); err != nil {
		return nil, nil, nil, err
	}
	prevStateFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/%s", prevBlock.Props().ImageHash, ts, StateFileName))
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

		// run container, passing the tx inputs
		// note: certain err's, here, should remove the tx from the pending tx pool
		nextState, err = sbSvc.Play(&sandbox.PlayConfig{
			ImageID:      tx.Props().ImageHash,
			Payload:      payload,
			InitialState: prevState,
		})

		if err != nil {
			return nil, nil, nil, err
		}

		//log.Printf("[miner] container new state: %s", string(nextState))
		nextStateFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/nextState.txt", prevBlock.Props().ImageHash, ts))
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

		nextStateHashBytes := hashutil.Hash(nextState)
		nextStateHash := hexutil.EncodeToString(nextStateHashBytes[:])
		log.Printf("[miner] state prev diff hash: %s", *diffStruct.Props().DiffHash)
		log.Printf("[miner] state current hash: %s", nextStateHash)
		nextStateStruct := statechain.New(&statechain.BlockProps{
			BlockNumber:       hexutil.EncodeUint64(prevBlockNumber),
			BlockTime:         hexutil.EncodeUint64(uint64(ts)),
			ImageHash:         prevBlock.Props().ImageHash,
			TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
			PrevBlockHash:     *prevBlock.Props().BlockHash,
			StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
			StateCurrentHash:  nextStateHash,
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
		log.Errorf("[miner] tx hash is nil for %v", tx.Props())
		return nil, nil, errors.New("nil tx hash")
	}

	if tx.Props().TxHash == nil {
		log.Errorf("[miner] tx hash is nil for %v", tx.Props())
		return nil, nil, errors.New("nil tx hash")
	}

	log.Printf("[miner] tx method %s", tx.Props().Method)

	// initial state
	nextState := tx.Props().Payload
	log.Printf("[miner] container initial state: %s", string(nextState))

	stateFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/%s", imageHash, ts, StateFileName))
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(stateFile.Name()) // clean up
	if err = stateFile.Close(); err != nil {
		return nil, nil, err
	}

	nextStateFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/nextState.txt", imageHash, ts))
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

	outPatchFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/diff.patch", imageHash, ts))
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err = outPatchFile.Close(); err != nil {
		return nil, nil, err
	}

	log.Printf("[miner] diffing the files:\ntmp state: %s\nnext state: %s\nout patch: %s", stateFile.Name(), nextStateFile.Name(), outPatchFile.Name())

	if err = diffing.Diff(stateFile.Name(), nextStateFile.Name(), outPatchFile.Name(), false); err != nil {
		return nil, nil, err
	}

	// build the diff struct
	diffData, err := ioutil.ReadFile(outPatchFile.Name())
	if err != nil {
		return nil, nil, err
	}

	log.Println(colorlog.Yellow("[miner] diff data from patch file: %s", string(diffData)))

	diffStruct := statechain.NewDiff(&statechain.DiffProps{
		Data: string(diffData),
	})
	if err := diffStruct.SetHash(); err != nil {
		return nil, nil, err
	}

	nextStateHashBytes := hashutil.Hash(nextState)
	nextStateHash := hexutil.EncodeToString(nextStateHashBytes[:])
	log.Printf("[miner] state prev diff hash: %s", *diffStruct.Props().DiffHash)
	log.Printf("[miner] state current hash: %s", nextStateHash)
	nextStateStruct := statechain.New(&statechain.BlockProps{
		BlockNumber:       hexutil.EncodeUint64(0),
		BlockTime:         hexutil.EncodeUint64(uint64(ts)),
		ImageHash:         imageHash,
		TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
		PrevBlockHash:     "0x",
		StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
		StateCurrentHash:  nextStateHash,
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
			log.Errorf("[miner] block is nil")
			ch <- errors.New("nil block")

			return
		}
		if block.Props().BlockHash == nil {
			log.Errorf("[miner] block hash is nil")
			ch <- errors.New("nil block hash")

			return
		}
		if ctx.Err() != nil {
			return
		}

		// gather the diffs
		diffs, err := GatherDiffs(ctx, p2pSvc, block)
		if err != nil {
			log.Errorf("[miner] error gathering diffs\n%v", err)
			ch <- err

			return
		}

		// apply the diffs to get the current state
		// TODO: get the genesis state of the block?
		var genesisState []byte
		imageHash := block.Props().ImageHash
		state, err := GenerateStateFromDiffs(ctx, imageHash, genesisState, diffs)
		if err != nil {
			log.Errorf("[miner] error reading state file\n%v", err)
			ch <- err

			return
		}
		if ctx.Err() != nil {
			log.Errorf("[miner] received context error; %s", err)
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
			log.Errorf("[miner] channel error; %s", err)
			return nil, err

		case []byte:
			bytes, _ := v.([]byte)

			return bytes, nil

		default:
			log.Errorf("[miner] received unknown message of type %T\n%v", v, v)
			return nil, errors.New("received message of unknown type")

		}

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GatherDiffs ...
func GatherDiffs(ctx context.Context, p2pSvc p2p.Interface, block *statechain.Block) ([]*statechain.Diff, error) {
	var diffs []*statechain.Diff

	// gather the diffs
	diffCID, err := p2p.GetCIDByHash(block.Props().StatePrevDiffHash)
	if err != nil {
		log.Errorf("[miner] err getting diff cid by hash\n%v", err)
		return nil, err
	}

	diff, err := p2pSvc.GetStatechainDiff(diffCID)
	if err != nil {
		log.Errorf("[miner] err getting diff by cid\n%v", err)
		return nil, err
	}

	log.Printf("[miner] state chain diff; %v", diff)

	// note: prepend
	diffs = append([]*statechain.Diff{diff}, diffs...)

	head := block
	for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
		if ctx.Err() != nil {
			log.Errorf("[miner] gather diffs context error; %v", err)
			return nil, ctx.Err()
		}

		prevStateCID, err := p2p.GetCIDByHash(head.Props().PrevBlockHash)
		if err != nil {
			log.Errorf("[miner] err getting statechain cid by hash\n%v", err)
			return nil, err
		}

		prevStateBlock, err := p2pSvc.GetStatechainBlock(prevStateCID)
		if err != nil {
			log.Errorf("[miner] err getting state chain block by cid\n%v", err)
			return nil, err
		}

		head = prevStateBlock

		diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
		if err != nil {
			log.Errorf("[miner] err getting diffCID by hash\n%v", err)
			return nil, err
		}

		diff, err := p2pSvc.GetStatechainDiff(diffCID)
		if err != nil {
			log.Errorf("[miner] err gitting diff by cid\n%v", err)
			return nil, err
		}

		// note: prepend
		diffs = append([]*statechain.Diff{diff}, diffs...)
	}

	if diffs == nil {
		log.Error("[miner] error; diffs is nil")
		return nil, errors.New("diffs is nil")
	}

	return diffs, nil
}

// GenerateStateFromDiffs ...
func GenerateStateFromDiffs(ctx context.Context, imageHash string, genesisState []byte, diffs []*statechain.Diff) ([]byte, error) {
	combinedDiff, err := generateCombinedDiffs(ctx, imageHash, diffs)
	if err != nil {
		log.Errorf("[miner] error generating combined diffs; %s", err)
		return nil, err
	}

	ts := time.Now().Unix()
	var fileNames []string
	defer cleanupFiles(&fileNames)

	tmpStateFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/%s", imageHash, ts, StateFileName))
	if err != nil {
		log.Errorf("[miner] error generating tmp state file; %s", err)
		return nil, err
	}

	fileNames = append(fileNames, tmpStateFile.Name())
	if _, err := tmpStateFile.Write(genesisState); err != nil {
		log.Errorf("[miner] error writing to genesis state to tmp state file; %s", err)
		return nil, err
	}
	if err := tmpStateFile.Close(); err != nil {
		log.Errorf("[miner] error closing tmp state file; %s", err)
		return nil, err
	}

	combinedPatchFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/combined.patch", imageHash, ts))
	if err != nil {
		log.Errorf("[miner] error creating combined patch file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, combinedPatchFile.Name())
	if _, err := combinedPatchFile.Write(combinedDiff); err != nil {
		log.Errorf("[miner] error writing combined diff to combined patch file; %s", err)
		return nil, err
	}
	if err := combinedPatchFile.Close(); err != nil {
		log.Errorf("[miner] error closing combined patch file; %s", err)
		return nil, err
	}

	// now apply the combined patch file to the state
	if err := diffing.Patch(combinedPatchFile.Name(), tmpStateFile.Name(), false, true); err != nil {
		log.Errorf("[miner] error applying combined patch file; %s", err)
		return nil, err
	}
	state, err := ioutil.ReadFile(tmpStateFile.Name())
	if err != nil {
		log.Errorf("[miner] error reading tmp state file; %s", err)
		return nil, err
	}

	return state, nil
}

func generateCombinedDiffs(ctx context.Context, imageHash string, diffs []*statechain.Diff) ([]byte, error) {
	ts := time.Now().Unix()
	var fileNames []string
	defer cleanupFiles(&fileNames)

	if diffs == nil || len(diffs) == 0 {
		return nil, errors.New("nil diffs")
	}

	combinedPatchFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/combined.patch", imageHash, ts))
	if err != nil {
		log.Errorf("[miner] error creating combined patch file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, combinedPatchFile.Name())
	if _, err := combinedPatchFile.Write([]byte(diffs[0].Props().Data)); err != nil {
		log.Errorf("[miner] error writing to genesis state to tmp state file; %s", err)
		return nil, err
	}
	if err := combinedPatchFile.Close(); err != nil {
		log.Errorf("[miner] error closing combined patch file; %s", err)
		return nil, err
	}

	tmpPatchFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/tmp.patch", imageHash, ts))
	if err != nil {
		log.Errorf("[miner] error creating tmp patch file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, tmpPatchFile.Name())
	if err := tmpPatchFile.Close(); err != nil {
		log.Errorf("[miner] error closing tmp patch file; %s", err)
		return nil, err
	}
	tmpPatchFile1, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/tmp1.patch", imageHash, ts))
	if err != nil {
		log.Errorf("[miner] error creating tmp patch file; %s", err)
		return nil, err
	}
	fileNames = append(fileNames, tmpPatchFile1.Name())
	if err := tmpPatchFile1.Close(); err != nil {
		log.Errorf("[miner] error closing tmp patch file 1; %s", err)
		return nil, err
	}

	for i := 1; i < len(diffs); i++ {
		if ctx.Err() != nil {
			log.Errorf("[miner] error diffing; %s", err)
			return nil, ctx.Err()
		}

		prevComb, err := ioutil.ReadFile(combinedPatchFile.Name())
		if err != nil {
			log.Errorf("[miner] error reading previous combined patch file; %s", err)
			return nil, err
		}
		if err := ioutil.WriteFile(tmpPatchFile.Name(), prevComb, os.ModePerm); err != nil {
			log.Errorf("[miner] error writing to tmp patch file; %s", err)
			return nil, err
		}

		if err := ioutil.WriteFile(tmpPatchFile1.Name(), []byte(diffs[i].Props().Data), os.ModePerm); err != nil {
			log.Errorf("[miner] error writing to tmp patch file; %s", err)
			return nil, err
		}

		// reading for debug logging
		combinedPathFileData, err := ioutil.ReadFile(combinedPatchFile.Name())
		if err != nil {
			return nil, err
		}
		// reading for debug logging
		tmpPatchFileData, err := ioutil.ReadFile(tmpPatchFile.Name())
		if err != nil {
			return nil, err
		}

		log.Printf("[miner] combining diffs\n%s\n%s\n%s\n%s\nout: %s", combinedPatchFile.Name(), string(combinedPathFileData), tmpPatchFile.Name(), string(tmpPatchFileData), combinedPatchFile.Name())

		if err := diffing.CombineDiff(tmpPatchFile.Name(), tmpPatchFile1.Name(), combinedPatchFile.Name()); err != nil {
			log.Errorf("[miner] error invoking diffing combined diff with combined patch file, tmp patch file and combined patch file; %s", err)
			return nil, err
		}

		// reading for debug logging
		combinedPathFileData, err = ioutil.ReadFile(combinedPatchFile.Name())
		if err != nil {
			return nil, err
		}

		log.Printf("[miner] combined diffs\n%s\n%s", combinedPatchFile.Name(), string(combinedPathFileData))
	}

	patch, err := ioutil.ReadFile(combinedPatchFile.Name())
	if err != nil {
		log.Errorf("[miner] error reading combined patch file; %s", err)
		return nil, err
	}

	return patch, nil
}

// VerifyMerkleTreeFromMinedBlock ...
func VerifyMerkleTreeFromMinedBlock(ctx context.Context, minedBlock *MinedBlock) (bool, error) {
	if minedBlock.NextBlock == nil {
		log.Errorf("[miner] verify mined block error - block data is nil")
		return false, nil
	}
	/*
		// note: mined block with no tx will have nil state chain blocks map
		if minedBlock.StatechainBlocksMap == nil {
			log.Errorf("[miner] verify mined block error - state chain blocks is nil")
			return false, nil
		}
	*/
	if minedBlock.MerkleTreesMap == nil {
		log.Errorf("[miner] verify mined block error - merkle tree map is nil")
		return false, nil
	}

	tree, ok := minedBlock.MerkleTreesMap[minedBlock.NextBlock.Props().StateBlocksMerkleHash]
	if !ok || tree == nil {
		log.Errorf("[miner] verify mined block error - merkle trees map %s is not ok or nil", minedBlock.NextBlock.Props().StateBlocksMerkleHash)
		return false, nil
	}
	if tree.Props().MerkleTreeRootHash == nil {
		log.Errorf("[miner] verify mined block error - merkle tree root hash is nil")
		return false, nil
	}

	tmpTree, err := merkle.New(&merkle.TreeProps{
		Hashes: tree.Props().Hashes,
		Kind:   merkle.StatechainBlocksKindStr,
	})
	if err != nil {
		log.Errorf("[miner] verify mined block error - new merkle tree error: %s", err)
		return false, err
	}
	if err := tmpTree.SetHash(); err != nil {
		log.Errorf("[miner] verify mined block error - set hash error: %s", err)
		return false, err
	}

	if *tmpTree.Props().MerkleTreeRootHash != *tree.Props().MerkleTreeRootHash {
		log.Errorf("[miner] verify mined block error - merkle tree root hash doesn't match; %s", *tmpTree.Props().MerkleTreeRootHash)
		return false, nil
	}

	if len(tmpTree.Props().Hashes) != len(minedBlock.StatechainBlocksMap) {
		log.Errorf("[miner] verify mined block error - tree hashes length doesn't match; %v", len(tmpTree.Props().Hashes))
		return false, nil
	}

	for _, hash := range tmpTree.Props().Hashes {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		statechainBlock, ok := minedBlock.StatechainBlocksMap[hash]
		if !ok || statechainBlock == nil {
			log.Errorf("[miner] verify mined block error - state chain block from map is nil or not ok")
			return false, nil
		}

		tmpHash, err := statechainBlock.CalculateHash()
		if err != nil {
			log.Errorf("[miner] verify mined block error - state chain calculate hash error: %s", err)
			return false, err
		}
		if hash != tmpHash {
			log.Errorf("[miner] verify mined block error - hash does not match %s", tmpHash)
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
			// note: we need this err, what if the network errs and we allow another deploy?
			prevStateBlock, err := p2pSvc.FetchMostRecentStateBlock(imageHash, prevBlock)
			if err != nil {
				return false, nil, nil, err
			}
			if prevStateBlock != nil {
				log.Errorf("[miner] prev state block exists image hash %s", imageHash)
				return false, tx, nil, ErrInvalidTx
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

// input is a map with a key on block hash
// return is a map with keys on the image hash
func groupStateBlocksByImageHash(stateBlocksMap map[string]*statechain.Block) (map[string][]*statechain.Block, error) {
	/*
		// note: mined block with no transactions will have nil state block map
		if stateBlocksMap == nil {
			log.Errorf("[miner] state blocks map is nil")
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
	err := fileutil.RemoveFiles(fileNames)
	if err != nil {
		log.Errorf("[miner] err cleaning up files; %v", err)
	}
}
