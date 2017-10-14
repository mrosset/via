package via

import (
	"testing"
)

func TestPathJoin(t *testing.T) {
	var (
		tmp    = Path("/tmp")
		expect = Path("/tmp/foo")
		got    = tmp.Join("foo")
	)

	if got != expect {
		t.Errorf("expect: %s got %s", expect, got)
	}
}

func TestPathJoinS(t *testing.T) {
	var (
		tmp    = Path("/tmp")
		expect = Path("/tmp/foo")
		got    = tmp.JoinS("foo")
	)
	if got != expect {
		t.Errorf("expect: %s got %s", expect, got)
	}
}

func TestToUnix(t *testing.T) {
	var (
		expect = "/c/msys64/usr/bin"
		got = Path("c:\\msys64\\usr\\bin").ToUnix()
	)
	if got != expect {
		t.Errorf("expect: %s got %s", expect, got)
	}
}
