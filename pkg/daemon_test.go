package via

import (
	"testing"
	"time"
)

func TestBuild(t *testing.T) {
	t.SkipNow()
	go StartDaemon(testConfig)
	time.Sleep(time.Second * 1)
	_, err := Connect()
	if err != nil {
		t.Fatal(err)
	}
}
