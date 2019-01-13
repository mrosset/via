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
		is   = is.New(t)
		p, _ = NewPlan("emacs")
		got  = NewBatch(testConfig)
	)
	got.Walk(p)
	is.Equal(len(got.Plans), 81)
}

func TestBatchInstall(t *testing.T) {
	var (
		p, _   = NewPlan("make")
		got    = NewBatch(testConfig)
		expect = join(testConfig.Repo, "repo", p.PackageFile())
	)
	defer os.RemoveAll(testConfig.Repo)
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
