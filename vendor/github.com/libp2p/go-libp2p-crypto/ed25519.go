package crypto

import (
	"bytes"
	"fmt"
	"io"

	// TODO: change back to libp2p for pr
	pb "github.com/libp2p/go-libp2p-crypto/pb"

	"github.com/agl/ed25519"
	extra "github.com/agl/ed25519/extra25519"
	proto "github.com/gogo/protobuf/proto"
)

// Ed25519PrivateKey is an ed25519 private key
type Ed25519PrivateKey struct {
	sk *[64]byte
	pk *[32]byte
}

// Ed25519PublicKey is an ed25519 public key
type Ed25519PublicKey struct {
	k *[32]byte
}

// GenerateEd25519Key generate a new ed25519 private and public key pair
func GenerateEd25519Key(src io.Reader) (PrivKey, PubKey, error) {
	pub, priv, err := ed25519.GenerateKey(src)
	if err != nil {
		return nil, nil, err
	}

	return &Ed25519PrivateKey{
			sk: priv,
			pk: pub,
		},
		&Ed25519PublicKey{
			k: pub,
		},
		nil
}

// Bytes marshals an ed25519 private key to protobuf bytes
func (k *Ed25519PrivateKey) Bytes() ([]byte, error) {
	pbmes := new(pb.PrivateKey)
	typ := pb.KeyType_Ed25519
	pbmes.Type = &typ

	buf := make([]byte, 96)
	copy(buf, k.sk[:])
	copy(buf[64:], k.pk[:])
	pbmes.Data = buf
	return proto.Marshal(pbmes)
}

// Equals compares two ed25519 private keys
func (k *Ed25519PrivateKey) Equals(o Key) bool {
	edk, ok := o.(*Ed25519PrivateKey)
	if !ok {
		return false
	}

	return bytes.Equal((*k.sk)[:], (*edk.sk)[:]) && bytes.Equal((*k.pk)[:], (*edk.pk)[:])
}

// GetPublic returns an ed25519 public key from a private key
func (k *Ed25519PrivateKey) GetPublic() PubKey {
	return &Ed25519PublicKey{k.pk}
}

// Sign returns a signature from an input message
func (k *Ed25519PrivateKey) Sign(msg []byte) ([]byte, error) {
	out := ed25519.Sign(k.sk, msg)
	return (*out)[:], nil
}

// ToCurve25519 returns the private key's curve
func (k *Ed25519PrivateKey) ToCurve25519() *[32]byte {
	var sk [32]byte
	extra.PrivateKeyToCurve25519(&sk, k.sk)
	return &sk
}

// Bytes returns a ed25519 public key as protobuf bytes
func (k *Ed25519PublicKey) Bytes() ([]byte, error) {
	pbmes := new(pb.PublicKey)
	typ := pb.KeyType_Ed25519
	pbmes.Type = &typ
	pbmes.Data = (*k.k)[:]
	return proto.Marshal(pbmes)
}

// Equals compares two ed25519 public keys
func (k *Ed25519PublicKey) Equals(o Key) bool {
	edk, ok := o.(*Ed25519PublicKey)
	if !ok {
		return false
	}

	return bytes.Equal((*k.k)[:], (*edk.k)[:])
}

// Verify checks a signature agains the input data
func (k *Ed25519PublicKey) Verify(data []byte, sig []byte) (bool, error) {
	var asig [64]byte
	copy(asig[:], sig)
	return ed25519.Verify(k.k, data, &asig), nil
}

// ToCurve25519 returns the public key's curve
func (k *Ed25519PublicKey) ToCurve25519() (*[32]byte, error) {
	var pk [32]byte
	success := extra.PublicKeyToCurve25519(&pk, k.k)
	if !success {
		return nil, fmt.Errorf("Error converting ed25519 pubkey to curve25519 pubkey")
	}
	return &pk, nil
}

// UnmarshalEd25519PublicKey returns a public key from input bytes
func UnmarshalEd25519PublicKey(data []byte) (PubKey, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("expect ed25519 public key data size to be 32")
	}

	var pub [32]byte
	copy(pub[:], data)

	return &Ed25519PublicKey{
		k: &pub,
	}, nil
}

// UnmarshalEd25519PrivateKey returns a private key from input bytes
func UnmarshalEd25519PrivateKey(data []byte) (PrivKey, error) {
	if len(data) != 96 {
		return nil, fmt.Errorf("expected ed25519 data size to be 96")
	}
	var priv [64]byte
	var pub [32]byte
	copy(priv[:], data)
	copy(pub[:], data[64:])

	return &Ed25519PrivateKey{
		sk: &priv,
		pk: &pub,
	}, nil
}
