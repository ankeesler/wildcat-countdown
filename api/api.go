// Package api provides an http.Handler for the wildcat-countdown API.
package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//go:generate mockgen -destination mock_api/mock_api.go github.com/ankeesler/wildcat-countdown/api IntervalHolder

// IntervalHolder is an object that can handle the setting and getting of an interval.
type IntervalHolder interface {
	// SetInterval should set the provided interval on the object.
	SetInterval(interval time.Duration) error
	// GetInterval should get the interval from the object.
	GetInterval() time.Duration
}

// API is a object that can run the wildcat-countdown web service.
type API struct {
	intervalHolder IntervalHolder

	mux *http.ServeMux
}

// New returns an instance of an API.
func New(intervalHolder IntervalHolder) *API {
	a := &API{intervalHolder: intervalHolder}

	a.mux = http.NewServeMux()
	a.mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Go 'Cats!\n"))
	})
	a.mux.HandleFunc("/api/interval", a.handleInterval)

	return a
}

// HandleFunc hands off the mux'ing of an HTTP request to a configured http.ServeMux.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *API) handleInterval(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPut && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet {
		interval := a.intervalHolder.GetInterval()
		w.Write([]byte(fmt.Sprintf("interval = %s\n", time.Duration(interval).String())))
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

	a.intervalHolder.SetInterval(time.Duration(interval))
	w.Write([]byte(fmt.Sprintf("interval set to %s\n", time.Duration(interval).String())))
}
