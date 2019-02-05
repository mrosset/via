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
			config:  testConfig,
		}
	)

	os.MkdirAll(testConfig.Repo, 0755)
	file.Touch(plan.PackagePath())

	test{
		Expect: nil,
		Got:    NewInstaller(testConfig, plan).VerifyCid(),
	}.equals(t.Errorf)

	plan.Cid = ""
	file.Touch(plan.PackagePath())

	test{
		Expect: "verify-0.0.1 Plans CID does not match tarballs got QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH",
		Got:    NewInstaller(testConfig, plan).Install().Error(),
	}.equals(t.Errorf)

}

func fixmeTestBuild(t *testing.T) {
	var (
		files = []string{
			"testdata/cache/src/hello-2.9.tar.gz",
			"testdata/cache/bld/hello-2.9/a.out",
			"testdata/cache/pkg/hello-2.9/opt/via/bin/a.out",
		}
	)
	if err := BuildSteps(testConfig, testPlan); err != nil {
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
