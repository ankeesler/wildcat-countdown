package api_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/api"
)

func TestAPIInterval(t *testing.T) {
	address := "127.0.0.1:12345"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	callbackInterval := time.Duration(0)
	callback := func(interval time.Duration) {
		callbackInterval = interval
	}
	api := api.New(listener, callback)

	errChan := make(chan error)
	if err := api.Start(errChan); err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("http://%s/api/interval", address)
	buf := bytes.NewBuffer([]byte("10000000000"))
	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("received response payload:", string(data))

	if rsp.StatusCode != http.StatusNoContent {
		t.Errorf("wanted %d, got %d", http.StatusNoContent, rsp.StatusCode)
	}

	if callbackInterval != time.Second*10 {
		t.Errorf("wanted %d, got %d", time.Second*10, callbackInterval)
	}
}
