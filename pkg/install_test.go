package via

import (
	. "github.com/mrosset/test"
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
	PackagePath(testConfig, plan).Touch()

	Tests{
		{
			Name:   "verify cid",
			Expect: nil,
			Got:    NewInstaller(testConfig, plan).VerifyCid(),
		},
	}.Equals(t)

	plan.Cid = ""
	PackagePath(testConfig, plan).Touch()

	Test{
		Name:   "verify empty cid",
		Expect: "verify-0.0.1 Plans CID does not match tarballs got QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH",
		Got:    NewInstaller(testConfig, plan).Install().Error(),
	}.Equals(t)

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
