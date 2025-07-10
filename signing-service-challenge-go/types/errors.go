package types

import (
	"errors"
)

var (
	ErrUnknownSigningAlgorithm = errors.New("unknown signing algorithm")
	ErrDeviceNotFound          = errors.New("device with given ID does not exist")
	ErrDeviceAlreadyExists     = errors.New("device with given ID already exist")
)
