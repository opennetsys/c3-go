package statechain

import (
	"encoding/json"

	"github.com/c3systems/c3/common/coder"
	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"

	"github.com/c3systems/merkletree"
)

// New ...
func New(props *BlockProps) *Block {
	if props == nil {
		return &Block{}
	}

	return &Block{
		props: *props,
	}
}

// Props ...
func (b Block) Props() BlockProps {
	return b.props
}

// Serialize ...
func (b *Block) Serialize() ([]byte, error) {
	return coder.Serialize(b.props)
}

// Deserialize ...
func (b *Block) Deserialize(data []byte) error {
	var tmpProps BlockProps
	if err := coder.Deserialize(data, &tmpProps); err != nil {
		return err
	}

	b.props = tmpProps
	return nil
}

// SerializeString ...
func (b Block) SerializeString() (string, error) {
	bytes, err := b.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

// DeserializeString ...
func (b *Block) DeserializeString(hexStr string) error {
	if b == nil {
		return ErrNilBlock
	}

	str, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return b.Deserialize([]byte(str))
}

// CalculateHash ...
func (b Block) CalculateHash() (string, error) {
	bytes, err := b.CalculateHashBytes()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

// CalculateHashBytes ...
func (b Block) CalculateHashBytes() ([]byte, error) {
	tmpBlock := Block{
		props: BlockProps{
			BlockNumber:       b.props.BlockNumber,
			BlockTime:         b.props.BlockTime,
			ImageHash:         b.props.ImageHash,
			TxHash:            b.props.TxHash,
			PrevBlockHash:     b.props.PrevBlockHash,
			StatePrevDiffHash: b.props.StatePrevDiffHash,
			StateCurrentHash:  b.props.StateCurrentHash,
		},
	}

	bytes, err := tmpBlock.Serialize()
	if err != nil {
		return nil, err
	}

	hashedBytes := hashing.Hash(bytes)
	return hashedBytes[:], nil
}

// Equals ...
func (b Block) Equals(other merkletree.Content) (bool, error) {
	bHash, err := b.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	oHash, err := other.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	return string(bHash) == string(oHash), nil
}

// SetHash ...
func (b *Block) SetHash() error {
	if b == nil {
		return ErrNilBlock
	}

	hash, err := b.CalculateHash()
	if err != nil {
		return err
	}

	b.props.BlockHash = &hash

	return nil
}

// MarshalJSON ...
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.props)
}

// UnmarshalJSON ...
func (b *Block) UnmarshalJSON(data []byte) error {
	var props BlockProps
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	b.props = props

	return nil
}
