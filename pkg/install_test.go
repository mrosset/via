package via

import (
	"github.com/mrosset/util/file"
	"os"
	"testing"
)

func TestBuild(t *testing.T) {
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

func TestInstaller(t *testing.T) {
	var (
		in    = NewInstaller(testConfig, testPlan)
		files = []string{"testdata/root/opt/via/bin/a.out"}
	)
	defer os.RemoveAll("testdata/root")
	if err := in.Install(); err != nil {
		t.Error(err)
	}
	for _, expect := range files {
		if !file.Exists(expect) {
			t.Errorf("expected %s file got %v", expect, false)
		}
	}
}
