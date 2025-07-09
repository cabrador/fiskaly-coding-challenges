package api

type CreateSignatureDeviceRequest struct {
	Algorithm string `json:"algorithm,omitempty"`
	Label     string `json:"label"`
}

type SignTransactionRequest struct {
	DeviceID       string `json:"deviceId,omitempty"`
	DataToBeSigned string `json:"data_to_be_signed,omitempty"`
}
