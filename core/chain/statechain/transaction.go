package statechain

import (
	"crypto/ecdsa"
	"encoding/json"

	"github.com/c3systems/c3/common/coder"
	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
	"github.com/c3systems/merkletree"
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
func (tx *Transaction) Props() TransactionProps {
	return tx.props
}

// Serialize ...
func (tx *Transaction) Serialize() ([]byte, error) {
	return coder.Serialize(tx.props)
}

// Deserialize ...
func (tx *Transaction) Deserialize(data []byte) error {
	var tmpProps TransactionProps
	if err := coder.Deserialize(data, &tmpProps); err != nil {
		return err
	}

	tx.props = tmpProps
	return nil
}

// SerializeString ...
func (tx *Transaction) SerializeString() (string, error) {
	data, err := tx.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(data)), nil
}

// DeserializeString ...
func (tx *Transaction) DeserializeString(hexStr string) error {
	if tx == nil {
		return ErrNilTx
	}

	str, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return tx.Deserialize([]byte(str))
}

// CalculateHash ...
func (tx *Transaction) CalculateHash() (string, error) {
	bytes, err := tx.CalculateHashBytes()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
}

// CalculateHashBytes ...
func (tx *Transaction) CalculateHashBytes() ([]byte, error) {
	tmpTx := Transaction{
		props: TransactionProps{
			ImageHash: tx.props.ImageHash,
			Method:    tx.props.Method,
			Payload:   tx.props.Payload,
			From:      tx.props.From,
		},
	}

	data, err := tmpTx.Serialize()
	if err != nil {
		return nil, err
	}

	hashedBytes := hashing.Hash(data)
	return hashedBytes[:], nil
}

// Equals ...
func (tx *Transaction) Equals(other merkletree.Content) (bool, error) {
	tHash, err := tx.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	oHash, err := other.CalculateHashBytes()
	if err != nil {
		return false, err
	}

	return string(tHash) == string(oHash), nil
}

// SetHash ...
func (tx *Transaction) SetHash() error {
	if tx == nil {
		return ErrNilTx
	}

	hash, err := tx.CalculateHash()
	if err != nil {
		return err
	}

	tx.props.TxHash = &hash

	return nil
}

// CalcSig ...
func (tx *Transaction) CalcSig(priv *ecdsa.PrivateKey) (*TxSig, error) {
	hash, err := tx.CalculateHash()
	if err != nil {
		return nil, err
	}

	r, s, err := c3crypto.Sign(priv, []byte(hash))
	if err != nil {
		return nil, err
	}

	return &TxSig{
		R: hexutil.EncodeBigInt(r),
		S: hexutil.EncodeBigInt(s),
	}, nil
}

// SetSig ...
func (tx *Transaction) SetSig(priv *ecdsa.PrivateKey) error {
	if tx == nil {
		return ErrNilTx
	}

	sig, err := tx.CalcSig(priv)
	if err != nil {
		return err
	}

	tx.props.Sig = sig

	return nil
}

// MarshalJSON ...
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(tx.props)
}

// UnmarshalJSON ...
func (tx *Transaction) UnmarshalJSON(data []byte) error {
	var props TransactionProps
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	tx.props = props

	return nil
}
