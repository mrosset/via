package via

import (
	"testing"
)

func TestEnv_KeyValue(t *testing.T) {
	var (
		env = Env{
			"OS":   "linux",
			"TERM": "dumb",
		}
	)

	test{
		Expect: []string{"OS=linux", "TERM=dumb"},
		Got:    env.KeyValue(),
	}.equals(t)
}
