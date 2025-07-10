package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
	"sync"
)

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		lock: sync.Mutex{},
		db:   make(map[string]*types.SignatureDevice),
	}
}

// InMemoryDatabase is a simple in-memory thread-safe implementation of the Database interface.
type InMemoryDatabase struct {
	lock sync.Mutex
	db   map[string]*types.SignatureDevice
}

func (d *InMemoryDatabase) GetSignatureDevice(id string) (*types.SignatureDevice, error) {
	d.lock.Lock()
	device, exists := d.db[id]
	d.lock.Unlock()
	if exists {
		return device, nil
	}
	return nil, types.ErrDeviceNotFound
}

func (d *InMemoryDatabase) CreateSignatureDevice(device *types.SignatureDevice) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	// We must not overwrite an existing device.
	_, exists := d.db[device.ID]
	if exists {
		return types.ErrDeviceAlreadyExists
	}
	d.db[device.ID] = device
	return nil
}

func (d *InMemoryDatabase) UpdateSignatureDevice(updatedDevice *types.SignatureDevice) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if _, exist := d.db[updatedDevice.ID]; !exist {
		return types.ErrDeviceNotFound
	}
	d.db[updatedDevice.ID] = updatedDevice
	return nil
}

func (d *InMemoryDatabase) GetAllSignatureDevices() []*types.SignatureDevice {
	d.lock.Lock()
	defer d.lock.Unlock()
	if len(d.db) == 0 {
		return []*types.SignatureDevice{}
	}
	devices := make([]*types.SignatureDevice, len(d.db))
	idx := 0
	for _, device := range d.db {
		devices[idx] = device
		idx++
	}

	return devices
}
func (d *InMemoryDatabase) GetDeviceSignatures(id string) ([][]byte, error) {
	d.lock.Lock()
	device, exists := d.db[id]
	d.lock.Unlock()
	if !exists {
		return nil, types.ErrDeviceNotFound
	}
	if len(device.PreviousSignatures) == 0 {
		return [][]byte{}, nil
	}

	signatures := make([][]byte, len(device.PreviousSignatures))
	idx := 0
	for _, signature := range device.PreviousSignatures {
		signatures[idx] = signature
		idx++
	}
	return signatures, nil
}
