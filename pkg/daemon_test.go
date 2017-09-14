package via

import (
	"testing"
	"time"
)

func TestBuild(t *testing.T) {
	go StartDaemon()
	time.Sleep(time.Second * 1)
	_, err := Connect()
	if err != nil {
		t.Fatal(err)
	}
}
