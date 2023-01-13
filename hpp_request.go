package hpp

import (
	"errors"
)

// HppRequest represents all of the contents of an Hpp request
type HppRequest struct {
    server     *Server
	pid        uint32
	payload    []byte
	rmcRequest RMCRequest
}

// Server returns the request server
func (request *HppRequest) Server() *Server {
	return request.server
}

// SetPID sets the clients NEX PID
func (request *HppRequest) SetPID(pid uint32) {
	request.pid = pid
}

// PID returns the clients NEX PID
func (request *HppRequest) PID() uint32 {
	return request.pid
}

// SetPayload sets the request payload
func (request *HppRequest) SetPayload(payload []byte) {
	request.payload = payload
}

// Payload returns the request payload
func (request *HppRequest) Payload() []byte {
	return request.payload
}

// RMCRequest returns the request RMC request
func (request *HppRequest) RMCRequest() RMCRequest {
	return request.rmcRequest
}

// NewHppRequest returns a new Hpp request
func NewHppRequest(server *Server, data []byte) (*HppRequest, error) {
	hppRequest := HppRequest{server: server}

	if data != nil {
		hppRequest.payload = data

		rmcRequest := NewRMCRequest()
		err := rmcRequest.FromBytes(data)
		if err != nil {
			return &HppRequest{}, errors.New("[HppRequest] Error parsing RMC request: " + err.Error())
		}

		hppRequest.rmcRequest = rmcRequest
	}

	return &hppRequest, nil
}
