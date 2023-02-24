// Package errors provides utilities to create and handle all kinds errors.
// It's recommended to use this package to replace any use of standard package "errors".
//
// The error handler design is imspired by Go2 error handling proposal.
// See https://github.com/golang/proposal/blob/master/design/go2draft-error-handling-overview.md.
package errors

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// HandlerError is an error containing several sub errors.
type HandlerError interface {
	error

	Unwrap() []error

	// KeyError returns the key error which will be considered as the major reason that causes error.
	// When Shana reports error message to client, the key error will be used as error message.
	KeyError() error
}

type handlerError struct {
	errors []error
	msg    atomic.Pointer[string]
}

var _ HandlerError = new(handlerError)

func newHandlerError(errors []error) *handlerError {
	return &handlerError{
		errors: errors,
	}
}

func (he *handlerError) Error() string {
	if msg := he.msg.Load(); msg != nil {
		return *msg
	}

	builder := &strings.Builder{}
	builder.WriteString(he.errors[0].Error())

	for _, err := range he.errors[1:] {
		builder.WriteRune('\n')
		builder.WriteString(err.Error())
	}

	msg := builder.String()
	he.msg.Store(&msg)
	return msg
}

func (he *handlerError) Unwrap() []error {
	return he.errors
}

func (he *handlerError) KeyError() error {
	return he.errors[0]
}

// Handle recovers from panic and returns error to the target if panic is an error.
func Handle(target *error) {
	r := recover()

	if r == nil {
		return
	}

	if err, ok := r.(error); ok {
		*target = err
		return
	}

	*target = New(fmt.Sprintf("errors: uncaught panic: %v", r))
}
