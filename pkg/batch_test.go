package via

import (
	"github.com/cheekybits/is"
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
	is.Equal(len(got.Plans), 80)
}

func TestBatchDownload(t *testing.T) {
	var (
		p, _ = NewPlan("devel")
		got  = NewBatch(config)
	)
	got.Walk(p)
	errors := got.Download()
	if len(errors) != 0 {
		t.Error(errors)
	}
}

func TestBatchInstall(t *testing.T) {
	var (
		p, _ = NewPlan("devel")
		got  = NewBatch(config)
	)
	got.Walk(p)
	errors := got.Install()
	if len(errors) != 0 {
		t.Error(errors)
	}

}
