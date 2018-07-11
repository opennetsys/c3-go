package miner

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/merkle"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/diffing"
	"github.com/c3systems/c3/core/p2p"
)

// New returns a new service
func New(props *Props) (*Service, error) {
	if props == nil {
		return nil, errors.New("props are required")
	}

	statechainBlocksMap := make(map[string]*statechain.Block)
	transactionsMap := make(map[string]*statechain.Transaction)
	diffsMap := make(map[string]*statechain.Diff)
	merkleTreesMap := make(map[string]*merkle.Tree)

	s := &Service{
		props: *props,
		minedBlock: &MinedBlock{
			NextBlock:           nil,
			PreviousBlock:       props.PreviousBlock,
			mut:                 sync.Mutex{},
			StatechainBlocksMap: statechainBlocksMap,
			TransactionsMap:     transactionsMap,
			DiffsMap:            diffsMap,
			MerkleTreesMap:      merkleTreesMap,
		},
	}

	nextBlock, err := s.bootstrapNextBlock()
	if err != nil {
		return nil, err
	}
	s.minedBlock.NextBlock = nextBlock

	return s, nil
}

// Props returns the props
func (s Service) Props() Props {
	return s.props
}

// SpawnMiner ...
func (s Service) SpawnMiner() error {
	// TODO: reward ourselves with some coin
	go func() {
		var (
			err error
		)

		switch s.props.Async {
		case true:
			err = s.buildMainchainBlockAsync()
			if err != nil {
				s.props.Channel <- err
				return
			}

		default:
			err = s.buildMainchainBlock()
			if err != nil {
				s.props.Channel <- err
				return
			}
		}

		if s.minedBlock == nil || s.minedBlock.NextBlock == nil {
			s.props.Channel <- errors.New("built a nil block")
			return
		}

		s.props.Channel <- s.minedBlock
	}()

	return nil
}

func (s Service) buildMainchainBlockAsync() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	// 1. gather tx's
	// TODO: only choose high value tx's to mine
	txsMap := BuildTxsMap(s.props.PendingTransactions)

	// 2. apply txs
	for imageHash, transactions := range txsMap {
		wg.Add(1)
		go func(iHash string, txs []*statechain.Transaction) {
			defer wg.Done()
			if s.props.IsValid == nil || *s.props.IsValid == false {
				err = errors.New("miner is invalid")
				return
			}

			if err1 := s.buildNextStates(iHash, txs); err1 != nil {
				// err = err1 note: don't do this, we'll just skip this image hash
				log.Printf("[miner] err mining state block for hash %s transactions %v: %v", iHash, txs, err1)
				return
			}
		}(imageHash, transactions)
	}
	wg.Wait()

	if err != nil {
		return err
	}

	// 3. mine main block
	return s.mineBlock()
}

func (s Service) buildMainchainBlock() error {
	// 1. gather tx's
	// TODO: only choose high value tx's to mine
	txsMap := BuildTxsMap(s.props.PendingTransactions)

	// 2. apply txs
	for imageHash, transactions := range txsMap {
		if s.props.IsValid == nil || *s.props.IsValid == false {
			return errors.New("miner is invalid")
		}

		if err := s.buildNextStates(imageHash, transactions); err != nil {
			log.Printf("[miner] err mining state block for hash %s transactions %v: %v", imageHash, transactions, err)
			continue
		}
	}

	// 3. mine main block
	return s.mineBlock()
}

func (s Service) mineBlock() error {
	// TODO: timeout?
	if err := s.generateMerkle(); err != nil {
		return err
	}

	for {
		if s.props.IsValid == nil || *s.props.IsValid == false {
			return errors.New("miner is invalid")
		}

		hash, _, err := s.generateHashAndNonce()
		if err != nil {
			return err
		}

		check, err := CheckHashAgainstDifficulty(hash, s.props.Difficulty)
		if err != nil {
			return err
		}

		if check {
			return s.minedBlock.NextBlock.SetHash()
		}
	}
}

