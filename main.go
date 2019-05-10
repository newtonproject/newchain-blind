package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/cryptoballot/fdh"
	"github.com/cryptoballot/rsablind"
	"gitlab.newtonproject.org/yangchenzhong/NewChainBlind/blind"
	"gitlab.newtonproject.org/yangchenzhong/NewChainBlind/cli"
	"log"
)

func blindSign(test string, data []byte, key *rsa.PrivateKey, pub *rsa.PublicKey) {
	hashed := fdh.Sum(crypto.SHA256, 256, data)
	fmt.Println(hex.EncodeToString(hashed))

	// customer, requestBlindSig
	blinded, unblinder, err := rsablind.Blind(pub, hashed)
	if err != nil {
		fmt.Println(err)
		return
	}

	// bank sign, writeBlindSig
	sig, err := rsablind.BlindSign(key, blinded)
	if err != nil {
		fmt.Println(err)
		return
	}

	// unblindSig, to
	unblindSig := rsablind.Unblind(pub, sig, unblinder)
	fmt.Println(hex.EncodeToString(unblindSig))

	// Check to make sure both the blinded and unblided data can be verified with the same signature
	if err := rsablind.VerifyBlindSignature(pub, hashed, unblindSig); err != nil {
		fmt.Errorf(test+": Failed to verify for unblinded signature: %v", err)
	}
	if err := rsablind.VerifyBlindSignature(pub, blinded, sig); err != nil {
		fmt.Errorf(test+": Failed to verify for blinded signature: %v", err)
	}

	// Check to make sure blind signing does not work when mismatched
	if err := rsablind.VerifyBlindSignature(pub, data, sig); err == nil {
		fmt.Errorf(test + ": Faulty Verification for mismatched signature 1")
	}
	if err := rsablind.VerifyBlindSignature(pub, blinded, unblindSig); err == nil {
		fmt.Errorf(test + ": Faulty Verification for mismatched signature 2")
	}
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

var (
	ErrKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey    = errors.New("Key is not a valid RSA private key")
	ErrNotRSAPublicKey     = errors.New("Key is not a valid RSA public key")
)

// Parse PEM encoded PKCS1 or PKCS8 private key
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

func main() {
	cli.NewCLI().Execute()
	return

	// openssl genrsa -out key.pem 4096
	// openssl rsa -in key.pem -pubout -out pubkey.pem

	key, err := blind.LoadPEMPrivKeyFile("key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(hex.EncodeToString(key.N.Bytes()))

	pub, err := blind.LoadPEMPublicKeyFile("pubkey.pem")
	if err != nil {
		fmt.Println(err)
		return
	}

	c := 2048
	data := make([]byte, c)
	_, err = rand.Read(data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(key.Public())
	fmt.Println(hex.EncodeToString(key.N.Bytes()))

	// Do it twice to make sure we are also testing using cached precomputed values.
	blindSign("TestBlindSignBig", data, key, pub)
	blindSign("TestBlindSignBig", data, key, pub)
}
