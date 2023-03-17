package web

import (
	"bytes"
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

func routeServeChainedHandlers(r *Route) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		crw, ok := rw.(*customResponseWriter)
		if !ok {
			crw = newCRW(rw, http.StatusOK)
		}

		for _, handler := range r.Handlers {
			if crw.written && !r.FallThroughPostResponse {
				break
			}
			handler(crw, req)
		}
	}
}

// init does all the initializations required for the route
func (r *Route) init() error {
	if r.initialized {
		return nil
	}
	r.initialized = true

	r.parseURIWithParams()
	r.serve = defaultRouteServe(r)
	return nil
}

func defaultRouteServe(r *Route) http.HandlerFunc {
	if len(r.Handlers) > 1 {
		return routeServeChainedHandlers(r)
	}
	// when there is only 1 handler, the custom response writer does not
	// have to check if the answer is already written or enabled
	return r.Handlers[0]
}

// matchPath matches requestURI with a route URI pattern
func (r *Route) matchPath(requestURI string) (bool, map[string]string) {
	p := bytes.NewBufferString(r.Pattern)
	if r.TrailingSlash {
		p.WriteString("/")
	} else {
		if requestURI[len(requestURI)-1] == '/' {
			return false, nil
		}
	}

	if r.Pattern == requestURI || p.String() == requestURI {
		return true, nil
	}

	return r.matchWithWildcard(requestURI)
}

func (r *Route) matchWithWildcard(requestURI string) (bool, map[string]string) {
	// if r.fragments is empty, it means that there are no variables in the URI template,
	// hence there is no point in checking it
	if len(r.fragments) == 0 {
		return false, nil
	}

	params := make(map[string]string, r.paramsCount)
	uriFragments := strings.Split(requestURI, "/")[1:]
	fragmentsLastIdx := len(r.fragments) - 1
	fragmentIdx := 0
	uriParameter := make([]string, 0, len(uriFragments))

	for idx, fragment := range uriFragments {
		// if the part is empty, it means it is the end of the URI with a slash.
		if fragment == "" {
			break
		}

		if fragmentIdx > fragmentsLastIdx {
			return false, nil
		}

		currentFragment := r.fragments[fragmentIdx]
		if !currentFragment.isVariable && currentFragment.fragment != fragment {
			return false, nil
		}

		uriParameter = append(uriParameter, fragment)
		if currentFragment.isVariable {
			params[currentFragment.fragment] = strings.Join(uriParameter, "/")
		}

		if !currentFragment.hasWildcard {
			uriParameter = make([]string, 0, len(uriFragments)-idx)
			fragmentIdx++
			continue
		}

		nextIdx := fragmentIdx + 1
		if nextIdx > fragmentsLastIdx {
			continue
		}
		nextPart := r.fragments[nextIdx]

		// If the URI has more fragments/parameters after the wildcard,
		// then the fragment immediately following the wildcard cannot be a variable or another wildcard.
		if !nextPart.isVariable && nextPart.fragment == fragment {
			// remove the last added 'part' from the parameters, because it is part of the static URI
			params[currentFragment.fragment] = strings.Join(uriParameter[:len(uriParameter)-1], "/")
			uriParameter = make([]string, 0, len(uriFragments)-idx)
			fragmentIdx += 2
		}
	}

	if len(params) != r.paramsCount {
		return false, nil
	}

	return true, params
}
