package via

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	test = "zlib"
	repo = "testdata/repo"
	plan *Plan
)

func init() {
	Verbose(false)
	os.MkdirAll(repo, 0700)
	config.Repo = repo
	var err error
	plan, err = FindPlan(test)
	if err != nil {
		log.Fatal(err)
	}
}

func Testlint(t *testing.T) {
	if err := Lint(); err != nil {
		t.Fatal(err)
	}
}

func TestBuildsteps(t *testing.T) {
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

var expectFiles = []string{
	"usr/local/via/include/zconf.h",
	"usr/local/via/include/zlib.h",
	"usr/local/via/lib/libz.a",
	"usr/local/via/lib/pkgconfig/zlib.pc",
	"usr/local/via/share/man/man3/zlib.3",
}

func TestPackage(t *testing.T) {
	pfile := join(config.Repo, plan.PackageFile())
	got, err := ReadPackManifest(pfile)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got.Files, expectFiles) {
		fmt.Println("expect")
		printSlice(expectFiles)
		fmt.Println("got")
		printSlice(got.Files)
	}
}

func TestrepoCreate(t *testing.T) {
	err := RepoCreate()
	if err != nil {
		t.Error(err)
	}
}

func TestrepoSync(t *testing.T) {
	err := PlanSync()
	if err != nil {
		t.Error(err)
	}
}

func TestExpand(t *testing.T) {
	p, err := FindPlan("bash")
	if err != nil {
		t.Error(err)
	}
	e := "http://mirrors.kernel.org/gnu/bash/bash-4.2.tar.gz"
	g := p.Expand("Url")
	if e != g {
		t.Errorf("expected %s got %s", e, g)
	}
}

// Clean up test directories
func TestFinal(t *testing.T) {
	os.RemoveAll(repo)
}
