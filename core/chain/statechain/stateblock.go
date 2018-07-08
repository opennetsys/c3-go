package statechain

import (
	"encoding/hex"
	"encoding/json"

	"github.com/c3systems/c3/common/hashing"

	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	mh "github.com/multiformats/go-multihash"
)

// New ...
func New(props *StateBlockProps) *Block {
	if props == nil {
		return &Block{}
	}

	return &Block{
		props: *props,
	}
}

// Props ...
func (b Block) Props() StateBlockProps {
	return b.props
}

// Serialize ...
func (b Block) Serialize() ([]byte, error) {
	return json.Marshal(b.props)
}

// Deserialize ...
func (b *Block) Deserialize(bytes []byte) error {
	var tmpProps StateBlockProps
	if err := json.Unmarshal(bytes, &tmpProps); err != nil {
		return err
	}

	b.props = tmpProps
	return nil
}

// SerializeString ...
func (b Block) SerializeString() (string, error) {
	bytes, err := json.Marshal(b.props)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// DeserializeString ...
func (b *Block) DeserializeString(str string) error {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return err
	}

	return b.Deserialize(bytes)
}

// CID ...
// TODO: implement attributeName enum
func (b Block) CID(attributeName string) (*cid.Cid, error) {
	nd, err := cbor.WrapObject(struct {
		BlockNumber   string
		ImageHash     string
		AttributeName string
	}{
		BlockNumber:   b.props.BlockNumber,
		ImageHash:     b.props.ImageHash,
		AttributeName: attributeName,
	}, mh.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}

// VerifyBlock verifies a block
// TODO: everything
func VerifyBlock(block *Block) (bool, error) {
	return false, nil
}

// Hash ...
func (b Block) Hash() (string, error) {
	if b.props.BlockHash != nil {
		return *b.props.BlockHash, nil
	}

	bytes, err := b.Serialize()
	if err != nil {
		return "", err
	}

	return hashing.HashToHexString(bytes, nil)
}

// BuildNextState ...
// TODO: everything...
func BuildNextState(imageHash string, transactions []*Transaction) (*Block, error) {
	return nil, nil
}
