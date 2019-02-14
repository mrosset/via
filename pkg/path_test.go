package via

import "testing"

func TestPath_String(t *testing.T) {
	tests{
		{
			Expect: "testdata/plans/core/hello.json",
			Got:    NewPlanContext(testConfig, testPlan).PlanFilePath(),
		},
	}.equals(t)
}
