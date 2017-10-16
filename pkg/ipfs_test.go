package via

import (
	"testing"
)

func TestIpfs(t *testing.T) {
	expect := "QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V"
	got, err := AddR("./testdata/ipfs")
	if err != nil {
		t.Fatal(err)
	}
	if expect != got {
		t.Errorf("expect cid %s got %s", expect, got)
	}
}
