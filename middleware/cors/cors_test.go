package cors

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestCORSEmptyconfig(t *testing.T) {
	port := "9696"
	routes := getRoutes()
	routes = append(routes, AddOptionsHandlers(nil)...)
	router, err := setup(port, routes)
	if err != nil {
		t.Error(err.Error())
		return
	}
	router.Use(CORS(&Config{TimeoutSecs: 50}))
	router.SetupMiddleware()

	url := fmt.Sprintf("http://localhost:%s/hello", port)
	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	router.ServeHTTP(w, req)
	body, _ := ioutil.ReadAll(w.Body)
	str := string(body)
	if str != "hello" {
		t.Errorf(
			"Expected body '%s', got '%s'",
			"hello",
			str,
		)
	}

	if w.Header().Get(headerMethods) != defaultAllowMethods {
		t.Errorf(
			"Expected header %s to be '%s', got '%s'",
			headerMethods,
			defaultAllowMethods,
			w.Header().Get(headerMethods),
		)
	}
	if w.Header().Get(headerCreds) != "true" {
		t.Errorf(
			"Expected header %s to be 'true', got '%s'",
			headerCreds,
			w.Header().Get(headerCreds),
		)
	}
	if w.Header().Get(headerAccessControlAge) != "50" {
		t.Errorf(
			"Expected '%s' to be '50', got '%s'",
			headerAccessControlAge,
			w.Header().Get(headerAccessControlAge),
		)
	}

	if w.Header().Get(headerAllowHeaders) != allowHeaders {
		t.Errorf(
			"Expected '%s' to be '%s', got '%s'",
			headerAllowHeaders,
			allowHeaders,
			w.Header().Get(headerAllowHeaders),
		)
	}

	// check OPTIONS method
	w = httptest.NewRecorder()
	req = httptest.NewRequest(
		http.MethodOptions,
		url,
		nil,
	)
	router.ServeHTTP(w, req)
	body, _ = ioutil.ReadAll(w.Body)
	str = string(body)
	if str != "" {
		t.Errorf(
			"Expected empty body, got '%s'",
			str,
		)
	}
}
