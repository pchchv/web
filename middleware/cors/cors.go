/*
The cors package sets the appropriate CORS (https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) response headers and allows for customization. The following settings are allowed:
  - provide a list of allowed domains
  - provide a list of headers
  - set the maximum age of CORS headers.

The list of allowed methods is as follows
*/
package cors

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/pchchv/web"
)

const (
	headerCreds            = "Access-Control-Allow-Credentials"
	allowHeaders           = "Accept,Content-Type,Content-Length,Accept-Encoding,Access-Control-Request-Headers,"
	headerOrigin           = "Access-Control-Allow-Origin"
	headerMethods          = "Access-Control-Allow-Methods"
	headerReqHeaders       = "Access-Control-Request-Headers"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAccessControlAge = "Access-Control-Max-Age"
)

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

func allowedOriginsRegex(allowedOrigins ...string) []regexp.Regexp {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	} else {
		// If "*" is one of the allowed domains,
		// i.e. all domains, then the other values are ignored
		for _, val := range allowedOrigins {
			val = strings.TrimSpace(val)

			if val == "*" {
				allowedOrigins = []string{"*"}
				break
			}
		}
	}

	allowedOriginRegex := make([]regexp.Regexp, 0, len(allowedOrigins))
	for _, ao := range allowedOrigins {
		parts := strings.Split(ao, ":")
		str := strings.TrimSpace(parts[0])
		if str == "" {
			continue
		}

		if str == "*" {
			allowedOriginRegex = append(
				allowedOriginRegex,
				*(regexp.MustCompile(".+")),
			)
			break
		}

		regStr := fmt.Sprintf(`^(http)?(https)?(:\/\/)?(.+\.)?%s(:[0-9]+)?$`, str)

		allowedOriginRegex = append(
			allowedOriginRegex,
			// Allow any port number of the specified domain
			*(regexp.MustCompile(regStr)),
		)
	}
	return allowedOriginRegex
}

// AddOptionsHandlers adds an OPTIONS handler for all routes.
// The response body will be empty for all new added handlers
func AddOptionsHandlers(routes []*web.Route) []*web.Route {
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {}
	if len(routes) == 0 {
		return []*web.Route{
			{
				Name:          "cors",
				Pattern:       "/:w*",
				Method:        http.MethodOptions,
				TrailingSlash: true,
				Handlers:      []http.HandlerFunc{dummyHandler},
			},
		}
	}

	list := make([]*web.Route, 0, len(routes))
	list = append(list, routes...)

	for _, r := range routes {
		list = append(list, &web.Route{
			Name:          fmt.Sprintf("%s-CORS", r.Name),
			Method:        http.MethodOptions,
			Pattern:       r.Pattern,
			TrailingSlash: true,
			Handlers:      []http.HandlerFunc{dummyHandler},
		})
	}
	return list
}

// Middleware allows the user to use this middleware without the web
func Middleware(allowedOriginRegex []regexp.Regexp, corsTimeout, allowedMethods, allowedHeaders string) web.Middleware {
	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		reqOrigin := getReqOrigin(req)
		allowed := allowedOrigin(reqOrigin, allowedOriginRegex)

		if !allowed {
			// If CORS fails, the corresponding headers are not set.
			// But execution is allowed.
			next(rw, req)
			return
		}

		// Set the appropriate response headers needed for CORS
		rw.Header().Set(headerOrigin, reqOrigin)
		rw.Header().Set(headerAccessControlAge, corsTimeout)
		rw.Header().Set(headerCreds, "true")
		rw.Header().Set(headerMethods, allowedMethods)
		rw.Header().Set(headerAllowHeaders, allowedHeaders+req.Header.Get(headerReqHeaders))

		if req.Method == http.MethodOptions {
			web.SendHeader(rw, http.StatusOK)
			return
		}
		next(rw, req)
	}
}

// CORS is a single CORS middleware that can be applied to the entire application at once
func CORS(cfg *Config) web.Middleware {
	if cfg == nil {
		cfg = new(Config)
		// 30 minutes
		cfg.TimeoutSecs = 30 * 60
	}

	allowedOrigins := cfg.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = allowedDomains()
	}

	allowedOriginRegex := allowedOriginsRegex(allowedOrigins...)
	allowedmethods := allowedMethods(cfg.Routes)
	allowedHeaders := allowedHeaders(cfg.AllowedHeaders)
	corsTimeout := fmt.Sprintf("%d", cfg.TimeoutSecs)

	return Middleware(
		allowedOriginRegex,
		corsTimeout,
		allowedmethods,
		allowedHeaders,
	)
}
