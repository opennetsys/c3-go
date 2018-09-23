package mainchain

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/c3systems/c3-go/common/coder"
	"github.com/c3systems/c3-go/common/hashutil"
	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/c3systems/merkletree"
)

// New ...
func New(props *Props) *Block {
	if props == nil {
		return &Block{
			props: Props{
				ImageHash: ImageHash,
			},
		}
	}

	props.ImageHash = ImageHash
	return &Block{
		props: *props,
	}
}

// Props ...
func (b *Block) Props() Props {
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
		return nil
	}

	props, err := BuildBlockPropsFromBytes(bytes)
	if err != nil {
		return err
	}

	b.props = *props

	return nil
}

// SerializeString ...
func (b *Block) SerializeString() (string, error) {
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
func (b *Block) CalculateHash() (string, error) {
	bts, err := b.CalculateHashBytes()
	if err != nil {
		return "", err
	}

	return strings.ToLower(hexutil.EncodeToString(bts)), nil
}

// CalculateHashBytes ...
func (b *Block) CalculateHashBytes() ([]byte, error) {
	tmpBlock := Block{
		props: Props{
			BlockNumber:           b.props.BlockNumber,
			BlockTime:             b.props.BlockTime,
			ImageHash:             b.props.ImageHash,
			StateBlocksMerkleHash: b.props.StateBlocksMerkleHash,
			PrevBlockHash:         b.props.PrevBlockHash,
			Nonce:                 b.props.Nonce,
			Difficulty:            b.props.Difficulty,
			MinerAddress:          b.props.MinerAddress,
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
func (b *Block) Equals(other merkletree.Content) (bool, error) {
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

// SetMinerSig ...
func (b *Block) SetMinerSig(sig *MinerSig) error {
	b.props.MinerSig = sig

	return nil
}

// MarshalJSON ...
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.props)
}

// UnmarshalJSON ...
func (b *Block) UnmarshalJSON(data []byte) error {
	var props Props
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	b.props = props

	return nil
}

// BuildCoderFromBlock ...
func BuildCoderFromBlock(b *Block) *coder.MainchainBlock {
	tmp := &coder.MainchainBlock{
		BlockNumber:           b.props.BlockNumber,
		BlockTime:             b.props.BlockTime,
		ImageHash:             b.props.ImageHash,
		StateBlocksMerkleHash: b.props.StateBlocksMerkleHash,
		PrevBlockHash:         b.props.PrevBlockHash,
		Nonce:                 b.props.Nonce,
		Difficulty:            b.props.Difficulty,
		MinerAddress:          b.props.MinerAddress,
	}

	// note: is there a better way to handle nil with protobuff?
	if b.props.BlockHash != nil {
		tmp.BlockHash = *b.props.BlockHash
	}
	if b.props.MinerSig != nil {
		tmp.MinerSig = &coder.MinerSig{
			R: b.props.MinerSig.R,
			S: b.props.MinerSig.S,
		}
	}

	return tmp
}

// BuildBlockPropsFromBytes ...
func BuildBlockPropsFromBytes(data []byte) (*Props, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	c, err := BuildCoderFromBytes(data)
	if err != nil {
		return nil, err
	}

	return BuildBlockPropsFromCoder(c)
}

// BuildCoderFromBytes ...
func BuildCoderFromBytes(data []byte) (*coder.MainchainBlock, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	tmp := new(coder.MainchainBlock)
	if err := tmp.Unmarshal(data); err != nil {
		return nil, err
	}
	if tmp == nil {
		return nil, errors.New("nil output")
	}

	return tmp, nil
}

// BuildBlockPropsFromCoder ...
func BuildBlockPropsFromCoder(tmp *coder.MainchainBlock) (*Props, error) {
	if tmp == nil {
		return nil, errors.New("nil coder")
	}

	props := &Props{
		BlockNumber:           tmp.BlockNumber,
		BlockTime:             tmp.BlockTime,
		ImageHash:             tmp.ImageHash,
		StateBlocksMerkleHash: tmp.StateBlocksMerkleHash,
		PrevBlockHash:         tmp.PrevBlockHash,
		Nonce:                 tmp.Nonce,
		Difficulty:            tmp.Difficulty,
		MinerAddress:          tmp.MinerAddress,
	}
	// note: is there any better way of checking forn nil with protobuf?
	if tmp.BlockHash != "" {
		s := tmp.BlockHash
		props.BlockHash = &s
	}
	if tmp.MinerSig != nil {
		props.MinerSig = &MinerSig{
			R: tmp.MinerSig.R,
			S: tmp.MinerSig.S,
		}
	}

	return props, nil
}
