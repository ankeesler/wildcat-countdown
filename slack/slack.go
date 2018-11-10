package slack

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Send will post a message to the provided Slack webhook.
func Send(url, message string) error {
	data := strings.NewReader(fmt.Sprintf(`{"text":"%s"}`, message))
	rsp, err := http.Post(url, "application/json", data)
	if err != nil {
		return err
	}

	log.Println("sent", message, "to slack, got back", rsp.StatusCode)

	return nil
}
