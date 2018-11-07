package main // import "github.com/ankeesler/wildcat-countdown"

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("hey")

	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(address, nil))
}
