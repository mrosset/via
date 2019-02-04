package via

import (
	"os"
	"testing"
)

const EXPECT_GOT_FMT = "expect '%v' got '%v'"

func init() {
	Verbose(false)
}

type test struct {
	Expect string
	Got    string
}

func (vt test) equals(fn func(format string, arg ...interface{})) {
	if vt.Expect == "" && vt.Got == "" {
		fn("expect and got will always be equal")
	}
	if vt.Expect != vt.Got {
		fn(EXPECT_GOT_FMT, vt.Expect, vt.Got)
	}
}

func equals(expect, got string, fn func(format string, arg ...interface{})) {
	test{
		Expect: expect,
		Got:    got,
	}.equals(fn)
}

func TestTestType(t *testing.T) {
	test{
		Expect: "foo",
		Got:    "foo",
	}.equals(t.Errorf)
}

func TestLint(t *testing.T) {
	if err := Lint(); err != nil {
		t.Fatal(err)
	}
}

func TestCreate(t *testing.T) {
	var (
		expect = "2.9"
	)
	os.Remove(testPlan.Path())
	err := Create(testPlan.Expand().Url, "core")
	if err != nil {
		t.Error(err)
	}
	_, err = NewPlan(config, testPlan.Name)
	if err != nil {
		t.Error(err)
	}
	got := testPlan.Version
	if expect != testPlan.Version {

		t.Errorf("expected '%s' got '%s'", expect, got)
	}
	os.Remove(testPlan.Path())
}

func TestRepoCreate(t *testing.T) {
	t.Parallel()
	err := RepoCreate(config)
	if err != nil {
		t.Error(err)
	}
}

func TestReadelf(t *testing.T) {
	t.Parallel()
	err := Readelf(join(cache.Packages(), "ccache-3.1.7/bin/ccache"))
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
