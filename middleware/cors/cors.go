/*
The cors package sets the appropriate CORS (https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) response headers and allows for customization. The following settings are allowed:
  - provide a list of allowed domains
  - provide a list of headers
  - set the maximum age of CORS headers.

The list of allowed methods is as follows
*/
package cors

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/pchchv/web"
)

const allowHeaders = "Accept,Content-Type,Content-Length,Accept-Encoding,Access-Control-Request-Headers,"

var defaultAllowMethods = "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS"

// Config holds all the configurations that are available to configure this middleware
type Config struct {
	TimeoutSecs    int
	Routes         []*web.Route
	AllowedOrigins []string
	AllowedHeaders []string
}

func allowedDomains() []string {
	// The domains mentioned here are default
	domains := []string{"*"}
	return domains
}

func getReqOrigin(r *http.Request) string {
	return r.Header.Get("Origin")
}

func allowedHeaders(headers []string) string {
	if len(headers) == 0 {
		return allowHeaders
	}

	allowedHeaders := strings.Join(headers, ",")
	if allowedHeaders[len(allowedHeaders)-1] != ',' {
		allowedHeaders += ","
	}
	return allowedHeaders
}

func allowedMethods(routes []*web.Route) string {
	if len(routes) == 0 {
		return defaultAllowMethods
	}

	methods := make([]string, 0, len(routes))
	for _, r := range routes {
		found := false
		for _, m := range methods {
			if m == r.Method {
				found = true
				break
			}
		}
		if found {
			continue
		}
		methods = append(methods, r.Method)
	}
	sort.Strings(methods)
	return strings.Join(methods, ",")
}

func allowedOrigin(reqOrigin string, allowedOriginRegex []regexp.Regexp) bool {
	for _, o := range allowedOriginRegex {
		// Set the appropriate response headers needed for CORS
		if o.MatchString(reqOrigin) || reqOrigin == "" {
			return true
		}
	}
	return false
}
