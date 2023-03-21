package main

import (
	"net/http"

	"github.com/pchchv/web"
)

func ParamHandler(w http.ResponseWriter, r *http.Request) {
	// Web context
	wctx := web.Context(r)
	// URI parameters, map[string]string
	params := wctx.Params()
	// route, the web.Route which is executing this request
	route := wctx.Route
	web.R200(
		w,
		map[string]interface{}{
			"route_name":    route.Name,
			"route_pattern": route.Pattern,
			"params":        params,
			"chained":       r.Header.Get("chained"),
		},
	)
}

func InvalidJSONHandler(w http.ResponseWriter, r *http.Request) {
	web.R200(w, make(chan int))
}
