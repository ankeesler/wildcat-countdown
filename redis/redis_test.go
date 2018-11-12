package redis_test

import (
	"os"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/redis"
	"github.com/tedsuo/ifrit"
)

func TestBasic(t *testing.T) {
	config := redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}
	r := redis.New(time.Millisecond*500, callback, &config)

	proc := ifrit.Invoke(r)
	<-proc.Ready()

	time.Sleep(time.Millisecond * 250)

	chanEmpty(t, called)

	time.Sleep(time.Millisecond * 350)

	rxTimeout(t, 1, called)

	time.Sleep(time.Second + time.Millisecond*10)

	rxTimeout(t, 2, called)

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestSetIntervalShorter(t *testing.T) {
	config := redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}
	r := redis.New(time.Second, callback, &config)

	proc := ifrit.Invoke(r)
	<-proc.Ready()

	time.Sleep(time.Millisecond * 350)

	chanEmpty(t, called)

	r.SetInterval(time.Millisecond * 250)

	rxTimeout(t, 1, called)

	time.Sleep(time.Millisecond*500 + time.Millisecond*10)

	rxTimeout(t, 2, called)

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}
}

func rxTimeout(t *testing.T, count int, c <-chan struct{}) {
	got := 0
	for count > 0 {
		select {
		case <-time.After(time.Second):
			t.Fatalf("expected %d, got %d", count, got)
		case <-c:
			got++
			if got == count {
				return
			}
		}
	}
}

func chanEmpty(t *testing.T, c <-chan struct{}) {
	select {
	case <-c:
		t.Fatal("expected channel to be empty")
	default:
		// yay!
	}
}

func assertInterval(
	t *testing.T,
	r *redis.Redis,
	expected time.Duration) {
	actual := r.GetInterval()
	if actual != expected {
		t.Errorf("wanted %s, got %s", expected.String(), actual.String())
	}
}
