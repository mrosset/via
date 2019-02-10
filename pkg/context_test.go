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
	got := NewPlanContext(testConfig, plan).PackagePath()
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, expect, got)
	}

	plan.Cid = "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"

	expect = "testdata/repo/QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH.tar.gz"
	got = NewPlanContext(testConfig, plan).PackagePath()
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, expect, got)

	}
}
