package rpc

import (
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
	"github.com/c3systems/c3-go/core/p2p"
	log "github.com/sirupsen/logrus"
)

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

func (s *RPC) getStateBlockByHash(hash string) *statechain.Block {
	cid, err := p2p.GetCIDByHash(hash)
	if err != nil {
		log.Fatal(err)
	}

	block, err := s.p2p.GetStatechainBlock(cid)
	if err != nil {
		log.Fatal(err)
	}

	return block
}
