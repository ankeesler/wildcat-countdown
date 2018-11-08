// Package runner provides a dependency injection point for the main functionality
// of the wildcat-countdown app.
package runner

import "log"

//go:generate mockgen -destination mock_runner/mock_runner.go github.com/ankeesler/wildcat-countdown/runner API,Periodic

// API is an interface to describe a type that can spin up a service.
type API interface {
	// Start is called when the API should begin to run. It uses the errChan argument
	// to communicate if the API exits.
	Start(errChan chan<- error) error
}

// Periodic is an interface to describe an object that will kick off a periodic function
// to be run after every provided time interval.
type Periodic interface {
	// Start is called when the Periodic object should start calling the provided
	// function.
	Start() error
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

// Run will kick off the Runner's functionality. It will Start() the API and
// then Start() the Periodic.
func (r *Runner) Run() error {
	errChan := make(chan error)
	if err := r.api.Start(errChan); err != nil {
		return err
	}

	go func() {
		if err := <-errChan; err != nil {
			log.Fatal("API errored:", err)
		}
	}()

	if err := r.periodic.Start(); err != nil {
		return err
	}

	return nil
}
