package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pchchv/web"
	"github.com/pchchv/web/extensions/sse"
	"github.com/pchchv/web/middleware/accesslog"
	"github.com/pchchv/web/middleware/cors"
)

var lastModified = time.Now().Format(http.TimeFormat)

func chain(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("chained", "true")
}

func routegroupMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("routegroup", "true")
	next(w, r)
}

// errLogger is a middleware which will log all errors returned/set by a handler
func errLogger(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)

	err := web.GetError(r)
	if err != nil {
		// log only server errors
		if web.ResponseStatus(w) > 499 {
			log.Println("errorLogger:", err.Error())
		}
	}
}

func getRoutes(sse *sse.SSE) []*web.Route {
	return []*web.Route{
		{
			Name:          "root",
			Method:        http.MethodGet,
			Pattern:       "/",
			Handlers:      []http.HandlerFunc{HomeHandler},
			TrailingSlash: true,
		},
		{
			Name:          "matchall",
			Method:        http.MethodGet,
			Pattern:       "/matchall/:wildcard*",
			Handlers:      []http.HandlerFunc{ParamHandler},
			TrailingSlash: true,
		},
		{
			Name:                    "api",
			Method:                  http.MethodGet,
			Pattern:                 "/api/:param",
			Handlers:                []http.HandlerFunc{chain, ParamHandler},
			TrailingSlash:           true,
			FallThroughPostResponse: true,
		},
		{
			Name:          "invalidjson",
			Method:        http.MethodGet,
			Pattern:       "/invalidjson",
			Handlers:      []http.HandlerFunc{InvalidJSONHandler},
			TrailingSlash: true,
		},
		{
			Name:          "error-setter",
			Method:        http.MethodGet,
			Pattern:       "/error-setter",
			Handlers:      []http.HandlerFunc{ErrorSetterHandler},
			TrailingSlash: true,
		},
		{
			Name:          "original-responsewriter",
			Method:        http.MethodGet,
			Pattern:       "/original-responsewriter",
			Handlers:      []http.HandlerFunc{OriginalResponseWriterHandler},
			TrailingSlash: true,
		},
		{
			Name:          "static",
			Method:        http.MethodGet,
			Pattern:       "/static/:w*",
			Handlers:      []http.HandlerFunc{StaticFilesHandler},
			TrailingSlash: true,
		},
		{
			Name:          "sse",
			Method:        http.MethodGet,
			Pattern:       "/sse/:clientID",
			Handlers:      []http.HandlerFunc{SSEHandler(sse)},
			TrailingSlash: true,
		},
	}
}

func setup() (*web.Router, *sse.SSE) {
	port := strings.TrimSpace(os.Getenv("HTTP_PORT"))
	if port == "" {
		port = "8080"
	}
	cfg := &web.Config{
		Host:         "",
		Port:         port,
		HTTPSPort:    "9595",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 1 * time.Hour,
		CertFile:     "./certs/localhost.crt",
		KeyFile:      "./certs/localhost.decrypted.key",
	}

	web.GlobalLoggerConfig(
		nil, nil,
		web.LogCfgDisableDebug,
	)

	routeGroup := web.NewRouteGroup("/v7.0.0", false)
	routeGroup.Add(web.Route{
		Name:     "router-group-prefix-v7.0.0_api",
		Method:   http.MethodGet,
		Pattern:  "/api/:param",
		Handlers: []http.HandlerFunc{chain, ParamHandler},
	})
	routeGroup.Use(routegroupMiddleware)

	sseService := sse.New()
	sseService.OnRemoveClient = func(ctx context.Context, clientID string, count int) {
		log.Printf("\nClient %q removed, active client(s): %d\n", clientID, count)
	}
	sseService.OnCreateClient = func(ctx context.Context, client *sse.Client, count int) {
		log.Printf("\nClient %q added, active client(s): %d\n", client.ID, count)
	}

	routes := getRoutes(sseService)
	routes = append(routes, routeGroup.Routes()...)

	router := web.NewRouter(cfg, routes...)
	router.UseOnSpecialHandlers(accesslog.AccessLog)
	router.Use(
		errLogger,
		cors.CORS(nil),
		accesslog.AccessLog,
	)

	return router, sseService
}

func main() {
	router, sseService := setup()
	clients := []*sse.Client{}
	sseService.OnCreateClient = func(ctx context.Context, client *sse.Client, count int) {
		clients = append(clients, client)
	}
	// broadcast server time to all SSE listeners
	go func() {
		retry := time.Millisecond * 500
		for {
			now := time.Now().Format(time.RFC1123Z)
			sseService.Broadcast(sse.Message{
				Data:  now + fmt.Sprintf(" (%d)", sseService.ActiveClients()),
				Retry: retry,
			})
			time.Sleep(time.Second)
		}
	}()

	go router.StartHTTPS()
	router.Start()
}
