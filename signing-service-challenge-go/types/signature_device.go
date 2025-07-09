package types

// SignatureDevice represents a device that can sign data using a specific signing algorithm.
type SignatureDevice struct {
	ID                 string
	Algorithm          SigningAlgorithm
	Label              string
	Counter            uint32
	PkPem              []byte
	PreviousSignatures map[uint32][]byte // counter -> signature mapping
}
