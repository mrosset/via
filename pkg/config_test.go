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

func TestBranchs(t *testing.T) {
	var (
		c = &Config{
			Plans: "../plans",
			Repo:  "../publish",
		}
		expect = "x86_64-via-linux-gnu"
		got    = c.PlanBranch()
	)
	// Test Plans
	if expect != got {
		t.Errorf("expected plan branch '%s' got '%s'.", expect, got)
	}

	// Test Repo
	got = c.RepoBranch()
	if expect != got {
		t.Errorf("expected repo branch '%s' got '%s'.", expect, got)
	}
}

func TestConfig(t *testing.T) {
	if config == nil {
		t.Errorf("config is nil")
	}
}
