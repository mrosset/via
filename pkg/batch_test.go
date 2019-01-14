package via

import (
	"github.com/cheekybits/is"
	"github.com/mrosset/util/file"
	"os"
	"testing"
)

func TestBatchAdd(t *testing.T) {
	is := is.New(t)
	d := NewBatch(config)
	d.Add(testPlan)
	is.Equal(d.Plans["plan"], testPlan)
}

func TestBatchWalk(t *testing.T) {
	var (
		p, _   = NewPlan("ccache")
		got    = NewBatch(testConfig)
		expect = 3
	)
	got.Walk(p)
	if len(got.Plans) != 3 {
		t.Errorf("expect %d depends got %d", expect, len(got.Plans))
	}
}

func TestBatchInstall(t *testing.T) {
	var (
		p, _   = NewPlan("ccache")
		got    = NewBatch(testConfig)
		expect = join(testConfig.Repo, "repo", p.PackageFile())
	)
	defer os.RemoveAll(testConfig.Repo)
	defer os.RemoveAll(testConfig.Root)
	defer os.RemoveAll(testConfig.DB.Installed())
	got.Walk(p)
	errors := got.Install()
	if len(errors) != 0 {
		t.Error(errors)
	}
	got.MarkDone()
	if !file.Exists(expect) {
		t.Errorf("expect: %s got: %v", expect, false)
	}
}
