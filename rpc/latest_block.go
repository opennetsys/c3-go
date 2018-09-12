package rpc

// Ping ...
import (
	pb "github.com/c3systems/c3-go/rpc/pb"
	log "github.com/sirupsen/logrus"
)

// latestBlock ...
func (s *RPC) latestBlock() *pb.LatestBlockResponse {
	headBlock, err := s.mempool.GetHeadBlock()
	if err != nil {
		log.Fatal(err)
	}

	return &pb.LatestBlockResponse{
		Data: headBlock.Props().BlockNumber,
	}
}
