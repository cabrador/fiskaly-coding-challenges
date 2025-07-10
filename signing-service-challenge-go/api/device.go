package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/types"
	"io"
	"net/http"
)

// TODO: REST endpoints ...

func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			fmt.Sprintf("Failed to read request body: %s", err.Error()),
		})
		return
	}
	unmarshalled := CreateSignatureDeviceRequest{}
	if err = json.Unmarshal(body, &unmarshalled); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			fmt.Sprintf("Incorrect request format: %s", err.Error()),
		})
		return
	}
	device, err := s.deviceService.Create(types.NewSignatureDevice{
		Algorithm: unmarshalled.Algorithm,
		Label:     unmarshalled.Label,
	})
	if err != nil {
		if errors.Is(err, types.ErrUnknownSigningAlgorithm) {
			WriteErrorResponse(response, http.StatusBadRequest, []string{
				err.Error(),
			})
		} else if errors.Is(err, types.ErrDeviceAlreadyExists) {
			WriteErrorResponse(response, http.StatusConflict, []string{
				err.Error(),
			})
		} else {
			WriteInternalError(response, request.URL.Path, err)
		}
		return
	}
	WriteAPIResponse(response, http.StatusCreated, device)
}

type SignTransactionResponse struct {
	Signature  []byte `json:"signature"`
	SignedData []byte `json:"signed_data"`
}

func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			fmt.Sprintf("Failed to read request body: %s", err.Error()),
		})
		return
	}
	unmarshalled := SignTransactionRequest{}
	if err = json.Unmarshal(body, &unmarshalled); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			fmt.Sprintf("Incorrect request format: %s", err.Error()),
		})
		return
	}
	signature, signedData, err := s.deviceService.SignUsingDevice(unmarshalled.DeviceID, []byte(unmarshalled.DataToBeSigned))
	if err != nil {
		if errors.Is(err, types.ErrDeviceNotFound) {
			WriteErrorResponse(response, http.StatusNotFound, []string{
				err.Error(),
			})

		} else {
			WriteInternalError(response, request.URL.Path, err)
		}
		return
	}
	WriteAPIResponse(response, http.StatusCreated, SignTransactionResponse{
		Signature:  signature,
		SignedData: signedData,
	})
}

func (s *Server) Devices(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	all := s.deviceService.GetAll()
	WriteAPIResponse(response, http.StatusOK, all)
}

func (s *Server) DeviceSignatures(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	signatures, err := s.deviceService.GetDeviceSignatures(request.PathValue("id"))
	if err != nil {
		if errors.Is(err, types.ErrDeviceNotFound) {
			WriteErrorResponse(response, http.StatusNotFound, []string{
				err.Error(),
			})
		} else {
			WriteInternalError(response, request.URL.Path, err)
		}
		return
	}
	WriteAPIResponse(response, http.StatusOK, signatures)
}
