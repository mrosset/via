package via

import (
	"github.com/str1ngs/util/json"
	"testing"
)

var (
	testPlan = &Plan{
		Name:         "plan",
		Version:      "1.0",
		Url:          "{{.Mirror}}/{{.Name}}-{{.Version}}.tar.gz",
		BuildInStage: true,
		Package:      []string{"cp a.out $PKGDIR/"},
		Files:        []string{"a.out"},
		Mirror:       "http://mirrors.kernel.org/gnu",
	}
)

func TestTemplate(t *testing.T) {
	var (
		expect = testPlan.Url
		got    string
	)
	err := json.Execute(testPlan)
	if err != nil {
		t.Error(err)
	}
	got = testPlan.template.Url
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestFindPlan(t *testing.T) {
	var (
		expect = "devel"
		got    = ""
	)
	plan, err := NewPlan("devel")
	if err != nil {
		t.Error(err)
	}
	got = plan.Name
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestGetUrl(t *testing.T) {
	var (
		expect = "http://mirrors.kernel.org/gnu/plan-1.0.tar.gz"
		got    = testPlan.Url
	)
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestBuildDir(t *testing.T) {
	var (
		expect = "/home/strings/via/cache/stg/plan-1.0"
		got    = testPlan.BuildDir()
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
