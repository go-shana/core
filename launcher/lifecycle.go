package launcher

import (
	"context"

	"github.com/go-shana/core/errors"
	"github.com/go-shana/core/internal/lifecycle"
)

// Func is a callback function in life cycle.
type Func func(context.Context) error

var errInvalidCallback = errors.New("invalid lifecycle callback")

// OnConnect registers f as a connecting service callback.
// All connecting service callbacks will be called before any start-up callbacks and
// after all config data initialized.
func OnConnect(f Func) {
	if f == nil {
		errors.Throw(errInvalidCallback)
		return
	}

	lifecycle.OnConnect.AddFunc(lifecycle.Func(f))
}

// OnStart registers f as a start-up callback.
// All start-up callbacks will be called before service starts serving and
// after all clients have connected.
func OnStart(f Func) {
	if f == nil {
		errors.Throw(errInvalidCallback)
		return
	}

	lifecycle.OnStart.AddFunc(lifecycle.Func(f))
}

// OnShutdown registers f as an shutdown callback.
// All shutdown callbacks will be called after service is gracefully shutdown and about to exit.
//
// It's not guaranteed to be called as a service can be killed by SIGKILL.
func OnShutdown(f Func) {
	if f == nil {
		errors.Throw(errInvalidCallback)
		return
	}

	lifecycle.OnShutdown.AddFunc(lifecycle.Func(f))
}
