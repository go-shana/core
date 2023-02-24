package rpc

import (
	"context"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-shana/core/errors"
	"github.com/go-shana/core/internal/rpc"
	"github.com/huandu/xstrings"
)

type HandlerFunc[Request, Response any] func(ctx context.Context, req *Request) (resp *Response, err error)

type funcInfo struct {
	Package string
	Name    string
}

// Export exports a method to RPC.
// The name of method will be converted to kebab-case.
func Export[Request, Response any](method HandlerFunc[Request, Response]) {
	val := reflect.ValueOf(method)
	info := parseFuncInfo(val)
	compiled := Compile(method)
	name := xstrings.ToKebabCase(info.Name)
	register(info.Package, info.Name, name, reflect.ValueOf(compiled))
}

// ExportName exports a method to RPC with specified name.
func ExportName[Request, Response any](name string, method HandlerFunc[Request, Response]) {
	val := reflect.ValueOf(method)
	info := parseFuncInfo(val)
	compiled := Compile(method)
	register(info.Package, info.Name, name, reflect.ValueOf(compiled))
}

func parseFuncInfo(val reflect.Value) *funcInfo {
	errors.Assert(val.Kind() == reflect.Func)

	pc := val.Pointer()
	f := runtime.FuncForPC(pc)
	errors.Assert(f != nil)
	fullName := f.Name()

	idx := strings.LastIndex(fullName, ".")
	errors.Assert(idx > 0)
	pkg := fullName[:idx]
	name := fullName[idx+1:]

	return &funcInfo{
		Package: pkg,
		Name:    name,
	}
}

func register(pkg, funcName, name string, val reflect.Value) {
	registry := rpc.DefaultRegistry()
	handler := &rpc.Handler{
		Package:  pkg,
		Name:     name,
		FuncName: funcName,
		Func:     val,
	}
	registry.Register(handler)
}
