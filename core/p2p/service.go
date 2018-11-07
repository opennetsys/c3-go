package p2p

import (
	"context"
	"errors"

	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/merkle"
	"github.com/c3systems/c3-go/core/chain/statechain"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	bitswap "github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-bitswap/network"
	bserv "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
)

// Ensure the struct implements the interface
var _ Interface = (*Service)(nil)

var latestMainchainBlockKey = []byte("latestMainchainBlock")

// New ...
func New(props *Props) (*Service, error) {
	var err error

	once.Do(func() {
		if props == nil {
			err = errors.New("props cannot be nil")
			return
		}
		if props.Host == nil || props.BlockStore == nil || props.Router == nil {
			err = errors.New("host, blockstore and router are required")
			return
		}

		bsnet := network.NewFromIpfsHost(props.Host, props.Router)
		bswap := bitswap.New(context.Background(), bsnet, props.BlockStore)

		// Bitswap only fetches blocks from other nodes, to fetch blocks from
		// either the local cache, or a remote node, we can wrap it in a
		// 'blockservice'
		bservice := bserv.New(props.BlockStore, bswap)

		service = &Service{
			props:        *props,
			peersOrLocal: bservice,
			local:        props.BlockStore,
		}
	})

	return service, err
}

// Props ...
func (s Service) Props() Props {
	return s.props
}

// Set ...
func (s Service) Set(v interface{}) (*cid.Cid, error) {
	return Put(s.peersOrLocal, v)
}

// SetMainchainBlock ...
// note: this function does not do any validation!
func (s Service) SetMainchainBlock(block *mainchain.Block) (*cid.Cid, error) {
	return PutMainchainBlock(s.peersOrLocal, block)
}

// SetStatechainBlock ...
func (s Service) SetStatechainBlock(block *statechain.Block) (*cid.Cid, error) {
	return PutStatechainBlock(s.peersOrLocal, block)
}

// SetStatechainTransaction ...
func (s Service) SetStatechainTransaction(tx *statechain.Transaction) (*cid.Cid, error) {
	return PutStatechainTransaction(s.peersOrLocal, tx)
}

// SetStatechainDiff ...
func (s Service) SetStatechainDiff(d *statechain.Diff) (*cid.Cid, error) {
	return PutStatechainDiff(s.peersOrLocal, d)
}

// SetMerkleTree ..
func (s Service) SetMerkleTree(tree *merkle.Tree) (*cid.Cid, error) {
	return PutMerkleTree(s.peersOrLocal, tree)
}

// SetBytes ...
func (s Service) SetBytes(b []byte) (*cid.Cid, error) {
	return PutBytes(s.peersOrLocal, b)
}

// SetLatestBlock ...
func (s Service) SetLatestBlock(block *mainchain.Block) (*cid.Cid, error) {
	return PutLatestBlock(s.peersOrLocal, block)
}

//// SaveLocal ...
//func (s Service) SaveLocal(v interface{}) (*cid.Cid, error) {
//return Put(s.local, v)
//}

//// SaveLocalMainchainBlock ...
//// note: this function does not do any validation!
//func (s Service) SaveLocalMainchainBlock(block *mainchain.Block) (*cid.Cid, error) {
//return PutMainchainBlock(s.local, block)
//}

//// SaveLocalStatechainBlock ...
//func (s Service) SaveLocalStatechainBlock(block *statechain.Block) (*cid.Cid, error) {
//return PutStatechainBlock(s.local, block)
//}

//// SaveLocalStatechainTransaction ...
//func (s Service) SaveLocalStatechainTransaction(tx *statechain.Transaction) (*cid.Cid, error) {
//return PutStatechainTransaction(s.local, tx)
//}

//// SaveLocalStatechainDiff ...
//func (s Service) SaveLocalStatechainDiff(d *statechain.Diff) (*cid.Cid, error) {
//return PutStatechainDiff(s.local, d)
//}

// note: cannot do generic get bc need to know the type to deserialize into
// Get ...
//func (s Service) Get(c *cid.Cid) (interface{}, error) {
//return Fetch(s.peersOrLocal, c)
//}

// GetMainchainBlock ...
func (s Service) GetMainchainBlock(c *cid.Cid) (*mainchain.Block, error) {
	return FetchMainchainBlock(s.peersOrLocal, c)
}

// GetStatechainBlock ...
func (s Service) GetStatechainBlock(c *cid.Cid) (*statechain.Block, error) {
	return FetchStateChainBlock(s.peersOrLocal, c)
}

// GetStatechainTransaction ...
func (s Service) GetStatechainTransaction(c *cid.Cid) (*statechain.Transaction, error) {
	return FetchStateChainTransaction(s.peersOrLocal, c)
}

// GetStatechainDiff ...
func (s Service) GetStatechainDiff(c *cid.Cid) (*statechain.Diff, error) {
	return FetchStateChainDiff(s.peersOrLocal, c)
}

// GetMerkleTree ...
func (s Service) GetMerkleTree(c *cid.Cid) (*merkle.Tree, error) {
	return FetchMerkleTree(s.peersOrLocal, c)
}

// GetBytes ...
func (s Service) GetBytes(c *cid.Cid) ([]byte, error) {
	return FetchBytes(s.peersOrLocal, c)
}

// GetLatestBlock ...
func (s Service) GetLatestBlock() (*mainchain.Block, error) {
	c, err := GetBytesCID(latestMainchainBlockKey)
	if err != nil {
		return nil, err
	}
	return FetchLatestBlock(s.peersOrLocal, c)
}

