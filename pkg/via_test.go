package via

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var (
	test  = "zlib"
	troot = "troot"
	trepo = "trepo"
)

func clean() {
	os.RemoveAll(troot)
	os.RemoveAll(trepo)
}

func init() {
	clean()
	os.MkdirAll(troot, 0700)
	os.MkdirAll(trepo, 0700)
	config.Root = troot
	config.Repo = trepo
	Verbose(false)
}

func Testlint(t *testing.T) {
	if err := Lint(); err != nil {
		t.Fatal(err)
	}
}

func Testbuildsteps(t *testing.T) {
	plan, err := FindPlan(test)
	if err != nil {
		t.Fatal(err)
	}
	if err := BuildSteps(plan); err != nil {
		t.Fatal(err)
	}
}

var hwant = []string{
	"troot",
	"troot/usr",
	"troot/usr/local",
	"troot/usr/local/via",
	"troot/usr/local/via/db",
	"troot/usr/local/via/db/installed",
	"troot/usr/local/via/db/installed/zlib",
	"troot/usr/local/via/db/installed/zlib/manifest.json",
	"troot/usr/local/via/include",
	"troot/usr/local/via/include/zconf.h",
	"troot/usr/local/via/include/zlib.h",
	"troot/usr/local/via/lib",
	"troot/usr/local/via/lib/libz.a",
	"troot/usr/local/via/lib/libz.so",
	"troot/usr/local/via/lib/libz.so.1",
	"troot/usr/local/via/lib/libz.so.1.2.7",
	"troot/usr/local/via/lib/pkgconfig",
	"troot/usr/local/via/lib/pkgconfig/zlib.pc",
	"troot/usr/local/via/share",
	"troot/usr/local/via/share/man",
	"troot/usr/local/via/share/man/man3",
	"troot/usr/local/via/share/man/man3/zlib.3",
}

func Testinstall(t *testing.T) {
	err := Install(test)
	if err != nil {
		t.Error(err)
	}
	hgot, err := walkPath(troot)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(hwant, hgot) {
		fmt.Println("want:")
		printSlice(hwant)
		fmt.Println("got:")
		printSlice(hgot)
		t.Fail()
	}
}

func Testremove(t *testing.T) {
	err := Remove(test)
	if err != nil {
		t.Error(err)
	}
}

func Testreadelf(t *testing.T) {
	err := Readelf(join(cache.Pkgs(), "ccache-3.1.7/bin/ccache"))
	if err != nil {
		t.Error(err)
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
	expect := "http://mirrors.kernel.org/gnu/bash/bash-4.2.tar.gz"
	got := p.Expand("Url")
	if got != expect {
		t.Errorf("expected %s got %s", got, expect)
	}
}