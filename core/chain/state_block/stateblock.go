package stateblock

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

type Props struct {
	BlockHash         *string `json:"blockHash,omitempty"`
	BlockNumber       string  `json:"blockNumber"`
	BlockTime         string  `json:"blockTime"` // unix timestamp
	ImageHash         string  `json:"imageHash"`
	TxsHash           string  `json:"txsHash"`
	StatePrevDiffHash string  `json:"statePrevDiffHash"`
	StateCurrentHash  string  `json:"stateCurrentHash"`
}

// Block ...
type Block struct {
	props Props
}

func New(props *Props) *Block {
	if props == nil {
		return &Block{}
	}

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
