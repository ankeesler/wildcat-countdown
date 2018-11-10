package main // import "github.com/ankeesler/wildcat-countdown"

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ankeesler/wildcat-countdown/api"
	"github.com/ankeesler/wildcat-countdown/messager"
	"github.com/ankeesler/wildcat-countdown/periodic"
	"github.com/ankeesler/wildcat-countdown/runner"
	"github.com/ankeesler/wildcat-countdown/slack"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("hello from wildcat-countdown")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("must specify PORT env var!")
	}

	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	periodic := periodic.New(time.Minute, sendSlackMessage)

	api := api.New(listener, periodic)

	runner := runner.New(api, periodic)
	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c)
	<-c
}

func sendSlackMessage() {
	url := os.Getenv("SLACK_URL")
	if url == "" {
		log.Fatal("ERROR:", "must specify SLACK_URL!")
	}

	if err := slack.Send(url, messager.New()); err != nil {
		log.Println("ERROR:", err)
	}
	log.Println("just sent message to slack!")
}
