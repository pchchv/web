package web

import (
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
