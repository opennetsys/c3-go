package rpc

// Ping ...
import pb "github.com/c3systems/c3-go/rpc/pb"

// ping ...
func (s *RPC) ping() *pb.PingResponse {
	return &pb.PingResponse{
		Data: "pong",
	}
}
