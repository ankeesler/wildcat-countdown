package runner_test

import (
	"errors"
	"os"
	"testing"

	"github.com/ankeesler/wildcat-countdown/api/mock_api"
	"github.com/ankeesler/wildcat-countdown/runner"
	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock_api.NewMockAPI(ctrl)
	api.EXPECT().Start()

	r := runner.New(api)
	if err := r.Run(os.Stdout); err != nil {
		t.Fatal(err)
	}
}

func TestAPIStartFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedError := errors.New("some api error")
	api := mock_api.NewMockAPI(ctrl)
	api.EXPECT().Start().Return(expectedError)

	r := runner.New(api)
	if err := r.Run(os.Stdout); err != expectedError {
		t.Fatalf("wanted %v, got %v", expectedError, err)
	}
}
