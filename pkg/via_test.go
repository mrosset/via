package via

import (
	"bytes"
	"github.com/mrosset/util/file"
	"os"
	"os/exec"
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
	var (
		path = "testdata/plans/files.json"
	)
	defer os.Remove(path)

	test{
		Expect: nil,
		Got:    RepoCreate(testConfig),
	}.equals(t.Errorf)

	test{
		Expect: true,
		Got:    file.Exists(path),
	}.equals(t.Errorf)

}

func TestReadelf(t *testing.T) {
	t.Parallel()
	var (
		out = "testdata/a.out"
	)
	defer os.Remove(out)
	bin, err := exec.LookPath("gcc")
	if err != nil {
		t.Fatal(err)
	}
	gcc := &exec.Cmd{
		Path:  bin,
		Args:  []string{"gcc", "-o", out, "-xc", "-"},
		Stdin: bytes.NewBufferString("int main() {}\n"),
	}
	if err := gcc.Start(); err != nil {
		t.Fatal(err)
	}
	if err = Readelf(out); err != nil {
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
