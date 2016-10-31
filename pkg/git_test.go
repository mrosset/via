package via

import (
	"github.com/mrosset/util/file"
	"os"
	"testing"
)

func TestClone(t *testing.T) {
	var (
		expect = "testdata/git/README"
		gitd   = "testdata/git"
	)
	defer os.RemoveAll(gitd)
	if err := Clone(gitd, "https://github.com/mrosset/gur"); err != nil {
		t.Error(err)
	}
	if !file.Exists(expect) {
		t.Errorf("exected %s but file does not exist", expect)
	}
	expect = "master"
	got, err := Branch(gitd)
	if err != nil {
		t.Fatal(err)
	}
	if expect != got {
		t.Logf("expect '%s' got '%s'", expect, got)
	}
}

func TestSubBranch(t *testing.T) {
	var (
		expect = "linux-x86_64"
		got    = ""
	)
	got, err := Branch("../publish")
	if err != nil {
		t.Fatal(err)
	}
	if expect != got {
		t.Logf("expect '%s' got '%s'", expect, got)
	}
}
