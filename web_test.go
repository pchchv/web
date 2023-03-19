package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkRouter(b *testing.B) {
	GlobalLoggerConfig(nil, nil, LogCfgDisableDebug, LogCfgDisableInfo, LogCfgDisableWarn)
	t := &testing.T{}
	router, err := setup(t, "1595")
	if err != nil {
		b.Error(err)
		return
	}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/a/b/-/c/~/d/./e", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, r)
		if w.Result().StatusCode != http.StatusOK {
			b.Error("expected status 200, got", w.Result().StatusCode)
			return
		}
	}
}

func TestStart(t *testing.T) {
	t.Parallel()
	router, _ := setup(t, "9696")
	go router.Start()
	time.Sleep(time.Second * 2)
	err := router.Shutdown()
	if err != nil {
		t.Fatal(err)
	}
}
func TestStartHTTPS(t *testing.T) {
	t.Parallel()
	router, _ := setup(t, "8443")
	go router.StartHTTPS()
	time.Sleep(time.Second * 2)
	err := router.ShutdownHTTPS()
	if err != nil {
		t.Fatal(err)
	}
}
