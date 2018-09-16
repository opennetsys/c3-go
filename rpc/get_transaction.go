package rpc

// Ping ...
import (
	pb "github.com/c3systems/c3-go/rpc/pb"
	log "github.com/sirupsen/logrus"
)

// getTransaction ...
func (s *RPC) getTransaction(params []string) *pb.TransactionResponse {
	headBlock, err := s.mempool.GetHeadBlock()
	if err != nil {
		log.Fatal(err)
	}

	_ = headBlock

	// TODO

	/*
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
		}
	*/

	return &pb.TransactionResponse{
	/*
		TxHash:
		ImageHash:
		Method:
		Payload:
		From:
		Sig:
	*/
	}
}
