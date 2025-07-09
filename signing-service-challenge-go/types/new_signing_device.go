package types

type NewSignatureDevice struct {
	Algorithm string `json:"algorithm,omitempty"`
	Label     string `json:"label"`
}
