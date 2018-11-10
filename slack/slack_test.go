package slack_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/slack"
)

func TestSlack(t *testing.T) {
	message := fmt.Sprintf("this is a test running at %s", time.Now().String())
	err := slack.Send(
		"https://hooks.slack.com/services/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",
		message,
	)
	if err != nil {
		t.Fatal(err)
	}
}
