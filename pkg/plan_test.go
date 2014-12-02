package via

import (
	"testing"
)

var (
	testPlan = &Plan{
		Name:         "plan",
		Version:      "1.0",
		Url:          "http://foo.com/plan-1.0.tar.gz",
		BuildInStage: true,
		Package:      []string{"cp a.out $PKGDIR/"},
		Files:        []string{"a.out"},
	}
)

func TestExpand(t *testing.T) {
	var (
		expect = "http://foo.com/plan-1.0.tar.gz"
		got    = testPlan.GetUrl()
	)
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestBuildDir(t *testing.T) {
	var (
		expect = "/home/strings/via/cache/stg/plan-1.0"
		got    = testPlan.GetBuildDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
func TestStageDir(t *testing.T) {
	var (
		expect = "/home/strings/via/cache/stg/plan-1.0"
		got    = testPlan.GetStageDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
