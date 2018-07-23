package chain

import (
	"errors"

	mainchain "github.com/c3systems/c3-go/core/chain/mainchain"
	statechain "github.com/c3systems/c3-go/core/chain/statechain"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	cid "github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
)

// New ...
func New(props *Props) (*Service, error) {
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}

	return &Service{
		props: *props,
	}, nil
}

// Props ...
func (s Service) Props() Props {
	return s.props
}

// TODO: implement methods

// AddMainBlock ...
func (s Service) AddMainBlock(block *mainchain.Block) *cid.Cid {
	return nil
}

// PendingTransactions ...
func (s Service) PendingTransactions() []*statechain.Transaction {
	return nil
}

// MainHead ...
func (s Service) MainHead() (*mainchain.Block, error) {
	return nil, nil
}

// StateHead ...
func (s Service) StateHead(hash string) (*statechain.Block, error) {
	return nil, nil
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
