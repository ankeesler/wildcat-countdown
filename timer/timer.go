// Package timer contains an object that call periodically send to a channel, even
// across runtime restarts.
package timer

import (
	"log"
	"time"

	"code.cloudfoundry.org/clock"
)

//go:generate mockgen -destination mock_timer/mock_timer.go github.com/ankeesler/wildcat-countdown/timer Store

// Store is an object that can Get()/Set() key/value pairs from a persistent store.
type Store interface {
	// Get retrieves a value associated with a key from the persistent store. If
	// there is no key in the store, it should return (nil, nil).
	Get(key string) (interface{}, error)
	// Set associated a key with a value in the persistent store.
	Set(key string, value interface{}) error
}

// Timer is an object that can deliver signals on a channel to indicate when
// a certain amount of time has passed. It uses a persistant store to function
// properly across runtime restarts. It is totally thread-safe.
type Timer struct {
	store                Store
	clock                clock.Clock
	fireTime             time.Time
	fireChan, cancelChan chan struct{}
	flag                 flag
}

// New instantiates a new Timer.
func New(store Store, clock clock.Clock) *Timer {
	return &Timer{
		store:      store,
		clock:      clock,
		fireChan:   make(chan struct{}),
		cancelChan: make(chan struct{}, 1),
	}
}

// Reset tells the timer to begin its countdown at a certain time in the future.
// The fireTime parameter is just a default - the Timer will try to overwrite it with
// the value from its store. Once the value from the store is used, it will be
// erased.
//
// Reset should be called from the same goroutine!
func (t *Timer) Reset(fireTime time.Time) error {
	if t.flag.get() {
		t.cancelChan <- struct{}{}
		t.flag.wait()
	}

	interval := t.getInterval(fireTime)
	timer := t.clock.NewTimer(interval)
	log.Printf("reseting to fire in %s (fireTime = %s)", interval, fireTime)

	t.flag.set(true)
	go func() {
		select {
		case <-timer.C():
			t.store.Set("time", -1)
			t.fireChan <- struct{}{}
		case <-t.cancelChan:
			if !timer.Stop() {
				<-timer.C()
			}
		}
		t.flag.set(false)
	}()
	return nil
}

// C returns a channel that can be waited on by a caller. The channel will be
// filled with a value when the Timer expires.
func (t *Timer) C() <-chan struct{} {
	return t.fireChan
}

func (t *Timer) getInterval(fireTime time.Time) time.Duration {
	if val, err := t.store.Get("time"); err != nil {
		log.Println("warning:", err)
	} else if val != nil {
		if interval, ok := val.(time.Duration); !ok {
			log.Printf("warning: cannot cast val (%+v) to time.Duration", val)
		} else if interval != -1 {
			return interval
		}
	}
	return fireTime.Sub(t.clock.Now())
}
