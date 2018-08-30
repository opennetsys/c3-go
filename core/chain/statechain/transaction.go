package statechain

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"

	"github.com/c3systems/c3-go/common/c3crypto"
	"github.com/c3systems/c3-go/common/coder"
	"github.com/c3systems/c3-go/common/hashutil"
	"github.com/c3systems/c3-go/common/hexutil"
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
	tmp, err := BuildCoderFromTransaction(tx)
	if err != nil {
		return nil, err
	}

	bytes, err := tmp.Marshal()
	if err != nil {
		return nil, err
	}

	return coder.AppendCode(bytes), nil
}

// Deserialize ...
func (tx *Transaction) Deserialize(data []byte) error {
	if data == nil {
		return errors.New("nil bytes")
	}
	if tx == nil {
		return errors.New("nil tx")
	}

	_, bytes, err := coder.StripCode(data)
	if err != nil {
		return err
	}

	props, err := BuildTransactionPropsFromBytes(bytes)
	if err != nil {
		return err
	}
	tx.props = *props

	return nil
}

// SerializeString ...
func (tx *Transaction) SerializeString() (string, error) {
	data, err := tx.Serialize()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(data), nil
}

// DeserializeString ...
func (tx *Transaction) DeserializeString(hexStr string) error {
	if tx == nil {
		return ErrNilTx
	}

	b, err := hexutil.DecodeString(hexStr)
	if err != nil {
		return err
	}

	return tx.Deserialize(b)
}

// CalculateHash ...
func (tx *Transaction) CalculateHash() (string, error) {
	bts, err := tx.CalculateHashBytes()
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(bts), nil
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

	hashedBytes := hashutil.Hash(data)
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

// BuildCoderFromTransaction ...
func BuildCoderFromTransaction(tx *Transaction) (*coder.Transaction, error) {
	// note: is there a better way to handle interfaces in protobufs? google.protobuf.Any? bson?
	payloadBytes, err := json.Marshal(tx.props.Payload)
	if err != nil {
		return nil, err
	}

	tmp := &coder.Transaction{
		ImageHash: tx.props.ImageHash,
		Method:    tx.props.Method,
		Payload:   payloadBytes,
		From:      tx.props.From,
	}

	// note: is there a better way to handle nil with protobuff?
	if tx.props.TxHash != nil {
		tmp.TxHash = *tx.props.TxHash
	}
	if tx.props.Sig != nil {
		tmp.Sig = &coder.TxSig{
			R: tx.props.Sig.R,
			S: tx.props.Sig.S,
		}
	}

	return tmp, nil
}

// BuildTransactionPropsFromBytes ...
func BuildTransactionPropsFromBytes(data []byte) (*TransactionProps, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	c, err := BuildTransactionCoderFromBytes(data)
	if err != nil {
		return nil, err
	}

	return BuildTransactionPropsFromCoder(c)
}

// BuildTransactionCoderFromBytes ...
func BuildTransactionCoderFromBytes(data []byte) (*coder.Transaction, error) {
	if data == nil {
		return nil, errors.New("nil bytes")
	}

	tmp := new(coder.Transaction)
	if err := tmp.Unmarshal(data); err != nil {
		return nil, err
	}
	if tmp == nil {
		return nil, errors.New("nil output")
	}

	return tmp, nil
}

// BuildTransactionPropsFromCoder ...
func BuildTransactionPropsFromCoder(tmp *coder.Transaction) (*TransactionProps, error) {
	if tmp == nil {
		return nil, errors.New("nil coder")
	}

	props := &TransactionProps{
		ImageHash: tmp.ImageHash,
		Method:    tmp.Method,
		From:      tmp.From,
	}

	if tmp.Payload != nil {
		var v []byte
		if err := json.Unmarshal(tmp.Payload, &v); err != nil {
			return nil, err
		}
		props.Payload = v
	}

	// note: is there a better way to handle nil with protobuff?
	if tmp.TxHash != "" {
		s := tmp.TxHash
		props.TxHash = &s
	}
	if tmp.Sig != nil && tmp.Sig.R != "" && tmp.Sig.S != "" {
		sig := &TxSig{
			R: tmp.Sig.R,
			S: tmp.Sig.S,
		}

		props.Sig = sig
	}

	return props, nil
}
