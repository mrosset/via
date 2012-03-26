package via

import (
	"testing"
	"util"
)

var (
	tests = []string{"ccache"}
	turl  = "http://libtorrent.rakshasa.no/downloads/rtorrent-0.8.9.tar.gz"
)

func init() {
	util.Verbose = false
}

func TestBuildSteps(t *testing.T) {
	for _, test := range tests {
		plan, err := ReadPlan(test)
		if err != nil {
			t.Fatal(err)
		}
		if err := BuildSteps(plan); err != nil {
			t.Fatal(err)
		}
		if err := Install(test); err != nil {
			t.Fatal(err)
		}
		if err := List(test); err != nil {
			t.Fatal(err)
		}
		if err := Remove(test); err != nil {
			t.Error(err)
		}
	}
}

func TestCreate(t *testing.T) {
	err := Create(turl)
	if err != nil {
		t.Error(err)
	}
}
