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

type cidObj struct {
	BlockNumber string `json:"blockNumber"`
	ImageHash   string `json:"imageHash"`
}

type Props struct {
	TxsHash              string  `json:"txsHash"`
	ImageHash            string  `json:"imageHash"`
	StatePrevDiffHash    string  `json:"statePrevDiffHash"`
	StateGenesisDiffHash string  `json:"stateGenesisDiffHash"`
	StateCurrentHash     string  `json:"stateCurrentHash"`
	BlockNumber          string  `json:"blockNumber"`
	TimeStamp            string  `json:"timeStamp"` // unix timestamp
	Nonce                string  `json:"nonce"`
	Hash                 *string `json:"hash,omitempty"`
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

// FromBytes ...
func (b *Block) FromBytes(bytes []byte) error {
	var tmpProps Props
	if err := json.Unmarshal(bytes, &tmpProps); err != nil {
		return err
	}

	b.props = tmpProps
	return nil
}

// CID ...
func (b Block) CID() (*cid.Cid, error) {
	nd, err := cbor.WrapObject(cidObj{
		BlockNumber: b.props.BlockNumber,
		ImageHash:   b.props.ImageHash,
	}, mh.SHA2_256, -1)
	if err != nil {
		return nil, err
	}

	return nd.Cid(), nil
}

// Hash ...
// note: should be mined?
func (b Block) Hash() (string, error) {
	if b.props.Hash != nil {
		return *b.props.Hash, nil
	}

	bytes, err := b.Serialize()
	if err != nil {
		return "", err
	}

	shaSum := sha256.Sum256(bytes)
	return hex.EncodeToString(shaSum[:]), nil
}
