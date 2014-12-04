package via

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var (
	test          = "sed"
	repo          = "testdata/repo"
	expectDepends = []string{"glibc"}
	expectFiles   = []string{
		"a.out",
	}
)

func init() {
	Verbose(false)
	os.MkdirAll(repo, 0700)
	config.Repo = repo
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

func TestReadelf(t *testing.T) {
	err := Readelf(join(cache.Pkgs(), "ccache-3.1.7/bin/ccache"))
	if err != nil {
		t.Error(err)
	}
}

func TestPackage(t *testing.T) {
	Clean(testPlan.Name)
	if err := BuildSteps(testPlan); err != nil {
		t.Fatal(err)
	}
	pfile := testPlan.PackagePath()
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
	return
	err := PlanSync()
	if err != nil {
		t.Error(err)
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
