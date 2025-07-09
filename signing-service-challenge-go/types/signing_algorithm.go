package types

type SigningAlgorithm string

const (
	ECC SigningAlgorithm = "ECC"
	RSA SigningAlgorithm = "RSA"
)

func IsAllowedSigningAlgorithm(algorithm SigningAlgorithm) bool {
	switch algorithm {
	case ECC, RSA:
		return true
	default:
		return false
	}
}
