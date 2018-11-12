package slack_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/slack"
)

type m struct{}

func (m m) Message() (string, error) {
	message := fmt.Sprintf("this is a test running at %s", time.Now().String())
	return message, nil
}

func TestSlack(t *testing.T) {
	if os.Getenv("SLACK_URL") == "" {
		t.Fatal("must set SLACK_URL env var!")
	}

	client := slack.New(m{})
	err := client.Send(os.Getenv("SLACK_URL"))
	if err != nil {
		t.Fatal(err)
	}
}
