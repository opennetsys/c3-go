package rpc

import (
	pb "github.com/c3systems/c3-go/rpc/pb"
)

// invokeMethod ...
func (s *RPC) invokeMethod(params []string) *pb.InvokeMethodResponse {
	// TODO
	txHash := ""

	return &pb.InvokeMethodResponse{
		TxHash: txHash,
	}
}
