package jsonrpc2

type JsonRpc2 interface {
	RegisterMethod(string, any) error
	DeregisterMethod(string, any) error
	ProcessSingleRequest(Request) []byte
	ProcessBatchRequest([]Request) []byte
}

func NewJsonRpc2() JsonRpc2 {
	rpc := new(Manager)
	rpc.methods = make(calls)

	return rpc
}
