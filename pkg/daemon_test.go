package via

import (
	"testing"
)

func init() {
	go Listen()
}

func TestRpc(t *testing.T) {
	c, dc := NewRpcClient()
	c.Start()
	defer c.Stop()

	ok, err := dc.Call("ping", nil)
	if err != nil {
		t.Error(err)
	}
	if !ok.(bool) {
		t.Errorf("server return unexpected response %v", ok)
	}
}
