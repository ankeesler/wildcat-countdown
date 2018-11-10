package periodic_test

import (
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/periodic"
)

func TestStart(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Millisecond*100, callback)

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 250)

	rxTimeout(t, 2, called)
}

func TestSetInterval(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Second, callback)

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	p.SetInterval(time.Millisecond * 100)

	time.Sleep(time.Millisecond * 250)

	rxTimeout(t, 2, called)
}

func TestSetIntervalBetweenExpire(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Second*3, callback)

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 500)

	p.SetInterval(time.Second)

	chanEmpty(t, called)

	time.Sleep(time.Millisecond * 750)

	rxTimeout(t, 1, called)
}

func TestSetIntervalAfterExpire(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Second, callback)

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 300)

	p.SetInterval(time.Millisecond * 250)

	rxTimeout(t, 1, called)
}

func TestBeforeStart(t *testing.T) {
	p := periodic.New(time.Second, func() {})
	err := p.SetInterval(time.Millisecond)
	if err == nil {
		t.Error("expected error, but got none")
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
