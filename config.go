package web

import "time"

// Config is used to read the application configuration from a json file
type Config struct {
	// Host is the host on which the server is listening
	Host string `json:"host,omitempty"`
	// Port is the port number on which the server should listen to HTTP requests
	Port string `json:"port,omitempty"`

	// CertFile is the path to TLS/SSL certificate file required for HTTPS
	CertFile string `json:"certFile,omitempty"`
	// KeyFile is the path to the certificate private key file
	KeyFile string `json:"keyFile,omitempty"`
	// HTTPSPort is the port number on which the server should listen to HTTP requests
	HTTPSPort string `json:"httpsPort,omitempty"`

	// ReadTimeout is the maximum length of time for which the server will read the request
	ReadTimeout time.Duration `json:"readTimeout,omitempty"`
	// WriteTimeout is the maximum time for which the server will try to respond to the request
	WriteTimeout time.Duration `json:"writeTimeout,omitempty"`

	// InsecureSkipVerify is the HTTP certificate verification
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`

	// ShutdownTimeout is the duration during which the preferential shutdown will be completed
	ShutdownTimeout time.Duration

	// ReverseMiddleware, if true,
	// will change the execution order of the middleware from the order it was added.
	// e.g. router.Use(m1,m2), m2 will be executed first if ReverseMiddleware is true
	ReverseMiddleware bool
}
