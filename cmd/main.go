package main

import (
	"net/http"
	"time"
)

var lastModified = time.Now().Format(http.TimeFormat)

func chain(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("chained", "true")
}

func routegroupMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("routegroup", "true")
	next(w, r)
}

func main() {
}
