package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	deviceService DeviceService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceService DeviceService) *Server {
	return &Server{
		listenAddress: listenAddress,
		deviceService: deviceService,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))
	mux.Handle("/api/v0/create-signature-device", http.HandlerFunc(s.CreateSignatureDevice))
	mux.Handle("/api/v0/sign-transaction", http.HandlerFunc(s.SignTransaction))
	mux.Handle("/api/v0/devices", http.HandlerFunc(s.Devices))
	mux.Handle("/api/v0/device-signs/{id}", http.HandlerFunc(s.DeviceSignatures))

	// TODO: register further HandlerFuncs here ...

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter, url string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	log.Printf("Internal error: \nPath: %s\nError msg: %v\n", url, err)
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w, "", err)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w, "", err)
	}

	w.Write(bytes)
}
