package snapshot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/c3systems/c3-go/common/hashutil"
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/core/diffing"
	"github.com/c3systems/c3-go/core/p2p"
	"github.com/c3systems/c3-go/core/sandbox"
	"github.com/c3systems/c3-go/node/store"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// TODO

// WORK IN PROGRESS

// Service ...
type Service struct {
	Mempool store.Interface
	P2P     p2p.Interface
	Sandbox sandbox.Interface
}

// Config ...
type Config struct {
	Mempool store.Interface
	P2P     p2p.Interface
}

// New ...
func New(cfg *Config) *Service {
	sb := sandbox.New(nil)

	return &Service{
		P2P:     cfg.P2P,
		Mempool: cfg.Mempool,
		Sandbox: sb,
	}
}

const (
	// StateFileName ...
	StateFileName string = "state.txt"
)

// Snapshot ...
func (s *Service) Snapshot(imageHash string, stateBlockNumber int) error {
	headBlock, err := s.Mempool.GetHeadBlock()
	if err != nil {
		return err
	}
	prevStateBlock, err := s.P2P.FetchMostRecentStateBlock(imageHash, &headBlock)
	if err != nil {
		return err
	}

	if prevStateBlock == nil {
		return errors.New("state block not found")
	}

	var (
		newDiffs            []*statechain.Diff
		newStatechainBlocks []*statechain.Block
		fileNames           []string
	)

	defer cleanupFiles(&fileNames)

	ts := time.Now().Unix()

	state := []byte(``)

	stateFile, err := makeTempFile(fmt.Sprintf("%s/%v/%s", imageHash, ts, StateFileName))
	if err != nil {
		return err
	}
	fileNames = append(fileNames, stateFile.Name())
	if _, err = stateFile.Write(state); err != nil {
		return err
	}
	if err = stateFile.Close(); err != nil {
		return err
	}

	nextStateFile, err := makeTempFile(fmt.Sprintf("%s/%v/nextState.txt", imageHash, ts))
	if err != nil {
		return err
	}
	fileNames = append(fileNames, nextStateFile.Name()) // clean up
	if err = nextStateFile.Close(); err != nil {
		return err
	}

	patchFile, err := makeTempFile(fmt.Sprintf("%s/%v/diff.patch", imageHash, ts))
	if err != nil {
		return err
	}
	fileNames = append(fileNames, patchFile.Name())
	if err = patchFile.Close(); err != nil {
		return err
	}

	spew.Dump(prevStateBlock)

	runningBlockNumber, err := hexutil.DecodeUint64(prevStateBlock.Props().BlockNumber)
	if err != nil {
		return err
	}
	runningBlockHash := *prevStateBlock.Props().BlockHash // note: already checked nil pointer, above
	latestState := state

	//tx := statechain.NewTransaction(&statechain.TransactionProps{ })

	// apply state to container and start running transactions

	var nextState []byte

	log.Printf("[miner] setting docker container initial state to %q", string(state))

	// run container, passing the tx inputs
	nextState, err = s.Sandbox.Play(&sandbox.PlayConfig{
		ImageID:      imageHash,
		Payload:      []byte(""),
		InitialState: latestState,
	})

	if err != nil {
		log.Errorf("[miner] error running container for image hash: %s; error: %s", imageHash, err)
		return err
	}

	log.Printf("[miner] container new state: %s", string(nextState))

	if err := ioutil.WriteFile(nextStateFile.Name(), nextState, os.ModePerm); err != nil {
		return err
	}

	if err = diffing.Diff(stateFile.Name(), nextStateFile.Name(), patchFile.Name(), false); err != nil {
		return err
	}

	// build the diff struct
	diffData, err := ioutil.ReadFile(patchFile.Name())
	if err != nil {
		return err
	}

	diffStruct := statechain.NewDiff(&statechain.DiffProps{
		Data: string(diffData),
	})
	if err := diffStruct.SetHash(); err != nil {
		return err
	}

	nextStateHashBytes := hashutil.Hash(nextState)
	nextStateHash := hexutil.EncodeToString(nextStateHashBytes[:])
	log.Printf("[miner] state prev diff hash: %s", *diffStruct.Props().DiffHash)
	log.Printf("[miner] state current hash: %s", nextStateHash)

	runningBlockNumber++
	nextStateStruct := statechain.New(&statechain.BlockProps{
		BlockNumber: hexutil.EncodeUint64(runningBlockNumber),
		BlockTime:   hexutil.EncodeUint64(uint64(ts)),
		ImageHash:   imageHash,
		//TxHash:            *tx.Props().TxHash, // note: checked for nil pointer, above
		PrevBlockHash:     runningBlockHash,
		StatePrevDiffHash: *diffStruct.Props().DiffHash, // note: used setHash, above so it would've erred
		StateCurrentHash:  nextStateHash,
	})
	if err := nextStateStruct.SetHash(); err != nil {
		return err
	}
	runningBlockHash = *nextStateStruct.Props().BlockHash

	newDiffs = append(newDiffs, diffStruct)
	newStatechainBlocks = append(newStatechainBlocks, nextStateStruct)

	// get ready for the next loop
	latestState = nextState

	if err := ioutil.WriteFile(stateFile.Name(), nextState, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func cleanupFiles(fileNames *[]string) {
	if fileNames == nil {
		return
	}

	for idx := range *fileNames {
		if err := os.Remove((*fileNames)[idx]); err != nil {
			log.Errorf("[miner] err cleaning up file %s", (*fileNames)[idx])
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