func (s Service) generateMerkle() error {
	var hashes []string

	s.minedBlock.mut.Lock()
	defer s.minedBlock.mut.Unlock()
	for _, statechainBlock := range s.minedBlock.StatechainBlocksMap {
		if statechainBlock == nil {
			return errors.New("nil block")
		}
		if statechainBlock.Props().BlockHash == nil {
			return errors.New("nil block hash")
		}

		hashes = append(hashes, *statechainBlock.Props().BlockHash)
	}

	tree, err := merkle.New(&merkle.TreeProps{
		Hashes: hashes,
		Kind:   merkle.StatechainBlocksKindStr,
	})
	if err != nil {
		return err
	}
	if tree == nil {
		return errors.New("nil tree")
	}
	if tree.Props().MerkleTreeRootHash == nil {
		return errors.New("nil merkle root hash")
	}

	s.minedBlock.MerkleTreesMap[*tree.Props().MerkleTreeRootHash] = tree

	nextProps := s.minedBlock.NextBlock.Props()
	nextProps.StateBlocksMerkleHash = *tree.Props().MerkleTreeRootHash
	nextBlock := mainchain.New(&nextProps)
	s.minedBlock.NextBlock = nextBlock

	return nil
}

func (s Service) generateHashAndNonce() (string, string, error) {
	nonce, err := s.generateNonce()
	if err != nil {
		return "", "", err
	}

	tmpProps := s.minedBlock.NextBlock.Props()
	tmpProps.Nonce = nonce
	tmpBlock := mainchain.New(&tmpProps)

	hash, err := tmpBlock.CalculateHash()
	return hash, nonce, err
}

func (s Service) generateNonce() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

func (s Service) bootstrapNextBlock() (*mainchain.Block, error) {
	nextProps := new(mainchain.Props)

	prevProps := s.props.PreviousBlock.Props()
	if prevProps.BlockHash == nil {
		return nil, errors.New("previous block's block hash is nil")
	}

	// note: checked for nil block hash, above
	nextProps.PrevBlockHash = *prevProps.BlockHash
	prevBlockHeight, err := hexutil.DecodeUint64(prevProps.BlockNumber)
	if err != nil {
		return nil, err
	}
	nextProps.BlockNumber = hexutil.EncodeUint64(prevBlockHeight + 1)
	nextProps.BlockTime = hexutil.EncodeUint64(uint64(time.Now().Unix()))
	nextProps.Difficulty = hexutil.EncodeUint64(s.props.Difficulty)

	return mainchain.New(nextProps), nil
}

