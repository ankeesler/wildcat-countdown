package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	t.Run("Get", testGet)
	t.Run("Interval", testInterval)
	t.Run("Slack", testSlack)
}

func testGet(t *testing.T) {
	s := hitAPI(t, http.MethodGet, "", "")
	if s != "Go 'Cats!\n" {
		t.Errorf("wanted \"Go 'Cats\\n\", got \"%s\"", s)
	}
}

func testInterval(t *testing.T) {
	s := hitAPI(t, http.MethodGet, "api/interval", "")
	if s != "interval = 10m0s\n" { // default timeout is 10 minutes!
		t.Errorf("wanted \"interval = 10m0s\n\", got \"%s\"", s)
	}

	s = hitAPI(t, http.MethodPut, "api/interval", "3600000000000")
	if s != "interval set to 1h0m0s\n" {
		t.Errorf("wanted \"interval set to 1h0m0s\n\", got \"%s\"", s)
	}

	s = hitAPI(t, http.MethodGet, "api/interval", "")
	if s != "interval = 1h0m0s\n" {
		t.Errorf("wanted \"interval = 1h0m0s\n\", got \"%s\"", s)
	}
}

func testSlack(t *testing.T) {
	// fake slack
	type stuff struct {
		method  string
		payload []byte
	}
	fakeSlackStuff := make(chan stuff, 1)
	go func() {
		var handler http.HandlerFunc
		handler = func(_ http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			fakeSlackStuff <- stuff{method: r.Method, payload: data}
		}
		if err := http.ListenAndServe("localhost:12346", handler); err != nil {
			t.Fatal(err)
		}
	}()

	s := hitAPI(t, http.MethodPut, "api/interval", "1000000000")
	t.Log(s)

	time.Sleep(time.Second * 1)

	stuf := <-fakeSlackStuff
	if stuf.method != http.MethodPost {
		t.Fatalf("wanted %s, got %s", http.MethodPost, stuf.method)
	}

	payload := make(map[string]string)
	if err := json.Unmarshal(stuf.payload, &payload); err != nil {
		t.Fatal(err)
	}
	if value, ok := payload["text"]; !ok {
		t.Errorf("expected 'text' key in payload (%s)", string(stuf.payload))
	} else if value != "Go 'Cats!" {
		t.Log(value)
	}
}

func setup(t *testing.T) func() {
	dir, err := ioutil.TempDir("", "wildcat-countdown-test")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(dir, "wildcat-countdown")

	if err := exec.Command("go", "build", "-o", binary, ".").Run(); err != nil {
		t.Fatal(err)
	}

	if os.Getenv("SLACK_URL") == "" {
		t.Fatal("must set SLACK_URL env var!")
	}

	cmd := exec.Command(binary)
	cmd.Env = []string{
		"PORT=12345",
		"SLACK_URL=http://localhost:12346", // fake slack...hehehe
	}
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	return func() {
		t.Log("STDOUT:", stdout.String())
		t.Log("STDERR:", stderr.String())
		cmd.Process.Signal(os.Kill)
		os.RemoveAll(dir)
	}
}

func hitAPI(t *testing.T, method, path, payload string) string {
	url := fmt.Sprintf("http://:12345/%s", path)
	buf := strings.NewReader(payload)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	return string(data)
}
