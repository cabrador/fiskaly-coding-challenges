package api

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
)

type DeviceService interface {
	// Get retrieves a device by its ID.
	Get(id string) (*types.SignatureDevice, error)
	// Create adds a new device to the system.
	Create(device types.NewSignatureDevice) (*types.SignatureDevice, error)
	// SignUsingDevice generates a signature for the given data using the specified device ID.
	// It returns the signature and the signed data.
	SignUsingDevice(deviceID string, data []byte) ([]byte, []byte, error)
	// GetAll retrieves all signature devices.
	GetAll() []*types.SignatureDevice
	// GetDeviceSignatures retrieves all signatures associated with a signature device by its ID.
	GetDeviceSignatures(deviceID string) ([][]byte, error)
}
