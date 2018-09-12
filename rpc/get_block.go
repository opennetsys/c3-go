package rpc

// Ping ...
import (
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/p2p"
	pb "github.com/c3systems/c3-go/rpc/pb"
	log "github.com/sirupsen/logrus"
)

// getBlock ...
func (s *RPC) getBlock(params []string) *pb.BlockResponse {
	headBlock, err := s.mempool.GetHeadBlock()
	if err != nil {
		log.Fatal(err)
	}

	wantBlockNumber, err := hexutil.DecodeFloat64(params[0])
	if err != nil {
		log.Fatal(err)
	}

	headBlockNumber, err := hexutil.DecodeFloat64(headBlock.Props().BlockNumber)
	if wantBlockNumber > headBlockNumber || wantBlockNumber <= 0 {
		return nil
	}

	currentBlock := headBlock
	for {
		blockNumber, err := hexutil.DecodeFloat64(currentBlock.Props().BlockNumber)
		if err != nil {
			log.Fatal(err)
		}
		if blockNumber == 0 {
			return nil
		}
		if blockNumber != wantBlockNumber {
			prevBlockHash := headBlock.Props().PrevBlockHash
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
		}
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
