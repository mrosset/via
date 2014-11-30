package via

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	test = "sed"
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

func TestLint(t *testing.T) {
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
	"usr/local/via/bin/sed",
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

func TestRepoSync(t *testing.T) {
	err := PlanSync()
	if err != nil {
		t.Error(err)
	}
}

func TestExpand(t *testing.T) {
	var (
		plan, _ = FindPlan("make")
		expect  = "http://mirrors.kernel.org/gnu/make/make-4.1.tar.bz2"
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

// Clean up test directories
func TestFinal(t *testing.T) {
	os.RemoveAll(repo)
}
