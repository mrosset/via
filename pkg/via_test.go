package via

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func init() {
	Verbose(false)
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
