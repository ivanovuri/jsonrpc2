package jsonrpc2

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/ivanovuri/jsonrpc2/inspect"
)

type rpcMethods map[string]any

type JsonRpc2 interface {
	RegisterCall(string, any) error
	Run(string, string) error
}

func NewJsonRpc2() JsonRpc2 {
	rpc := new(rpcManager)
	rpc.methods = make(rpcMethods)

	return rpc
}

type rpcManager struct {
	methods rpcMethods
}

// RegisterCall implements JsonRpc2
func (m *rpcManager) RegisterCall(name string, method interface{}) error {
	if _, ok := m.methods[name]; ok {
		return fmt.Errorf("identity [%s] already registered", name)
	}

	m.methods[name] = method

	return nil
}

// Run implements JsonRpc2
func (m *rpcManager) Run(addr, path string) error {
	log.Printf("Starting server on %s", addr+path)

	http.HandleFunc(path, m.RpcHandler)
	return http.ListenAndServe(addr, nil)
}

func (m *rpcManager) RpcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := DecodeRequest(r.Body)

	if callableMethod, rpcErr := m.GetMethod(request.Method); rpcErr != nil {
		w.Write(ErrorReply(request.Id,
			int(rpcErr.Code()),
			rpcErr.Message()))
	} else {
		var in []reflect.Value
		var err error

		if inspect.NamedCall(callableMethod) {
			in, err = inspect.ParseNamedParams(request.Params, callableMethod)
			if err != nil {
				w.Write(ErrorReply(request.Id,
					InvalidParamsCode.Code(),
					InvalidParamsCode.Message()))
				return
			}
		} else {
			in, err = inspect.ParsePositionalParams(request.Params, callableMethod)
			if err != nil {
				w.Write(ErrorReply(request.Id,
					InvalidParamsCode.Code(),
					InvalidParamsCode.Message()))
				return
			}
		}

		result := inspect.ExecuteMethod(callableMethod, in)
		response := MakeResponse(request, result)
		w.Write(response)
	}
}

func (m *rpcManager) GetMethod(methodName string) (any, RpcError) {
	if _, ok := m.methods[methodName]; ok {
		call := m.methods[methodName]
		return call, nil
	} else {
		return nil, MethodNotFoundCode
	}
}
