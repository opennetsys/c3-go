// +build unit

package c3crypto

import (
	"crypto/ecdsa"
	"os"
	"reflect"
	"testing"
)

func TestNewKeyPair(t *testing.T) {
	priv, pub, err := NewKeyPair()
	if err != nil {
		t.Errorf("received non nil err\n%v", err)
	}
	if priv == nil {
		t.Error("received nil private key")
	}
	if pub == nil {
		t.Error("received nil public key")
	}
}

func TestNewPrivateKey(t *testing.T) {
	priv, err := NewPrivateKey()
	if err != nil {
		t.Errorf("received non nil err\n%v", err)
	}
	if priv == nil {
		t.Error("received nil private key")
	}
}

func TestGetPublicKey(t *testing.T) {
	priv, pub, err := NewKeyPair()
	if err != nil {
		t.Errorf("received non nil err generating key pair\n%v", err)
	}
	if err != nil {
		t.Errorf("received non nil err\n%v", err)
	}
	if priv == nil {
		t.Errorf("received nil private key")
	}
	if pub == nil {
		t.Errorf("received nil public key")
	}

	expectedIfc := priv.Public()
	expectedPub, ok := expectedIfc.(*ecdsa.PublicKey)
	if !ok {
		t.Errorf("could not generate expected public key")
	}

	if !reflect.DeepEqual(*expectedPub, *pub) {
		t.Errorf("expected %v, received %v", *expectedPub, *pub)
	}
}

func TestSignAndVerify(t *testing.T) {
	inputs := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("foo bar"),
	}

	priv, pub, err := NewKeyPair()
	if err != nil {
		t.Errorf("received non nil err generating key pair\n%v", err)
	}

	if priv == nil || pub == nil {
		t.Errorf("received nil priv or pub key\n%v\n%v", priv, pub)
	}

	// note: is there a better way to test signing and verifying?
	for idx, in := range inputs {
		r, s, err := Sign(priv, in)
		if err != nil {
			t.Errorf("test %d got non nil err signing\n%v", idx+1, err)
		}

		if r == nil || s == nil {
			t.Errorf("test %d got nil r or s\n%v\n%v", idx+1, r, s)
		}

		ver, err := Verify(pub, in, r, s)
		if err != nil {
			t.Errorf("test %d received non nil err verifying\n%v", idx+1, err)
		}

		if !ver {
			t.Errorf("test %d didn't verify", idx+1)
		}
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	inputs := [][]byte{
		[]byte("foo"),
		[]byte("bar"),
		[]byte("foo bar"),
	}

	priv, pub, err := NewKeyPair()
	if err != nil {
		t.Errorf("received non nil err generating key pair\n%v", err)
	}

	if priv == nil || pub == nil {
		t.Errorf("received nil priv or pub key\n%v\n%v", priv, pub)
	}

	// note: is there a better way to test encrypting and decrypting?
	for idx, in := range inputs {
		e, err := Encrypt(pub, in)
		if err != nil {
			t.Errorf("test %d got non nil err encrypting\n%v", idx+1, err)
		}
		if e == nil {
			t.Errorf("test %d got nil encrypted bytes\n%v", idx+1, e)
		}

		d, err := Decrypt(priv, e)
		if err != nil {
			t.Errorf("test %d received non nil err decrypting\n%v", idx+1, err)
		}
		if d == nil {
			t.Errorf("test %d received nil decrypted text", idx+1)
		}

		if string(d) != string(in) {
			t.Errorf("test %d expected %s received %s", idx+1, string(in), string(d))
		}
	}
}

func TestSerializeAndDeserializePrivateKey(t *testing.T) {
	priv, err := NewPrivateKey()
	if err != nil {
		t.Errorf("received non nil err\n%v", err)
	}
	if priv == nil {
		t.Error("priv is nil")
	}

	privBytes, err := SerializePrivateKey(priv)
	if err != nil {
		t.Errorf("received non nil err\n%v", err)
	}

	priv2, err := DeserializePrivateKey(privBytes)
	if err != nil {
		t.Errorf("received non nil err\n%v", err)
	}
	if priv2 == nil {
		t.Error("priv2 is nil")
	}

	if !reflect.DeepEqual(*priv, *priv2) {
		t.Errorf("expected %v\nreceived %v", *priv, *priv2)
	}
}

