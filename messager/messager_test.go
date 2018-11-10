package messager_test

import (
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/messager"
)

func TestMessage(t *testing.T) {
	targetDate := time.Now().Add(time.Hour * 24 * 15) // 15 days
	targetDate = targetDate.Add(time.Minute)          // wiggleroom
	m := messager.New(targetDate)

	msg, err := m.Message()
	if err != nil {
		t.Fatal(err)
	}

	expected := "15 days until reunion!!!"
	if msg != expected {
		t.Errorf("got \"%s\", wanted \"%s\"", msg, expected)
	}
}
