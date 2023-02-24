// Package httpjson provides utilities to handle Shana-opinioned HTTP JSON requests.
package httpjson

import (
	"context"
	"net/http"

	"github.com/go-shana/core/rpc"
)

// Server is a HTTP JSON server.
type Server struct {
	config *Config
	server *http.Server
}

var _ rpc.Server = new(Server)

// NewServer creates a new HTTP JSON server.
func NewServer(config *Config) *Server {
	router := NewRouter(config)
	return &Server{
		config: config,
		server: &http.Server{
			Addr:    config.Addr,
			Handler: router,
		},
	}
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context) error {
	return s.server.ListenAndServe()
}

// Shutdown stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
