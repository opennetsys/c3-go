package redisstore

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/statechain"
	loghooks "github.com/c3systems/c3-go/log/hooks"
	redis "github.com/gomodule/redigo/redis"
)

const (
	transactionsMembersName = "transactions"
	blocksMembersName       = "blocks"
)

// Props ...
type Props struct {
	Pool *redis.Pool
}

// Service ...
type Service struct {
	props     Props
	headBlock *mainchain.Block // note: don't use a pointer bc we don't want it being modified after being passed
}

// New ...
func New(props *Props) (*Service, error) {
	// 1. check props
	if props == nil {
		return nil, errors.New("props cannot be nil")
	}
	if props.Pool == nil {
		return nil, errors.New("pool is required")
	}

	// 2. ping db
	c := props.Pool.Get()
	defer c.Close()
	_, err := c.Do("PING")

	// 3. return service
	return &Service{
		props: *props,
	}, err
}

// Props ...
func (s Service) Props() Props {
	return s.props
}

// HasTx ...
func (s Service) HasTx(hash string) (bool, error) {
	c := s.props.Pool.Get()
	defer c.Close()

	return redis.Bool(c.Do("EXISTS", buildKey(hash)))
}

// GetTx ...
func (s Service) GetTx(hash string) (*statechain.Transaction, error) {
	c := s.props.Pool.Get()
	defer c.Close()

	bytesStr, err := redis.Strings(c.Do("GET", hash))
	if err != nil {
		return nil, err
	}
	if len(bytesStr) == 0 {
		return nil, nil
	}

	var tx statechain.Transaction
	err = tx.DeserializeString(bytesStr[0])
	return &tx, err
}

// GetTxs ...
func (s Service) GetTxs(hashes []string) ([]*statechain.Transaction, error) {
	c := s.props.Pool.Get()
	defer c.Close()

	keys := buildKeys(hashes)
	// get many keys in a single MGET, ask redigo for []string result
	bytesStrs, err := redis.Strings(c.Do("MGET", keys))
	if err != nil {
		return nil, err
	}

	var txs []*statechain.Transaction
	for _, bytesStr := range bytesStrs {
		var tx statechain.Transaction
		if err := tx.DeserializeString(bytesStr); err != nil {
			return nil, err
		}

		txs = append(txs, &tx)
	}

	return txs, nil
}

// RemoveTx ...
func (s Service) RemoveTx(hash string) error {
	c := s.props.Pool.Get()
	defer c.Close()

	key := buildKey(hash)
	_, err := c.Do("DEL", key)
	if err != nil {
		return err
	}

	_, err = c.Do("SREM", key)
	return err
}

// RemoveTxs ...
func (s Service) RemoveTxs(hashes []string) error {
	c := s.props.Pool.Get()
	defer c.Close()

	if len(hashes) == 0 {
		return nil
	}

	keys := buildKeys(hashes)
	k := make([]interface{}, len(keys))
	for i, v := range k {
		k[i] = v
	}

	_, err := c.Do("DEL", k...)
	if err != nil {
		return err
	}

	_, err = c.Do("SREM", k...)

	return err
}

// AddTx ...
func (s Service) AddTx(tx *statechain.Transaction) error {
	if tx == nil {
		return errors.New("cannot add a nil transaction")
	}

	hash, err := tx.CalculateHash()
	if err != nil {
		return err
	}

	bytesStr, err := tx.SerializeString()
	if err != nil {
		return err
	}

	c := s.props.Pool.Get()
	defer c.Close()
	_, err = c.Do("SET", buildKey(hash), bytesStr)
	if err != nil {
		return err
	}

	_, err = c.Do("SADD", transactionsMembersName, hash)

	return err
}

// GatherPendingTransactions ...
func (s Service) GatherPendingTransactions() ([]*statechain.Transaction, error) {
	log.Println("[redismempool] gathering pending transactions")
	c := s.props.Pool.Get()
	defer c.Close()

	keys, err := redis.Strings(c.Do("SMEMBERS", transactionsMembersName))
	if err != nil {
		return nil, err
	}

	bytesStrs, err := redis.Strings(c.Do("MGET", keys))
	if err != nil {
		return nil, err
	}

	txs := []*statechain.Transaction{}
	if len(bytesStrs) == 0 {
		return txs, nil
	}

	for _, bytesStr := range bytesStrs {
		var tx statechain.Transaction
		if len(bytesStr) == 0 {
			continue
		}
		if err := tx.DeserializeString(bytesStr); err != nil {
			return nil, err
		}

		txs = append(txs, &tx)
	}

	return txs, nil
}

// GetHeadBlock ...
func (s *Service) GetHeadBlock() (mainchain.Block, error) {
	if s.headBlock == nil {
		return mainchain.Block{}, errors.New("no headblock")
	}

	return *s.headBlock, nil
}

// SetHeadBlock ...
func (s *Service) SetHeadBlock(block *mainchain.Block) error {
	s.headBlock = block
	return nil
}

// SetPendingMainchainBlock ...
func (s *Service) SetPendingMainchainBlock(block *mainchain.Block) error {
	if block == nil {
		return errors.New("block is nil")
	}

	if block.Props().BlockHash == nil {
		return errors.New("block hash is nil")
	}

	bytesStr, err := block.SerializeString()
	if err != nil {
		return err
	}

	hash := *block.Props().BlockHash
	c := s.props.Pool.Get()
	defer c.Close()
	_, err = c.Do("SET", buildKey(hash), bytesStr)
	if err != nil {
		return err
	}

	_, err = c.Do("SADD", blocksMembersName, hash)

	return nil
}

// GetPendingMainchainBlocks ...
func (s *Service) GetPendingMainchainBlocks() ([]*mainchain.Block, error) {
	c := s.props.Pool.Get()
	defer c.Close()

	hashes, err := redis.Strings(c.Do("SMEMBERS", blocksMembersName))
	if err != nil {
		return nil, err
	}

	keys := buildKeys(hashes)
	// get many keys in a single MGET, ask redigo for []string result
	bytesStrs, err := redis.Strings(c.Do("MGET", keys))
	if err != nil {
		return nil, err
	}

	var blks []*mainchain.Block
	for _, bytesStr := range bytesStrs {
		if bytesStr == "" {
			continue
		}
		var blk mainchain.Block
		if err := blk.DeserializeString(bytesStr); err != nil {
			return nil, err
		}

		blks = append(blks, &blk)
	}

	return blks, nil
}

// RemovePendingMainchainBlock ...
func (s *Service) RemovePendingMainchainBlock(blockHash string) error {
	c := s.props.Pool.Get()
	defer c.Close()

	if blockHash == "" {
		return nil
	}

	key := buildKey(blockHash)
	_, err := c.Do("DEL", key)
	if err != nil {
		return err
	}

	_, err = c.Do("SREM", key)
	return err
}

// RemovePendingMainchainBlocks ...
func (s *Service) RemovePendingMainchainBlocks(blockHashes []string) error {
	c := s.props.Pool.Get()
	defer c.Close()

	if len(blockHashes) == 0 {
		return nil
	}

	keys := buildKeys(blockHashes)
	k := make([]interface{}, len(keys))
	for i, v := range k {
		k[i] = v
	}
	_, err := c.Do("DEL", k...)
	if err != nil {
		return err
	}

	_, err = c.Do("SREM", k...)

	return err
}

func init() {
	log.AddHook(loghooks.ContextHook{})
}
