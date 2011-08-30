package via

import (
	"fmt"
	"os"
	"testing"
)

var (
	tests    = []string{"bash", "ncdu"}
	testArch = "x86_64"
	testRoot = "./tmp"
)

func TestFindPlan(t *testing.T) {
	for _, test := range tests {
		expected := test
		plan, err := FindPlan(expected)
		checkError(t, err)
		if plan.Name != expected {
			t.Errorf("exected %s for Name got %s", expected, plan.Name)
		}
	}
}

func TestPackage(t *testing.T) {
	for _, test := range tests {
		err := Package(test, testArch)
		checkError(t, err)
	}
}

func TestUnPack(t *testing.T) {
	for _, test := range tests {
		plan, err := FindPlan(test)
		checkError(t, err)
		err = Unpack(testRoot, PkgAbsFile(plan, testArch))
		checkError(t, err)
	}
}

func TestUpdateRepo(t *testing.T) {
	err := UpdateRepo(testArch)
	checkError(t, err)
}

func TestLoadRepo(t *testing.T) {
	rep, err := LoadRepo(testArch)
	checkError(t, err)
	for _, m := range rep.Manifests {
		fmt.Printf("%-10.10s %s\n", m.Meta.Name, m.Meta.Tarball)
	}
}

func checkError(t *testing.T, err os.Error) {
	if err != nil {
		t.Error(err)
	}
}
