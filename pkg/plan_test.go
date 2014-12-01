package via

import (
	"testing"
)

var (
	tpp = &Plan{
		Name:    "plan",
		Version: "1.0",
	}
)

func TestBuildDir(t *testing.T) {
	var (
		expect = "/home/strings/via/cache/bld/plan-1.0"
		got    = tpp.GetBuildDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
func TestStageDir(t *testing.T) {
	var (
		expect = "/home/strings/via/cache/stg/plan-1.0"
		got    = tpp.GetStageDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
