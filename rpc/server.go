package rpc

import "context"

// Server represents a long running server.
type Server interface {
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
