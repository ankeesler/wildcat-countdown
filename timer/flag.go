package timer

import "sync"

// flag is a simple concurrent boolean that can be read or written.
type flag struct {
	rwMutex sync.RWMutex
	wg      sync.WaitGroup
	value   bool
}

func (c *flag) get() bool {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.value
}

func (c *flag) set(value bool) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	c.value = value
	if value {
		c.wg.Add(1)
	} else {
		c.wg.Add(-1)
	}
}

// wait until the flag value is false
func (c *flag) wait() {
	c.wg.Wait()
}
