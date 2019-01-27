package via

import (
	"os"
	"testing"
)

func init() {
	srcpath = "testdata/git"
}

func TestClone(t *testing.T) {
	var (
		expect = "aarch64-via-linux-gnu"
	)
	defer os.RemoveAll(srcpath)
	config, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.Branch != expect {
		t.Errorf("expects %s -> got %s", expect, config.Branch)
	}
}
