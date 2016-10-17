package git

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

func shell(path string) {
	sh := exec.Command("sh", "-c", path)
	sh.Stdout = os.Stdout
	sh.Stdin = os.Stdin
	if err := sh.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	shell("scripts/init")
}

func TestBranch(t *testing.T) {
	var (
		expect = "linux-x86_64"
		got    = ""
	)
	got, err := Branch("testdata/.git")
	if err != nil {
		t.Fatal(err)
	}
	if expect != got {
		t.Logf("expect '%s' got '%s'", expect, got)
	}
}

func TestBranchFail(t *testing.T) {
	var (
		expect = "nobranch"
		got    = ""
	)
	got, err := Branch("testdata/.git")
	if err != nil {
		t.Fatal(err)
	}
	if expect == got {
		t.Errorf("expect '%s' got '%s'", expect, got)
	}
}

func TestCleanup(t *testing.T) {
	err := os.RemoveAll("testdata")
	if err != nil {
		t.Error(err)
	}
}
