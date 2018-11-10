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
	client := slack.New(m{})
	err := client.Send(os.Getenv("SLACK_URL"))
	if err != nil {
		t.Fatal(err)
	}
}
