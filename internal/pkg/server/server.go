package server

import (
	"net/http"
	"time"
)

const (
	readTimeout       = 10 * time.Second
	readHeaderTimeout = 2 * time.Second
	writeTimeout      = 5 * time.Second
	idleTimeout       = 30 * time.Second
)

func DefaultServer() *http.Server {
	return &http.Server{
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}
