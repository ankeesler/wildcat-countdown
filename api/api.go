// Package api provides the web service functionality for the wildcat-countdown app.
// At a cursory level, it spins up a web server and registers HTTP handlers for that
// server.
package api

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

//go:generate mockgen -destination mock_api/mock_api.go github.com/ankeesler/wildcat-countdown/api IntervalSetter

// IntervalSetter is an object that can handle the setting of an interval.
type IntervalSetter interface {
	// SetInterval should set the provided interval on the object.
	SetInterval(interval time.Duration) error
}

// API is a object that can run the wildcat-countdown web service on a net.Listener.
type API struct {
	listener       net.Listener
	intervalSetter IntervalSetter
}

// New returns an instance of an API configured with a net.Listener on which to run
// its service.
func New(listener net.Listener, intervalSetter IntervalSetter) *API {
	return &API{listener: listener, intervalSetter: intervalSetter}
}

// Start will simply register the necessary handlers, start the server, and return
// asynchronously. If there is no error, then the server is running happily. The
// errChan argument is filled with the error that the server returns if it returns.
func (a *API) Start(errChan chan<- error) error {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("Go 'Cats!\n"))
		})
		mux.HandleFunc("/api/interval", a.handleInterval)
		errChan <- http.Serve(a.listener, mux)
	}()
	return nil
}

func (a *API) handleInterval(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "could not read request body:", err)
		return
	}

	interval, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "could not convert payload to unsigned integer:", err)
		return
	}

	a.intervalSetter.SetInterval(time.Duration(interval))
	w.WriteHeader(http.StatusNoContent)
}
