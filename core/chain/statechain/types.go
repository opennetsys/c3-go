package statechain

import "errors"

var (
	// ErrNoHash ...
	ErrNoHash = errors.New("no hash present")
	// ErrNilTx ...
	ErrNilTx = errors.New("transaction is nil")
	// ErrNoSig ...
	ErrNoSig = errors.New("no signature present")
	// ErrInvalidFromAddress ...
	ErrInvalidFromAddress = errors.New("from address is not valid")
	// ErrNilBlock ...
	ErrNilBlock = errors.New("block is nil")
)

// TxSig ...
type TxSig struct {
	R string `json:"r"`
	S string `json:"s"`
}

// TransactionsMap is a list of transactions by image hashes
type TransactionsMap map[string][]*Transaction

// TransactionProps ...
type TransactionProps struct {
	TxHash    *string     `json:"txHash,omitempty"`
	ImageHash string      `json:"imageHash"`
	Method    string      `json:"method"`
	Payload   interface{} `json:"payload"`
	From      string      `json:"from"`
	Sig       *TxSig      `json:"txSig,omitempty"`
}

// Transaction ...
type Transaction struct {
	props TransactionProps
}

// BlockProps ...
type BlockProps struct {
	BlockHash         *string  `json:"blockHash,omitempty"`
	BlockNumber       string   `json:"blockNumber"`
	BlockTime         string   `json:"blockTime"` // unix timestamp
	ImageHash         string   `json:"imageHash"`
	TxsMerkleHash     string   `json:"txsMerkleHash"`
	TxHashes          []string `json:"txHashes"`
	StatePrevDiffHash string   `json:"statePrevDiffHash"`
	StateCurrentHash  string   `json:"stateCurrentHash"`
}

// Block ...
type Block struct {
	props BlockProps
}
