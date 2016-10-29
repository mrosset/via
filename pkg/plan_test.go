package via

import (
	"testing"
)

var (
	testPlan = &Plan{
		Name:         "plan",
		Version:      "1.0",
		Url:          "http://mirrors.kernel.org/gnu/plan-$Version.tar.gz",
		BuildInStage: true,
		Package:      []string{"cp a.out $PKGDIR/"},
		Files:        []string{"a.out"},
		Group:        "core",
	}
)

/*
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
*/

func TestExpand(t *testing.T) {
	var (
		p = &Plan{
			Name:    "plan",
			Version: "1.0",
			EUrl:    "http://mirrors.kernel.org/gnu/plan-{{.Version}}.tar.gz",
		}
		expect = "http://mirrors.kernel.org/gnu/plan-1.0.tar.gz"
		got    = ""
	)
	got = p.EUrl.Expand(p)
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestFindPlan(t *testing.T) {
	var (
		expect = "sed"
		got    = ""
	)
	plan, err := NewPlan("sed")
	if err != nil {
		t.Fatal(err)
	}
	got = plan.Name
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestGetUrl(t *testing.T) {
	var (
		expect = "http://mirrors.kernel.org/gnu/plan-1.0.tar.gz"
		got    = testPlan.ExpandField("Url")
	)
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestBuildDir(t *testing.T) {
	var (
		expect = "/home/strings/via_cache/stg/plan-1.0"
		got    = testPlan.BuildDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
func TestStageDir(t *testing.T) {
	var (
		expect = "/home/strings/via_cache/stg/plan-1.0"
		got    = testPlan.GetStageDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
