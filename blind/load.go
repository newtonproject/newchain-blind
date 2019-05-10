package blind

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	ErrNotECPublicKey  = errors.New("Key is not a valid ECDSA public key")
	ErrNotECPrivateKey = errors.New("Key is not a valid ECDSA private key")
)

// LoadPEMPrivKeyFile reads a PEM encoded RSA private key from the file name.
// A new RSA private key will be returned if no error.
func LoadPEMPrivKeyFile(name string) (*rsa.PrivateKey, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return LoadPEMPrivKey(file)
}

// LoadPEMPrivKey reads a PEM encoded RSA private key from the io.Reader.
// A new RSA private key will be returned if no error.
func LoadPEMPrivKey(reader io.Reader) (*rsa.PrivateKey, error) {
	block, err := loadPem(reader)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// LoadPEMPublicKeyFile reads a PEM encoded RSA publicKey key from the file name.
// A new RSA public key will be returned if no error.
func LoadPEMPublicKeyFile(name string) (*rsa.PublicKey, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return LoadPEMPublicKey(file)
}

// LoadPEMPublicKey reads a PEM encoded RSA public key from the io.Reader.
// A new RSA public key will be returned if no error.
func LoadPEMPublicKey(reader io.Reader) (*rsa.PublicKey, error) {
	block, err := loadPem(reader)
	if err != nil {
		return nil, err
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotECPublicKey
	}

	return pkey, nil
}

// LoadEncryptedPEMPrivKey decrypts the PEM encoded private key using the
// password provided returning a RSA private key. If the PEM data is invalid,
// or unable to decrypt an error will be returned.
func LoadEncryptedPEMPrivKey(reader io.Reader, password []byte) (*rsa.PrivateKey, error) {
	block, err := loadPem(reader)
	if err != nil {
		return nil, err
	}

	decryptedBlock, err := x509.DecryptPEMBlock(block, password)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(decryptedBlock)
}

func loadPem(reader io.Reader) (*pem.Block, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	if block == nil {
		// pem.Decode will set block to nil if there is no PEM data in the input
		// the second parameter will contain the provided bytes that failed
		// to be decoded.
		return nil, fmt.Errorf("no valid PEM data provided")
	}

	return block, nil
}
