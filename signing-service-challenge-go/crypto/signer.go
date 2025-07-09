package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
)

var VerificationFailedError = errors.New("signature verification failed")

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	Verify(data []byte, signature []byte) error
}

// NewSigner creates a new Signer based on the provided SignatureDevice.
func NewSigner(algorithm types.SigningAlgorithm, pkPem []byte) (Signer, error) {
	switch algorithm {
	case types.RSA:
		pair, err := unmarshalRSA(pkPem)
		if err != nil {
			return nil, fmt.Errorf("failed to create RSA key pair from PEM: %w", err)
		}
		return &RSASigner{pair: pair}, nil
	case types.ECC:
		pair, err := unmarshalECC(pkPem)
		if err != nil {
			return nil, fmt.Errorf("failed to create ECC key pair from PEM: %w", err)
		}
		return &ECCSigner{pair: pair}, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// GenerateNewPair generates a new key pair based on the specified signing algorithm
// and returns the public and private keys in PEM format.
func GenerateNewPair(algorithm types.SigningAlgorithm) ([]byte, []byte, error) {
	switch algorithm {
	case types.RSA:
		pair, err := generateRSA()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate RSA private key: %w", err)
		}
		return marshalRSA(*pair)
	case types.ECC:
		pair, err := generateECC()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate ECC private key: %w", err)
		}
		return marshalECC(*pair)
	default:
		return nil, nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

type RSASigner struct {
	pair *RSAKeyPair
}

func (r RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	return r.pair.Private.Sign(rand.Reader, hashed[:], crypto.SHA256)
}

func (r RSASigner) Verify(data []byte, signature []byte) error {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(r.pair.Public, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Printf("failed to verify signature: %v\n", err)
		return VerificationFailedError
	}
	return nil
}

type ECCSigner struct {
	pair *ECCKeyPair
}

func (r ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha512.Sum384(dataToBeSigned)
	return r.pair.Private.Sign(rand.Reader, hashed[:], crypto.SHA384)
}

func (r ECCSigner) Verify(data []byte, signature []byte) error {
	hashed := sha512.Sum384(data)
	if !ecdsa.VerifyASN1(r.pair.Public, hashed[:], signature) {
		return VerificationFailedError
	}
	return nil
}
