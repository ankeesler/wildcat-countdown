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
	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
)

const (
	reunionTimeRFC3339 = "2019-06-07T00:00:00-08:00"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() | log.Lshortfile)
	log.Println("hello from wildcat-countdown")
	printAppDetails()

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

func printAppDetails() {
	appEnv, err := cfenv.Current()
	if err != nil {
		log.Println("NOTE:", "cannot get current cfenv:", err)
		return
	}

	log.Println("ID:", appEnv.ID)
	log.Println("Index:", appEnv.Index)
	log.Println("Name:", appEnv.Name)
	log.Println("Host:", appEnv.Host)
	log.Println("Port:", appEnv.Port)
	log.Println("Version:", appEnv.Version)
	log.Println("Home:", appEnv.Home)
	log.Println("MemoryLimit:", appEnv.MemoryLimit)
	log.Println("WorkingDir:", appEnv.WorkingDir)
	log.Println("TempDir:", appEnv.TempDir)
	log.Println("User:", appEnv.User)
	log.Println("Services:", appEnv.Services)
}
