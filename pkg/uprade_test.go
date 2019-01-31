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

func TestUpgrader(t *testing.T) {
	u := NewUpgrader(config)
	if u.Upgrade() != nil {
		t.Error("Upgrader is not nil")
	}
	u.Check()
	if err := u.Upgrade(); err != nil {
		t.Error(err)
	}
}
