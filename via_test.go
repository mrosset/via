package via

import (
	"testing"
	"util"
)

var tests = []string{"ccache"}

func init() {
	util.Verbose = false
}

func TestBuildSteps(t *testing.T) {
	for _, test := range tests {
		plan, err := ReadPlan(test)
		if err != nil {
			t.Error(err)
		}
		err = BuildSteps(plan)
		if err != nil {
			t.Error(err)
		}
	}
}
