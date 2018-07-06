package mainchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/c3systems/c3/core/chain/statechain"
	//cbor "gx/ipfs/QmRVSCwQtW1rjHCay9NqKXDwbtKTgDcN4iY7PrpSqfKM5D/go-ipld-cbor"
	//mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
)

// ImageHash is the main chain identifier
// The main chain does not have an image (i.e. the image hash is nil).
// The hex encoded, sha256 hash of a nil bytes array is ImageHash
// https://play.golang.org/p/33_3vY6XyjD
const ImageHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Props ...
type Props struct {
	BlockHash       *string `json:"blockHash,omitempty"`
	BlockNumber     string  `json:"blockNumber"`
	BlockTime       string  `json:"blockTime"` // unix timestamp
	ImageHash       string  `json:"imageHash"`
	StateBlocksHash string  `json:"stateBlocksHash"`
	PrevBlockHash   string  `json:"prevBlockHash"`
	Nonce           string  `json:"nonce"`
	Difficulty      string  `json:"difficulty"`
}

// Block ...
type Block struct {
	props Props
}

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
	return json.Marshal(b.props)
}

// Deserialize ...
func (b *Block) Deserialize(bytes []byte) error {
	var tmpProps Props
	if err := json.Unmarshal(bytes, &tmpProps); err != nil {
		return err
	}

	b.props = tmpProps
	return nil
}

// SerializeString ...
func (b Block) SerializeString() (string, error) {
	bytes, err := json.Marshal(b.props)
	return hex.EncodeToString(bytes), err
}

// DeserializeString ...
func (b *Block) DeserializeString(str string) error {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return err
	}

	return b.Deserialize(bytes)
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

	shaSum := sha256.Sum256(bytes)
	return hex.EncodeToString(shaSum[:]), nil
}

// VerifyBlock verifies a block
// TODO: everything
func VerifyBlock(block *Block) (bool, error) {
	return false, nil
}

// HashProps ...
func HashProps(props Props) (string, error) {
	if props.BlockHash != nil {
		return *props.BlockHash, nil
	}

	bytes, err := json.Marshal(props)
	if err != nil {
		return "", err
	}

	shaSum := sha256.Sum256(bytes)
	return hex.EncodeToString(shaSum[:]), nil
}

// NewFromStateBlocks ...
// TODO: everything...
func NewFromStateBlocks(stateBlocks []*statechain.Block) (*Block, error) {
	return nil, nil
}
