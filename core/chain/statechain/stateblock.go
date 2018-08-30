package statechain

import (
	"encoding/json"
	"errors"

	"github.com/c3systems/c3-go/common/coder"
	"github.com/c3systems/c3-go/common/hashutil"
	"github.com/c3systems/c3-go/common/hexutil"

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
	tmp := BuildCoderFromBlock(b)
	bytes, err := tmp.Marshal()
	if err != nil {
		return nil, err
	}

	return coder.AppendCode(bytes), nil
}

// Deserialize ...
func (b *Block) Deserialize(data []byte) error {
	if data == nil {
		return errors.New("nil bytes")
	}
	if b == nil {
		return errors.New("nil block")
	}

	_, bytes, err := coder.StripCode(data)
	if err != nil {
		return err
	}

	props, err := BuildBlockPropsFromBytes(bytes)
	if err != nil {
		return err
	}
	b.props = *props

	return nil
}

// SerializeString ...
func (b Block) SerializeString() (string, error) {
	bts, err := b.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(bts), nil
}

// DeserializeString ...
func (b *Block) DeserializeString(hexStr string) error {
	if b == nil {
		return ErrNilBlock
	}

	bts, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return b.Deserialize(bts)
}

// CalculateHash ...
func (b Block) CalculateHash() (string, error) {
	bts, err := b.CalculateHashBytes()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(bts), nil
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

	hashedBytes := hashutil.Hash(bytes)
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

// BuildCoderFromBlock ...
func BuildCoderFromBlock(b *Block) *coder.StatechainBlock {
	tmp := &coder.StatechainBlock{
		BlockNumber:       b.props.BlockNumber,
		BlockTime:         b.props.BlockTime,
		ImageHash:         b.props.ImageHash,
		TxHash:            b.props.TxHash,
		PrevBlockHash:     b.props.PrevBlockHash,
		StatePrevDiffHash: b.props.StatePrevDiffHash,
		StateCurrentHash:  b.props.StateCurrentHash,
	}

	// note: is there a better way to handle nil with protobuff?
	if b.props.BlockHash != nil {
		tmp.BlockHash = *b.props.BlockHash
	}

	return tmp
}

// BuildBlockPropsFromBytes ...
func BuildBlockPropsFromBytes(data []byte) (*BlockProps, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	c, err := BuildBlockCoderFromBytes(data)
	if err != nil {
		return nil, err
	}

	return BuildBlockPropsFromCoder(c)
}

// BuildBlockCoderFromBytes ...
func BuildBlockCoderFromBytes(data []byte) (*coder.StatechainBlock, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	tmp := new(coder.StatechainBlock)
	if err := tmp.Unmarshal(data); err != nil {
		return nil, err
	}

	return tmp, nil
}

// BuildBlockPropsFromCoder ...
func BuildBlockPropsFromCoder(tmp *coder.StatechainBlock) (*BlockProps, error) {
	if tmp == nil {
		return nil, errors.New("nil bytes")
	}

	props := &BlockProps{
		BlockNumber:       tmp.BlockNumber,
		BlockTime:         tmp.BlockTime,
		ImageHash:         tmp.ImageHash,
		TxHash:            tmp.TxHash,
		PrevBlockHash:     tmp.PrevBlockHash,
		StatePrevDiffHash: tmp.StatePrevDiffHash,
		StateCurrentHash:  tmp.StateCurrentHash,
	}
	// note: is there any better way of checking forn nil with protobuf?
	if tmp.BlockHash != "" {
		s := tmp.BlockHash
		props.BlockHash = &s
	}

	return props, nil
}
