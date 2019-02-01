package upstream

import (
	"github.com/mrosset/util/console"
	"github.com/mrosset/via/pkg"
	"testing"
)

func TestUpstreamError(t *testing.T) {
	if err := Upstream(); err != nil {
		t.Error(err)
	}
}

func TestParseName(t *testing.T) {
	var (
		names = map[string]string{
			"bash":               "bash-4.4.12.tar.gz",
			"Bash":               "Bash-4.4.12.tar.gz",
			"one-two":            "one-two-4.4.12.tar.gz",
			"one-two-three":      "one-two-three-4.4.12.tar.gz",
			"one-two-three-four": "one-two-three-four-4.4.12.tar.gz",
		}
	)

	for expect, n := range names {
		got := ParseName(n)
		if expect != got {
			t.Errorf("expect '%s' got '%s'", expect, got)
		}
	}
}

func TestParseVersion(t *testing.T) {
	var (
		vers = map[string]string{
			"4.4.12":        "bash-4.4.12.tar.gz",
			"4.4":           "Bash-4.4.tar.gz",
			"4.444":         "Bash-4.444.tar.gz",
			"4.444.444":     "Bash-4.444.444.tar.gz",
			"4.444.444.555": "Bash-4.444.444.555.tar.gz",
		}
	)

	for expect, v := range vers {
		got := ParseVersion(v)
		if expect != got {
			t.Errorf("expect '%s' got '%s'", expect, got)
		}
	}
}

func OTestEachPlanFile(t *testing.T) {
	plans, err := via.GetPlans()
	if err != nil {
		t.Error(err)
	}
	for _, p := range plans {
		if p.Cid == "" {

		}
		// file := path.Base(p.Expand().Url)
		//console.Println(file, ParseName(file), ParseVersion(file))
	}
	console.Flush()
}
