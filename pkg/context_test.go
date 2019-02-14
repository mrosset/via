package via

import (
	"testing"
)

func TestContextPackagePath(t *testing.T) {
	var (
		plan = &Plan{
			Name:    "testplan",
			Version: "1.0.1",
		}
	)

	expect := "testdata/repo/testplan-1.0.1-linux-x86_64.tar.gz"
	got := PackagePath(testConfig, plan)
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}

	plan.Cid = "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"

	expect = "testdata/repo/QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH.tar.gz"
	got = PackagePath(testConfig, plan)
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)

	}
}

func TestBuildDir(t *testing.T) {
	var (
		expect = "testdata/cache/bld/hello-2.9"
		got    = NewPlanContext(testConfig, testPlan).BuildDir()
	)
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}
}

func TestSourcPath(t *testing.T) {
	var (
		expect = "testdata/cache/src/hello-2.9.tar.gz"
		got    = NewPlanContext(testConfig, testPlan).SourcePath()
	)
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}
}

func TestStageDir(t *testing.T) {
	var (
		expect = "testdata/cache/stg/hello-2.9"
		got    = NewPlanContext(testConfig, testPlan).StageDir()
	)
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}
}
