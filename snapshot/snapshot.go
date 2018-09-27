package snapshot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/c3systems/c3-go/common/fileutil"
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/common/stringutil"
	"github.com/c3systems/c3-go/core/docker"
	"github.com/c3systems/c3-go/core/miner"
	"github.com/c3systems/c3-go/core/p2p"
	"github.com/c3systems/c3-go/core/sandbox"
	"github.com/c3systems/c3-go/node/store"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

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

// Snapshot ...
func (s *Service) Snapshot(imageHash string, stateBlockNumber int) (string, error) {
	// stateFileName ...
	var stateFileName = "state.txt"
	headBlock, err := s.Mempool.GetHeadBlock()
	if err != nil {
		return "", err
	}
	prevStateBlock, err := s.P2P.FetchMostRecentStateBlock(imageHash, &headBlock)
	if err != nil {
		return "", err
	}

	if prevStateBlock == nil {
		return "", errors.New("state block not found")
	}

	var fileNames []string
	defer fileutil.RemoveFiles(&fileNames)

	ts := time.Now().Unix()

	minr, err := miner.New(&miner.Props{
		Context:             context.Background(),
		PreviousBlock:       &headBlock,
		P2P:                 s.P2P,
		PendingTransactions: nil,
	})
	if err != nil {
		return "", err
	}

	diffs, err := minr.GatherDiffs(prevStateBlock)
	if err != nil {
		return "", err
	}

	// debug logs
	for i := range diffs {
		log.Printf("diff %v\n%s", i, diffs[i].Props().Data)
	}

	var genesisState []byte
	state, err := miner.GenerateStateFromDiffs(context.Background(), imageHash, genesisState, diffs)
	if err != nil {
		return "", err
	}

	sta, err := stringutil.CompactJSON(state)
	if err != nil {
		return "", err
	}

	var st map[string]string
	err = json.Unmarshal(sta, &st)
	if err != nil {
		return "''", err
	}

	// debug logs
	for key, value := range st {
		k, _ := hexutil.DecodeString(key)
		v, _ := hexutil.DecodeString(value)
		log.Printf("[sandbox] state k/v %s=>%s", string(k), string(v))
	}

	stateFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/%s", imageHash, ts, stateFileName))
	if err != nil {
		return "", err
	}
	fileNames = append(fileNames, stateFile.Name())
	if _, err = stateFile.Write(state); err != nil {
		return "", err
	}
	if err = stateFile.Close(); err != nil {
		return "", err
	}

	nextStateFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/nextState.txt", imageHash, ts))
	if err != nil {
		return "", err
	}
	fileNames = append(fileNames, nextStateFile.Name()) // clean up
	if err = nextStateFile.Close(); err != nil {
		return "", err
	}

	patchFile, err := fileutil.CreateTempFile(fmt.Sprintf("%s/%v/diff.patch", imageHash, ts))
	if err != nil {
		return "", err
	}
	fileNames = append(fileNames, patchFile.Name())
	if err = patchFile.Close(); err != nil {
		return "", err
	}

	spew.Dump(prevStateBlock)

	latestState := state

	log.Printf("[miner] setting docker container initial state to %q", string(state))

	committedImageID, err := s.Sandbox.CommitPlay(&sandbox.PlayConfig{
		ImageID:      imageHash,
		Payload:      []byte(""),
		InitialState: latestState,
	})
	if err != nil {
		return "", err
	}

	if committedImageID == "" {
		return "", errors.New("expected image ID")
	}

	return docker.ShortImageID(committedImageID), nil
}
