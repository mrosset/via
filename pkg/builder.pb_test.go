package via

import (
	"testing"
)

func TestRpcBuilder(t *testing.T) {
	go func() {
		if err := Listen(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := ClientRequestBuild("make", true); err != nil {
		t.Error(err)
	}
}
