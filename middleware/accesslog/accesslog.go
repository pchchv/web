/*
The accesslogs package provides a simple middleware for access logs.
The logs have the following format:
< timestamp> < HTTP request method> <full URL including request string parameters> <execution duration> <HTTP response status code>.
*/
package accesslog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pchchv/web"
)

// AccessLog is a middleware that prints an access log to stdout
func AccessLog(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, req)
	end := time.Now()

	web.LOGHANDLER.Info(
		fmt.Sprintf(
			"%s %s %s %d",
			req.Method,
			req.URL.String(),
			end.Sub(start).String(),
			web.ResponseStatus(rw),
		),
	)
}
