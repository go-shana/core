package httpjson

import (
	"fmt"
	"net/http"
	"sort"

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

	lines := make([]string, 0, len(uris))

	for _, uri := range uris {
		lines = append(lines, fmt.Sprintf("GET|POST\t%-[1]*[2]s => %[3]s", maxLen+len(parent)+1, uri[0], uri[1]))
	}

	sort.Strings(lines)

	for _, line := range lines {
		// TODO: use logger instead of fmt.
		fmt.Println(line)
	}

	for path, tree := range root.subRoutes {
		printRouteTree(tree, parent+"/"+path)
	}
}
