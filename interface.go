package jsonrpc2

type JsonRpc2_ interface {
	RegisterMethod(string, any) error
	DeregisterMethod(string, any) error
	ProcessRequest(Request) []byte // To be deleted
	ProcessSingleRequest(Request) []byte
	ProcessBatchRequest([]Request) []byte
}

func NewJsonRpc2_() JsonRpc2_ {
	rpc := new(Manager)
	rpc.methods = make(calls)

	return rpc
}
