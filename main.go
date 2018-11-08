package main // import "github.com/ankeesler/wildcat-countdown"

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ankeesler/wildcat-countdown/api"
	"github.com/ankeesler/wildcat-countdown/periodic"
	"github.com/ankeesler/wildcat-countdown/runner"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("hello from wildcat-countdown")

	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	api := api.New(listener)
	periodic := periodic.New(time.Hour*24, func() {
		fmt.Println("hello, tuna")
	})

	runner := runner.New(api, periodic)
	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c)
	<-c
}
