package statechain

import (
	"crypto/ecdsa"

	"github.com/c3systems/c3/common/hexutil"
	"github.com/c3systems/c3/core/c3crypto"
)

// BuildNextState ...
// TODO: everything...
func BuildNextState(imageHash string, transactions []*Transaction) (*Block, error) {
	return nil, nil
}

// VerifyBlock verifies a block
// TODO: everything
func VerifyBlock(block *Block) (bool, error) {
	return false, nil
}

// VerifyTransaction ...
func VerifyTransaction(tx *Transaction) (bool, error) {
	// note: we hash the message and then sign the hash
	// TODO: check the image hash exists?
	// TODO: check for blank inputs?
	if tx == nil {
		return false, ErrNilTx
	}

	// 1. tx must have a hash
	if tx.props.TxHash == nil {
		return false, ErrNoHash
	}

	// 2. tx must have a sig
	if tx.props.Sig == nil {
		return false, ErrNoSig
	}

	// 3. verify the hash
	tmpProps := TransactionProps{
		ImageHash: tx.props.ImageHash,
		Method:    tx.props.Method,
		Payload:   tx.props.Payload,
		From:      tx.props.From,
	}
	tmpTx := Transaction{
		props: tmpProps,
	}

	tmpHash, err := tmpTx.CalcHash()
	if err != nil {
		return false, err
	}
	// note: already checked for nil hash
	if *tx.props.TxHash != tmpHash {
		return false, nil
	}

	// 4. the sig must verify
	pub, err := PubFromTx(tx)
	if err != nil {
		return false, err
	}

	// note: checked for nil sig, above
	r, err := hexutil.DecodeBigInt(tx.props.Sig.R)
	if err != nil {
		return false, err
	}
	s, err := hexutil.DecodeBigInt(tx.props.Sig.S)
	if err != nil {
		return false, err
	}

	return c3crypto.Verify(pub, []byte(*tx.props.TxHash), r, s)
}

// PubFromTx ...
func PubFromTx(tx *Transaction) (*ecdsa.PublicKey, error) {
	if tx == nil {
		return nil, ErrNilTx
	}

	pubStr, err := hexutil.DecodeString(tx.props.From)
	if err != nil {
		return nil, err
	}
	pub, err := c3crypto.DeserializePublicKey([]byte(pubStr))
	if err != nil {
		return nil, err
	}
	if pub == nil {
		return nil, ErrInvalidFromAddress
	}

	return pub, nil

}
