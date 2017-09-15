package via

import (
	"github.com/mrosset/util/file"
	"testing"
	"time"
)

func TestBuild(t *testing.T) {
	t.SkipNow()
	if file.Exists(SOCKET_FILE) {
		t.Fatalf("%s: exists", SOCKET_FILE)
	}
	go StartDaemon()
	time.Sleep(time.Second * 2)
	_, err := Connect()
	if err != nil {
		t.Fatal(err)
	}
}
