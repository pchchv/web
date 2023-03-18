package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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
