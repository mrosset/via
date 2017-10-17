package via

import (
	"fmt"
	"testing"
)

func TestIpfs(t *testing.T) {
	var (
		got string
		err error
	)
	if got, err = Store("testdata/ipfs"); err != nil {
		t.Error(err)
	}
	fmt.Println("GOT", got)
}
