package timer_test

import (
	"errors"
	"runtime"
	"testing"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/ankeesler/wildcat-countdown/timer"
	"github.com/ankeesler/wildcat-countdown/timer/mock_timer"
	"github.com/golang/mock/gomock"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_timer.NewMockStore(ctrl)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).Return(nil)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).Return(nil)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).Return(nil)

	clock := fakeclock.NewFakeClock(time.Now())

	timeout := func() time.Time {
		return clock.Now().Add(time.Millisecond * 500)
	}

	timer := timer.New(store, clock)
	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}

	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())
	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())

	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}

	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())
	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())

	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}

	clock.Increment(time.Millisecond * 600)
	rxTimeout(t, 1, timer.C())

	clock.Increment(time.Millisecond * 600)
	chanEmpty(t, timer.C())
}

func TestStore(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_timer.NewMockStore(ctrl)
	store.EXPECT().Get("time").Return(time.Millisecond*250, nil)
	store.EXPECT().Set("time", -1).Return(nil)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).Return(nil)

	clock := fakeclock.NewFakeClock(time.Now())

	timeout := func() time.Time {
		return clock.Now().Add(time.Millisecond * 500)
	}

	timer := timer.New(store, clock)
	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}

	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())
	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())

	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}

	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())
	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())
}

func TestStoreFailure(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_timer.NewMockStore(ctrl)
	store.EXPECT().Get("time").Return(nil, errors.New("some error"))
	store.EXPECT().Set("time", -1).Return(nil)

	clock := fakeclock.NewFakeClock(time.Now())

	timeout := func() time.Time {
		return clock.Now().Add(time.Millisecond * 500)
	}

	timer := timer.New(store, clock)
	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}

	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())
	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())
}

func TestReset(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_timer.NewMockStore(ctrl)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).Return(nil)

	clock := fakeclock.NewFakeClock(time.Now())

	timer := timer.New(store, clock)
	if err := timer.Reset(clock.Now().Add(time.Millisecond * 500)); err != nil {
		t.Fatal(err)
	}

	beforeIncrement := clock.Now()
	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())

	if err := timer.Reset(beforeIncrement.Add(time.Millisecond * 750)); err != nil {
		t.Fatal(err)
	}
	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())
	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())
}

func TestFireInPast(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_timer.NewMockStore(ctrl)
	store.EXPECT().Get("time").Return(time.Now().Add(-time.Second).Sub(time.Now()), nil)
	store.EXPECT().Set("time", -1).Return(nil)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).Return(nil)

	clock := fakeclock.NewFakeClock(time.Now())

	timeout := func() time.Time {
		return clock.Now().Add(time.Millisecond * 500)
	}

	timer := timer.New(store, clock)
	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}

	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())
	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())

	timeout = func() time.Time {
		return time.Now().Add(-time.Minute)
	}

	if err := timer.Reset(timeout()); err != nil {
		t.Fatal(err)
	}
	clock.Increment(time.Millisecond * 300)
	rxTimeout(t, 1, timer.C())
	clock.Increment(time.Millisecond * 300)
	chanEmpty(t, timer.C())
}

func TestOverlap(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// we may or may not get the first call to store.Set(), so we say it
	// should be called 0 or 1 times
	store := mock_timer.NewMockStore(ctrl)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).MinTimes(0).MaxTimes(1).Return(nil)
	store.EXPECT().Get("time").Return(nil, nil)
	store.EXPECT().Set("time", -1).MinTimes(0).MaxTimes(1).Return(nil)

	clock := fakeclock.NewFakeClock(time.Now())

	timer := timer.New(store, clock)

	timer.Reset(clock.Now().Add(time.Millisecond * 500))
	clock.Increment(time.Millisecond * 600)
	// don't drain timer.C()

	timer.Reset(clock.Now().Add(time.Millisecond * 100))
	clock.Increment(time.Millisecond * 200)

	rxTimeout(t, 1, timer.C())
	chanEmpty(t, timer.C())
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

func callerLine() int {
	_, _, line, ok := runtime.Caller(2)
	if !ok {
		return -1
	}

	return line
}
