package miner

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/merkle"
	"github.com/c3systems/c3/core/chain/statechain"
	"github.com/c3systems/c3/core/diffing"
	"github.com/c3systems/c3/core/p2p"
	"github.com/c3systems/c3/core/sandbox"
	methodTypes "github.com/c3systems/c3/core/types/methods"
	"github.com/c3systems/c3/logger"
	"github.com/c3systems/merkletree"
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

	log.Printf("[miner] build mainchain block async; tx count: %v", len(txsMap))

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

	log.Printf("[miner] build mainchain block; tx count: %v", len(txsMap))

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
		log.Printf("[miner] error mining block; %s", err)
		return err
	}

	for {
		if s.props.IsValid == nil || *s.props.IsValid == false {
			log.Println("[miner] miner is invalid")
			return errors.New("miner is invalid")
		}

		hash, nonce, err := s.generateHashAndNonce()
		if err != nil {
			log.Printf("[miner] error generating hash and nonce; %s", err)
			return err
		}

		check, err := CheckHashAgainstDifficulty(hash, s.props.Difficulty)
		if err != nil {
			log.Printf("[miner] error checking hash against difficulty; %s", err)
			return err
		}

		nextProps := s.minedBlock.NextBlock.Props()
		nextProps.Nonce = nonce
		nextBlock := mainchain.New(&nextProps)
		s.minedBlock.NextBlock = nextBlock

		if check {
			log.Println("[miner] difficulty checks out")
			return s.minedBlock.NextBlock.SetHash()
		}
	}
}

