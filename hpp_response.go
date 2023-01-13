package hpp

// HppResponse represents all of the contents of an Hpp response
type HppResponse struct {
	pid        uint32
	payload    []byte
	rmcResponse RMCResponse
}

// PID returns the clients NEX PID
func (response *HppResponse) PID() uint32 {
	return response.pid
}

// SetPayload sets the response payload
func (request *HppResponse) SetPayload(payload []byte) {
	request.payload = payload
}

// Payload returns the response payload
func (request *HppResponse) Payload() []byte {
	return request.payload
}

// RMCResponse returns the response RMC response
func (request *HppResponse) RMCResponse() RMCResponse {
	return request.rmcResponse
}

// NewHppResponse returns a new Hpp response
func NewHppResponse(rmcResponse RMCResponse, pid uint32) (*HppResponse) {
	hppResponse := HppResponse{}

	data := rmcResponse.Bytes()
	hppResponse.payload = data
	hppResponse.pid = pid

	return &hppResponse
}
