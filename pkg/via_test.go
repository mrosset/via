package via

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

//revive:disable
const EXPECT_GOT_FMT = "%s: expect '%s' got '%s'"

//revive:enable

func init() {
	Verbose(false)
}

type test struct {
	Name   string
	Expect interface{}
	Got    interface{}
}

type tests []test

func (ts tests) equals(t *testing.T) {
	for _, test := range ts {
		test.equals(t)
	}
}

func (vt test) equals(t *testing.T) bool {

	if !reflect.DeepEqual(vt.Expect, vt.Got) {
		t.Errorf(EXPECT_GOT_FMT, vt.Name, vt.Expect, vt.Got)
		return false
	}
	return true
}

func TestTestType(t *testing.T) {
	test{
		Expect: "foo",
		Got:    "foo",
	}.equals(t)
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
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}

}
