package miner

import (
	"errors"

	"github.com/c3systems/c3-go/common/coder"
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/c3-go/core/chain/mainchain"
	"github.com/c3systems/c3-go/core/chain/merkle"
	"github.com/c3systems/c3-go/core/chain/statechain"
)

// Serialize ...
func (m *MinedBlock) Serialize() ([]byte, error) {
	tmp, err := BuildCoderFromBlock(m)
	if err != nil {
		return nil, err
	}

	bytes, err := tmp.Marshal()
	if err != nil {
		return nil, err
	}

	return coder.AppendCode(bytes), nil
}

// Deserialize ...
func (m *MinedBlock) Deserialize(data []byte) error {
	if data == nil {
		return errors.New("nil bytes")
	}
	if m == nil {
		return errors.New("nil block")
	}

	_, bytes, err := coder.StripCode(data)
	if err != nil {
		return err
	}

	b, err := BuildBlockFromBytes(bytes)
	if err != nil {
		return err
	}
	if b != nil {
		*m = *b // note: ignore the sync copy linting error
	}

	return nil
}

// SerializeString ...
func (m *MinedBlock) SerializeString() (string, error) {
	b, err := m.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(b), nil
}

// DeserializeString ...
func (m *MinedBlock) DeserializeString(hexStr string) error {
	if m == nil {
		return ErrNilBlock
	}

	b, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return m.Deserialize(b)
}

// BuildCoderFromBlock ...
func BuildCoderFromBlock(m *MinedBlock) (*coder.MinedBlock, error) {
	tmp := &coder.MinedBlock{}
	if m.NextBlock != nil {
		tmp.NextBlock = mainchain.BuildCoderFromBlock(m.NextBlock)
	}
	if m.PreviousBlock != nil {
		tmp.PreviousBlock = mainchain.BuildCoderFromBlock(m.PreviousBlock)
	}
	if m.StatechainBlocksMap != nil && len(m.StatechainBlocksMap) > 0 {
		for k, v := range m.StatechainBlocksMap {
			if tmp.StatechainBlocksMap == nil {
				tmp.StatechainBlocksMap = make(map[string]*coder.StatechainBlock)
			}

			tmp.StatechainBlocksMap[k] = statechain.BuildCoderFromBlock(v)
		}
	}
	if m.TransactionsMap != nil && len(m.TransactionsMap) > 0 {
		for k, v := range m.TransactionsMap {
			if tmp.TransactionsMap == nil {
				tmp.TransactionsMap = make(map[string]*coder.Transaction)
			}

			c, err := statechain.BuildCoderFromTransaction(v)
			if err != nil {
				return nil, err
			}

			tmp.TransactionsMap[k] = c
		}
	}
	if m.DiffsMap != nil && len(m.DiffsMap) > 0 {
		for k, v := range m.DiffsMap {
			if tmp.DiffsMap == nil {
				tmp.DiffsMap = make(map[string]*coder.Diff)
			}

			tmp.DiffsMap[k] = statechain.BuildCoderFromDiff(v)
		}
	}
	if m.MerkleTreesMap != nil && len(m.MerkleTreesMap) > 0 {
		for k, v := range m.MerkleTreesMap {
			if tmp.MerkleTreesMap == nil {
				tmp.MerkleTreesMap = make(map[string]*coder.MerkleTree)
			}

			tmp.MerkleTreesMap[k] = merkle.BuildCoderFromTree(v)
		}
	}

	return tmp, nil
}

// BuildBlockFromBytes ...
func BuildBlockFromBytes(data []byte) (*MinedBlock, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	c, err := BuildCoderFromBytes(data)
	if err != nil {
		return nil, err
	}

	return BuildBlockFromCoder(c)
}

// BuildCoderFromBytes ...
func BuildCoderFromBytes(data []byte) (*coder.MinedBlock, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	tmp := new(coder.MinedBlock)
	if err := tmp.Unmarshal(data); err != nil {
		return nil, err
	}
	if tmp == nil {
		return nil, errors.New("nil output")
	}

	return tmp, nil
}

// BuildBlockFromCoder ...
func BuildBlockFromCoder(tmp *coder.MinedBlock) (*MinedBlock, error) {
	if tmp == nil {
		return nil, errors.New("nil coder")
	}

	block := new(MinedBlock)
	if tmp.NextBlock != nil {
		blockProps, err := mainchain.BuildBlockPropsFromCoder(tmp.NextBlock)
		if err != nil {
			return nil, err
		}

		block.NextBlock = mainchain.New(blockProps)
	}

	if tmp.PreviousBlock != nil {
		blockProps, err := mainchain.BuildBlockPropsFromCoder(tmp.PreviousBlock)
		if err != nil {
			return nil, err
		}

		block.PreviousBlock = mainchain.New(blockProps)
	}

	if tmp.StatechainBlocksMap != nil && len(tmp.StatechainBlocksMap) > 0 {
		for k, v := range tmp.StatechainBlocksMap {
			if block.StatechainBlocksMap == nil {
				block.StatechainBlocksMap = make(map[string]*statechain.Block)
			}

			props, err := statechain.BuildBlockPropsFromCoder(v)
			if err != nil {
				return nil, err
			}

			block.StatechainBlocksMap[k] = statechain.New(props)
		}
	}

	if tmp.TransactionsMap != nil && len(tmp.TransactionsMap) > 0 {
		for k, v := range tmp.TransactionsMap {
			if block.TransactionsMap == nil {
				block.TransactionsMap = make(map[string]*statechain.Transaction)
			}

			props, err := statechain.BuildTransactionPropsFromCoder(v)
			if err != nil {
				return nil, err
			}

			tx := statechain.NewTransaction(props)
			block.TransactionsMap[k] = tx
		}
	}

	if tmp.DiffsMap != nil && len(tmp.DiffsMap) > 0 {
		for k, v := range tmp.DiffsMap {
			if block.DiffsMap == nil {
				block.DiffsMap = make(map[string]*statechain.Diff)
			}

			props, err := statechain.BuildDiffPropsFromCoder(v)
			if err != nil {
				return nil, err
			}

			d := statechain.NewDiff(props)
			block.DiffsMap[k] = d
		}
	}

	if tmp.MerkleTreesMap != nil && len(tmp.MerkleTreesMap) > 0 {
		for k, v := range tmp.MerkleTreesMap {
			if block.MerkleTreesMap == nil {
				block.MerkleTreesMap = make(map[string]*merkle.Tree)
			}

			props, err := merkle.BuildTreePropsFromCoder(v)
			if err != nil {
				return nil, err
			}

			t, err := merkle.New(props)
			if err != nil {
				return nil, err
			}
			block.MerkleTreesMap[k] = t
		}
	}

	return block, nil
}
