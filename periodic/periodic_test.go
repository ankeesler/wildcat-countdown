package periodic_test

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/periodic"
	"github.com/ankeesler/wildcat-countdown/periodic/mock_periodic"
	"github.com/golang/mock/gomock"
	"github.com/tedsuo/ifrit"
)

func TestRun(t *testing.T) {
	t.Parallel()

	cleanup, called, p := setup(t, time.Millisecond*100, 1)
	defer cleanup()

	assertInterval(t, p, time.Millisecond*100)

	proc := run(p)

	time.Sleep(time.Millisecond * 250)

	stop(t, proc)

	rxTimeout(t, 2, called)
}

func TestSetInterval(t *testing.T) {
	t.Parallel()

	cleanup, called, p := setup(t, time.Second, 1)
	defer cleanup()

	assertInterval(t, p, time.Second)

	proc := run(p)

	p.SetInterval(time.Millisecond * 100)

	time.Sleep(time.Millisecond * 250)

	stop(t, proc)

	rxTimeout(t, 2, called)
}

func TestSetIntervalBetweenExpire(t *testing.T) {
	t.Parallel()

	cleanup, called, p := setup(t, time.Second*3, 1)
	defer cleanup()

	assertInterval(t, p, time.Second*3)

	proc := run(p)

	time.Sleep(time.Millisecond * 500)

	p.SetInterval(time.Second)

	chanEmpty(t, called)

	time.Sleep(time.Millisecond * 750)

	stop(t, proc)

	rxTimeout(t, 1, called)
}

func TestSetIntervalAfterExpire(t *testing.T) {
	t.Parallel()

	cleanup, called, p := setup(t, time.Second, 1)
	defer cleanup()

	assertInterval(t, p, time.Second)

	proc := run(p)

	time.Sleep(time.Millisecond * 300)

	p.SetInterval(time.Millisecond * 250)

	stop(t, proc)

	rxTimeout(t, 1, called)
}

func TestSetIntervalLonger(t *testing.T) {
	t.Parallel()

	cleanup, called, p := setup(t, time.Millisecond*500, 1)
	defer cleanup()

	assertInterval(t, p, time.Millisecond*500)

	proc := run(p)

	time.Sleep(time.Millisecond * 300)

	p.SetInterval(time.Millisecond * 750)

	chanEmpty(t, called)

	time.Sleep(time.Millisecond * 300)

	chanEmpty(t, called)

	time.Sleep(time.Millisecond * 300)

	stop(t, proc)

	rxTimeout(t, 1, called)
}

func TestBeforeRun(t *testing.T) {
	t.Parallel()

	cleanup, _, p := setup(t, time.Millisecond, 0)
	defer cleanup()

	err := p.SetInterval(time.Millisecond)
	if err == nil {
		t.Error("expected error, but got none")
	}
}

func TestRestartNotExpired(t *testing.T) {
	t.Parallel()

	cleanup, called, p := setup(t, time.Second, 2)
	defer cleanup()

	proc := run(p)

	time.Sleep(time.Millisecond * 500)

	chanEmpty(t, called)

	stop(t, proc)

	time.Sleep(time.Millisecond * 250)

	proc = run(p)

	chanEmpty(t, called)

	time.Sleep(time.Millisecond * 250)

	stop(t, proc)

	rxTimeout(t, 1, called)
}

func TestRestartExpired(t *testing.T) {
	t.Parallel()

	cleanup, called, p := setup(t, time.Second, 2)
	defer cleanup()

	proc := run(p)

	time.Sleep(time.Millisecond * 500)

	chanEmpty(t, called)

	stop(t, proc)

	time.Sleep(time.Millisecond * 750)

	proc = run(p)

	stop(t, proc)

	rxTimeout(t, 1, called)
}

func setup(t *testing.T, interval time.Duration, mockSetup int) (cleanup func(), called chan struct{}, p *periodic.Periodic) {
	ctrl := gomock.NewController(t)

	storeClient := mock_periodic.NewMockStoreClient(ctrl)
	switch mockSetup {
	case 0:
		// don't set anything up
	case 1:
		storeClient.EXPECT().Get("time").Return(nil, nil)
	case 2:
		storeClient.EXPECT().Get("time").Return(nil, nil)
		storeClient.EXPECT().Get("time").Return(time.Now().Add(time.Millisecond*500), nil)
	}

	called = make(chan struct{}, 10)
	callback := func() {
		called <- struct{}{}
	}

	cleanup = func() {
		ctrl.Finish()
		close(called)
	}

	p = periodic.New(storeClient, interval, callback)

	return
}

func rxTimeout(t *testing.T, count int, c <-chan struct{}) {
	got := 0
	for count > 0 {
		select {
		case <-time.After(time.Second):
			t.Fatalf("expected %d, got %d (from line %d)", count, got, callerLine())
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
		t.Fatalf("expected channel to be empty (from line %d)", callerLine())
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
		t.Errorf("wanted %s, got %s (from line %d)",
			expected.String(),
			actual.String(),
			callerLine())
	}
}

func callerLine() int {
	_, _, line, ok := runtime.Caller(2)
	if !ok {
		return -1
	}

	return line
}

func run(runner ifrit.Runner) ifrit.Process {
	proc := ifrit.Invoke(runner)
	return proc
}

func stop(t *testing.T, proc ifrit.Process) {
	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}
}
