package git

import (
	"github.com/mrosset/util/file"
	"os"
	"testing"
)

func TestClone(t *testing.T) {
	var (
		expect = "testdata/config.json"
	)
	defer os.RemoveAll("testdata")
	if err := Clone("testdata", "https://github.com/mrosset/plans"); err != nil {
		t.Error(err)
	}
	if !file.Exists(expect) {
		t.Errorf("exected %s but file does not exist", expect)
	}
}

func TestBranch(t *testing.T) {
	var (
		expect = "field_expansion"
		got    = ""
	)
	got, err := Branch("../../../via")
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
	got, err := Branch("../../publish")
	if err != nil {
		t.Fatal(err)
	}
	if expect != got {
		t.Logf("expect '%s' got '%s'", expect, got)
	}
}
