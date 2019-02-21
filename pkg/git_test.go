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

func TestCloneBranch(t *testing.T) {
	t.Parallel()
	var (
		gitd = Path("testdata/plans.git")
	)
	defer gitd.RemoveAll()
	test{
		Expect: nil,
		Got: CloneBranch(
			gitd,
			"../plans",
			"x86_64-via-linux-gnu-release",
		),
	}.equals(t)
	got, err := Branch(gitd)
	if err != nil {
		t.Fatal(err)
	}
	test{
		Expect: "x86_64-via-linux-gnu-release",
		Got:    got,
	}.equals(t)
}
