package rpc

import "reflect"

// Handler represents a RPC handler.
type Handler struct {
	Package  string // Package name.
	Name     string // API name.
	FuncName string // Function name.

	// Func must be a function with signature:
	//
	//	func(ctx context.Context, req *Request) (resp *Response, err error)
	Func reflect.Value
}
