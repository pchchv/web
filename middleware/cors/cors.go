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

	"github.com/pchchv/web"
)

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
