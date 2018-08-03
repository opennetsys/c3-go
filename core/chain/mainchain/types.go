package mainchain

import (
	"errors"
)

var (
	// ErrNilBlock ...
	ErrNilBlock = errors.New("block is nil")
	// GenesisBlockHash ...
	GenesisBlockHash = "0x90770ba62574b68607ebd14a41dd2eb9df3c537df11b68647faea34a88315e49"
	// GenesisBlock ...
	GenesisBlock = Block{
		props: Props{
			BlockHash:             &GenesisBlockHash,
			BlockNumber:           "0x0",
			BlockTime:             "0x0",
			ImageHash:             ImageHash,
			StateBlocksMerkleHash: "0x",
			PrevBlockHash:         "0x",
			Nonce:                 "0x",
			Difficulty:            "0x0",
			MinerAddress:          "0x",
			MinerSig:              nil,
		},
	}
)

// ImageHash is the main chain identifier
// The main chain does not have an image (i.e. the image hash is nil).
// The hex encoded, sha2_256 hash of a nil bytes array is ImageHash
// https://play.golang.org/p/69Z8ot5uly5
const ImageHash = "0xc672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a"

// MinerSig ...
type MinerSig struct {
	R string `json:"r"`
	S string `json:"s"`
}

// Props ...
type Props struct {
	BlockHash             *string   `json:"blockHash,omitempty" rlp:"nil"`
	BlockNumber           string    `json:"blockNumber"`
	BlockTime             string    `json:"blockTime"` // unix timestamp
	ImageHash             string    `json:"imageHash"`
	StateBlocksMerkleHash string    `json:"stateBlocksMerkleHash"`
	PrevBlockHash         string    `json:"prevBlockHash"`
	Nonce                 string    `json:"nonce"`
	Difficulty            string    `json:"difficulty"`
	MinerAddress          string    `json:"minerAddress"`
	MinerSig              *MinerSig `json:"minerSig,omitempty" rlp:"nil"`
}

// Block ...
type Block struct {
	props Props
}
