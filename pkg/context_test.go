package via

import (
	"testing"
)

func TestContextPackagePath(t *testing.T) {
	var (
		config = Config{
			OS:   "linux",
			Arch: "x86_64",
			Repo: "testdata/repo",
		}
		plan = &Plan{
			Name:    "testplan",
			Version: "1.0.0",
			Cid:     "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH",
		}
		expect = "testdata/repo/QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH.tar.gz"
	)

	ctx := NewViaContext(config, plan)
	if got := ctx.PackagePath(); expect != got {
		t.Errorf(EXPECT_GOT_FMT, expect, got)
	}
}
