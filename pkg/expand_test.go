package via

import (
	"testing"
)

func TestPlanExpand(t *testing.T) {
	var (
		p = &Plan{
			Name:    "plan",
			Version: "1.0",
			Url:     "http://mirrors.kernel.org/gnu/plan-{{.Version}}.tar.gz",
		}
		expect = "http://mirrors.kernel.org/gnu/plan-1.0.tar.gz"
		got    = ""
	)
	got = Expand(p, p.Url)
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}
