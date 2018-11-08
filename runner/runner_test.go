package runner_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/runner"
	"github.com/ankeesler/wildcat-countdown/runner/mock_runner"
	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_runner.NewMockAPI(ctrl)
	api.EXPECT().Start()

	periodic := mock_runner.NewMockPeriodic(ctrl)
	periodic.EXPECT().Start(time.Second*5, gomock.Any())

	r := runner.New(api, periodic)

	buf := bytes.NewBuffer([]byte{})
	if err := r.Run(time.Second*5, buf); err != nil {
		t.Fatal(err)
	}
}

func TestAPIStartFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedError := errors.New("some api error")
	api := mock_runner.NewMockAPI(ctrl)
	api.EXPECT().Start().Return(expectedError)

	periodic := mock_runner.NewMockPeriodic(ctrl)

	r := runner.New(api, periodic)

	buf := bytes.NewBuffer([]byte{})
	if err := r.Run(time.Second*5, buf); err != expectedError {
		t.Fatalf("wanted %v, got %v", expectedError, err)
	}
}

func TestPeriodicStartFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_runner.NewMockAPI(ctrl)
	api.EXPECT().Start()

	expectedError := errors.New("some periodic error")
	periodic := mock_runner.NewMockPeriodic(ctrl)
	periodic.EXPECT().Start(time.Second*5, gomock.Any()).Return(expectedError)

	r := runner.New(api, periodic)

	buf := bytes.NewBuffer([]byte{})
	if err := r.Run(time.Second*5, buf); err != expectedError {
		t.Fatalf("wanted %v, got %v", expectedError, err)
	}
}
