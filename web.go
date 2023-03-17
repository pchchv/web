/*
The web package is a lightweight framework for building web applications.
It has a multiplexer, a mechanism for connecting middleware, and its own context management.
The main goal of web is to get as far away from the developer as possible,
i.e. it doesn't force you to build your application according to any particular template,
but just helps you make all the trivial things faster and easier.
e.g.
1. Getting named URI parameters.
2. Multiplexer for regex-matching URIs and the like.
3. Implementation of special application-level configurations or any similar objects into the request context as required.
*/
package web

import (
	"crypto/tls"
	"net/http"
)

const wgoCtxKey = ctxkey("webcontext")

var supportedHTTPMethods = []string{
	http.MethodOptions,
	http.MethodHead,
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

// ctxkey is a custom string type to store the Web context within the HTTP request context.
type ctxkey string

// ContextPayload is a WebContext.
// A new ContextPayload instance is injected inside the context object of each request.
type ContextPayload struct {
	Route     *Route
	Err       error
	URIParams map[string]string
}

// Params returns the URI parameters of the corresponding route.
func (cp *ContextPayload) Params() map[string]string {
	return cp.URIParams
}

func (cp *ContextPayload) reset() {
	cp.Route = nil
	cp.Err = nil
}

// SetError sets the value of err in context.
func (cp *ContextPayload) SetError(err error) {
	cp.Err = err
}

// Error returns the error set within the context.
func (cp *ContextPayload) Error() error {
	return cp.Err
}

// Context returns the ContextPayload injected inside the HTTP request context.
func Context(r *http.Request) *ContextPayload {
	return r.Context().Value(wgoCtxKey).(*ContextPayload)
}

// SetError is an auxiliary function for setting an error in the web context
func SetError(r *http.Request, err error) {
	ctx := Context(r)
	ctx.SetError(err)
}

// GetError is an auxiliary function to get the error from the web context
func GetError(r *http.Request) error {
	return Context(r).Error()
}

// ResponseStatus returns the response status code.
// This only works if http.ResponseWriter is not wrapped in
// another response writer before calling ResponseStatus.
func ResponseStatus(rw http.ResponseWriter) int {
	crw, ok := rw.(*customResponseWriter)
	if !ok {
		return http.StatusOK
	}
	return crw.statusCode
}

// OriginalResponseWriter returns the Go response record stored in
// the custom web response record.
func OriginalResponseWriter(rw http.ResponseWriter) http.ResponseWriter {
	crw, ok := rw.(*customResponseWriter)
	if !ok {
		return nil
	}
	return crw.ResponseWriter
}

func (router *Router) setupServer() {
	cfg := router.config
	router.httpsServer = &http.Server{
		Addr:         "",
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		},
	}
	router.httpServer = &http.Server{
		Addr:         "",
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	router.SetupMiddleware()
}

// SetupMiddleware initializes all middleware added with "Use".
// This function does not need to be called explicitly if router.Start() or router.StartHTTPS() is used.
// Instead, if the router is passed to an external server, the SetupMiddleware function should be called.
func (router *Router) SetupMiddleware() {
	// load middleware for all routes
	for _, routes := range router.allHandlers {
		for _, route := range routes {
			route.setupMiddleware(router.config.ReverseMiddleware)
		}
	}
}

// StartHTTPS starts the server with HTTPS enabled
func (router *Router) StartHTTPS() {
	cfg := router.config
	if cfg.CertFile == "" {
		LOGHANDLER.Fatal("No certificate provided for HTTPS")
	}

	if cfg.KeyFile == "" {
		LOGHANDLER.Fatal("No key file provided for HTTPS")
	}

	router.setupServer()

	host := cfg.Host
	if len(cfg.HTTPSPort) > 0 {
		host += ":" + cfg.HTTPSPort
	}
	router.httpsServer.Addr = host

	LOGHANDLER.Info("HTTPS server, listening on", router.httpsServer.Addr)
	err := router.httpsServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	if err != nil && err != http.ErrServerClosed {
		LOGHANDLER.Error("HTTPS server exited with error:", err.Error())
	}
}
