package via

import (
	"github.com/mrosset/util/file"
	"os"
	"testing"
)

func TestInstallerCidVerifiy(t *testing.T) {
	var (
		plan = &Plan{
			Name:    "verify",
			Version: "0.0.1",
			Cid:     "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH",
		}
	)

	// Travis does not have a ipfs instance remove this once we
	// have offline Cid verfications
	if os.Getenv("TRAVIS") != "" {
		return
	}

	testConfig.Repo.Ensure()

	file.Touch(PackagePath(testConfig, testPlan))

	test{
		Expect: nil,
		Got:    NewInstaller(testConfig, plan).VerifyCid(),
	}.equals(t)

	plan.Cid = ""
	file.Touch(PackagePath(testConfig, testPlan))

	test{
		Expect: "verify-0.0.1 Plans CID does not match tarballs got QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH",
		Got:    NewInstaller(testConfig, plan).Install().Error(),
	}.equals(t)

}

func fixmeTestBuild(t *testing.T) {
	var (
		files = []string{
			"testdata/cache/src/hello-2.9.tar.gz",
			"testdata/cache/bld/hello-2.9/a.out",
			"testdata/cache/pkg/hello-2.9/opt/via/bin/a.out",
		}
	)
	ctx := NewPlanContext(testConfig, testPlan)
	if err := BuildSteps(ctx); err != nil {
		t.Error(err)
	}
	for _, expect := range files {
		if !file.Exists(expect) {
			t.Errorf("expected %s file got %v", expect, false)
		}
	}
}

func fixmeTestInstaller(t *testing.T) {
	var (
		files = []string{"testdata/root/opt/via/bin/a.out"}
	)
	defer os.RemoveAll("testdata/root")
	if err := NewInstaller(testConfig, testPlan).Install(); err != nil {
		t.Error(err)
	}
	for _, expect := range files {
		if !file.Exists(expect) {
			t.Errorf("expected %s file got %v", expect, false)
		}
	}
}
