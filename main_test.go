package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	t.Run("Get", testGet)
}

func testGet(t *testing.T) {
	rsp, err := http.Get("http://:12345")
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	s := string(data)
	if s != "Go 'Cats!\n" {
		t.Errorf("wanted \"Go 'Cats\\n\", got \"%s\"", s)
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

	cmd := exec.Command(binary)
	cmd.Env = []string{"PORT=12345"}
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
