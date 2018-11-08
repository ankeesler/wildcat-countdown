package periodic

import "time"

// Periodic is a object that can call a function after every certain time interval
// has passed.
type Periodic struct {
	interval time.Duration
	callback func()
}

// New instantiates a new Periodic.
func New(interval time.Duration, callback func()) *Periodic {
	return &Periodic{interval: interval, callback: callback}
}

// Start begins the periodic calling of a function after every time interval.
func (p *Periodic) Start() error {
	go func() {
		ticker := time.NewTicker(p.interval)
		for {
			<-ticker.C
			p.callback()
		}
	}()
	return nil
}
