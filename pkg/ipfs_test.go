package via

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestAddApi(t *testing.T) {
	var (
		expect = "QmbT3ShooYM8DnvWzLStH9nkDSjDxx1KhcGpg9RhSMcGdh"
	)
	got, err := AddR("testdata/make")
	if err != nil {
		t.Error(err)
	}
	if got != expect {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestIpfsVersion(t *testing.T) {
	is := is.New(t)
	got, err := IpfsVersion()
	is.Nil(err)
	is.Equal("0.4.12-dev", got)
}

func TestAdd(t *testing.T) {
	var (
		expect = "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB"
	)
	got, err := Add("testdata/ipfs/readme")
	if err != nil {
		t.Error(err)
	}
	if got != expect {
		t.Errorf("expect %v got %v", expect, got)
	}
}

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
