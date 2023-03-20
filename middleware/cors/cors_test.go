package cors

import (
	"net/http"
	"time"

	"github.com/pchchv/web"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`hello`))
}

func getRoutes() []*web.Route {
	return []*web.Route{
		{
			Name:     "hello",
			Pattern:  "/hello",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{handler},
		},
	}
}

func setup(port string, routes []*web.Route) (*web.Router, error) {
	cfg := &web.Config{
		Port:            "9696",
		ReadTimeout:     time.Second * 1,
		WriteTimeout:    time.Second * 1,
		ShutdownTimeout: time.Second * 10,
		CertFile:        "tests/ssl/server.crt",
		KeyFile:         "tests/ssl/server.key",
	}
	router := web.NewRouter(cfg, routes...)

	return router, nil
}
