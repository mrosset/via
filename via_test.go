package via

import (
	"fmt"
	"github.com/str1ngs/util"
	"testing"
)

var (
	test = "ccache"
	turl = "http://libtorrent.rakshasa.no/downloads/rtorrent-0.8.9.tar.gz"
)

func init() {
	Verbose(true)
	util.Verbose = false
}

func Testbuildsteps(t *testing.T) {
	plan, err := ReadPlan(test)
	if err != nil {
		t.Fatal(err)
	}
	if err := BuildSteps(plan); err != nil {
		t.Fatal(err)
	}
}

func TestPackage(t *testing.T) {
	plan, err := ReadPlan(test)
	if err != nil {
		t.Fatal(err)
	}
	if err := BuildSteps(plan); err != nil {
		t.Fatal(err)
	}
}

func TestReadelf(t *testing.T) {
	err := Readelf(join(cache.Pkgs(), "ccache-3.1.7/bin/ccache"))
	if err != nil {
		t.Error(err)
	}
}

func ExampleDepends() {
	fmt.Println(Depends("bash", "/", []string{"bin/bash"}))
	// output:
	// [readline ncurses]
}
