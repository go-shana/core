// Package httpjson provides utilities to handle Shana-opinioned HTTP JSON requests.
package httpjson

import (
	"context"
	"errors"
	"fmt"
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
	addr := fmt.Sprintf("%v:%v", config.IP, config.Port)
	return &Server{
		config: config,
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context) error {
	// TODO: use logger instead of fmt.
	fmt.Printf("Server is starting at address %v\n", s.server.Addr)

	err := s.server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

// Shutdown stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
