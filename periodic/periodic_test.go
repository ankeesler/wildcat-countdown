package periodic_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/periodic"
)

func TestStart(t *testing.T) {
	lock := sync.Mutex{}
	called := 0
	callback := func() {
		lock.Lock()
		defer lock.Unlock()
		called++
	}

	p := periodic.New(time.Millisecond*100, callback)

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 250)

	lock.Lock()
	got := called
	lock.Unlock()
	if got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}
