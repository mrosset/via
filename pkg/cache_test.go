package via

import (
	"testing"
)

func TestCachePaths(t *testing.T) {
	tests{
		{
			Expect: "testdata/cache/builds/hello-2.9",
			Got:    NewPlanContext(testConfig, testPlan).BuildDir(),
		},
		{
			Expect: "testdata/cache/stages/hello-2.9",
			Got:    NewPlanContext(testConfig, testPlan).StageDir(),
		},
	}.equals(t)
}
