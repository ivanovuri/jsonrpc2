# Golang JSON-RPC 2.0 Server using reflect library

This library is an HTTP server implementation of the [JSON-RPC 2.0 Specification](http://www.jsonrpc.org/specification). Not fully specifications compliant yet. Batch requests are coming soon.

Heavily driven by reflection, which allows avoiding writing additional logic for custom rpc methods.


### Quickstart

```golang
package main

import (
	"fmt"

	"github.com/ivanovuri/jsonrpc2"
)

func main() {
	manager := jsonrpc2.NewJsonRpc2()

	manager.RegisterCall("substract", Substract)
	manager.RegisterCall("add", AddTwoInts)
	manager.RegisterCall("pa", PositionalAdd)
	// Returning error in case of adding same method twice
	if err := manager.RegisterCall("count_formula", CountFormula); err != nil {
		fmt.Println(err)
	}
	// Method will not be added here
	if err := manager.RegisterCall("count_formula", CountFormula); err != nil {
		fmt.Println(err)
	}

	manager.Run(":8888", "/")
}

func Substract(a, b int) int {
	return a - b
}

func AddTwoInts(a, b int) (int, int) {
	return a + b, 2
}

func CountFormula(x, y, z float32) float32 {
	if y == 0 {
		return 0
	}
	return (24+x)/y - x*y/(2.4+z)
}

type PositionalAddParamsStructure struct {
	A int
	B int
}

func PositionalAdd(params PositionalAddParamsStructure) int {
	return (params.A + 3)
}
```

Transport independent using of JSON-RPC 2.0 library.
Same as earlier, but server logic should be created manually. “Underscore” object (NewJsonRpc2_) is used here for distinction.

```golang
func HttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := jsonrpc2.DecodeRequest(r.Body)

	rpcProcessor := jsonrpc2.NewJsonRpc2_()
	rpcProcessor.RegisterMethod("substract", Substract)
	rpcProcessor.RegisterMethod("add", AddTwoInts)
	if err := rpcProcessor.RegisterMethod("pa", PositionalAdd); err != nil {
		fmt.Println(err)
	}

	w.Write(rpcProcessor.ProcessRequest(request))
}

func main() {
	log.Printf("Starting server on localhost:8888/")

	http.HandleFunc("/", HttpHandler)
	if err := http.ListenAndServe("localhost:8888", nil); err != nil {
		panic(err)
	}
}
```

When defining your own registered methods with the rpc server, it is important to consider both named and positional parameters per the specification.

Rpc call with positional parameters:
```
curl --header "Content-Type: application/json" -d '{
    "id": "asd",
    "jsonrpc": "2.0",
    "method": "add",
    "params": [
        20,
        10
    ]
}' 'http://localhost:8888/'
```

Rpc call with named parameters:
```
curl --header "Content-Type: application/json" -d '{
    "id": "asd",
    "jsonrpc": "2.0",
    "method": "pa",
    "params": {
        "A": 5,
        "B": 5
    }
}' 'http://localhost:8888/'
```