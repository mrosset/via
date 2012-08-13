package via

import (
	"github.com/str1ngs/util"
	"testing"
)

var (
	test = "zlib"
	turl = "http://libtorrent.rakshasa.no/downloads/rtorrent-0.8.9.tar.gz"
)

func init() {
	Verbose(true)
	util.Verbose = false
}

func TestLint(t *testing.T) {
	if err := Lint(); err != nil {
		t.Fatal(err)
	}
}

func TestBuildsteps(t *testing.T) {
	plan, err := FindPlan(test)
	if err != nil {
		t.Fatal(err)
	}
	if err := BuildSteps(plan); err != nil {
		t.Fatal(err)
	}
}

func TestPackage(t *testing.T) {
	plan, err := FindPlan(test)
	if err != nil {
		t.Fatal(err)
	}
	if err := BuildSteps(plan); err != nil {
		t.Fatal(err)
	}
}

func TestInstall(t *testing.T) {
	config.Root = "tmp"
	err := Install(test)
	if err != nil {
		t.Error(err)
	}
	walkPath("tmp")
}

func TestRemove(t *testing.T) {
	config.Root = "tmp"
	err := Remove(test)
	if err != nil {
		t.Error(err)
	}
}

func TestReadelf(t *testing.T) {
	err := Readelf(join(cache.Pkgs(), "ccache-3.1.7/bin/ccache"))
	if err != nil {
		t.Error(err)
	}
}
