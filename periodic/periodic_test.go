package periodic_test

import (
	"os"
	"testing"
	"time"

	"github.com/tedsuo/ifrit"

	"github.com/ankeesler/wildcat-countdown/periodic"
)

func TestRun(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Millisecond*100, callback)

	assertInterval(t, p, time.Millisecond*100)

	proc := ifrit.Invoke(p)

	time.Sleep(time.Millisecond * 250)

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}

	rxTimeout(t, 2, called)
}

func TestSetInterval(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Second, callback)

	assertInterval(t, p, time.Second)

	proc := ifrit.Invoke(p)

	p.SetInterval(time.Millisecond * 100)

	time.Sleep(time.Millisecond * 250)

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}

	rxTimeout(t, 2, called)
}

func TestSetIntervalBetweenExpire(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Second*3, callback)

	assertInterval(t, p, time.Second*3)

	proc := ifrit.Invoke(p)

	time.Sleep(time.Millisecond * 500)

	p.SetInterval(time.Second)

	chanEmpty(t, called)

	time.Sleep(time.Millisecond * 750)

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}

	rxTimeout(t, 1, called)
}

func TestSetIntervalAfterExpire(t *testing.T) {
	called := make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	p := periodic.New(time.Second, callback)

	assertInterval(t, p, time.Second)

	proc := ifrit.Invoke(p)

	time.Sleep(time.Millisecond * 300)

	p.SetInterval(time.Millisecond * 250)

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}

	rxTimeout(t, 1, called)
}

func TestBeforeRun(t *testing.T) {
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

func assertInterval(
	t *testing.T,
	p *periodic.Periodic,
	expected time.Duration) {
	actual := p.GetInterval()
	if actual != expected {
		t.Errorf("wanted %s, got %s", expected.String(), actual.String())
	}
}
