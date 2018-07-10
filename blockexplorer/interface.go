package blockexplorer

import (
	"github.com/c3systems/c3/core/chain/mainchain"
	"github.com/c3systems/c3/core/chain/statechain"

	cid "github.com/c3systems/go-cid"
)

// Interface ...
type Interface interface {
	// Fetch CID's
	FetchCIDByImageHashAndBlockNumber(imageHash, blockNumber string) (*cid.Cid, error)

	// Fetch the block that includes a tx or stateblock, or tx's sent by a user
	FetchMainHashByStateBlockHash(hexHash string) (string, error)
	FetchStateBlockHashByTransactionHash(hexHash string) (string, error)
	FetchTransactionsBySenderAddress(address string, skip, limit uint64) ([]*cid.Cid, error)

	// Store
	StoreMainBlockMeta(block *mainchain.Block) error
	StoreStateBlockMeta(block *statechain.Block) error
	StoreTransactionMeta(tx *statechain.Transaction) error
	StoreMainBlockCID(hexHash string, c *cid.Cid) error
	StoreStateBlockCID(hexHash string, c *cid.Cid) error
	StoreTransactionCID(hexHash string, c *cid.Cid) error
}
