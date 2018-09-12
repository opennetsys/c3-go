package rpc

// Ping ...
import pb "github.com/c3systems/c3-go/rpc/pb"

func ping() *pb.PingResponse {
	return &pb.PingResponse{
		Data: "pong",
	}
}
