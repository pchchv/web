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

import "net/http"

const wgoCtxKey = ctxkey("webgocontext")

// ctxkey is a custom string type to store the WebGo context within the HTTP request context.
type ctxkey string

// ContextPayload is a WebgoContext.
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