func (s Service) generateMerkle() error {
	log.Println("[miner] generating merkle")
	var (
		hashes []string
		list   []merkletree.Content
	)

	s.minedBlock.mut.Lock()
	defer s.minedBlock.mut.Unlock()

	log.Printf("[miner] state chain blocks length; %v", len(s.minedBlock.StatechainBlocksMap))
	for _, statechainBlock := range s.minedBlock.StatechainBlocksMap {
		if statechainBlock == nil {
			log.Println("[miner] state chain block is nil")
			return errors.New("nil block")
		}
		if statechainBlock.Props().BlockHash == nil {
			log.Println("[miner] state chain block hash is nil")
			return errors.New("nil block hash")
		}

		hashes = append(hashes, *statechainBlock.Props().BlockHash)
		list = append(list, statechainBlock)
	}

	tree, err := merkle.BuildFromObjects(list, merkle.StatechainBlocksKindStr)
	//tree, err := merkle.New(&merkle.TreeProps{
	//Hashes: hashes,
	//Kind:   merkle.StatechainBlocksKindStr,
	//})
	if err != nil {
		log.Printf("[miner] error building merkle tree from objects; %s", list)
		return err
	}
	if tree == nil {
		log.Println("[miner] tree is nil")
		return errors.New("nil tree")
	}
	if tree.Props().MerkleTreeRootHash == nil {
		log.Println("[miner] merkle root hash is nil")
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

	nextProps.BlockNumber = hexutil.EncodeUint64(0)
	nextProps.BlockTime = hexutil.EncodeUint64(uint64(time.Now().Unix()))
	nextProps.Difficulty = hexutil.EncodeUint64(s.props.Difficulty)
	nextProps.MinerAddress = s.props.EncodedMinerAddress

	// previous block will be nil if first block
	if s.props.PreviousBlock != nil {
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
	}

	return mainchain.New(nextProps), nil
}

func (s Service) isGenesisTransaction(imageHash string, transactions []*statechain.Transaction) (bool, error) {
	for _, tx := range transactions {
		if tx.Props().Method == methodTypes.Deploy {
			prevStateBlock, err := s.props.P2P.FetchMostRecentStateBlock(imageHash, s.props.PreviousBlock)
			if err != nil {
				log.Printf("[miner] error fetching most recent state block for image hash %s %s", imageHash, err)
				return false, err
			}
			if prevStateBlock != nil {
				log.Printf("[miner] prev state block exists image hash %s", imageHash)
				return false, errors.New("prev state block exists; can't deploy")
			}

			return true, nil
		}
	}

	return false, nil
}

func (s Service) buildNextStates(imageHash string, transactions []*statechain.Transaction) error {
	log.Printf("[miner] build next state state; image hash: %s, tx count: %v", imageHash, len(transactions))

	var (
		diffs               []*statechain.Diff
		newDiffs            []*statechain.Diff
		newTxs              []*statechain.Transaction
		newStatechainBlocks []*statechain.Block
	)

	ts := time.Now().Unix()

	// tmp hack to make genesis block work
	// TODO: move genesis tx at the top of the list
	isGenesisTx, err := s.isGenesisTransaction(imageHash, transactions)
	if err != nil {
		return err
	}
	if isGenesisTx {
		// TODO: process other transactions as well
		tx := transactions[0]
		return s.buildGenesisStateBlock(imageHash, tx)
	}

	prevStateBlock, err := s.props.P2P.FetchMostRecentStateBlock(imageHash, s.props.PreviousBlock)
	if err != nil {
		log.Printf("[miner] error fetching most recent state block for image hash %s %s", imageHash, err)
		return err
	}
	if prevStateBlock != nil {
		log.Printf("[miner] prev state block exists image hash %s", imageHash)
		return errors.New("prev state block exists; can't deploy")
	}

	// gather the diffs
	diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
	if err != nil {
		log.Printf("[miner] error getting cid by hash for image hash %s", imageHash)
		return err
	}

	diff, err := s.props.P2P.GetStatechainDiff(diffCID)
	if err != nil {
		log.Printf("[miner] error getting statechain diff for image hash %s", imageHash)
		return err
	}
	// note: prepend
	diffs = append([]*statechain.Diff{diff}, diffs...)
	log.Printf("[miner] diffs %v", len(diffs))

	head := prevStateBlock
	for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
		prevStateCID, err := p2p.GetCIDByHash(head.Props().PrevBlockHash)
		if err != nil {
			log.Printf("[miner] error getting cid by hash for %s for image hash %s", head.Props().PrevBlockHash, imageHash)
			return err
		}
		log.Printf("[miner] got prev state cid for image hash %s", imageHash)

		prevStateBlock, err := s.props.P2P.GetStatechainBlock(prevStateCID)
		if err != nil {
			return err
		}
		head = prevStateBlock
		log.Printf("[miner] set head to prev state block for image hash %s", imageHash)

		diffCID, err := p2p.GetCIDByHash(prevStateBlock.Props().StatePrevDiffHash)
		if err != nil {
			log.Printf("[miner] error getting cid by hash for prev state block diff hash for image hash %s", imageHash)
			return err
		}
		diff, err := s.props.P2P.GetStatechainDiff(diffCID)
		if err != nil {
			log.Printf("[miner] error getting state chain diff for diff cid for image hash %s", imageHash)
			return err
		}
		// note: prepend
		diffs = append([]*statechain.Diff{diff}, diffs...)
	}

	log.Printf("[miner] total diffs %v", len(diffs))

	// apply the diffs to get the current state
	// TODO: get the genesis state of the block
	genesisState := ""
	tmpStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		log.Printf("[miner] error creating tmp state file %s", err)
		return err
	}
	log.Printf("[miner] tmp state file %s", tmpStateFile.Name())
	defer os.Remove(tmpStateFile.Name()) // clean up

	if _, err := tmpStateFile.Write([]byte(genesisState)); err != nil {
		log.Printf("[miner] error writing genesis state to tmp state file %s", err)
		return err
	}
	if err := tmpStateFile.Close(); err != nil {
		log.Printf("[miner] error closing tmp state file %s", err)
		return err
	}

	outPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/combined.txt", imageHash, ts))
	if err != nil {
		return err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err := outPatchFile.Close(); err != nil {
		return err
	}

	for i, diff := range diffs {
		tmpPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/patch.%d.txt", imageHash, ts, i))
		if err != nil {
			log.Printf("[miner] error creating tmp patch file %s", err)
			return err
		}
		log.Printf("[miner] tmp patch file %s", tmpPatchFile.Name())
		defer os.Remove(tmpPatchFile.Name()) // clean up

		if _, err := tmpPatchFile.Write([]byte(diff.Props().Data)); err != nil {
			log.Printf("[miner] error writing diff data to tmp patch file %s", err)
			return err
		}
		if err := tmpPatchFile.Close(); err != nil {
			log.Printf("[miner] error closing tmp patch file %s", err)
			return err
		}

		if err := diffing.CombineDiff(outPatchFile.Name(), tmpPatchFile.Name(), outPatchFile.Name()); err != nil {
			return err
		}
	}

	// now apply the combined patch file to the state
	if err := diffing.Patch(outPatchFile.Name(), false, true); err != nil {
		log.Printf("[miner] error diffing patch file %s", err)
		return err
	}
	state, err := ioutil.ReadFile(tmpStateFile.Name())
	if err != nil {
		log.Printf("[miner] error reading tmp state file %s", err)
		return err
	}

	log.Printf("[miner] state for image hash %s\nstate: %s", imageHash, string(state))
	headStateFileName := tmpStateFile.Name()
	runningBlockNumber, err := hexutil.DecodeUint64(prevStateBlock.Props().BlockNumber)
	if err != nil {
		return err
	}
	runningBlockHash := *prevStateBlock.Props().BlockHash // note: already checked nil pointer, above

	// apply state to container and start running transactions
	for i, tx := range transactions {
		if tx == nil {
			log.Printf("[miner] tx is nil for image hash %s", imageHash)
			return errors.New("nil tx")
		}
		if tx.Props().TxHash == nil {
			log.Printf("[miner] tx hash is nil for %v", tx.Props())
			return errors.New("nil tx hash")
		}

		var nextState []byte
		log.Printf("[miner] tx method %s", tx.Props().Method)

		if tx.Props().Method == methodTypes.InvokeMethod {
			payload := tx.Props().Payload

			var parsed []string
			if err := json.Unmarshal(payload, &parsed); err != nil {
				log.Printf("[miner] error unmarshalling json for image hash %s", imageHash)
				return err
			}

			inputsJSON, err := json.Marshal(struct {
				Method string   `json:"method"`
				Params []string `json:"params"`
			}{
				Method: parsed[0],
				Params: parsed[1:],
			})
			if err != nil {
				log.Printf("[miner] error marshalling json for image hash %s", imageHash)
				return err
			}

			log.Printf("[miner] invoking method %s for image hash %s", parsed[0], imageHash)

			// run container, passing the tx inputs
			sb := sandbox.NewSandbox(&sandbox.Config{})
			nextState, err = sb.Play(&sandbox.PlayConfig{
				ImageID:      tx.Props().ImageHash,
				Payload:      inputsJSON,
				InitialState: state,
			})

			if err != nil {
				log.Printf("[miner] error running container for image hash: %s; error: %s", tx.Props().ImageHash, err)
				return err
			}

			log.Printf("[miner] container new state: %s", string(nextState))
		}

		nextStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.%d.txt", imageHash, ts, i))
		if err != nil {
			return err
		}
		defer os.Remove(nextStateFile.Name()) // clean up

		if _, err := nextStateFile.Write(nextState); err != nil {
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

		nextStateHash := hexutil.EncodeBytes(nextState)
		runningBlockNumber++
		nextStateStruct := statechain.New(&statechain.BlockProps{
			BlockNumber:       hexutil.EncodeUint64(runningBlockNumber),
			BlockTime:         hexutil.EncodeUint64(uint64(ts)),
			ImageHash:         imageHash,
			TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
			PrevBlockHash:     runningBlockHash,
			StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
			StateCurrentHash:  string(nextStateHash),
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

// TODO: improve
func (s Service) buildGenesisStateBlock(imageHash string, tx *statechain.Transaction) error {
	log.Printf("[miner] building genesis state block for image hash %s", imageHash)

	ts := time.Now().Unix()
	var (
		newDiffs            []*statechain.Diff
		newTxs              []*statechain.Transaction
		newStatechainBlocks []*statechain.Block
	)

	if tx.Props().TxHash == nil {
		log.Printf("[miner] tx hash is nil for %v", tx.Props())
		return errors.New("nil tx hash")
	}

	if tx.Props().TxHash == nil {
		log.Printf("[miner] tx hash is nil for %v", tx.Props())
		return errors.New("nil tx hash")
	}

	log.Printf("[miner] tx method %s", tx.Props().Method)

	// initial state
	nextState := tx.Props().Payload
	log.Printf("[miner] container initial state: %s", string(nextState))

	nextStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		return err
	}
	defer os.Remove(nextStateFile.Name()) // clean up

	tmpStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/state.txt", imageHash, ts))
	if err != nil {
		return err
	}
	headStateFileName := tmpStateFile.Name()

	if _, err := nextStateFile.Write(nextState); err != nil {
		return err
	}
	if err := nextStateFile.Close(); err != nil {
		return err
	}

	outPatchFile, err := makeTempFile(fmt.Sprintf("%s/%v/combined.txt", imageHash, ts))
	if err != nil {
		return err
	}
	defer os.Remove(outPatchFile.Name()) // clean up
	if err := outPatchFile.Close(); err != nil {
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
		return err
	}

	newDiffs = append(newDiffs, diffStruct)
	newTxs = append(newTxs, tx)
	newStatechainBlocks = append(newStatechainBlocks, nextStateStruct)

	// write to the mined block
	s.minedBlock.mut.Lock()
	// note: they should all have same length
	for i := 0; i < len(newDiffs); i++ {
		s.minedBlock.DiffsMap[*newDiffs[i].Props().DiffHash] = newDiffs[i]
		s.minedBlock.TransactionsMap[*newTxs[i].Props().TxHash] = newTxs[i]
		s.minedBlock.StatechainBlocksMap[*newStatechainBlocks[i].Props().BlockHash] = newStatechainBlocks[i]
	}
	s.minedBlock.mut.Unlock()
	log.Printf("[miner] mined state block for image hash %s", imageHash)

	return nil
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

func init() {
	log.AddHook(logger.ContextHook{})
}
