package jsonrpc2

import (
	"encoding/json"
	"io"
	"reflect"
)

const protocolVersion = "2.0"

type RpcError interface {
	Code() int
	Message() string
}

type ErrorCode int

const (
	ParseErrorCode     ErrorCode = -32700
	InvalidRequestCode ErrorCode = -32600
	MethodNotFoundCode ErrorCode = -32601
	InvalidParamsCode  ErrorCode = -32602
	InternalErrorCode  ErrorCode = -32603
)

func (s ErrorCode) Code() int {
	return int(s)
}

func (s ErrorCode) Message() string {
	switch s {
	case ParseErrorCode:
		return "Parse error"
	case InvalidRequestCode:
		return "Invalid Request"
	case MethodNotFoundCode:
		return "Method not found"
	case InvalidParamsCode:
		return "Invalid params"
	case InternalErrorCode:
		return "Internal error"
	}

	return "Undefined error"
}

func (s ErrorCode) String() string {
	return s.Message()
}

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Id      any             `json:"id"`
}

type Response struct {
	Jsonrpc string       `json:"jsonrpc"`
	Result  any          `json:"result,omitempty"`
	Id      any          `json:"id"`
	Error   *ErrorObject `json:"error,omitempty"`
}

type ErrorObject struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func ErrorReply(rpcId any, errCode int, message string) []byte {
	errorResponse := Response{
		Jsonrpc: protocolVersion,
		Id:      rpcId,
		Error: &ErrorObject{
			Code:    errCode,
			Message: message,
		},
	}

	errorReply, _ := json.Marshal(errorResponse)

	return errorReply
}

func MakeResult(returnedValues []reflect.Value) json.RawMessage {
	r := make([]any, len(returnedValues))
	for k, v := range returnedValues {
		r[k] = v.Interface()
	}

	marshalledResult, err := json.Marshal(r)
	if err != nil {
		return marshalledResult
	}

	return marshalledResult
}

func MakeResponse(r Request, returnedValues []reflect.Value) []byte {
	if r.Id != nil {
		rpcResult := MakeResult(returnedValues)

		marshalledResponse, _ := json.Marshal(Response{
			Jsonrpc: protocolVersion,
			Result:  rpcResult,
			Id:      r.Id,
		})
		return marshalledResponse
	}
	return nil
}

func DecodeRequest(requestReader io.Reader) Request {
	incomingRequest := new(Request)
	json.NewDecoder(requestReader).Decode(incomingRequest)

	return *incomingRequest
}

func DecodeRequest2(requestReader io.Reader) (*Request, error) {
	incomingRequest := new(Request)
	if err := json.NewDecoder(requestReader).Decode(incomingRequest); err != nil {
		return nil, err
	}

	return incomingRequest, nil
}

func DecodeBatchRequests(requestReader io.Reader) (*[]Request, error) {
	incomingRequest := new([]Request)
	if err := json.NewDecoder(requestReader).Decode(incomingRequest); err != nil {
		return nil, err
	}

	return incomingRequest, nil
}
