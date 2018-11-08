package periodic

import "time"

// Periodic is a object that can call a function after every certain time interval
// has passed.
type Periodic struct {
}

// New instantiates a new Periodic.
func New() *Periodic {
	return &Periodic{}
}

// Start begins the periodic calling of a function after every time interval.
func (p *Periodic) Start(interval time.Duration, f func()) error {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			f()
		}
	}()
	return nil
}
