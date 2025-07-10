package domain

import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
	"go.uber.org/mock/gomock"
	"strings"
	"testing"
)

func Test_DeviceService_CreateNewSignatureDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testErr := errors.New("error")
	tests := []struct {
		name          string
		algorithm     string
		expectedError error
		setup         func(*MockDatabase)
	}{
		{
			name:          "Unknown Algorithm",
			algorithm:     "Unknown",
			expectedError: types.ErrUnknownSigningAlgorithm,
			setup: func(*MockDatabase) {
				// this test case returns early - no setup
			},
		},
		{
			name:          "Db Error",
			algorithm:     "ECC",
			expectedError: testErr,
			setup: func(db *MockDatabase) {
				db.EXPECT().CreateSignatureDevice(gomock.Any()).Return(testErr)
			},
		},
		{
			name:      "Successful Creation",
			algorithm: "RSA",
			setup: func(db *MockDatabase) {
				db.EXPECT().CreateSignatureDevice(gomock.Any()).Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := NewMockDatabase(ctrl)
			deviceService := NewDeviceService(db)
			test.setup(db)
			_, err := deviceService.Create(types.NewSignatureDevice{
				Algorithm: test.algorithm,
				Label:     "Label",
			})
			if test.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error %q, got nil", test.expectedError)
				}
				if !errors.Is(err, test.expectedError) {
					t.Fatalf("expected error %q, got %q", test.expectedError, err.Error())
				}

			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func Test_DeviceService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name             string
		deviceID         string
		expectedErrorMsg string
		setup            func(*MockDatabase)
	}{
		{
			name:             "Device Not Found",
			deviceID:         "nonexistent",
			expectedErrorMsg: "device not found",
			setup: func(db *MockDatabase) {
				db.EXPECT().GetSignatureDevice("nonexistent").Return(nil, errors.New("device not found"))
			},
		},
		{
			name:     "Successful Retrieval",
			deviceID: "valid-id",
			setup: func(db *MockDatabase) {
				db.EXPECT().GetSignatureDevice("valid-id").Return(&types.SignatureDevice{ID: "valid-id"}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := NewMockDatabase(ctrl)
			deviceService := NewDeviceService(db)
			test.setup(db)
			device, err := deviceService.Get(test.deviceID)
			if test.expectedErrorMsg != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", test.expectedErrorMsg)
				}
				if !strings.Contains(err.Error(), test.expectedErrorMsg) {
					t.Fatalf("expected error %q, got %q", test.expectedErrorMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if device.ID != test.deviceID {
				t.Fatalf("expected device ID %q, got %q", test.deviceID, device.ID)
			}
		})
	}
}

func Test_DeviceService_SignUsingDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// generate a valid pair for testing
	_, privatePem, err := crypto.GenerateNewPair(types.ECC)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		name           string
		device         types.SignatureDevice
		wantSignedData string
	}{
		{
			name: "Zero Counter",
			device: types.SignatureDevice{
				ID:                 "valid-id",
				Algorithm:          types.ECC,
				PkPem:              privatePem,
				PreviousSignatures: make(map[uint32][]byte),
				Counter:            0,
			}, wantSignedData: "0_test data_dmFsaWQtaWQ=",
		},
		{
			name: "Non-Zero Counter",
			device: types.SignatureDevice{
				ID:        "valid-id",
				Algorithm: types.ECC,
				PkPem:     privatePem,
				PreviousSignatures: map[uint32][]byte{
					0: []byte("previous-signature"),
				},
				Counter: 1,
			}, wantSignedData: "1_test data_cHJldmlvdXMtc2lnbmF0dXJl",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := NewMockDatabase(ctrl)
			db.EXPECT().GetSignatureDevice("valid-id").Return(&test.device, nil)
			db.EXPECT().UpdateSignatureDevice(gomock.Any()).Return(nil)
			deviceService := NewDeviceService(db)

			signature, signedData, err := deviceService.SignUsingDevice("valid-id", []byte("test data"))
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if signature == nil {
				t.Fatal("expected non-nil signature and signed data")
			}
			if signedData == nil {
				t.Fatalf("expected non-nil signed data")
			}
			if !strings.EqualFold(string(signedData), test.wantSignedData) {
				t.Fatalf("expected signed data to match, got %s", string(signedData))
			}
		})
	}

}
