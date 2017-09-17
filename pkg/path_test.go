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
		t.Errorf("expectd: %s got %s", expect, got)
	}
}

func TestPathJoinS(t *testing.T) {
	var (
		tmp    = Path("/tmp")
		expect = Path("/tmp/foo")
		got    = tmp.JoinS("foo")
	)
	if got != expect {

		t.Errorf("expectd: %s got %s", expect, got)
	}
}
