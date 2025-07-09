package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// RSAKeyPair is a DTO that holds RSA private and public keys.
type RSAKeyPair struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

// marshalRSA takes an RSAKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func marshalRSA(keyPair RSAKeyPair) ([]byte, []byte, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(keyPair.Private)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(keyPair.Public)

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodePublic := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodePublic, encodedPrivate, nil
}

// unmarshalRSA takes an encoded RSA private key and transforms it into a rsa.PrivateKey.
func unmarshalRSA(privateKeyBytes []byte) (*RSAKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
