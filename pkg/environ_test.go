package via

import (
	. "github.com/mrosset/test"
	"testing"
)

func TestEnv_KeyValue(t *testing.T) {
	var (
		env = Env{
			"OS":   "linux",
			"TERM": "dumb",
		}
	)

	Tests{
		{
			Expect: []string{"OS=linux", "TERM=dumb"},
			Got:    env.KeyValue(),
		},
		{
			Expect: "OS=linux",
			Got:    env.Value("OS"),
		},
	}.Equals(t)
}
