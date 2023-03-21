package main

import (
	"log"
	"net/http"
	"time"

	"github.com/pchchv/web"
	"github.com/pchchv/web/extensions/sse"
)

var lastModified = time.Now().Format(http.TimeFormat)

func chain(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("chained", "true")
}

func routegroupMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("routegroup", "true")
	next(w, r)
}

// errLogger is a middleware which will log all errors returned/set by a handler
func errLogger(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)

	err := web.GetError(r)
	if err != nil {
		// log only server errors
		if web.ResponseStatus(w) > 499 {
			log.Println("errorLogger:", err.Error())
		}
	}
}

func getRoutes(sse *sse.SSE) []*web.Route {
	return []*web.Route{
		{
			Name:          "root",
			Method:        http.MethodGet,
			Pattern:       "/",
			Handlers:      []http.HandlerFunc{HomeHandler},
			TrailingSlash: true,
		},
		{
			Name:          "matchall",
			Method:        http.MethodGet,
			Pattern:       "/matchall/:wildcard*",
			Handlers:      []http.HandlerFunc{ParamHandler},
			TrailingSlash: true,
		},
		{
			Name:                    "api",
			Method:                  http.MethodGet,
			Pattern:                 "/api/:param",
			Handlers:                []http.HandlerFunc{chain, ParamHandler},
			TrailingSlash:           true,
			FallThroughPostResponse: true,
		},
		{
			Name:          "invalidjson",
			Method:        http.MethodGet,
			Pattern:       "/invalidjson",
			Handlers:      []http.HandlerFunc{InvalidJSONHandler},
			TrailingSlash: true,
		},
		{
			Name:          "error-setter",
			Method:        http.MethodGet,
			Pattern:       "/error-setter",
			Handlers:      []http.HandlerFunc{ErrorSetterHandler},
			TrailingSlash: true,
		},
		{
			Name:          "original-responsewriter",
			Method:        http.MethodGet,
			Pattern:       "/original-responsewriter",
			Handlers:      []http.HandlerFunc{OriginalResponseWriterHandler},
			TrailingSlash: true,
		},
		{
			Name:          "static",
			Method:        http.MethodGet,
			Pattern:       "/static/:w*",
			Handlers:      []http.HandlerFunc{StaticFilesHandler},
			TrailingSlash: true,
		},
		{
			Name:          "sse",
			Method:        http.MethodGet,
			Pattern:       "/sse/:clientID",
			Handlers:      []http.HandlerFunc{SSEHandler(sse)},
			TrailingSlash: true,
		},
	}
}

func main() {
}
