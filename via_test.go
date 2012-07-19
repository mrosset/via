package via

import (
	"fmt"
	"github.com/str1ngs/util"
	"testing"
)

var (
	//tests = []string{"ccache", "eglibc"}
	tests = []string{"gcc"}
	turl  = "http://libtorrent.rakshasa.no/downloads/rtorrent-0.8.9.tar.gz"
)

func init() {
	util.Verbose = false
}

func ExampleDepends() {
	fmt.Println(Depends("bash", "/", []string{"bin/bash"}))
	// output:
	// [readline ncurses glibc]
}

func TestLint(t *testing.T) {
	err := Lint()
	if err != nil {
		t.Error(err)
	}
}

func TestStage(t *testing.T) {
	for _, test := range tests {
		plan, err := ReadPlan(test)
		if err != nil {
			t.Fatal(err)
		}
		if err := Stage(plan); err != nil {
			t.Fatal(err)
		}
	}
}
func Testbuildsteps(t *testing.T) {
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
		fmt.Printf(lfmt, "removing", plan.NameVersion())
		if err := Remove(test); err != nil {
			t.Error(err)
		}
	}
}

func Testcreate(t *testing.T) {
	err := Create(turl)
	if err != nil {
		t.Error(err)
	}
}
