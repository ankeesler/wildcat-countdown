package main // import "github.com/ankeesler/wildcat-countdown"

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ankeesler/wildcat-countdown/api"
	"github.com/ankeesler/wildcat-countdown/messager"
	"github.com/ankeesler/wildcat-countdown/periodic"
	"github.com/ankeesler/wildcat-countdown/slack"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
)

const (
	reunionTimeRFC3339 = "2019-06-07T00:00:00-08:00"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("hello from wildcat-countdown")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("must specify PORT env var!")
	}
	address := fmt.Sprintf(":%s", port)

	messager := wireMessager()
	periodic := wirePeriodic(messager)

	api := http_server.New(address, api.New(periodic))

	members := []grouper.Member{
		{Name: "periodic", Runner: periodic},
		{Name: "api", Runner: api},
	}

	grouper := grouper.NewParallel(os.Kill, members)
	process := ifrit.Invoke(grouper)

	c := make(chan os.Signal)
	signal.Notify(c)
	signal := <-c

	process.Signal(signal)
	if err := <-process.Wait(); err != nil {
		log.Fatal(err)
	}
}

func wireMessager() *messager.Messager {
	targetDate, err := time.Parse(time.RFC3339, reunionTimeRFC3339)
	if err != nil {
		log.Fatal(err)
	}
	messager := messager.New(targetDate)
	return messager
}

func wirePeriodic(messager *messager.Messager) *periodic.Periodic {
	client := slack.New(messager)
	periodic := periodic.New(time.Minute*10, func() {
		url := os.Getenv("SLACK_URL")
		if url == "" {
			log.Fatal("ERROR:", "must specify SLACK_URL!")
		}

		if err := client.Send(url); err != nil {
			log.Println("ERROR:", err)
		} else {
			log.Println("just sent message to slack!")
		}
	})
	return periodic
}
