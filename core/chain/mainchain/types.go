package mainchain

import "errors"

var (
	// ErrNilBlock ...
	ErrNilBlock = errors.New("block is nil")
)

// ImageHash is the main chain identifier
// The main chain does not have an image (i.e. the image hash is nil).
// The hex encoded, sha256 hash of a nil bytes array is ImageHash
// https://play.golang.org/p/33_3vY6XyjD
const ImageHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// BlockSig ...
type BlockSig struct {
	R string `json:"r"`
	S string `json:"s"`
}

// BlockProps ...
type BlockProps struct {
	BlockHash             *string   `json:"blockHash,omitempty"`
	BlockNumber           string    `json:"blockNumber"`
	BlockTime             string    `json:"blockTime"` // unix timestamp
	ImageHash             string    `json:"imageHash"`
	StateBlocksMerkleHash string    `json:"stateBlocksMerkleHash"`
	StateBlockHashes      []string  `json:"stateBlockHashes"`
	PrevBlockHash         string    `json:"prevBlockHash"`
	Nonce                 string    `json:"nonce"`
	Difficulty            string    `json:"difficulty"`
	MinerAddress          string    `json:"minerAddress"`
	MinerSig              *BlockSig `json:"blockSig,omitempty"`
}

// Block ...
type Block struct {
	props BlockProps
}
