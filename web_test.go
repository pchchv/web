package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
