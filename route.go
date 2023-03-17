package web

import (
	"net/http"
	"strings"
)

// Route defines a route for each API
type Route struct {
	// Name is unique identifier for the route
	Name string
	// Method is the HTTP request method/type
	Method string
	// Pattern is the URI pattern to match
	Pattern string
	// TrailingSlash if set to true, the URI will be matched with or without
	// a trailing slash. IMPORTANT: It does not redirect.
	TrailingSlash bool

	// FallThroughPostResponse if enabled will execute all the handlers even if a response was already sent to the client
	FallThroughPostResponse bool

	// Handlers is a slice of http.HandlerFunc which can be middlewares or anything else. Though only 1 of them will be allowed to respond to client.
	// subsequent writes from the following handlers will be ignored
	Handlers []http.HandlerFunc

	hasWildcard bool
	fragments   []uriFragment
	paramsCount int

	// skipMiddleware if true, middleware added using `router` will not be applied to this Route.
	// This is used only when a Route is set using the RouteGroup, which can have its own set of middleware
	skipMiddleware bool
	// middlewareList is used at the last stage, i.e. right before starting the server
	middlewarelist []Middleware

	initialized bool

	serve http.HandlerFunc
}

type uriFragment struct {
	isVariable  bool
	hasWildcard bool
	// fragment will be the key name, if it's a variable/named URI parameter
	fragment string
}

func (r *Route) parseURIWithParams() {
	// if there are no URI params, then there's no need to set route parts
	if !strings.Contains(r.Pattern, ":") {
		return
	}

	fragments := strings.Split(r.Pattern, "/")
	if len(fragments) == 1 {
		return
	}

	rFragments := make([]uriFragment, 0, len(fragments))
	for _, fragment := range fragments[1:] {
		hasParam := false
		hasWildcard := false

		if strings.Contains(fragment, ":") {
			hasParam = true
			r.paramsCount++
		}
		if strings.Contains(fragment, "*") {
			r.hasWildcard = true
			hasWildcard = true
		}

		key := strings.ReplaceAll(fragment, ":", "")
		key = strings.ReplaceAll(key, "*", "")
		rFragments = append(
			rFragments,
			uriFragment{
				isVariable:  hasParam,
				hasWildcard: hasWildcard,
				fragment:    key,
			})
	}
	r.fragments = rFragments
}

func (r *Route) setupMiddleware(reverse bool) {
	if reverse {
		for i := range r.middlewarelist {
			m := r.middlewarelist[i]
			srv := r.serve
			r.serve = func(rw http.ResponseWriter, req *http.Request) {
				m(rw, req, srv)
			}
		}
	} else {
		for i := len(r.middlewarelist) - 1; i >= 0; i-- {
			m := r.middlewarelist[i]
			srv := r.serve
			r.serve = func(rw http.ResponseWriter, req *http.Request) {
				m(rw, req, srv)
			}
		}
	}
	// clear the middleware list, because it is already configured for the route
	r.middlewarelist = nil
}
