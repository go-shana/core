package rpc

import (
	"strings"
)

type PackageMap map[string][]*Handler

type Registry struct {
	packages PackageMap
}

var defaultRegistry = &Registry{
	packages: PackageMap{},
}

// DefaultRegistry returns the default registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// Register registers a RPC handler.
func Register(handler *Handler) {
	defaultRegistry.Register(handler)
}

// Register registers a RPC handler.
func (r *Registry) Register(handler *Handler) {
	r.packages[handler.Package] = append(r.packages[handler.Package], handler)
}

// Handlers returns all registered handlers under a package.
func (r *Registry) Handlers(pkgPrefix string) (handlers []*Handler) {
	for pkg, hs := range r.packages {
		if pkgPrefix == "" || strings.HasPrefix(pkg, pkgPrefix) {
			if len(pkg) > len(pkgPrefix) && pkg[len(pkgPrefix)] != '/' {
				continue
			}

			handlers = append(handlers, hs...)
		}
	}

	return
}
