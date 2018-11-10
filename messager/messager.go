package messager

import (
	"fmt"
	"time"
)

// Messager is an object from which a caller can get a message.
type Messager struct {
	targetDate time.Time
}

// New instantiates a new Messager.
func New(targetDate time.Time) *Messager {
	return &Messager{targetDate: targetDate}
}

// Message returns a message created by a Getter.
func (m *Messager) Message() (string, error) {
	duration := m.targetDate.Sub(time.Now())
	days := duration / (time.Hour * 24)
	return fmt.Sprintf("%d days until reunion!!!", days), nil
}
