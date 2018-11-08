// Package runner provides a dependency injection point for the main functionality
// of the wildcat-countdown app.
package runner

import (
	"io"
)

//go:generate mockgen -destination mock_runner/mock_runner.go github.com/ankeesler/wildcat-countdown/runner API

// API is an interface to describe a type that can spin up a service.
type API interface {
	Start() error
}

// Runner provides a dependency injection point for the main functionality of the
// wildcat-countdown app.
type Runner struct {
	api API
}

// New creates a new Runner to be used for running wildcat-countdown app
// functionality.
func New(api API) *Runner {
	return &Runner{api: api}
}

// Run will kick off the Runner's functionality. It will write output to the
// out io.Writer argument.
func (r *Runner) Run(out io.Writer) error {
	if err := r.api.Start(); err != nil {
		return err
	}

	return nil
}
