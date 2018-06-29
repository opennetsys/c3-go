package mainblock

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	mh "github.com/multiformats/go-multihash"
	//cbor "gx/ipfs/QmRVSCwQtW1rjHCay9NqKXDwbtKTgDcN4iY7PrpSqfKM5D/go-ipld-cbor"
	//mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
)

// MainChainImageHash is the main chain identifier
// The main chain does not have an image (i.e. the image hash is nil).
// The hex encoded, sha256 hash of a nil bytes array is MainChainImageHash
// https://play.golang.org/p/33_3vY6XyjD
const MainChainImageHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

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

func New(props *Props) *Block {
	if props == nil {
		return &Block{
			props: Props{
				ImageHash: MainChainImageHash,
			},
		}
	}

	props.ImageHash = MainChainImageHash
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
