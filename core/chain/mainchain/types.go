package mainchain

import (
	"errors"
)

var (
	// ErrNilBlock ...
	ErrNilBlock = errors.New("block is nil")
)

// ImageHash is the main chain identifier
// The main chain does not have an image (i.e. the image hash is nil).
// The hex encoded, sha512_256 hash of a nil bytes array is ImageHash
// https://play.golang.org/p/69Z8ot5uly5
const ImageHash = "0xc672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a"

// BlockSig ...
type BlockSig struct {
	R string `json:"r"`
	S string `json:"s"`
}

// Props ...
type Props struct {
	BlockHash             *string   `json:"blockHash,omitempty"`
	BlockNumber           string    `json:"blockNumber"`
	BlockTime             string    `json:"blockTime"` // unix timestamp
	ImageHash             string    `json:"imageHash"`
	StateBlocksMerkleHash string    `json:"stateBlocksMerkleHash"`
	StateBlockHashes      []*string `json:"stateBlockHashes"`
	PrevBlockHash         string    `json:"prevBlockHash"`
	Nonce                 string    `json:"nonce"`
	Difficulty            string    `json:"difficulty"`
	MinerAddress          string    `json:"minerAddress"`
	MinerSig              *BlockSig `json:"blockSig,omitempty"`
}

// Block ...
type Block struct {
	props Props
}
