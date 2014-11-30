package github

import (
	"testing"
)

func TestAssests(t *testing.T) {
	var (
		expect = true
		got    = GetRelease().Assets.Contains("sed-4.2.1-linux-x86_64.tar.gz")
	)
	if expect != got {
		t.Fatalf("expect %v got %v", expect, got)
	}
}

func TestPush(t *testing.T) {
	err := PushAll()
	if err != nil {
		t.Fatal(err)
	}
}
