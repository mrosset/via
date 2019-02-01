package via

import (
	"testing"
)

// testdata/ipfs directory is generated using.
// ipfs get -o testdata/ipfs QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv

func fixmeTestAdd(t *testing.T) {
	var (
		expect = "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB"
	)
	got, err := IpfsAdd(testConfig, "testdata/ipfs/readme")
	if err != nil {
		t.Error(err)
	}
	if got != expect {
		t.Errorf("expect %s got %s", expect, got)
	}

}

func fixmeTestHashOnly(t *testing.T) {
	var (
		expect = "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB"
	)
	got, err := HashOnly(testConfig, "testdata/ipfs/readme")
	if err != nil {
		t.Error(err)
	}
	if got != expect {
		t.Errorf("expect %s got %s", expect, got)
	}
}
