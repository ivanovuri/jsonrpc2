package jsonrpc2

import (
	"fmt"
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

func (m *Manager) ProcessRequest(request Request) []byte {
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
		return MakeResponse(request, result)
	}
}

func (m *Manager) ProcessRequestEx(batch []Request) []byte {
	for _, v := range batch {
		fmt.Println(v.Id)
	}
	return nil
	// if callableMethod, rpcErr := m.methods.Get(request.Method); rpcErr != nil {
	// 	return ErrorReply(request.Id,
	// 		MethodNotFoundCode.Code(),
	// 		MethodNotFoundCode.Message())
	// } else {
	// 	var in []reflect.Value
	// 	var err error

	// 	if inspect.NamedCall(callableMethod) {
	// 		in, err = inspect.ParseNamedParams(request.Params, callableMethod)
	// 		if err != nil {
	// 			return ErrorReply(request.Id,
	// 				InvalidParamsCode.Code(),
	// 				InvalidParamsCode.Message())
	// 		}
	// 	} else {
	// 		in, err = inspect.ParsePositionalParams(request.Params, callableMethod)
	// 		if err != nil {
	// 			return ErrorReply(request.Id,
	// 				InvalidParamsCode.Code(),
	// 				InvalidParamsCode.Message())
	// 		}
	// 	}
	// 	result := inspect.ExecuteMethod(callableMethod, in)
	// 	return MakeResponse(request, result)
	// }
}
