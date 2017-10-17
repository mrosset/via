package via

import (
	"testing"
)

func TestIpfsStore(t *testing.T) {
	var (
		expect = "QmbT3ShooYM8DnvWzLStH9nkDSjDxx1KhcGpg9RhSMcGdh"
	)
	got, err := IpfsAdd("testdata/ipfs", false)
	if err != nil {
		t.Error(err)
	}
	if got != expect {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestIpfsStoreSingle(t *testing.T) {
	var (
		expect = "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB"
	)
	got, err := IpfsAdd("testdata/ipfs/readme", false)
	if err != nil {
		t.Error(err)
	}

	if got != expect {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestIpfsGet(t *testing.T) {
	err := IpfsGet("testdata", "QmbT3ShooYM8DnvWzLStH9nkDSjDxx1KhcGpg9RhSMcGdh")
	if err != nil {
		t.Fatal(err)
	}
}
