package git

import (
	"os"
	"testing"
)

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
