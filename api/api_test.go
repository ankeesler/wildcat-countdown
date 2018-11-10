package api_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/ankeesler/wildcat-countdown/api"
	"github.com/ankeesler/wildcat-countdown/api/mock_api"
)

func TestAPIInterval(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	intervalSetter := mock_api.NewMockIntervalSetter(ctrl)
	intervalSetter.EXPECT().SetInterval(time.Second * 10)

	address := "127.0.0.1:12345"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	api := api.New(listener, intervalSetter)

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

	if rsp.StatusCode != http.StatusOK {
		t.Errorf("wanted %d, got %d", http.StatusNoContent, rsp.StatusCode)
	}

	if string(data) != "interval set to 10s\n" {
		t.Errorf("wanted 'interval set to 10s\n', got %s", string(data))
	}
}
