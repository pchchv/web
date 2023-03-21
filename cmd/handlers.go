package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/pchchv/golog"
	"github.com/pchchv/web"
	"github.com/pchchv/web/extensions/sse"
)

// StaticFilesHandler is used to serve static files
func StaticFilesHandler(rw http.ResponseWriter, r *http.Request) {
	wctx := web.Context(r)
	// '..' is replaced to prevent directory traversal which could go out of static directory
	path := strings.ReplaceAll(wctx.Params()["w"], "..", "-")
	path = strings.ReplaceAll(path, "~", "-")

	rw.Header().Set("Last-Modified", lastModified)
	http.ServeFile(rw, r, fmt.Sprintf("./static/%s", path))
}

func OriginalResponseWriterHandler(w http.ResponseWriter, r *http.Request) {
	rw := web.OriginalResponseWriter(w)
	if rw == nil {
		web.Send(w, "text/html", "got nil", http.StatusPreconditionFailed)
		return
	}
	web.Send(w, "text/html", "success", http.StatusOK)
}

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

func pushCSS(pusher http.Pusher, r *http.Request, path string) {
	cssOpts := &http.PushOptions{
		Header: http.Header{
			"Accept-Encoding": r.Header["Accept-Encoding"],
			"Content-Type":    []string{"text/css; charset=UTF-8"},
		},
	}
	err := pusher.Push(path, cssOpts)
	if err != nil {
		web.LOGHANDLER.Error(err)
	}
}

func pushJS(pusher http.Pusher, r *http.Request, path string) {
	cssOpts := &http.PushOptions{
		Header: http.Header{
			"Accept-Encoding": r.Header["Accept-Encoding"],
			"Content-Type":    []string{"application/javascript"},
		},
	}
	err := pusher.Push(path, cssOpts)
	if err != nil {
		web.LOGHANDLER.Error(err)
	}
}

func pushHomepage(r *http.Request, w http.ResponseWriter) {
	pusher, ok := w.(http.Pusher)
	if !ok {
		return
	}

	cp, _ := r.Cookie("pusher")
	if cp != nil {
		return
	}

	cookie := &http.Cookie{
		Name:   "pusher",
		Value:  "css,js",
		MaxAge: 300,
	}
	http.SetCookie(w, cookie)
	pushCSS(pusher, r, "/static/css/main.css")
	pushCSS(pusher, r, "/static/css/normalize.css")
	pushJS(pusher, r, "/static/js/main.js")
	pushJS(pusher, r, "/static/js/sse.js")
}
