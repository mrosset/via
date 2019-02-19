package via

import (
	"path/filepath"
	"testing"
)

func TestCachePaths(t *testing.T) {
	tests{
		{
			Expect: filepath.Join(wd, "testdata/cache/builds/hello-2.9"),
			Got:    NewPlanContext(testConfig, testPlan).BuildDir(),
		},
		{
			Expect: filepath.Join(wd, "testdata/cache/stages/hello-2.9"),
			Got:    NewPlanContext(testConfig, testPlan).StageDir(),
		},
	}.equals(t)
}
