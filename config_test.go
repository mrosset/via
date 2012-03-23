package via

import (
	"testing"
)

func TestNilConfig(t *testing.T) {
	if config == nil {
		t.Errorf("config is nil")
	}
}
