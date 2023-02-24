package rpc

import (
	"context"

	"github.com/go-shana/core/errors"
	"github.com/go-shana/core/initer"
	"github.com/go-shana/core/validator"
)

// Compile wraps a method with request validator, initer and error handler.
// The returned method is used in RPC server to handle any request.
//
// Business code should never call Compile unless for test purpose, e.g.
// writing unit test cases for a RPC method.
func Compile[Request, Response any](method HandlerFunc[Request, Response]) HandlerFunc[Request, Response] {
	return func(ctx context.Context, req *Request) (resp *Response, err error) {
		defer errors.Handle(&err)

		validator.Validate(ctx, req)
		initer.Init(ctx, req)
		resp, err = method(ctx, req)
		return
	}
}
