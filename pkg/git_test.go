package via

import (
	"testing"
)

func TestClone(t *testing.T) {
	t.Parallel()
	var (
		gitd = Path("testdata/via-test")
	)
	defer gitd.RemoveAll()
	tests{
		{
			Expect: nil,
			Got:    gitd.Clone("https://github.com/mrosset/via-test"),
		},
		{
			Expect: true,
			Got:    gitd.Join("README.md").Exists(),
		},
	}.equals(t)
}

func TestCheckout(t *testing.T) {
}

func TestCloneBranch(t *testing.T) {
	var (
		gitd  = Path("testdata/plans.git")
		plans = Path("..").Join("plans")
	)
	tests{
		{
			Name:   "Plans exists",
			Expect: true,
			Got:    plans.Exists(),
		},
		{
			Name:   "CloneBranch",
			Expect: nil,
			Got: CloneBranch(
				gitd,
				plans.String(),
				"x86_64-via-linux-gnu-release",
			),
		},
	}.equals(t)
}

func TestBranch(t *testing.T) {
	var (
		gitd = Path("testdata/plans.git")
	)
	defer gitd.RemoveAll()
	got, err := Branch(gitd)
	tests{
		{
			Name:   "Branch error",
			Expect: nil,
			Got:    err,
		},
		{
			Name:   "Test branch",
			Expect: "x86_64-via-linux-gnu-release",
			Got:    got,
		},
	}.equals(t)
}
