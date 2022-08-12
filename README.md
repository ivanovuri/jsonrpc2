# Golang JSON-RPC 2.0 reflection driven library

This library is an implementation of the [JSON-RPC 2.0 Specification](http://www.jsonrpc.org/specification). Looks fully specifications compliant. HTTP Server code which was included earlier removed, so this can be used in independent manner.

Heavily driven by reflection, which allows avoiding writing additional logic for custom rpc methods. All you need is to create regular Go function and register it in repository with specific name.

The main disadvantage of "independence", the need to write code to implement the transport logic.
It is possible that major transports will be added soon to make it easier to use by client code, but the ability to use it in its raw form will remain unchanged.

### Quickstart

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ivanovuri/jsonrpc2"
)

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	rpcProcessor := jsonrpc2.NewJsonRpc2_()
	rpcProcessor.RegisterMethod("substract", Substract)
	rpcProcessor.RegisterMethod("add", AddTwoInts)
	if err := rpcProcessor.RegisterMethod("pa", PositionalAdd); err != nil {
		fmt.Println(err)
	}
	rpcProcessor.RegisterMethod("rsf", returnStructFn)

	incomingRequestData, _ := ioutil.ReadAll(r.Body)

	w.Header().Set("Content-Type", "application/json")

	if rSingle, err := jsonrpc2.DecodeSingleRequest(incomingRequestData); err != nil {
		if rBatch, _ := jsonrpc2.DecodeBatchRequests(incomingRequestData); err != nil {
			w.Write(rpcProcessor.ProcessBatchRequest(rBatch))
		} else {
			fmt.Println("Wtf, Junk input ;-)")
		}
	} else {
		w.Write(rpcProcessor.ProcessSingleRequest(*rSingle))
	}
}

func main() {
	log.Printf("Starting server on localhost:8888/")

	http.HandleFunc("/", HttpHandler)
	if err := http.ListenAndServe("localhost:8888", nil); err != nil {
		panic(err)
	}
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

type returnStruct struct {
	Inf string `json:"information"`
	Res int    `json:"result"`
}

func returnStructFn(a, b int) returnStruct {
	return returnStruct{
		Inf: "big Information string",
		Res: a + b,
	}
}
```
### Requests examples
Single RPC request with positional parameters:
```bash
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
Single RPC request with named parameters:
```bash
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
Batch RPC request:
```bash
curl --header "Content-Type: application/json" -d '[
  { "id": "asd3", "jsonrpc": "2.0", "method": "pa", "params": { "A": 5, "B": 5 }},
  { "jsonrpc": "2.0", "method": "add", "params": [ 20, 10 ]},
  { "id": "threeplustwo", "jsonrpc": "2.0", "method": "add", "params": [ 3, 2 ]}
]' 'http://localhost:8888/'
```
RPC to function returning struct:
```bash
curl --header "Content-Type: application/json" -d '{
    "id": "asd",
    "jsonrpc": "2.0",
    "method": "rsf",
    "params": [1,2]
}' 'http://localhost:8888/'
```