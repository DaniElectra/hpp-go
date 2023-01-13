# hpp-go
Go implementation of the NEX HPP protocol

### Other NEX libraries
[nex-go](https://github.com/PretendoNetwork/nex-go) - NEX PRUDP server library

[nex-protocols-go](https://github.com/PretendoNetwork/nex-protocols-go) - NEX protocol definitions

### Install

`go get github.com/PretendoNetwork/hpp-go`

### Usage note

This module is a barebones HPP server for use with titles using the Nintendo NEX library over HTTPS. You will need to provide a function which returns the NEX password of a given PID account:

```go
func(pid uint32) (string, uint32) {
    return password, errorCode
}
```

### Usage

```go
package main

import (
    "fmt"

    hpp "github.com/PretendoNetwork/hpp-go"
)

func main() {
    nexServer := hpp.NewServer()
    nexServer.SetAccessKey("76f26496")
    nexServer.SetPasswordFromPIDFunction(passwordFromPID)

    nexServer.On("Data", func(packet *hpp.HppRequest) {
        request := packet.RMCRequest()

        fmt.Println("== Swapdoodle - Hpp ==")
        fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
        fmt.Printf("Method ID: %#v\n", request.MethodID())
        fmt.Println("======================")
    })

    nexServer.Listen("hpp-001a2c00-l1.n.app.pretendo.cc", "path/to/server.crt", "path/to/server.key")
}
```
