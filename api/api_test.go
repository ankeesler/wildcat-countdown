package api_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ankeesler/wildcat-countdown/api"
	"github.com/ankeesler/wildcat-countdown/api/mock_api"
	"github.com/golang/mock/gomock"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

func TestAPIIntervalPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	intervalHolder := mock_api.NewMockIntervalHolder(ctrl)
	intervalHolder.EXPECT().SetInterval(time.Second * 10)

	address := "127.0.0.1:12345"
	api := http_server.New(address, api.New(intervalHolder))
	proc := ifrit.Invoke(api)

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

	if rsp.StatusCode != http.StatusOK {
		t.Errorf("wanted %d, got %d", http.StatusOK, rsp.StatusCode)
	}

	if string(data) != "interval set to 10s\n" {
		t.Errorf("wanted 'interval set to 10s\n', got %s", string(data))
	}

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestAPIIntervalGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	intervalHolder := mock_api.NewMockIntervalHolder(ctrl)
	intervalHolder.EXPECT().GetInterval().Return(time.Hour)

	address := "127.0.0.1:12345"
	api := http_server.New(address, api.New(intervalHolder))
	proc := ifrit.Invoke(api)

	url := fmt.Sprintf("http://%s/api/interval", address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	if rsp.StatusCode != http.StatusOK {
		t.Errorf("wanted %d, got %d", http.StatusOK, rsp.StatusCode)
	}

	if string(data) != "interval = 1h0m0s\n" {
		t.Errorf("wanted 'interval = 1h0m0s\n', got %s", string(data))
	}

	proc.Signal(os.Kill)
	if err := <-proc.Wait(); err != nil {
		t.Fatal(err)
	}
}