func TestSerializeAndDeserializePublicKey(t *testing.T) {
	_, pub, err := NewKeyPair()
	if err != nil {
		t.Errorf("received non nil err generating key pair\n%v", err)
	}
	if pub == nil {
		t.Error("pub is nil")
	}

	pubBytes, err := SerializePublicKey(pub)
	if err != nil {
		t.Errorf("received non nil err serializing\n%v", err)
	}

	pub2, err := DeserializePublicKey(pubBytes)
	if err != nil {
		t.Errorf("received non nil err deserializing\n%v", err)
	}
	if pub2 == nil {
		t.Error("priv2 is nil")
	}

	if !reflect.DeepEqual(*pub, *pub2) {
		t.Errorf("expected %v\nreceived %v", *pub, *pub2)
	}
}

func TestReadWritePrivateKeyToPEM(t *testing.T) {
	priv, err := NewPrivateKey()
	if err != nil {
		t.Error(err)
	}

	password := "password"
	fileName := "priv.test.pem"
	if err := WritePrivateKeyToPemFile(priv, &password, fileName); err != nil {
		t.Error(err)
	}
	defer os.Remove(fileName)

	priv2, err := ReadPrivateKeyFromPem(fileName, &password)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(priv, priv2) {
		t.Errorf("expected: %v\nreceived: %v", priv, priv2)
	}
}

func TestReadWritePublicKeyToPEM(t *testing.T) {
	_, pub, err := NewKeyPair()
	if err != nil {
		t.Error(err)
	}

	password := "password"
	fileName := "pub.test.pem"
	if err := WritePublicKeyToPemFile(pub, &password, fileName); err != nil {
		t.Error(err)
	}
	defer os.Remove(fileName)

	pub2, err := ReadPublicKeyFromPem(fileName, &password)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(pub, pub2) {
		t.Errorf("expected: %v\nreceived: %v", pub, pub2)
	}
}

func TestReadWritePairsToPEM(t *testing.T) {
	priv, pub, err := NewKeyPair()
	if err != nil {
		t.Error(err)
	}

	password := "password"
	fileName := "keys.pem"
	if err := WriteKeyPairToPem(priv, pub, &password, fileName); err != nil {
		t.Error(err)
	}
	defer os.Remove(fileName)

	priv2, pub2, err := ReadKeyPairFromPem(fileName, &password)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(pub, pub2) || !reflect.DeepEqual(priv, priv2) {
		t.Errorf("expected priv: %v\nexpected pub: %v\nreceived priv: %v\nreceived pub: %v", priv, pub, priv2, pub2)
	}
}

func TestDecodeAddress(t *testing.T) {
	_, pub, err := NewKeyPair()
	if err != nil {
		t.Error(err)
	}

	addr, err := EncodeAddress(pub)
	if err != nil {
		t.Error(err)
	}

	newPub, err := DecodeAddress(addr)

	if !reflect.DeepEqual(pub, newPub) {
		t.Errorf("public keys do not match; received %v", newPub)
	}
}

func TestEncodeAndDecodeAddress(t *testing.T) {
	_, pub, err := NewKeyPair()
	if err != nil {
		t.Error(err)
	}
	if pub == nil {
		t.Error("nil pub key")
	}

	addr, err := EncodeAddress(pub)
	if err != nil {
		t.Error(err)
	}
	t.Logf("addr %s", addr)

	pub1, err := DecodeAddress(addr)
	if err != nil {
		t.Error(err)
	}
	if pub1 == nil {
		t.Error("pub1 is nil")
	}

	if !reflect.DeepEqual(pub, pub1) {
		t.Errorf("expected %v\nreceived %v", *pub, *pub1)
	}
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}

	return nil
}
