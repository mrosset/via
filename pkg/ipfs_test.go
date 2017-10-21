package via

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestAdd(t *testing.T) {
	var (
		expect = "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB"
		is     = is.New(t)
	)
	got, err := Add("testdata/ipfs/readme")
	is.Nil(err)
	is.Equal(got, expect)
}

func TestIpfsStore(t *testing.T) {
	var (
		expect = "QmbT3ShooYM8DnvWzLStH9nkDSjDxx1KhcGpg9RhSMcGdh"
		is     = is.New(t)
	)
	got, err := IpfsAdd("testdata/ipfs", false)
	is.Nil(err)
	is.Equal(got, expect)
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
