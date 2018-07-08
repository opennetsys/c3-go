package mainchain

import (
	"encoding/json"

	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
)

// New ...
func New(props *BlockProps) *Block {
	if props == nil {
		return &Block{
			props: BlockProps{
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
func (b Block) Props() BlockProps {
	return b.props
}

// Serialize ...
func (b Block) Serialize() ([]byte, error) {
	return json.Marshal(b.props)
}

// Deserialize ...
func (b *Block) Deserialize(bytes []byte) error {
	if b == nil {
		return ErrNilBlock
	}

	var tmpProps BlockProps
	if err := json.Unmarshal(bytes, &tmpProps); err != nil {
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

// CalcHash ...
func (b Block) CalcHash() (string, error) {
	tmpBlock := Block{
		props: BlockProps{
			BlockNumber:           b.props.BlockNumber,
			BlockTime:             b.props.BlockTime,
			ImageHash:             b.props.ImageHash,
			StateBlocksMerkleHash: b.props.StateBlocksMerkleHash,
			StateBlockHashes:      b.props.StateBlockHashes,
			PrevBlockHash:         b.props.PrevBlockHash,
			Nonce:                 b.props.Nonce,
			Difficulty:            b.props.Difficulty,
			MinerAddress:          b.props.MinerAddress,
		},
	}

	bytes, err := tmpBlock.Serialize()
	if err != nil {
		return "", err
	}

	return hashing.HashToHexString(bytes), nil
}

// SetHash ...
func (b *Block) SetHash() error {
	if b == nil {
		return ErrNilBlock
	}

	hash, err := b.CalcHash()
	if err != nil {
		return err
	}

	b.props.BlockHash = &hash

	return nil
}
