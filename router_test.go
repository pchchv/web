package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testLogger struct {
	out bytes.Buffer
}

func (tl *testLogger) Debug(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Info(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Warn(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Error(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Fatal(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}

func testTable() []struct {
	Name      string
	TestType  string
	Path      string
	Method    string
	Want      interface{}
	WantErr   bool
	Err       error
	ParamKeys []string
	Params    []string
	Body      io.Reader
} {
	return []struct {
		Name      string
		TestType  string
		Path      string
		Method    string
		Want      interface{}
		WantErr   bool
		Err       error
		ParamKeys []string
		Params    []string
		Body      io.Reader
	}{
		{
			Name:     "Check root path without params",
			TestType: "checkpath",
			Path:     "/",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check root path without params - duplicate",
			TestType: "checkpath",
			Path:     "/",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - 1",
			TestType: "checkpath",
			Path:     "/a",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - 2",
			TestType: "checkpath",
			Path:     "/a/b",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - 3",
			TestType: "checkpath",
			Path:     "/a/b/-/c",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - 4",
			TestType: "checkpath",
			Path:     "/a/b/-/c/~/d",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - 5",
			TestType: "checkpath",
			Path:     "/a/b/-/c/~/d/./e",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - 5",
			TestType: "checkpathnotrailingslash",
			Path:     "/a/b/-/c/~/d/./e/notrail",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - OPTION",
			TestType: "checkpathnotrailingslash",
			Path:     "/a/b/-/c/~/d/./e",
			Method:   http.MethodOptions,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - HEAD",
			TestType: "checkpathnotrailingslash",
			Path:     "/a/b/-/c/~/d/./e",
			Method:   http.MethodHead,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - POST",
			TestType: "checkpathnotrailingslash",
			Path:     "/a/b/-/c/~/d/./e",
			Method:   http.MethodPost,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - PUT",
			TestType: "checkpathnotrailingslash",
			Path:     "/a/b/-/c/~/d/./e",
			Method:   http.MethodPut,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - PATCH",
			TestType: "checkpathnotrailingslash",
			Path:     "/a/b/-/c/~/d/./e",
			Method:   http.MethodPatch,
			WantErr:  false,
		},
		{
			Name:     "Check nested path without params - DELETE",
			TestType: "checkpathnotrailingslash",
			Path:     "/a/b/-/c/~/d/./e",
			Method:   http.MethodDelete,
			WantErr:  false,
		},
		{
			Name:      "Check with params - 1",
			TestType:  "checkparams",
			Path:      "/params/:a",
			Method:    http.MethodGet,
			ParamKeys: []string{"a"},
			Params:    []string{"hello"},
			WantErr:   false,
		},
		{
			Name:      "Check with params - 2",
			TestType:  "checkparams",
			Path:      "/params/:a/:b",
			Method:    http.MethodGet,
			ParamKeys: []string{"a", "b"},
			Params:    []string{"hello", "world"},
			WantErr:   false,
		},
		{
			Name:      "Check with wildcard",
			TestType:  "checkparams",
			Path:      "/wildcard/:a*",
			Method:    http.MethodGet,
			ParamKeys: []string{"a"},
			Params:    []string{"w1/hello/world/hi/there"},
			WantErr:   false,
		},
		{
			Name:      "Check with wildcard - 2",
			TestType:  "checkparams",
			Path:      "/wildcard2/:a*",
			Method:    http.MethodGet,
			ParamKeys: []string{"a"},
			Params:    []string{"w2/hello/world/hi/there/-/~/./again"},
			WantErr:   false,
		},
		{
			Name:      "Check with wildcard - 3",
			TestType:  "widlcardwithouttrailingslash",
			Path:      "/wildcard3/:a*",
			Method:    http.MethodGet,
			ParamKeys: []string{"a"},
			Params:    []string{"w3/hello/world/hi/there/-/~/./again/"},
			WantErr:   true,
		},
		{
			Name:      "Check with wildcard - 4",
			TestType:  "widlcardwithouttrailingslash",
			Path:      "/wildcard3/:a*",
			Method:    http.MethodGet,
			ParamKeys: []string{"a"},
			Params:    []string{"w4/hello/world/hi/there/-/~/./again"},
			WantErr:   false,
		},
		{
			Name:     "Check not implemented",
			TestType: "notimplemented",
			Path:     "/notimplemented",
			Method:   "HELLO",
			WantErr:  false,
		},
		{
			Name:     "Check not found",
			TestType: "notfound",
			Path:     "/notfound",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check chaining",
			TestType: "chaining",
			Path:     "/chained",
			Method:   http.MethodGet,
			WantErr:  false,
		},
		{
			Name:     "Check chaining",
			TestType: "chaining-nofallthrough",
			Path:     "/chained/nofallthrough",
			Method:   http.MethodGet,
			WantErr:  false,
		},
	}
}

func chainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("chained", "true")
}

func chainNoFallthroughHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("chained", "true")
	_, _ = w.Write([]byte(`yay, blocked!`))
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	wctx := Context(r)
	params := wctx.Params()
	R200(
		w,
		map[string]interface{}{
			"path":   r.URL.Path,
			"params": params,
		},
	)
}

func getRoutes(t *testing.T) []*Route {
	t.Helper()
	list := testTable()
	rr := make([]*Route, 0, len(list))
	for _, l := range list {
		switch l.TestType {
		case "checkpath", "checkparams", "checkparamswildcard":
			{
				rr = append(rr,
					&Route{
						Name:                    l.Name,
						Method:                  l.Method,
						Pattern:                 l.Path,
						TrailingSlash:           true,
						FallThroughPostResponse: false,
						Handlers:                []http.HandlerFunc{successHandler},
					},
				)
			}
		case "checkpathnotrailingslash", "widlcardwithouttrailingslash":
			{
				rr = append(rr,
					&Route{
						Name:                    l.Name,
						Method:                  l.Method,
						Pattern:                 l.Path,
						TrailingSlash:           false,
						FallThroughPostResponse: false,
						Handlers:                []http.HandlerFunc{successHandler},
					},
				)

			}

		case "chaining":
			{
				rr = append(
					rr,
					&Route{
						Name:                    l.Name,
						Method:                  l.Method,
						Pattern:                 l.Path,
						TrailingSlash:           false,
						FallThroughPostResponse: false,
						Handlers:                []http.HandlerFunc{chainHandler, successHandler},
					},
				)
			}
		case "chaining-nofallthrough":
			{
				{
					rr = append(
						rr,
						&Route{
							Name:                    l.Name,
							Method:                  l.Method,
							Pattern:                 l.Path,
							TrailingSlash:           false,
							FallThroughPostResponse: false,
							Handlers:                []http.HandlerFunc{chainHandler, chainNoFallthroughHandler, successHandler},
						},
					)
				}
			}
		}
	}
	return rr
}

func checkPath(req *http.Request, resp *httptest.ResponseRecorder) error {
	want := req.URL.EscapedPath()
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response, '%s'", err.Error())
	}

	body := struct {
		Data struct {
			Path   string
			Params map[string]string
		}
	}{}
	err = json.Unmarshal(rbody, &body)
	if err != nil {
		return fmt.Errorf(
			"json decode failed '%s', got response: '%s'",
			err.Error(),
			string(rbody),
		)
	}

	if want != body.Data.Path {
		return fmt.Errorf("wanted URI path '%s', got '%s'", want, body.Data.Path)
	}

	return nil
}

func checkPathWildCard(req *http.Request, resp *httptest.ResponseRecorder) error {
	want := req.URL.EscapedPath()
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response, '%s'", err.Error())
	}

	body := struct {
		Data struct {
			Path   string
			Params map[string]string
		}
	}{}
	err = json.Unmarshal(rbody, &body)
	if err != nil {
		return fmt.Errorf("json decode failed '%s', got response: '%s'", err.Error(), string(rbody))
	}

	if want != body.Data.Path {
		return fmt.Errorf("wanted URI path '%s', got '%s'", want, body.Data.Path)
	}

	if len(body.Data.Params) != 1 {
		return fmt.Errorf("expected no.of params: %d, got %d. response: '%s'", 1, len(body.Data.Params), string(rbody))
	}

	wantWildcardParamValue := ""
	parts := strings.Split(want, "/")[2:]
	wantWildcardParamValue = strings.Join(parts, "/")
	if body.Data.Params["a"] != wantWildcardParamValue {
		return fmt.Errorf(
			"wildcard value\nexpected: %s\ngot: %s",
			wantWildcardParamValue,
			body.Data.Params["a"],
		)
	}

	return nil
}