func (s Service) buildNextStates(imageHash string, transactions []*statechain.Transaction) error {
	// TODO: add miguel's code, here
	var (
		diffs []*statechain.Diff

		newDiffs            []*statechain.Diff
		newTxs              []*statechain.Transaction
		newStatechainBlocks []*statechain.Block
	)

	// fetch the most recent state block
	prevStateBlock, err := s.props.P2P.FetchMostRecentStateBlock(imageHash, s.props.PreviousBlock)
	if err != nil {
		return err
	}
	// TODO: upload a genesis block if nil
	if prevStateBlock == nil {
		return errors.New("could not find most recent state block")
	}
	if prevStateBlock.Props().BlockHash == nil {
		return errors.New("prev state block hash is nil")
	}

	// gather the diffs
	diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
	if err != nil {
		return err
	}
	diff, err := s.props.P2P.GetStatechainDiff(diffCID)
	if err != nil {
		return err
	}
	// note: prepend
	diffs = append([]*statechain.Diff{diff}, diffs...)

	head := prevStateBlock
	for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
		prevStateCID, err := p2p.GetCIDByHash(head.Props().PrevBlockHash)
		if err != nil {
			return err
		}

		prevStateBlock, err := s.props.P2P.GetStatechainBlock(prevStateCID)
		if err != nil {
			return err
		}
		head = prevStateBlock

		diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
		if err != nil {
			return err
		}
		diff, err := s.props.P2P.GetStatechainDiff(diffCID)
		if err != nil {
			return err
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
		return err
	}
	defer os.Remove(tmpStateFile.Name()) // clean up

	if _, err := tmpStateFile.Write([]byte(genesisState)); err != nil {
		return err
	}
	if err := tmpStateFile.Close(); err != nil {
		return err
	}

	outPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/combined.txt", imageHash, ts))
	if err != nil {
		return err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err := outPatchFile.Close(); err != nil {
		return err
	}

	for i, diff := range diffs {
		tmpPatchFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/patch.%d.txt", imageHash, ts, i))
		if err != nil {
			return err
		}
		defer os.Remove(tmpPatchFile.Name()) // clean up

		if _, err := tmpPatchFile.Write([]byte(diff.Props().Data)); err != nil {
			return err
		}
		if err := tmpPatchFile.Close(); err != nil {
			return err
		}

		if err := diffing.CombineDiff(outPatchFile.Name(), tmpPatchFile.Name(), outPatchFile.Name()); err != nil {
			return err
		}
	}

	// now apply the combined patch file to the state
	if err := diffing.Patch(outPatchFile.Name(), false, true); err != nil {
		return err
	}
	state, err := ioutil.ReadFile(tmpStateFile.Name())
	if err != nil {
		return err
	}

	log.Printf("state\n%s", string(state))
	headStateFileName := tmpStateFile.Name()
	runningBlockNumber, err := hexutil.DecodeUint64(prevStateBlock.Props().BlockNumber)
	if err != nil {
		return err
	}
	runningBlockHash := *prevStateBlock.Props().BlockHash // note: already checked nil pointer, above
	// TODO: apply state to container and start running transactions
	for i, tx := range transactions {
		if tx == nil {
			return errors.New("nil tx")
		}
		if tx.Props().TxHash == nil {
			return errors.New("nil tx hash")
		}

		nextState := ""
		nextStateFile, err := ioutil.TempFile("", fmt.Sprintf("%s/%v/state.%d.txt", imageHash, ts, i))
		if err != nil {
			return err
		}
		defer os.Remove(nextStateFile.Name()) // clean up

		if _, err := nextStateFile.Write([]byte(nextState)); err != nil {
			return err
		}
		if err := nextStateFile.Close(); err != nil {
			return err
		}

		if err := diffing.Diff(headStateFileName, nextStateFile.Name(), outPatchFile.Name(), false); err != nil {
			return err
		}
		headStateFileName = nextStateFile.Name()

		// build the diff struct
		diffData, err := ioutil.ReadFile(outPatchFile.Name())
		if err != nil {
			return err
		}

		diffStruct := statechain.NewDiff(&statechain.DiffProps{
			Data: string(diffData),
		})
		if err := diffStruct.SetHash(); err != nil {
			return err
		}

		nextStateHash := hexutil.EncodeString(nextState)
		runningBlockNumber++
		nextStateStruct := statechain.New(&statechain.BlockProps{
			BlockNumber:       hexutil.EncodeUint64(runningBlockNumber),
			BlockTime:         hexutil.EncodeUint64(uint64(ts)),
			ImageHash:         imageHash,
			TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
			PrevBlockHash:     runningBlockHash,
			StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
			StateCurrentHash:  nextStateHash,
		})
		if err := nextStateStruct.SetHash(); err != nil {
			return err
		}
		runningBlockHash = *nextStateStruct.Props().BlockHash

		newDiffs = append(newDiffs, diffStruct)
		newTxs = append(newTxs, tx)
		newStatechainBlocks = append(newStatechainBlocks, nextStateStruct)
	}

	// write to the mined block
	s.minedBlock.mut.Lock()
	// note: they should all have same length
	for i := 0; i < len(newDiffs); i++ {
		s.minedBlock.DiffsMap[*newDiffs[i].Props().DiffHash] = newDiffs[i]
		s.minedBlock.TransactionsMap[*newTxs[i].Props().TxHash] = newTxs[i]
		s.minedBlock.StatechainBlocksMap[*newStatechainBlocks[i].Props().BlockHash] = newStatechainBlocks[i]
	}
	s.minedBlock.mut.Unlock()

	return nil
}
