package upstream

import (
	"github.com/blang/semver"
	"testing"
)

type release struct {
	name    string
	url     string
	current string
	expect  string
}

func TestGnuUpstreamLatest(t *testing.T) {
	var (
		releases = []release{
			{
				"bash",
				"http://mirrors.kernel.org/gnu/bash/",
				"4.4",
				"5.0",
			},
			{
				"emacs",
				"http://mirrors.kernel.org/gnu/emacs/",
				"25.1",
				"26.1",
			},
		}
	)
	for _, r := range releases {
		sv, err := semver.ParseTolerant(r.current)
		if err != nil {
			t.Error(err)
		}
		got, err := GnuUpstreamLatest(r.name, r.url, sv)
		if err != nil {
			t.Error(err)
		}
		if r.expect != got {
			t.Errorf("expect '%s' got '%s'", r.expect, got)
		}
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
			"5.0":           "bash-5.0.tar.gz",
		}
	)

	for expect, v := range vers {
		got := ParseVersion(v)
		if expect != got {
			t.Errorf("expect '%s' got '%s'", expect, got)
		}
	}
}
