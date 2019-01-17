package via

import (
	"github.com/mrosset/util/file"
	"os"
	"testing"
)

func TestBuild(t *testing.T) {
	var (
		expectSrc = "testdata/cache/src/hello-2.9.tar.gz"
	)
	err := BuildSteps(testConfig, testPlan)
	if err != nil {
		t.Error(err)
	}
	if !file.Exists(expectSrc) {
		t.Errorf("expected %s file to exist", expectSrc)
	}
}

func TestInstaller(t *testing.T) {
	var (
		p, err = NewPlan("ccache")
		in     = NewInstaller(testConfig, p)
	)
	defer os.RemoveAll("testdata/root")
	defer os.RemoveAll("testdata/repo")
	if err != nil {
		t.Error(err)
	}
	err = in.Install()
	if err != nil {
		t.Error(err)
	}
}
