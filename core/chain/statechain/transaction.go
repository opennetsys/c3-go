package statechain

import (
	"encoding/hex"
	"encoding/json"

	"github.com/c3systems/c3/common/hashing"
)

// NewTransaction ...
func NewTransaction(props *TransactionProps) *Transaction {
	if props == nil {
		return &Transaction{}
	}

	return &Transaction{
		props: *props,
	}
}

// Props ...
func (tx Transaction) Props() TransactionProps {
	return tx.props
}

// Serialize ...
func (tx Transaction) Serialize() ([]byte, error) {
	return json.Marshal(tx.props)
}

// Deserialize ...
func (tx *Transaction) Deserialize(bytes []byte) error {
	var tmpProps TransactionProps
	if err := json.Unmarshal(bytes, &tmpProps); err != nil {
		return err
	}

	tx.props = tmpProps
	return nil
}

// SerializeString ...
func (tx Transaction) SerializeString() (string, error) {
	bytes, err := json.Marshal(tx.props)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// DeserializeString ...
func (tx *Transaction) DeserializeString(str string) error {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return err
	}

	return tx.Deserialize(bytes)
}

// Hash ...
func (tx Transaction) Hash() (string, error) {
	if tx.props.TxHash != nil {
		return *tx.props.TxHash, nil
	}

	bytes, err := tx.Serialize()
	if err != nil {
		return "", err
	}

	return hashing.HashToHexString(bytes), nil
}
