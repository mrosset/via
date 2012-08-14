package via

import (
	"testing"
)

func TestConfig(t *testing.T) {
	if config == nil {
		t.Errorf("config is nil")
	}
}
