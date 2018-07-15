package c3crypto

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogo/protobuf/proto"
	glc "github.com/libp2p/go-libp2p-crypto"
	pb "github.com/libp2p/go-libp2p-crypto/pb"
)

// note: these structs implement the go-libp2p-crypto private and public key interfaces
// https://godoc.org/github.com/ipfs/go-libp2p-crypto

type signature struct {
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

// WrappedPrivateKeyProps ...
type WrappedPrivateKeyProps struct {
	Priv *ecdsa.PrivateKey
}

// WrappedPrivateKey ...
type WrappedPrivateKey struct {
	props WrappedPrivateKeyProps
}

// WrappedPublicKeyProps ...
type WrappedPublicKeyProps struct {
	Pub *ecdsa.PublicKey
}

// WrappedPublicKey ...
type WrappedPublicKey struct {
	props WrappedPublicKeyProps
}

// NewWrappedKeyPair ...
func NewWrappedKeyPair(priv *ecdsa.PrivateKey) (*WrappedPrivateKey, *WrappedPublicKey) {
	if priv == nil {
		return nil, nil
	}

	wPriv := NewWrappedPrivateKey(&WrappedPrivateKeyProps{
		Priv: priv,
	})

	return wPriv, NewWrappedPublicKey(&WrappedPublicKeyProps{
		Pub: &priv.PublicKey,
	})
}

// NewWrappedPrivateKey ...
func NewWrappedPrivateKey(props *WrappedPrivateKeyProps) *WrappedPrivateKey {
	if props == nil {
		return nil
	}
	if props.Priv == nil {
		return nil
	}

	return &WrappedPrivateKey{
		props: *props,
	}
}

// NewWrappedPublicKey ...
func NewWrappedPublicKey(props *WrappedPublicKeyProps) *WrappedPublicKey {
	if props == nil {
		return nil
	}
	if props.Pub == nil {
		return nil
	}

	return &WrappedPublicKey{
		props: *props,
	}
}

// Bytes ...
func (wPriv *WrappedPrivateKey) Bytes() ([]byte, error) {
	return SerializePrivateKey(wPriv.props.Priv)
}

// Equals ...
func (wPriv *WrappedPrivateKey) Equals(other glc.Key) bool {
	if other == nil {
		return false
	}

	return reflect.DeepEqual(wPriv, other)
}

// Sign ...
func (wPriv *WrappedPrivateKey) Sign(data []byte) ([]byte, error) {
	r, s, err := Sign(wPriv.props.Priv, data)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&signature{
		R: r,
		S: s,
	})
}

// GetPublic ...
func (wPriv *WrappedPrivateKey) GetPublic() glc.PubKey {
	return NewWrappedPublicKey(&WrappedPublicKeyProps{
		Pub: &wPriv.props.Priv.PublicKey,
	})
}

// Bytes ...
func (wPub *WrappedPublicKey) Bytes() ([]byte, error) {
	pbmes := new(pb.PublicKey)
	typ := pb.KeyType_RSA
	pbmes.Type = &typ
	pbmes.Data = crypto.FromECDSAPub(wPub.props.Pub)
	return proto.Marshal(pbmes)
}

// Equals ...
func (wPub *WrappedPublicKey) Equals(other glc.Key) bool {
	if other == nil {
		return false
	}

	// note: hate using reflect here, any better way?
	return reflect.DeepEqual(wPub, other)
}

// Verify ...
func (wPub *WrappedPublicKey) Verify(data []byte, sig []byte) (bool, error) {
	var tmpSig signature
	if err := json.Unmarshal(sig, &tmpSig); err != nil {
		return false, err
	}

	return Verify(wPub.props.Pub, data, tmpSig.R, tmpSig.S)
}
