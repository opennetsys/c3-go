package rpc

// Ping ...
import (
	"github.com/c3systems/c3-go/common/hexutil"
	pb "github.com/c3systems/c3-go/rpc/pb"
)

// getStateblock ...
func (s *RPC) getStateblock(params []string) (*pb.StateBlockResponse, error) {
	headBlock, err := s.mempool.GetHeadBlock()
	if err != nil {
		return nil, ErrBlockNotFound
	}

	imageHash := params[0]

	wantStateBlockNumber, err := hexutil.DecodeInt(params[1])
	if err != nil {
		return nil, err
	}

	if wantStateBlockNumber <= 0 {
		return nil, ErrStateBlockNotFound
	}

	currentStateBlock, err := s.p2p.FetchMostRecentStateBlock(imageHash, &headBlock)
	if err != nil {
		return nil, err
	}

	if currentStateBlock == nil {
		return nil, ErrStateBlockNotFound
	}

	for {
		stateBlockNumber, err := hexutil.DecodeInt(currentStateBlock.Props().BlockNumber)
		if err != nil {
			return nil, err
		}
		if stateBlockNumber <= 0 {
			return nil, ErrStateBlockNotFound
		}

		if currentStateBlock == nil {
			return nil, ErrStateBlockNotFound
		}

		props := currentStateBlock.Props()
		stateBlockNumber, err = hexutil.DecodeInt(props.BlockNumber)
		if err != nil {
			return nil, err
		}

		if stateBlockNumber != wantStateBlockNumber {
			prevStateBlockHash := props.PrevBlockHash
			prevStateBlock := s.getStateBlockByHash(prevStateBlockHash)
			currentStateBlock = prevStateBlock
			continue
		}

		blockHash := props.BlockHash

		return &pb.StateBlockResponse{
			BlockHash:         *blockHash,
			BlockNumber:       props.BlockNumber,
			BlockTime:         props.BlockTime,
			ImageHash:         props.ImageHash,
			TxHash:            props.TxHash,
			PrevBlockHash:     props.PrevBlockHash,
			StatePrevDiffHash: props.StatePrevDiffHash,
			StateCurrentHash:  props.StateCurrentHash,
		}, nil
	}
}
