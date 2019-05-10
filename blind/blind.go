package blind

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/cryptoballot/fdh"
	"github.com/cryptoballot/rsablind"
)

func BlindCash(data []byte, pub *rsa.PublicKey) (blindedData []byte, unblinder []byte, err error) {
	hashed := fdh.Sum(crypto.SHA256, 256, data)
	return rsablind.Blind(pub, hashed)
}

func blindSign(test string, data []byte, key *rsa.PrivateKey, pub *rsa.PublicKey) {
	hashed := fdh.Sum(crypto.SHA256, 256, data)

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
