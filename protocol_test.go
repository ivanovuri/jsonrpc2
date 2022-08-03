package jsonrpc2

import (
	"encoding/json"
	"testing"
)

func TestRequest(t *testing.T) {
	jsRequest := new(Request)
	jsonPositionalRequest := []byte(`{
		"id": 30,
		"jsonrpc": "2.0",
		"method": "add",
		"params": [42, 123]
	}`)

	err := json.Unmarshal(jsonPositionalRequest, &jsRequest)
	if err != nil {
		t.Fatalf(`Unmarshal request failed with %s`, err)
	}

	if jsRequest.Method != "add" {
		t.Fatalf("%s != add", jsRequest.Method)
	}
}
