package accesslog

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pchchv/web"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`hello`))
}

func setup(port string) (*web.Router, error) {
	cfg := &web.Config{
		Port:            "9696",
		ReadTimeout:     time.Second * 1,
		WriteTimeout:    time.Second * 1,
		ShutdownTimeout: time.Second * 10,
		CertFile:        "tests/ssl/server.crt",
		KeyFile:         "tests/ssl/server.key",
	}
	router := web.NewRouter(cfg, &web.Route{
		Name:     "hello",
		Pattern:  "/hello",
		Method:   http.MethodGet,
		Handlers: []http.HandlerFunc{handler},
	})
	return router, nil
}

func TestAccessLog(t *testing.T) {
	stdout := bytes.NewBuffer([]byte(``))
	stderr := bytes.NewBuffer([]byte(``))
	web.GlobalLoggerConfig(stdout, stderr)
	port := "9696"
	router, err := setup(port)
	if err != nil {
		t.Error(err.Error())
		return
	}
	router.Use(AccessLog)
	router.SetupMiddleware()

	url := fmt.Sprintf("http://localhost:%s/hello", port)
	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	router.ServeHTTP(w, req)

	parts := strings.Split(stdout.String(), " ")
	if len(parts) != 7 {
		t.Errorf(
			"Expected log to have %d parts, got %d",
			7,
			len(parts),
		)
		return
	}

	if parts[0] != "Info" {
		t.Errorf("expected log type 'Info', got '%s'", parts[0])
	}

	if parts[3] != http.MethodGet {
		t.Errorf("expected HTTP method %s, got %s", http.MethodGet, parts[3])
	}

	if parts[4] != url {
		t.Errorf("expected HTTP full URL '%s', got '%s'", url, parts[4])
	}

	if parts[6][0:3] != "200" {
		t.Errorf("expected HTTP status code '%d', got '%s'", http.StatusOK, parts[6][0:3])
	}
}