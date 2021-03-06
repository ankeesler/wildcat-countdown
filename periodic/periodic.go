// Package periodic contains functionality to call a callback after an interval
// expires. It uses a per-golang runtime local timer in addition to a persistent
// storage client to do this. The persistent storage client is used so that the
// periodic timer can be used across reboots.
package periodic

import (
	"errors"
	"os"
	"sync/atomic"
	"time"
)

// Periodic is a object that can call a function after every certain time interval
// has passed. Note that the callback for a Periodic is called on a separate
// goroutine from the one on which the object was instantiated.
type Periodic struct {
	interval time.Duration
	callback func()

	started int32

	resetChan chan time.Duration
}

// New instantiates a new Periodic.
func New(interval time.Duration, callback func()) *Periodic {
	return &Periodic{
		interval:  interval,
		callback:  callback,
		resetChan: make(chan time.Duration),
	}
}

// Run begins the periodic calling of a function after every time interval.
func (p *Periodic) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	close(ready)
	atomic.StoreInt32(&p.started, 1)

	timer := time.NewTimer(p.interval)
	startTime := time.Now()
	for {
		select {
		case <-timer.C:
			p.callback()
			timer.Reset(p.interval)

		case newInterval := <-p.resetChan:
			if !timer.Stop() {
				<-timer.C
			}

			elapsed := time.Now().Sub(startTime)
			if elapsed > newInterval {
				p.callback()
				timer.Reset(newInterval)
			} else {
				timer.Reset(newInterval - elapsed)
			}
			p.interval = newInterval

		case <-signals:
			return nil
		}
		startTime = time.Now()
	}
}

// SetInterval resets this periodic function timer to the provided interval.
//
// If the new interval would have experied by now, the callback is fired, and the
// timer is reset.
// Else, the next time the timer will fire is the provided duration plus the start
// time, and then the timer will be reset to the provided duration.
//
// This function must be called after Start().
//
// This function is thread-safe.
func (p *Periodic) SetInterval(interval time.Duration) error {
	if atomic.LoadInt32(&p.started) == 0 {
		return errors.New("must call Start() first")
	}
	p.resetChan <- interval
	return nil
}

// GetInterval gets the current interval for this Periodic.
func (p *Periodic) GetInterval() time.Duration {
	return p.interval // this is racey, but I am lazy...
}
