// Package runner provides a dependency injection point for the main functionality
// of the wildcat-countdown app.
package runner

import (
	"io"
	"time"
)

//go:generate mockgen -destination mock_runner/mock_runner.go github.com/ankeesler/wildcat-countdown/runner API,Periodic

// API is an interface to describe a type that can spin up a service.
type API interface {
	// Start is called when the API should begin to run.
	Start() error
}

// Periodic is an interface to describe an object that will kick off a periodic function
// to be run after every provided time interval.
type Periodic interface {
	// Start is called when the Periodic object should start calling the provided function
	// after every provided time interval.
	Start(time.Duration, func()) error
}

// Runner provides a dependency injection point for the main functionality of the
// wildcat-countdown app.
type Runner struct {
	api      API
	periodic Periodic
}

// New creates a new Runner to be used for running wildcat-countdown app
// functionality.
func New(api API, periodic Periodic) *Runner {
	return &Runner{api: api, periodic: periodic}
}

// Run will kick off the Runner's functionality. It will call a periodic function
// after every interval and it will write output to the out io.Writer argument.
func (r *Runner) Run(interval time.Duration, out io.Writer) error {
	if err := r.api.Start(); err != nil {
		return err
	}

	if err := r.periodic.Start(interval, func() {}); err != nil {
		return err
	}

	return nil
}
