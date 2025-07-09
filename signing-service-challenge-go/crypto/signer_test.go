package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
	"testing"
)

func Test_Factories(t *testing.T) {
	t.Run("Unsupported Algorithm", func(t *testing.T) {
		_, _, err := GenerateNewPair("unsupported")
		if err == nil {
			t.Fatal("expected error for unsupported algorithm, got nil")
		}
		signer, err := NewSigner("unsupported", []byte{})
		if err == nil {
			t.Fatal("expected error for unsupported algorithm, got nil")
		}
		if signer != nil {
			t.Fatalf("expected nil signer for unsupported algorithm, got %T", signer)
		}
	})
	t.Run("RSA Signer", func(t *testing.T) {
		_, pkPem, err := GenerateNewPair(types.RSA)
		if err != nil {
			t.Fatalf("failed to create RSA signer: %v", err)
		}
		signer, err := NewSigner(types.RSA, pkPem)
		if _, ok := signer.(*RSASigner); !ok {
			t.Fatalf("expected RSASigner type, got different %T", signer)
		}
	})
	t.Run("ECC Signer", func(t *testing.T) {
		_, pkPem, err := GenerateNewPair(types.ECC)
		if err != nil {
			t.Fatalf("failed to create ECC signer: %v", err)
		}
		signer, err := NewSigner(types.ECC, pkPem)
		if _, ok := signer.(*ECCSigner); !ok {
			t.Fatalf("expected ECCSigner type, got different %T", signer)
		}
	})
}

func TestRSASigner_Sign(t *testing.T) {
	pair, err := generateRSA()
	if err != nil {
		t.Fatalf("failed to create generate RSA pair: %v", err)
	}
	signer := RSASigner{pair: pair}
	data := []byte("test data to be signed")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("failed to sign data: %v", err)
	}
	if len(signature) == 0 {
		t.Fatal("signature should not be empty")
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(signer.pair.Public, crypto.SHA256, hashed[:], signature)
	if err != nil {
		t.Errorf("signature verification failed: %v", err)
	}
}

func TestECCSigner_Sign(t *testing.T) {
	pair, err := generateECC()
	if err != nil {
		t.Fatalf("failed to create generate RSA pair: %v", err)
	}
	signer := ECCSigner{pair: pair}
	data := []byte("test data to be signed")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("failed to sign data: %v", err)
	}
	if len(signature) == 0 {
		t.Fatal("signature should not be empty")
	}
	hashed := sha512.Sum384(data)
	if !ecdsa.VerifyASN1(signer.pair.Public, hashed[:], signature) {
		t.Error("signature verification failed")
	}
}

func TestSigner_Verify(t *testing.T) {
	tests := []struct {
		name   string
		alg    types.SigningAlgorithm
		passes bool
	}{
		{
			name:   "RSA Success",
			alg:    types.RSA,
			passes: true,
		},
		{
			name:   "ECC Success",
			alg:    types.ECC,
			passes: true,
		},
		{
			name:   "RSA Fail",
			alg:    types.RSA,
			passes: false,
		},
		{
			name:   "ECC Fail",
			alg:    types.ECC,
			passes: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, pkPemm, err := GenerateNewPair(test.alg)
			if err != nil {
				t.Fatalf("failed to create signer: %v", err)
			}
			signer, err := NewSigner(test.alg, pkPemm)
			signature, err := signer.Sign([]byte("success"))
			if err != nil {
				t.Fatalf("failed to sign data: %v", err)
			}
			// botch the signature if the test is supposed to fail
			if !test.passes {
				signature[0] ^= 0xFF // flip the first byte to invalidate the signature
			}
			err = signer.Verify([]byte("success"), signature)
			if test.passes && err != nil {
				t.Errorf("expected verification to pass, but got error: %v", err)
			}
			if !test.passes && err == nil {
				t.Error("expected verification to fail, but it passed")
			}

		})
	}
}