// FetchMostRecentStateBlock ...
func (s Service) FetchMostRecentStateBlock(imageHash string, block *mainchain.Block) (*statechain.Block, error) {
	if block == nil {
		log.Errorf("[p2p] block is nil for image hash %s", imageHash)
		return nil, errors.New("block is nil")
	}

	if block.Props().BlockHash == nil {
		log.Errorf("[p2p] block hash is nil for image hash %s", imageHash)
		return nil, errors.New("block hash is nil")
	}

	log.Printf("[p2p] state block merkle hash is %s for image hash %s", block.Props().StateBlocksMerkleHash, imageHash)
	log.Printf("[p2p] fetch most recent state block starting with main chain block number %s for image hash %s", block.Props().BlockNumber, imageHash)

	merkleHash := block.Props().StateBlocksMerkleHash

	log.Printf("[p2p] block merkle hash for block %s is %q for image hash %s", *block.Props().BlockHash, merkleHash, imageHash)

	if merkleHash != hexutil.EncodeString("") {
		log.Printf("[p2p] merkle hash is empty; %s", block.Props().StateBlocksMerkleHash)

		// 1. search the current block
		treeCID, err := GetCIDByHash(block.Props().StateBlocksMerkleHash)
		if err != nil {
			log.Errorf("[p2p] error getting cid by hash for state block merkle hash %s for image hash %s", block.Props().StateBlocksMerkleHash, imageHash)
			return nil, err
		}

		log.Printf("[p2p] cid by hash for state block merkle hash is %s for image hash %s", treeCID, imageHash)

		tree, err := s.GetMerkleTree(treeCID)
		if err != nil {
			log.Errorf("[p2p] error getting merkle tree for tree cid %s for image hash %s", treeCID, imageHash)
			return nil, err
		}

		log.Printf("[p2p] tree hashes for for image hash %s; %v", imageHash, len(tree.Props().Hashes))

		// TODO: check kind
		// TODO: use go routines
		for _, stateBlockHash := range tree.Props().Hashes {
			stateBlockCID, err := GetCIDByHash(stateBlockHash)
			if err != nil {
				log.Errorf("[p2p] error getting cid by hash for state block hash %s for image hash %s", stateBlockHash, imageHash)
				return nil, err
			}

			stateBlock, err := s.GetStatechainBlock(stateBlockCID)
			if err != nil {
				log.Errorf("[p2p] error getting state chain block for state block cid %s for image hash %s", stateBlockCID, imageHash)
				return nil, err
			}

			if stateBlock.Props().ImageHash == imageHash {
				log.Printf("[p2p] state block image hash matches image hash %s", imageHash)
				return stateBlock, nil
			}
		}
	}

	log.Printf("[p2p] mainchain genesis block block number is %v", mainchain.GenesisBlock.Props().BlockNumber)

	// walk the mainchain
	head := block
	for head.Props().BlockNumber != mainchain.GenesisBlock.Props().BlockNumber {
		log.Printf("[p2p] checking head block number %v for image hash %s", head.Props().BlockNumber, imageHash)
		prevCID, err := GetCIDByHash(head.Props().PrevBlockHash)
		if err != nil {
			log.Errorf("[p2p] error getting cid by hash for prev block hash %s for image hash %s", head.Props().PrevBlockHash, imageHash)
			return nil, err
		}

		prevBlock, err := s.GetMainchainBlock(prevCID)
		if err != nil {
			log.Errorf("[p2p] error getting main chain block for prev cid %s for image hash %s", prevCID, imageHash)
			return nil, err
		}
		head = prevBlock

		merkleHash = head.Props().StateBlocksMerkleHash

		log.Printf("[p2p] block merkle hash for block %s is %q for image hash %s", *block.Props().BlockHash, merkleHash, imageHash)

		if merkleHash == hexutil.EncodeString("") {
			log.Warnf("[p2p] merkle hash for block %s is empty for image hash %s, continuing", *prevBlock.Props().BlockHash, imageHash)
			continue
		}

		treeCID, err := GetCIDByHash(prevBlock.Props().StateBlocksMerkleHash)
		if err != nil {
			log.Errorf("[p2p] error getting cid by hash for prev block state blocks merkle hash %s for image hash %s", prevBlock.Props().StateBlocksMerkleHash, imageHash)
			return nil, err
		}

		tree, err := s.GetMerkleTree(treeCID)
		if err != nil {
			log.Errorf("[p2p] error getting merkle tree by tree cid %s for image hash %s", treeCID, imageHash)
			return nil, err
		}

		// TODO: check kind
		// TODO: use go routines
		log.Printf("[p2p] tree hash count %v for image hash %v", len(tree.Props().Hashes), imageHash)
		for _, stateBlockHash := range tree.Props().Hashes {
			log.Printf("[p2p] getting cid for state chain block by state block hash %s", stateBlockHash)
			stateBlockCID, err := GetCIDByHash(stateBlockHash)
			if err != nil {
				log.Printf("[p2p] error getting cid by hash for state block hash %s for image hash %s", stateBlockHash, imageHash)
				return nil, err
			}

			log.Printf("[p2p] getting state chain block by state block hash %s and state block cid %s", stateBlockHash, stateBlockCID)
			stateBlock, err := s.GetStatechainBlock(stateBlockCID)
			if err != nil {
				log.Errorf("[p2p] error getting state chain block for state block cid %s for image hash %s", stateBlockCID, imageHash)
				return nil, err
			}

			if stateBlock.Props().ImageHash == imageHash {
				return stateBlock, nil
			}
		}
	}

	return nil, nil
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
