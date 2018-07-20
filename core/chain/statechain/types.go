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
	// ErrNilDiff ...
	ErrNilDiff = errors.New("diff is nil")
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
	TxHash    *string `json:"txHash,omitempty" rlp:"nil"`
	ImageHash string  `json:"imageHash"`
	Method    string  `json:"method"`
	Payload   []byte  `json:"payload"`
	From      string  `json:"from"`
	Sig       *TxSig  `json:"txSig,omitempty" rlp:"nil"`
}

// Transaction ...
type Transaction struct {
	props TransactionProps
}

// BlockProps ...
type BlockProps struct {
	BlockHash         *string `json:"blockHash,omitempty" rlp:"nil"`
	BlockNumber       string  `json:"blockNumber"`
	BlockTime         string  `json:"blockTime"` // unix timestamp
	ImageHash         string  `json:"imageHash"`
	TxHash            string  `json:"txHash"`
	PrevBlockHash     string  `json:"prevBlockHash"`
	StatePrevDiffHash string  `json:"statePrevDiffHash"`
	StateCurrentHash  string  `json:"stateCurrentHash"`
}

// Block ...
type Block struct {
	props BlockProps
}

// DiffProps ...
// note @miguelmota: any better system than simply storing as a string?
type DiffProps struct {
	DiffHash *string `json:"diffHash,omitempty" rlp:"nil"`
	// what's the best way to store a diff?
	Data string `json:"data"`
}

// Diff ...
type Diff struct {
	props DiffProps
}
