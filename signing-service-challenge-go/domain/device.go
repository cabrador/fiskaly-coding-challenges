package domain

import (
	"encoding/base64"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

// NewDeviceService creates a new DeviceService instance with the provided database.
func NewDeviceService(db Database) *DeviceService {
	return &DeviceService{
		db:      db,
		builder: strings.Builder{},
	}
}

type DeviceService struct {
	db      Database
	builder strings.Builder
}

// Get retrieves a device by its ID from the database.
func (d *DeviceService) Get(id string) (*types.SignatureDevice, error) {
	return d.db.GetSignatureDevice(id)
}

// Create adds a new device to the database.
func (d *DeviceService) Create(device types.NewSignatureDevice) (*types.SignatureDevice, error) {
	if !types.IsAllowedSigningAlgorithm(types.SigningAlgorithm(device.Algorithm)) {
		return nil, fmt.Errorf("%w: %s", types.ErrUnknownSigningAlgorithm, device.Algorithm)
	}

	// The probability of hitting an existing UUID is close to zero
	// nevertheless it should still be handled in real scenario.
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate device id: %v", err)
	}

	_, privatePem, err := crypto.GenerateNewPair(types.SigningAlgorithm(device.Algorithm))
	if err != nil {
		return nil, fmt.Errorf("failed to generate signer: %v", err)
	}
	newDevice := &types.SignatureDevice{
		ID:                 id.String(),
		Algorithm:          types.SigningAlgorithm(device.Algorithm),
		Label:              device.Label,
		Counter:            0,
		PkPem:              privatePem,
		PreviousSignatures: make(map[uint32][]byte),
	}

	if err = d.db.CreateSignatureDevice(newDevice); err != nil {
		return nil, fmt.Errorf("failed to save device into the db: %w", err)
	}

	return newDevice, nil
}

func (d *DeviceService) SignUsingDevice(deviceID string, data []byte) ([]byte, []byte, error) {
	signingDevice, err := d.Get(deviceID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get device: %w", err)
	}

	// RESET the builder to avoid appending to next call!!!
	defer d.builder.Reset()
	d.builder.WriteString(strconv.Itoa(int(signingDevice.Counter)))
	d.builder.WriteString("_")
	d.builder.Write(data)
	d.builder.WriteString("_")
	// First sign with this device?
	var prev []byte
	if signingDevice.Counter == 0 {
		prev = []byte(signingDevice.ID)
	} else {
		prev = signingDevice.PreviousSignatures[signingDevice.Counter-1]
	}
	d.builder.WriteString(base64.StdEncoding.EncodeToString(prev))
	toBeSigned := []byte(d.builder.String())
	signer, err := crypto.NewSigner(signingDevice.Algorithm, signingDevice.PkPem)
	if err != nil {
		return nil, nil, err
	}
	signature, err := signer.Sign(toBeSigned)
	if err != nil {
		return nil, nil, err
	}

	// Update the device with the new signature and increment the counter
	signingDevice.PreviousSignatures[signingDevice.Counter] = signature
	signingDevice.Counter++
	// Update the device in the database
	if err = d.db.UpdateSignatureDevice(signingDevice); err != nil {
		return nil, nil, err
	}

	return signature, toBeSigned, nil
}

func (d *DeviceService) GetAll() []*types.SignatureDevice {
	return d.db.GetAllSignatureDevices()
}

func (d *DeviceService) GetDeviceSignatures(deviceID string) ([][]byte, error) {
	return d.db.GetDeviceSignatures(deviceID)
}
