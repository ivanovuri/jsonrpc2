package jsonrpc2

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetMethodError(t *testing.T) {
	manager := rpcManager{
		methods: make(rpcMethods),
	}

	manager.RegisterCall("test", func(a int) int { return a * a })

	_, eCode := manager.GetMethod("test3")
	if eCode.Code() != int(MethodNotFoundCode) {
		t.Fatalf("Error code is incorrect %d != %d",
			eCode.Code(),
			MethodNotFoundCode)
	}
}

func init() {
	s := NewJsonRpc2()
	s.RegisterCall("add", addTwoInts)
	s.Run(":61001", "/api/v1/rpc")
}

func addTwoInts(a, b int) int {
	return a + b
}

func TestRpcCallWithPositionalParamter(t *testing.T) {
	var result struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  int    `json:"result"`
		Id      int    `json:"id"`
	}

	body := `{ "id": 30, "jsonrpc": "2.0", "method": "add", "params": [40, 203] }`
	buf := bytes.NewBuffer([]byte(body))

	resp, err := http.Post("http://localhost:61001/api/v1/rpc", "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}

	rdr := bufio.NewReader(resp.Body)
	dec := json.NewDecoder(rdr)
	dec.Decode(&result)

	if result.Id != 30 {
		t.Fatalf("Expected Id to be 30 got %d", result.Id)
	}

	if result.Jsonrpc != "2.0" {
		t.Fatalf("Expected Jsonrpc to be 2.0 got %s", result.Jsonrpc)
	}

	if result.Result != 243 {
		t.Fatalf("Expected Jsonrpc to be 243 got %d", result.Result)
	}
}
