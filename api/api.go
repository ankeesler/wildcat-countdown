// Package api provides the web service functionality for the wildcat-countdown app.
// At a cursory level, it spins up a web server and registers HTTP handlers for that
// server.
package api

import (
	"log"
	"net"
	"net/http"
)

//go:generate mockgen -destination mock_api/mock_api.go github.com/ankeesler/wildcat-countdown/api API

// API is an interface that describes how the core object of this package functions.
type API interface {
	// Start will simply start the server, register the necessary handlers, and return
	// asynchronously. If there is no error, then the server is running happily.
	Start() error
}

// New returns an instance of the web service that this package provides.
func New(listener net.Listener) API {
	return &api{listener: listener}
}

type api struct {
	listener net.Listener
}

func (a *api) Start() error {
	go func() {
		log.Fatal(http.Serve(a.listener, nil))
	}()
	return nil
}
