package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
)

//go:generate mockgen -source database.go -destination database_mock.go -package domain

// Database interface defines the methods required for device persistence.
type Database interface {
	// GetSignatureDevice retrieves a device by its ID.
	GetSignatureDevice(id string) (*types.SignatureDevice, error)
	// GetAllSignatureDevices retrieves all signature devices.
	GetAllSignatureDevices() []*types.SignatureDevice
	// GetDeviceSignatures retrieves all signatures associated with a signature device by its ID.
	GetDeviceSignatures(id string) ([][]byte, error)
	// CreateSignatureDevice adds a new signature device to the database.
	CreateSignatureDevice(device *types.SignatureDevice) error
	// UpdateSignatureDevice updates an existing signature device in the database.
	UpdateSignatureDevice(updatedDevice *types.SignatureDevice) error
}
