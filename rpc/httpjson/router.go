package httpjson

import (
	"fmt"
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

	printRouteTree(root, "")

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

func printRouteTree(root *routeTree, parent string) {
	uris := make([][]string, 0, len(root.handlers))
	maxLen := 0

	for path, handler := range root.handlers {
		if len(path) > maxLen {
			maxLen = len(path)
		}

		uris = append(uris, []string{parent + "/" + path, handler.handler.FuncName})
	}

	for _, uri := range uris {
		// TODO: use logger instead of fmt.
		fmt.Printf("%[1]*[2]s => %[3]s\n", maxLen, uri[0], uri[1])
	}

	for path, tree := range root.subRoutes {
		printRouteTree(tree, parent+"/"+path)
	}
}
