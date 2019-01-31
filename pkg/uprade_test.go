package via

import (
	"testing"
)

func TestUpgraderNil(t *testing.T) {
	var (
		expect = "config is nil"
	)
	defer func() {
		if got := recover(); got != expect {
			t.Errorf("expect '%s' got '%s'", expect, got)
		}
	}()
	NewUpgrader(nil)
}
