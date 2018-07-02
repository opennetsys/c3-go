package statechain

// TransactionsMap is a list of transactions by image hashes
type TransactionsMap map[string][]*Transaction

// TransactionProps
type TransactionProps struct {
	TxHash  *string
	Method  string
	Payload interface{}
}

// Transaction ...
type Transaction struct {
	props Props
}

// StateBlockProps
type StateBlockProps struct {
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
