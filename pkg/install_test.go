package via

import (
	"os"
	"testing"
)

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
