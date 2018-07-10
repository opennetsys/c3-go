package mainchain

import (
	"bytes"
	"encoding/gob"

	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
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
func (b Block) Props() Props {
	return b.props
}

// Serialize ...
func (b Block) Serialize() ([]byte, error) {
	data := new(bytes.Buffer)
	err := gob.NewEncoder(data).Encode(b.props)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

// Deserialize ...
func (b *Block) Deserialize(data []byte) error {
	if b == nil {
		return ErrNilBlock
	}

	// note: may not be able to recreate in other languages?
	var tmpProps Props
	byts := bytes.NewBuffer(data)
	if err := gob.NewDecoder(byts).Decode(&tmpProps); err != nil {
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
		props: Props{
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
