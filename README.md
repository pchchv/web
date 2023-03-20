# **web** [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/web)](https://goreportcard.com/report/github.com/pchchv/web) [![Go Reference](https://pkg.go.dev/badge/github.com/pchchv/web.svg)](https://pkg.go.dev/github.com/pchchv/web) [![GitHub license](https://img.shields.io/github/license/pchchv/web.svg)](https://github.com/pchchv/web/blob/master/LICENSE)

Web is a minimalist router for [Go](https://golang.org) to create web applications (server-side) without third-party dependencies. Web will always be compatible with the standard Go library; HTTP handlers have the same signature as [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc).

## Router

Web has a simplified, linear path-matching router and supports [URI](https://developer.mozilla.org/en-US/docs/Glossary/URI) definition with the following patterns:

1. `/api/users` - URI with no dynamic values
2. `/api/users/:userID`
   - URI with a named parameter, `userID`
   - If TrailingSlash is true, a URI ending in '/' will be accepted, see [sample](https://github.com/pchchv/web#sample).
3. `/api/users/:misc*`
   - Named URI parameter `misc`, with a wildcard suffix '\*'
   - This matches everything after `/api/users`. e.g. `/api/users/a/b/c/d`

If there are multiple handlers corresponding to the same URI, the request will only be handled by the first encountered handler.
Refer to [sample](https://github.com/pchchv/web#sample) to see how routes are configured. You can access the named URI parameters with the `Context` function.

Note: Web Context **not** available inside special handlers.

```golang
func helloWorld(w http.ResponseWriter, r *http.Request) {
	wctx := web.Context(r)
	// URI paramaters, map[string]string
	params := wctx.Params()
	// route, the web.Route which is executing this request
	route := wctx.Route
	web.R200(
		w,
		fmt.Sprintf(
			"Route name: '%s', params: '%s'",
			route.Name,
			params,
		),
	)
}
```

## Handler chaining

Handler chaining allows to execute multiple handlers for a given route. Chaining execution can be set to run even after the handler has written a response to an HTTP request by setting `FallThroughPostResponse` to `true` (see [sample](https://github.com/pchchv/web/blob/master/cmd/main.go)).

## Middleware

Web [middlware](https://godoc.org/github.com/pchchv/web#Middleware) allows you to wrap all routes with middleware as opposed to a handler chain. The router exposes the [Use](https://godoc.org/github.com/pchchv/web#Router.Use) and [UseOnSpecialHandlers](https://godoc.org/github.com/pchchv/web#Router.UseOnSpecialHandlers) methods to add Middleware to the router.

NotFound and NotImplemented are considered `special` handlers. The `web.Context(r)` inside special handlers will return `nil`.

You can add any number of intermediate programs to the router, the execution order of the intermediate programs will be [LIFO](<https://en.wikipedia.org/wiki/Stack_(abstract_data_type)>) (Last In First Out). E.g.:

```golang
func main() {
	router.Use(accesslog.AccessLog, cors.CORS(nil))
	router.Use(<more middleware>)
}
```

First **_CorsWrap_** will be executed, then **_AccessLog_**.

## Error handling

Web context has 2 methods for [set](https://github.com/pchchv/web/blob/master/web.go) and [get](https://github.com/pchchv/web/blob/master/web.go) errors in the request context. This allows the Web to implement a single middleware where errors returned in the HTTP handler can be handled. [set error](https://github.com/pchchv/web/blob/master/cmd/main.go), [get error](https://github.com/pchchv/web/blob/master/cmd/main.go).

## Helper functions

Web provides several helper functions. When using `Send` or `SendResponse` the response is wrapped in [response struct](https://github.com/pchchv/web/blob/master/responses.go) Web and serialized as JSON.

```json
{
  "data": "<any valid JSON payload>",
  "status": "<HTTP status code, of type integer>"
}
```

Using `SendError`, the response is wrapped in [error response struct](https://github.com/pchchv/web/blob/master/responses.go) Web and serialized as JSON.

```json
{
  "errors": "<any valid JSON payload>",
  "status": "<HTTP status code, of type integer>"
}
```

## HTTPS ready

The HTTPS server can be easily started by providing a key and a cert file. You can also have both HTTP and HTTPS servers running side by side.

Start HTTPS server

```golang
cfg := &web.Config{
	Port: "80",
	HTTPSPort: "443",
	CertFile: "/path/to/certfile",
	KeyFile: "/path/to/keyfile",
}
router := web.NewRouter(cfg, routes()...)
router.StartHTTPS()
```

Starting both HTTP & HTTPS server

```golang
cfg := &web.Config{
	Port: "80",
	HTTPSPort: "443",
	CertFile: "/path/to/certfile",
	KeyFile: "/path/to/keyfile",
}

router := web.NewRouter(cfg, routes()...)
go router.StartHTTPS()
router.Start()
```

## Graceful shutdown

Graceful shutdown allows you to shut down the server without affecting live connections/clients connected to the server. Any new connection request after initiating the shutdown will be ignored.  
Sample:

```golang
func main() {
	osSig := make(chan os.Signal, 5)

	cfg := &web.Config{
		Host:            "",
		Port:            "8080",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    60 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}
	router := web.NewRouter(cfg, routes()...)

	go func() {
		<-osSig
		// Initiate HTTP server shutdown
		err := router.Shutdown()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("shutdown complete")
			os.Exit(0)
		}

		// If HTTPS server running
		err := router.ShutdownHTTPS()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("shutdown complete")
			os.Exit(0)
		}
	}()

	go func(){
		time.Sleep(time.Second*15)
		signal.Notify(osSig, os.Interrupt, syscall.SIGTERM)
	}()

	router.Start()
}
```

## Logging

Web exposes a singleton & global scoped logger variable [LOGHANDLER](https://godoc.org/github.com/pchchv/web#Logger) with which you can plug in your custom logger by implementing the [Logger](https://godoc.org/github.com/pchchv/web#Logger) interface.

### Configuring the default Logger

The default logger uses the standard Go `log.Logger` library with `os.Stdout` for debugging and information logs and `os.Stderr` for warnings, errors, fatal events as io.Writers by default. You can set io.Writer and also disable certain log types with `GlobalLoggerConfig(stdout, stderr, cfgs...)`.

## Server-Sent Events

MDN has very good documentation on what [SSE (Server-Sent Events)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) is.

## Usage

A fully functional sample is available [here] (https://github.com/pchchv/web/blob/master/cmd/main.go).