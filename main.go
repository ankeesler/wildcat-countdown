package main // import "github.com/ankeesler/wildcat-countdown"

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/ankeesler/wildcat-countdown/api"
	"github.com/ankeesler/wildcat-countdown/runner"
)

func main() {
	log.Println("hello from wildcat-countdown")

	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	api := api.New(listener)
	runner := runner.New(api)
	if err := runner.Run(os.Stdout); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c)
	<-c
}
