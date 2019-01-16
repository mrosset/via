package via

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestAdd(t *testing.T) {
	var (
		is     = is.New(t)
		expect = "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB"
	)
	cid, err := IpfsAdd(testConfig, "testdata/ipfs/readme")
	is.Nil(err)
	is.Equal(cid, expect)
}

func TestHashOnly(t *testing.T) {
	var (
		is     = is.New(t)
		expect = "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB"
	)
	cid, err := HashOnly(testConfig, "testdata/ipfs/readme")
	is.Nil(err)
	is.Equal(cid, expect)
}
