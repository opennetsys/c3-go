package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/c3systems/c3-go/core/chain/statechain"
	pb "github.com/c3systems/c3-go/rpc/pb"
	"github.com/davecgh/go-spew/spew"
)

// invokeMethod ...
func (s *RPC) invokeMethod(params []string) (*pb.InvokeMethodResponse, error) {
	txstr := params[0]
	fmt.Printf("[rpc] invokeMethod received payload: %s\n", txstr)
	// TODO: pass root as json instead of json as first array value
	tx := new(statechain.Transaction)
	err := json.Unmarshal([]byte(txstr), tx)
	if err != nil {
		return nil, err
	}

	spew.Dump(tx)

	go s.node.HandleReceiptOfStatechainTransaction(tx)

	resp, err := s.node.BroadcastTransaction(tx)
	if err != nil {
		return nil, err
	}

	txHash := *resp.TxHash
	return &pb.InvokeMethodResponse{
		TxHash: txHash,
	}, nil
}
