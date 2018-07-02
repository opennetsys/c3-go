package statechain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func New(props *TransactionProps) *Transaction {
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
func (tx Transcation) Serialize() ([]byte, error) {
	return json.Marshal(tx.props)
}

// Deserialize ...
func (tx *Transaction) Deserialize(bytes []byte) error {
	var tmpProps Props
	if err := json.Unmarshal(bytes, &tmpProps); err != nil {
		return err
	}

	tx.props = tmpProps
	return nil
}

// SerializeString ...
func (tx Transaction) SerializeString() (string, error) {
	return hex.EncodeToString(json.Marshal(tx.props))
}

// DeserializeString ...
func (tx *Transaction) DeserializeString(str string) error {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	return tx.Deserialize(bytes)
}

// Hash ...
func (tx Transaction) Hash() (string, error) {
	if tx.props.BlockHash != nil {
		return *tx.props.BlockHash, nil
	}

	bytes, err := tx.Serialize()
	if err != nil {
		return "", err
	}

	shaSum := sha256.Sum256(bytes)
	return hex.EncodeToString(shaSum[:]), nil
}
