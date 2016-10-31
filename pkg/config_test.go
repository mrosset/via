package via

import (
	"testing"
)

func TestConfigExpand(t *testing.T) {
	var (
		c = &Config{
			Prefix: "/usr/local/via",
			Branch: "linux-x86_64",
			Binary: "https://bitbucket.org/strings/publish/raw/{{.Branch}}/repo",
			Env: map[string]string{
				"C_INCLUDE_PATH": "{{.Prefix}}/include",
			},
		}
		expect = "https://bitbucket.org/strings/publish/raw/linux-x86_64/repo"
		got    = ""
	)

	got = c.Expand().Binary
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}

	expect = "/usr/local/via/include"
	got = c.Expand().Env["C_INCLUDE_PATH"]
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestPlanBranch(t *testing.T) {
	var (
		c = &Config{
			Plans: "../plans",
			Repo:  "../publish/repo",
		}
		expect = "linux-x86_64"
		got    = c.PlanBranch()
	)
	if expect != got {
		t.Errorf("expected '%s' got '%s'.", expect, got)
	}
}

func TestRepoBranch(t *testing.T) {
	var (
		c = &Config{
			Repo: "../publish",
		}
		expect = "linux-x86_64"
		got    = c.RepoBranch()
	)
	if expect != got {
		t.Errorf("expected '%s' got '%s'.", expect, got)
	}
}

func TestConfig(t *testing.T) {
	if config == nil {
		t.Errorf("config is nil")
	}
}
