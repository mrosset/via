package upstream

import (
	"github.com/cheekybits/is"
	"github.com/mrosset/util/console"
	"github.com/mrosset/via/pkg"
	"testing"
)

func TestUpstream(t *testing.T) {
	var (
		is = is.New(t)
	)
	is.NoErr(Upstream())
}

func TestParseName(t *testing.T) {
	var (
		is    = is.New(t)
		names = map[string]string{
			"bash":               "bash-4.4.12.tar.gz",
			"Bash":               "Bash-4.4.12.tar.gz",
			"one-two":            "one-two-4.4.12.tar.gz",
			"one-two-three":      "one-two-three-4.4.12.tar.gz",
			"one-two-three-four": "one-two-three-four-4.4.12.tar.gz",
		}
	)

	for expect, v := range names {
		got := ParseName(v)
		is.Equal(expect, got)
	}
}

func TestParseVersion(t *testing.T) {
	var (
		is   = is.New(t)
		vers = map[string]string{
			"4.4.12":        "bash-4.4.12.tar.gz",
			"4.4":           "Bash-4.4.tar.gz",
			"4.444":         "Bash-4.444.tar.gz",
			"4.444.444":     "Bash-4.444.444.tar.gz",
			"4.444.444.555": "Bash-4.444.444.555.tar.gz",
		}
	)

	for expect, v := range vers {
		is.Equal(expect, ParseVersion(v))
	}

}

func TestEachPlanFile(t *testing.T) {
	var (
		is = is.New(t)
	)
	plans, err := via.GetPlans()
	is.NoErr(err)
	for _, p := range plans {
		if p.Url == "" {
			continue
		}
		// file := path.Base(p.Expand().Url)
		//console.Println(file, ParseName(file), ParseVersion(file))
	}
	console.Flush()
}
