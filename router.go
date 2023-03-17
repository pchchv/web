package web

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"sync"
)

var (
	validHTTPMethods = []string{
		http.MethodOptions,
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}
	crwPool = &sync.Pool{
		New: func() interface{} {
			return new(customResponseWriter)
		},
	}
)

// Router is the HTTP router
type Router struct {
	optHandlers    []*Route
	headHandlers   []*Route
	getHandlers    []*Route
	postHandlers   []*Route
	putHandlers    []*Route
	patchHandlers  []*Route
	deleteHandlers []*Route
	allHandlers    map[string][]*Route

	// NotFound is the generic handler for 404 resource not found response
	NotFound http.HandlerFunc

	// NotImplemented is the generic handler for 501 method not implemented
	NotImplemented http.HandlerFunc

	// config has all the app config
	config *Config

	// httpServer is the server handler for the active HTTP server
	httpServer *http.Server
	// httpsServer is the server handler for the active HTTPS server
	httpsServer *http.Server
}

// Middleware is the signature of WebGo's middleware
type Middleware func(http.ResponseWriter, *http.Request, http.HandlerFunc)

// customResponseWriter is a custom HTTP response writer
type customResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	written       bool
	headerWritten bool
}

func newCRW(rw http.ResponseWriter, rCode int) *customResponseWriter {
	crw := crwPool.Get().(*customResponseWriter)
	crw.ResponseWriter = rw
	crw.statusCode = rCode
	return crw
}

// WriteHeader is an implementation of an interface for getting the
// HTTP response code and adding it to the user's response writer.
func (crw *customResponseWriter) WriteHeader(code int) {
	if crw.headerWritten {
		return
	}

	crw.headerWritten = true
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Write is an implementation of an interface to respond to an HTTP request,
// but checks to see if a response has already been sent.
func (crw *customResponseWriter) Write(body []byte) (int, error) {
	crw.WriteHeader(crw.statusCode)
	crw.written = true
	return crw.ResponseWriter.Write(body)
}

// Flush calls http.Flusher to clean/flush the buffer.
func (crw *customResponseWriter) Flush() {
	if rw, ok := crw.ResponseWriter.(http.Flusher); ok {
		rw.Flush()
	}
}

// Hijack implements the http.Hijacker interface.
func (crw *customResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := crw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}

	return nil, nil, errors.New("unable to create hijacker")
}

// CloseNotify implements the http.CloseNotifier interface
func (crw *customResponseWriter) CloseNotify() <-chan bool {
	if n, ok := crw.ResponseWriter.(http.CloseNotifier); ok {
		return n.CloseNotify()
	}
	return nil
}

func (crw *customResponseWriter) Push(target string, opts *http.PushOptions) error {
	if n, ok := crw.ResponseWriter.(http.Pusher); ok {
		return n.Push(target, opts)
	}
	return errors.New("pusher not implemented")
}

func (crw *customResponseWriter) reset() {
	crw.statusCode = 0
	crw.written = false
	crw.headerWritten = false
	crw.ResponseWriter = nil
}

func releaseCRW(crw *customResponseWriter) {
	crw.reset()
	crwPool.Put(crw)
}

// discoverRoute returns the correct 'route', for the given request
func discoverRoute(path string, routes []*Route) (*Route, map[string]string) {
	for _, route := range routes {
		if ok, params := route.matchPath(path); ok {
			return route, params
		}
	}
	return nil, nil
}

// methodRoutes returns the list of Routes handling the HTTP method given the request
func (rtr *Router) methodRoutes(method string) (routes []*Route) {
	switch method {
	case http.MethodOptions:
		return rtr.optHandlers
	case http.MethodHead:
		return rtr.headHandlers
	case http.MethodGet:
		return rtr.getHandlers
	case http.MethodPost:
		return rtr.postHandlers
	case http.MethodPut:
		return rtr.putHandlers
	case http.MethodPatch:
		return rtr.patchHandlers
	case http.MethodDelete:
		return rtr.deleteHandlers
	}

	return nil
}
