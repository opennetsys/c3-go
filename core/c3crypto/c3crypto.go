package c3crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"math/big"

	"github.com/obscuren/ecies"
)

var (
	// ErrNilPrivateKey ...
	ErrNilPrivateKey = errors.New("private key is nil")
	// ErrNilPublicKey ...
	ErrNilPublicKey = errors.New("public key is nil")
	// ErrNilData ...
	ErrNilData = errors.New("data is nil")
	// ErrNilSigParams ...
	ErrNilSigParams = errors.New("r and s params cannot be nil")
	// ErrGeneratingECIESPublicKey ...
	ErrGeneratingECIESPublicKey = errors.New("could not generate ecies public key from ecda public key")
	// ErrGeneratingECIESPrivateKey ...
	ErrGeneratingECIESPrivateKey = errors.New("could not generate ecies private key from ecdsa private key")
	// ErrGeneratingECDSAPublicKey ...
	ErrGeneratingECDSAPublicKey = errors.New("could not generate ecdsa public key from ecdsa private key")
	// ErrNotECDSAPubKey ...
	ErrNotECDSAPubKey = errors.New("not an ecdsa public key")
	// Curve ...
	Curve = elliptic.P256()
	//Curve = elliptic.P384() // note: causes a shared key params are too big err
	//Curve = elliptic.P521() // note: causes an out of range panic at ecies.go L106!!
)

// NewKeyPair ...
func NewKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	priv, err := NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	pub, err := GetPublicKey(priv)
	return priv, pub, err
}

// NewPrivateKey ...
func NewPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(Curve, rand.Reader)
}

// GetPublicKey ...
func GetPublicKey(priv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	if priv == nil {
		return nil, ErrNilPrivateKey
	}

	return &priv.PublicKey, nil
}

// Sign ...
func Sign(priv *ecdsa.PrivateKey, data []byte) (*big.Int, *big.Int, error) {
	if data == nil {
		return nil, nil, ErrNilData
	}
	if priv == nil {
		return nil, nil, ErrNilPrivateKey
	}

	return ecdsa.Sign(rand.Reader, priv, data)
}

// Verify ...
func Verify(pub *ecdsa.PublicKey, data []byte, r, s *big.Int) (bool, error) {
	if pub == nil {
		return false, ErrNilPublicKey
	}
	if data == nil {
		return false, ErrNilData
	}
	if r == nil || s == nil {
		return false, ErrNilSigParams
	}

	return ecdsa.Verify(pub, data, r, s), nil
}

// Encrypt ...
func Encrypt(pub *ecdsa.PublicKey, data []byte) ([]byte, error) {
	if pub == nil {
		return nil, ErrNilPublicKey
	}
	if data == nil {
		return nil, ErrNilData
	}

	eciesPubKey := ecies.ImportECDSAPublic(pub)
	if eciesPubKey == nil {
		return nil, ErrGeneratingECIESPublicKey
	}

	return ecies.Encrypt(rand.Reader, eciesPubKey, data, nil, nil)

}

// Decrypt ...
func Decrypt(priv *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	if priv == nil {
		return nil, ErrNilPrivateKey
	}
	if data == nil {
		return nil, ErrNilData
	}

	eciesPrivKey := ecies.ImportECDSA(priv)
	if eciesPrivKey == nil {
		return nil, ErrGeneratingECIESPrivateKey
	}

	return eciesPrivKey.Decrypt(rand.Reader, data, nil, nil)
}

// SerializePrivateKey ...
func SerializePrivateKey(priv *ecdsa.PrivateKey) ([]byte, error) {
	if priv == nil {
		return nil, ErrNilPrivateKey
	}

	return x509.MarshalECPrivateKey(priv)
}

// SerializePublicKey ...
func SerializePublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	if pub == nil {
		return nil, ErrNilPublicKey
	}

	return x509.MarshalPKIXPublicKey(pub)
}

// DeserializePrivateKey ...
func DeserializePrivateKey(privBytes []byte) (*ecdsa.PrivateKey, error) {
	return x509.ParseECPrivateKey(privBytes)
}

// DeserializePublicKey ...
func DeserializePublicKey(pubBytes []byte) (*ecdsa.PublicKey, error) {
	pubIfc, err := x509.ParsePKIXPublicKey(pubBytes)
	if err != nil {
		return nil, err
	}

	pub, ok := pubIfc.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrNotECDSAPubKey
	}

	return pub, nil
}
