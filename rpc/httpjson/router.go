package httpjson

import (
	"net/http"

	"github.com/go-shana/core/internal/rpc"
)

// Router is a HTTP JSON router.
type Router struct {
	root *routeTree
}

var _ http.Handler = new(Router)

// NewRouter creates a new HTTP JSON router.
func NewRouter(config *Config) *Router {
	pkgPrefix := config.PkgPrefix
	registry := rpc.DefaultRegistry()
	handlers := registry.Handlers(pkgPrefix)
	root := parseRoute(config, handlers)

	return &Router{
		root: root,
	}
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handlerFunc := r.root.Lookup(req.URL.Path)

	if handlerFunc == nil {
		http.NotFound(w, req)
		return
	}

	handlerFunc(w, req)
}
