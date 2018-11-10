package slack_test

import (
	"fmt"
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
	err := slack.Send(
		"https://hooks.slack.com/services/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",
		m{},
	)
	if err != nil {
		t.Fatal(err)
	}
}
