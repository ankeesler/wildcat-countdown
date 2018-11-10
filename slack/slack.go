package slack

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Messager is an object from which a caller can get a message.
type Messager interface {
	// Message should return a message and an error, if an error occurred.
	Message() (string, error)
}

// Send will post a message (via the Messager) to the provided Slack webhook.
func Send(url string, messager Messager) error {
	message, err := messager.Message()
	if err != nil {
		return err
	}

	data := strings.NewReader(fmt.Sprintf(`{"text":"%s"}`, message))
	rsp, err := http.Post(url, "application/json", data)
	if err != nil {
		return err
	}

	log.Println("sent", message, "to slack, got back", rsp.StatusCode)

	return nil
}
