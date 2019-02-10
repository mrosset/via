package via

import (
	"testing"
)

const EXPECT_GOT_FMT = "expect '%v' got '%v'"

func init() {
	Verbose(false)
}

type test struct {
	Expect interface{}
	Got    interface{}
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
		files  = RepoFiles{"glibc": []string{"libc.so.6"}}
		expect = "glibc"
		got    = files.Owns("libc.so.6")
	)
	if expect != got {
		t.Errorf(EXPECT_GOT_FMT, expect, got)
	}

}
