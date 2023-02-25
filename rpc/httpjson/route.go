package httpjson

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/go-shana/core/data"
	"github.com/go-shana/core/errors"
	"github.com/go-shana/core/initer"
	"github.com/go-shana/core/internal/global"
	"github.com/go-shana/core/internal/rpc"
	"github.com/go-shana/core/validator"
)

// routeTree is a HTTP JSON route.
type routeTree struct {
	handlers  routeHandlerMap
	subRoutes routeMap
}

type routeHandler struct {
	handler     *rpc.Handler
	handlerFunc http.HandlerFunc
}

type routeHandlerMap map[string]*routeHandler

type routeMap map[string]*routeTree

func newRoute() *routeTree {
	return &routeTree{
		handlers:  routeHandlerMap{},
		subRoutes: routeMap{},
	}
}

// Lookup finds a handler by path.
func (r *routeTree) Lookup(path string) http.HandlerFunc {
	if path == "" {
		return nil
	}

	paths := strings.Split(path, "/")
	last := len(paths) - 1

	for ; last > 0; last-- {
		if paths[last] != "" {
			break
		}
	}

	handlerName := paths[last]
	route := r

	for _, path := range paths[:last] {
		if path == "" {
			continue
		}

		route = route.subRoutes[path]

		if route == nil {
			return nil
		}
	}

	handler := route.handlers[handlerName]

	if handler == nil {
		return nil
	}

	return handler.handlerFunc
}

func parseRoute(config *Config, handlers []*rpc.Handler) (root *routeTree) {
	pkgPrefix := config.PkgPrefix
	root = newRoute()

	for _, handler := range handlers {
		paths := parsePackage(pkgPrefix, handler.Package)
		fn := parseHandlerFunc(config, handler)
		r := root
		m := root.subRoutes
		ok := false

		for _, path := range paths {
			r, ok = m[path]

			if !ok {
				r = newRoute()
				m[path] = r
			}

			m = r.subRoutes
		}

		r.handlers[handler.Name] = &routeHandler{
			handler:     handler,
			handlerFunc: fn,
		}
	}

	sonic.Pretouch(reflect.TypeOf(Response{}))
	return
}

// parsePackage stripes pkgPrefix from pkg and separate pkg with '/'.
func parsePackage(pkgPrefix, pkg string) []string {
	path := pkg[len(pkgPrefix):]
	start := 0
	end := len(path)

	for i := 0; i < end; i++ {
		r := path[i]

		if r != '/' {
			break
		}

		start++
	}

	for i := end - 1; i >= 0; i-- {
		r := path[i]

		if r != '/' {
			break
		}

		end--
	}

	striped := path[start:end]

	if striped == "" {
		return nil
	}

	return strings.Split(striped, "/")
}

// Response is the Shana-opinioned HTTP JSON protocol response.
type Response struct {
	Code    any        `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Error   string     `json:"error,omitempty"`
	Debug   *DebugInfo `json:"_debug,omitempty"` // Only valid when debug is enabled.
	Data    any        `json:"data,omitempty"`
}

// DebugInfo contains more debug information.
type DebugInfo struct {
	FuncName string   `json:"funcName,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

func parseHandlerFunc(config *Config, handler *rpc.Handler) http.HandlerFunc {
	fn := handler.Func
	fnType := fn.Type()
	reqType := fnType.In(1).Elem()
	respType := fnType.Out(0).Elem()
	debug := global.Debug()

	sonic.Pretouch(reqType)
	sonic.Pretouch(respType)

	handleFunc := func(w http.ResponseWriter, r *http.Request) (respVal reflect.Value, err error) {
		defer errors.Handle(&err)

		reqVal := reflect.New(reqType)
		req := reqVal.Interface()
		respHeader := w.Header()

		respHeader.Set("Content-Type", "application/json; charset=utf-8")

		if debug {
			setCORSHeaders(respHeader)
		}

		switch r.Method {
		case http.MethodGet:
			errors.Check(unmarshalQueryString(reqVal, r.URL))

		case http.MethodPost:
			errors.Check(unmarshalQueryString(reqVal, r.URL))

			if r.Body != nil {
				dec := sonic.ConfigDefault.NewDecoder(r.Body)
				errors.Check(dec.Decode(req))
				r.Body.Close()
			}

		default:
			w.WriteHeader(405)
			return
		}

		// TODO: add more context information.
		ctx := r.Context()

		errors.Check(validator.Validate(ctx, req))
		errors.Check(initer.Init(ctx, req))

		ret := fn.Call([]reflect.Value{reflect.ValueOf(ctx), reqVal})
		errors.Assert(len(ret) == 2)
		respVal = ret[0]
		errVal := ret[1]

		if errVal.IsValid() && !errVal.IsNil() {
			err = errVal.Interface().(error)
		}

		return
	}

	return func(w http.ResponseWriter, r *http.Request) {
		respVal, err := handleFunc(w, r)
		resp := &Response{}

		var code any
		var key error
		var errs []error
		var msg string

		if err == nil && !respVal.IsValid() {
			err = errors.New("httpjson: invalid response value")
		}

		if err != nil {
			if he, ok := err.(errors.HandlerError); ok {
				key = he.KeyError()
				msg = key.Error()

				if debug {
					errs = he.Unwrap()
				}
			} else {
				key = err
				msg = err.Error()

				if debug {
					if e := errors.Unwrap(err); e != nil {
						errs = []error{e}
					}
				}
			}

			if codeFunc := reflect.ValueOf(key).MethodByName("Code"); codeFunc.IsValid() {
				if t := codeFunc.Type(); t.Kind() == reflect.Func && t.NumIn() == 0 && t.NumOut() == 1 {
					if ret := codeFunc.Call(nil)[0]; ret.IsValid() {
						code = ret.Interface()
					}
				}
			}
		}

		if respVal.IsValid() {
			resp.Data = respVal.Interface()
		}

		if code == nil {
			resp.Error = msg
		} else {
			resp.Code = code
			resp.Message = msg
		}

		if debug {
			errStrs := make([]string, len(errs))

			for i, e := range errs {
				errStrs[i] = e.Error()
			}

			resp.Debug = &DebugInfo{
				FuncName: handler.FuncName,
				Errors:   errStrs,
			}
		}

		enc := sonic.ConfigDefault.NewEncoder(w)
		enc.SetEscapeHTML(false)

		if debug {
			enc.SetIndent("", "  ")
		}

		enc.Encode(resp)
	}
}

func setCORSHeaders(header http.Header) {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Headers", "Content-Type")
	header.Set("Access-Control-Allow-Methods", "GET, POST")
	header.Set("Access-Control-Max-Age", "86400")
}

func unmarshalQueryString(ptr reflect.Value, u *url.URL) (err error) {
	defer errors.Handle(&err)
	values := u.Query()
	val := ptr.Elem()
	errors.Assert(val.Kind() == reflect.Struct)

	raw := data.RawData{}

	for k, vs := range values {
		l := len(vs)
		if l == 0 {
			continue
		}

		raw[k] = json.Number(vs[l-1])
	}

	d := data.Make(raw)
	dec := &data.Decoder{
		TagName: "json",
	}
	errors.Check(dec.Decode(d, ptr.Interface()))

	return
}
