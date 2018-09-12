package c3crypto

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/c3systems/c3-go/common/hexutil"
	"github.com/obscuren/ecies"
)

const (
	// EncryptionCipher ...
	EncryptionCipher = x509.PEMCipherAES256
	// PrivateKeyPEMType ...
	PrivateKeyPEMType = "ECDSA PRIVATE KEY"
	// PublicKeyPEMType ...
	PublicKeyPEMType = "ECDSA PUBLIC KEY"
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
	// ErrInvalidPublicKey ...
	ErrInvalidPublicKey = errors.New("invalid public key")
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
	//Curve = secp256k1.S256() // note: what ethereum uses
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

// EncodePrivateKeyToPem ...
// note password is optional
func EncodePrivateKeyToPem(priv *ecdsa.PrivateKey, password *string) (*pem.Block, error) {
	bytes, err := SerializePrivateKey(priv)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  PrivateKeyPEMType,
		Bytes: bytes,
	}

	if password != nil {
		return x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(*password), EncryptionCipher)
	}

	return block, nil
}

// EncodePublicKeyToPem ...
// note: password is optional
func EncodePublicKeyToPem(pub *ecdsa.PublicKey, password *string) (*pem.Block, error) {
	bytes, err := SerializePublicKey(pub)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  PublicKeyPEMType,
		Bytes: bytes,
	}

	if password != nil {
		return x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(*password), EncryptionCipher)
	}

	return block, nil
}

// DecodePemToPrivateKey ...
func DecodePemToPrivateKey(block *pem.Block, password *string) (*ecdsa.PrivateKey, error) {
	if password == nil && x509.IsEncryptedPEMBlock(block) {
		return nil, errors.New("pem block is encrypted but no password provided")
	}
	if password != nil && !x509.IsEncryptedPEMBlock(block) {
		return nil, errors.New("a password was provided but pem block is not encrypted")
	}

	if password != nil {
		bytes, err := x509.DecryptPEMBlock(block, []byte(*password))
		if err != nil {
			return nil, err
		}

		return DeserializePrivateKey(bytes)
	}

	return DeserializePrivateKey(block.Bytes)
}

// DecodePemToPublicKey ...
func DecodePemToPublicKey(block *pem.Block, password *string) (*ecdsa.PublicKey, error) {
	if password == nil && x509.IsEncryptedPEMBlock(block) {
		return nil, errors.New("pem block is encrypted but no password provided")
	}
	if password != nil && !x509.IsEncryptedPEMBlock(block) {
		return nil, errors.New("a password was provided but pem block is not encrypted")
	}

	if password != nil {
		bytes, err := x509.DecryptPEMBlock(block, []byte(*password))
		if err != nil {
			return nil, err
		}

		return DeserializePublicKey(bytes)
	}

	return DeserializePublicKey(block.Bytes)
}

// WritePrivateKeyToPemFile ...
func WritePrivateKeyToPemFile(priv *ecdsa.PrivateKey, password *string, fileName string) error {
	block, err := EncodePrivateKeyToPem(priv, password)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("[c3crypto] err creating file %s", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err := pem.Encode(w, block); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return f.Sync()
}

// WritePublicKeyToPemFile ...
func WritePublicKeyToPemFile(pub *ecdsa.PublicKey, password *string, fileName string) error {
	block, err := EncodePublicKeyToPem(pub, password)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err := pem.Encode(w, block); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return f.Sync()
}

// ReadPrivateKeyFromPem ...
func ReadPrivateKeyFromPem(fileName string, password *string) (*ecdsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)

	return DecodePemToPrivateKey(block, password)
}

// ReadPublicKeyFromPem ...
func ReadPublicKeyFromPem(fileName string, password *string) (*ecdsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)

	return DecodePemToPublicKey(block, password)
}

// WriteKeyPairToPem ...
func WriteKeyPairToPem(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, password *string, fileName string) error {
	privBlock, err := EncodePrivateKeyToPem(priv, password)
	if err != nil {
		return err
	}

	pubBlock, err := EncodePublicKeyToPem(pub, password)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err := pem.Encode(w, privBlock); err != nil {
		return err
	}
	if err := pem.Encode(w, pubBlock); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return f.Sync()
}

// ReadKeyPairFromPem ...
func ReadKeyPairFromPem(fileName string, password *string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	privBlock, rest := pem.Decode(bytes)
	pubBlock, _ := pem.Decode(rest)

	priv, err := DecodePemToPrivateKey(privBlock, password)
	if err != nil {
		return nil, nil, err
	}
	pub, err := DecodePemToPublicKey(pubBlock, password)
	if err != nil {
		return nil, nil, err
	}

	return priv, pub, nil
}

// EncodeAddress ...
func EncodeAddress(pub *ecdsa.PublicKey) (string, error) {
	bytes, err := SerializePublicKey(pub)
	if err != nil {
		return "", err
	}

	return hexutil.EncodeToString(bytes), nil
}

// DecodeAddress [for now] decodes public address hex to ECDSA public key.
// Public key are treated as public address at the moment.
func DecodeAddress(address string) (*ecdsa.PublicKey, error) {
	byteStr, err := hexutil.RemovePrefix(address)
	if err != nil {
		return nil, err
	}

	bytes, err := hexutil.DecodeBytes([]byte(byteStr))
	if err != nil {
		return nil, err
	}

	pub, err := DeserializePublicKey(bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

// PublicKeyToBytes ...
func PublicKeyToBytes(pub *ecdsa.PublicKey) ([]byte, error) {
	if pub == nil {
		return nil, ErrInvalidPublicKey
	}

	return elliptic.Marshal(Curve, pub.X, pub.Y), nil
}

// PublicKeyFromBytes ...
func PublicKeyFromBytes(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(Curve, pub)
	if x == nil {
		return nil, ErrInvalidPublicKey
	}

	return &ecdsa.PublicKey{Curve: Curve, X: x, Y: y}, nil
}
