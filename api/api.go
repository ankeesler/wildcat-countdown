// Package api provides the web service functionality for the wildcat-countdown app.
// At a cursory level, it spins up a web server and registers HTTP handlers for that
// server.
package api

import (
	"log"
	"net"
	"net/http"
)

// API is a object that can run the wildcat-countdown web service on a net.Listener.
type API struct {
	listener net.Listener
}

// New returns an instance of an API configured with a net.Listener on which to run
// its service.
func New(listener net.Listener) *API {
	return &API{listener: listener}
}

// Start will simply start the server, register the necessary handlers, and return
// asynchronously. If there is no error, then the server is running happily.
func (a *API) Start() error {
	go func() {
		log.Fatal(http.Serve(a.listener, nil))
	}()
	return nil
}
