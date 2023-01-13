// Package hpp implements an API for creating bare-bones
// NEX servers and clients and provides the underlying
// HPP implementation
//
// No NEX protocols are implemented in this package. For
// NEX protocols see https://github.com/PretendoNetwork/nex-protocols-go
//
// No PIA code is implemented in this package
package hpp

import (
	"bytes"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
)

// Server represents a PRUDP server
type Server struct {
	hppEventHandles           map[string][]func(*HppRequest)
	hppClientResponses        map[uint32](chan []byte)
	passwordFromPIDHandler    func(pid uint32) (string, uint32)
	accessKey                 string
	nexVersion                int
}

// Listen starts a NEX Hpp server on a given address
func (server *Server) Listen(address string, certFile string, keyFile string) {
	hppHandler := func(w http.ResponseWriter, req *http.Request) {
		pidValue := req.Header.Get("pid")
		accessKeySignature := req.Header.Get("signature1")
		passwordSignature := req.Header.Get("signature2")

		pid, err := strconv.Atoi(pidValue)
		if err != nil {
			logger.Error(fmt.Sprintf("[Hpp] Invalid PID - %s", pidValue))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		accessKeySignatureBytes, err := hex.DecodeString(accessKeySignature)
		if err != nil {
			logger.Error(fmt.Sprintf("[Hpp] Invalid access key signature - %s", accessKeySignature))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		passwordSignatureBytes, err := hex.DecodeString(passwordSignature)
		if err != nil {
			logger.Error(fmt.Sprintf("[Hpp] Invalid password signature - %s", passwordSignature))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rmcRequestString := req.FormValue("file")

		rmcRequestBytes := []byte(rmcRequestString)

		hppRequest, _ := NewHppRequest(server, rmcRequestBytes)
		hppRequest.SetPID(uint32(pid))

		generatedAccessKeySignature, err := GenerateAccessKeySignature(server.AccessKey(), rmcRequestBytes)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !bytes.Equal(generatedAccessKeySignature, accessKeySignatureBytes) {
			logger.Error("[Hpp] Access key calculated signature did not match")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pidPassword, errorCode := server.passwordFromPIDHandler(uint32(pid))
		if errorCode != 0 {
			rmcRequest := hppRequest.RMCRequest()
			callID := rmcRequest.CallID()

			errorResponse := NewRMCResponse(callID)
			errorResponse.SetError(errorCode)

			_, err = w.Write(errorResponse.Bytes())
			if err != nil {
				logger.Error(err.Error())
			}

			return
		}

		generatedPasswordSignature := GeneratePasswordSignature(uint32(pid), pidPassword, rmcRequestBytes)

		// When the password signature fails, the server returns error PythonCore::ValidationError
		if !bytes.Equal(generatedPasswordSignature, passwordSignatureBytes) {
			logger.Error("[Hpp] Password calculated signature did not match")
			rmcRequest := hppRequest.RMCRequest()
			callID := rmcRequest.CallID()

			validationErrorResponse := NewRMCResponse(callID)
			validationErrorResponse.SetError(Errors.PythonCore.ValidationError)

			_, err = w.Write(validationErrorResponse.Bytes())
			if err != nil {
				logger.Error(err.Error())
			}

			return
		}

		server.hppClientResponses[uint32(pid)] = make(chan []byte)

		server.Emit("Data", hppRequest)

		rmcResponseBytes := <- server.hppClientResponses[uint32(pid)]

		if len(rmcResponseBytes) > 0 {
			_, err = w.Write(rmcResponseBytes)
			if err != nil {
				logger.Error(err.Error())
			}
		}

		delete(server.hppClientResponses, uint32(pid))
	}

	http.HandleFunc("/hpp/", hppHandler)

	hppServer := &http.Server{
		Addr: address,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS11,
		},
	}

	logger.Success(fmt.Sprintf("Hpp server listening on address - %s", address))

	err := hppServer.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		panic(err)
	}
}

// On sets the data event handler
func (server *Server) On(event string, handler func(*HppRequest)) {
	server.hppEventHandles[event] = append(server.hppEventHandles[event], handler)
}

// Emit runs the given event handle
func (server *Server) Emit(event string, request *HppRequest) {
	eventName := server.hppEventHandles[event]
	for i := 0; i < len(eventName); i++ {
		handler := eventName[i]
		go handler(request)
	}
}

// NexVersion returns the server NEX version
func (server *Server) NexVersion() int {
	return server.nexVersion
}

// SetNexVersion sets the server NEX version
func (server *Server) SetNexVersion(nexVersion int) {
	server.nexVersion = nexVersion
}

// AccessKey returns the server access key
func (server *Server) AccessKey() string {
	return server.accessKey
}

// SetAccessKey sets the server access key
func (server *Server) SetAccessKey(accessKey string) {
	server.accessKey = accessKey
}

// SetPasswordFromPIDFunction sets the function for the server to get a NEX password using the PID
func (server *Server) SetPasswordFromPIDFunction(handler func(pid uint32) (string, uint32)) {
	server.passwordFromPIDHandler = handler
}

// Send writes data to client
func (server *Server) Send(response *HppResponse) {
	pid := response.PID()
	payload := response.Payload()
	server.hppClientResponses[pid] <- payload
}

// NewServer returns a new NEX server
func NewServer() *Server {
	server := &Server{
		hppEventHandles:       make(map[string][]func(*HppRequest)),
		hppClientResponses:    make(map[uint32](chan []byte)),
	}

	return server
}
