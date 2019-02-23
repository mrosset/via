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

	tests{
		{
			Expect: []string{"OS=linux", "TERM=dumb"},
			Got:    env.KeyValue(),
		},
		{
			Expect: "OS=linux",
			Got:    env.Value("OS"),
		},
	}.equals(t)
}
