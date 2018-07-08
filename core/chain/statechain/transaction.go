package statechain

import (
	"crypto/ecdsa"
	"encoding/json"

	"github.com/c3systems/c3/common/hashing"
	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
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
	if tx == nil {
		return ErrNilTx
	}

	var tmpProps TransactionProps
	if err := json.Unmarshal(bytes, &tmpProps); err != nil {
		return err
	}

	tx.props = tmpProps
	return nil
}

// SerializeString ...
func (tx Transaction) SerializeString() (string, error) {
	bytes, err := tx.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeString(string(bytes)), nil
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

// CalcHash ...
func (tx Transaction) CalcHash() (string, error) {
	tmpTx := Transaction{
		props: TransactionProps{
			ImageHash: tx.props.ImageHash,
			Method:    tx.props.Method,
			Payload:   tx.props.Payload,
			From:      tx.props.From,
		},
	}

	bytes, err := tmpTx.Serialize()
	if err != nil {
		return "", err
	}

	return hashing.HashToHexString(bytes), nil
}

// SetHash ...
func (tx *Transaction) SetHash() error {
	if tx == nil {
		return ErrNilTx
	}

	hash, err := tx.CalcHash()
	if err != nil {
		return err
	}

	tx.props.TxHash = &hash

	return nil
}

// CalcSig ...
func (tx Transaction) CalcSig(priv *ecdsa.PrivateKey) (*TxSig, error) {
	hash, err := tx.CalcHash()
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
