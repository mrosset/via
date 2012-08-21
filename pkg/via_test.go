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

func TestRepoCreate(t *testing.T) {
	err := RepoCreate()
	if err != nil {
		t.Error(err)
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
	"usr/local/via/lib/libz.so",
	"usr/local/via/lib/libz.so.1",
	"usr/local/via/lib/libz.so.1.2.7",
	"usr/local/via/lib/pkgconfig/zlib.pc",
	"usr/local/via/share/man/man3/zlib.3",
}

var expectDepends = []string{"glibc"}

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
	if !reflect.DeepEqual(got.Depends, expectDepends) {
		fmt.Println("expect")
		printSlice(expectDepends)
		fmt.Println("got")
		printSlice(got.Depends)
	}
}

func TestrepoSync(t *testing.T) {
	err := PlanSync()
	if err != nil {
		t.Error(err)
	}
}

func TestExpand(t *testing.T) {
	var (
		plan, _ = FindPlan("bash")
		expect  = "http://mirrors.kernel.org/gnu/bash/bash-4.2.tar.gz"
		got     = plan.Expand("Url")
	)
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestOwns(t *testing.T) {
	var (
		files, _ = ReadRepoFiles()
		expect   = "glibc"
		got      = files.Owns("libc.so.6")
	)
	if expect != got {
		t.Errorf("expected %s got %s.", expect, got)
	}
}

func TestSort(t *testing.T) {
	plans, err := NewPlanSlice()
	if err != nil {
		t.Error(err)
	}
	plans.SortSize().Print()
}

// Clean up test directories
func TestFinal(t *testing.T) {
	os.RemoveAll(repo)
}
