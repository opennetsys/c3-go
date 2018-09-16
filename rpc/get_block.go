package rpc

// Ping ...
import (
	"errors"

	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/p2p"
	pb "github.com/c3systems/c3-go/rpc/pb"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrBlockNotFound ...
	ErrBlockNotFound = errors.New("block not found")
)

// getBlock ...
func (s *RPC) getBlock(params []string) (*pb.BlockResponse, error) {
	headBlock, err := s.mempool.GetHeadBlock()
	if err != nil {
		return nil, err
	}

	wantBlockNumber, err := hexutil.DecodeInt(params[0])
	if err != nil {
		return nil, err
	}

	headBlockNumber, err := hexutil.DecodeInt(headBlock.Props().BlockNumber)

	if wantBlockNumber <= 0 {
		return nil, ErrBlockNotFound
	}

	if wantBlockNumber > headBlockNumber {
		return nil, ErrBlockNotFound
	}

	currentBlock := headBlock
	for {
		blockNumber, err := hexutil.DecodeInt(currentBlock.Props().BlockNumber)
		if err != nil {
			return nil, err
		}

		if blockNumber <= 0 || blockNumber > headBlockNumber {
			return nil, ErrBlockNotFound
		}
		if blockNumber != wantBlockNumber {
			prevBlockHash := currentBlock.Props().PrevBlockHash
			prevBlock := s.getBlockByHash(prevBlockHash)
			currentBlock = *prevBlock
			continue
		}

		props := currentBlock.Props()
		//sig := props.MinerSig
		sig := &pb.Signature{}
		blockHash := props.BlockHash

		return &pb.BlockResponse{
			BlockHash:             *blockHash,
			BlockNumber:           props.BlockNumber,
			BlockTime:             props.BlockTime,
			ImageHash:             props.ImageHash,
			StateBlocksMerkleHash: props.StateBlocksMerkleHash,
			PrevBlockHash:         props.PrevBlockHash,
			Nonce:                 props.Nonce,
			Difficulty:            props.Difficulty,
			MinerAddress:          props.MinerAddress,
			MinerSig:              sig,
		}, nil
	}
}

func (s *RPC) getBlockByHash(hash string) *mainchain.Block {
	cid, err := p2p.GetCIDByHash(hash)
	if err != nil {
		log.Fatal(err)
	}

	block, err := s.p2p.GetMainchainBlock(cid)
	if err != nil {
		log.Fatal(err)
	}

	return block
}
