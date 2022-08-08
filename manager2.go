package jsonrpc2

import (
	"encoding/json"
	"reflect"

	"github.com/ivanovuri/jsonrpc2/inspect"
)

type Manager struct {
	methods calls
}

func (m *Manager) RegisterMethod(name string, method any) error {
	return m.methods.Add(name, method)
}

func (m *Manager) DeregisterMethod(name string, method any) error {
	return m.methods.Remove(name)
}

func (m *Manager) ProcessRequest(request Request) *Response {
	if callableMethod, rpcErr := m.methods.Get(request.Method); rpcErr != nil {
		return ErrorReply(request.Id,
			MethodNotFoundCode.Code(),
			MethodNotFoundCode.Message())
	} else {
		var in []reflect.Value
		var err error

		if inspect.NamedCall(callableMethod) {
			in, err = inspect.ParseNamedParams(request.Params, callableMethod)
			if err != nil {
				return ErrorReply(request.Id,
					InvalidParamsCode.Code(),
					InvalidParamsCode.Message())
			}
		} else {
			in, err = inspect.ParsePositionalParams(request.Params, callableMethod)
			if err != nil {
				return ErrorReply(request.Id,
					InvalidParamsCode.Code(),
					InvalidParamsCode.Message())
			}
		}
		result := inspect.ExecuteMethod(callableMethod, in)
		replyObj := MakeReply(request, result)
		return replyObj
	}
}

func (m *Manager) ProcessSingleRequest(singleRequest Request) []byte {
	result, _ := json.Marshal(m.ProcessRequest(singleRequest))

	return result
}

func (m *Manager) ProcessBatchRequest(batch []Request) []byte {
	var replies []Response

	for _, singleRequest := range batch {
		processed := m.ProcessRequest(singleRequest)
		if processed != nil {
			replies = append(replies, *processed)
		}
	}

	result, _ := json.Marshal(replies)

	return result
}
