package via

import (
	"testing"
)

var (
	testPlan = &Plan{
		Name:         "plan",
		Version:      "1.0",
		Url:          "http://mirrors.kernel.org/gnu/plan-{{.Version}}tar.gz",
		BuildInStage: true,
		Package:      []string{"cp a.out $PKGDIR/"},
		Files:        []string{"a.out"},
		Group:        "core",
	}
)

func TestPlanExpand(t *testing.T) {
	var (
		p = &Plan{
			Name:    "plan",
			Version: "1.0",
			Url:     "http://mirrors.kernel.org/gnu/{{.Name}}-{{.Version}}.tar.gz",
		}
		expect = "http://mirrors.kernel.org/gnu/plan-1.0.tar.gz"
		got    = ""
	)

	got = p.Expand().Url
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestFindPlan(t *testing.T) {
	var (
		expect = &Plan{
			Name: "sed",
			Url:  "http://mirrors.kernel.org/gnu/sed/sed-{{.Version}}.tar.xz",
		}
	)
	got, err := NewPlan("sed")
	if err != nil {
		t.Fatal(err)
	}
	if expect.Name != got.Name || got.Url != expect.Url {
		t.Errorf("expected %s got %s", expect.Url, got.Url)
	}
}

func TestBuildDir(t *testing.T) {
	var (
		expect = "/home/mrosset/.cache/via/stg/plan-1.0"
		got    = testPlan.BuildDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
func TestStageDir(t *testing.T) {
	var (
		expect = "/home/mrosset/.cache/via/stg/plan-1.0"
		got    = testPlan.GetStageDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}
