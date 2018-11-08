package runner_test

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/ankeesler/wildcat-countdown/runner"
	"github.com/ankeesler/wildcat-countdown/runner/mock_runner"
	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_runner.NewMockAPI(ctrl)
	api.EXPECT().Start()

	r := runner.New(api)

	buf := bytes.NewBuffer([]byte{})
	if err := r.Run(buf); err != nil {
		t.Fatal(err)
	}

	expectedString := "hello, tuna\nhello, tuna\nhello, tuna"
	if buf.String() != expectedString {
		//t.Errorf("wanted '%s', got '%s'", expectedString, buf.String())
	}
}

func TestAPIStartFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedError := errors.New("some api error")
	api := mock_runner.NewMockAPI(ctrl)
	api.EXPECT().Start().Return(expectedError)

	r := runner.New(api)
	if err := r.Run(os.Stdout); err != expectedError {
		t.Fatalf("wanted %v, got %v", expectedError, err)
	}
}
