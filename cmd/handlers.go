package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/pchchv/golog"
	"github.com/pchchv/web"
	"github.com/pchchv/web/extensions/sse"
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

func SSEHandler(sse *sse.SSE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := web.Context(r).Params()
		r.Header.Set(sse.ClientIDHeader, params["clientID"])

		err := sse.Handler(w, r)
		if err != nil && !errors.Is(err, context.Canceled) {
			golog.Info("errorLogger:", err.Error())
			return
		}
	}
}

func ErrorSetterHandler(w http.ResponseWriter, r *http.Request) {
	err := errors.New("oh no, server error")
	web.SetError(r, err)

	web.R500(w, err.Error())
}

func InvalidJSONHandler(w http.ResponseWriter, r *http.Request) {
	web.R200(w, make(chan int))
}
