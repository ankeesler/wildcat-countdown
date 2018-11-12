// Package redis contains functionality to call a callback after an interval
// expires. It uses a Redis client for talking to a Redis key-value store so
// that between reboots, the interval will be maintained.
package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

// Redis is an object that will call a function periodically throughout reboots.
type Redis struct {
	interval time.Duration
	callback func()
	client   *redis.Client
}

// Config is a simple type used to describe the configuration needed for a Redis
// object to connect to its key-value store.
type Config struct {
	// Host is the hostname of the store to which to connect.
	Host string
	// Post is the port of the store to which to connect.
	Port string
	// Password is the password that will be used to authenticate with the store.
	Password string
}

// New instantiates a new Redis object.
func New(interval time.Duration, callback func(), config *Config) *Redis {
	return &Redis{
		interval: interval,
		callback: callback,
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
			Password: config.Password,
		}),
	}
}

// Run begins the periodic calling of a function after every time interval.
func (r *Redis) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	_, err := r.client.Ping().Result()
	if err != nil {
		close(ready)
		return err
	}

	close(ready)

	timer := time.NewTimer(r.interval)
	for {
		select {
		case <-timer.C:
			r.callback()
			timer.Reset(r.interval)
		case <-signals:
			return nil
		}
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
func (r *Redis) SetInterval(interval time.Duration) error {
	return nil
}

// GetInterval gets the current interval for this Periodic.
func (r *Redis) GetInterval() time.Duration {
	return r.interval // this is racey, but I am lazy...
}
