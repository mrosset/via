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
		got  = NewBatch(config)
	)
	got.Walk(p)
	is.Equal(len(got.Plans), 80)
}
