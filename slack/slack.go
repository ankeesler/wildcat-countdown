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

// Client is an object that can Send() message to Slack.
type Client struct {
	messager Messager
}

// New creates a new Client.
func New(messager Messager) *Client {
	return &Client{messager: messager}
}

// Send will post a message (via the Messager) to the provided Slack webhook.
func (c *Client) Send(url string) error {
	message, err := c.messager.Message()
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
