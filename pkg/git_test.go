package via

import (
	"testing"
)

func TestClone(t *testing.T) {
	t.Parallel()
	var (
		gitd = Path("testdata/git-test/clone")
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
	t.Parallel()
	var (
		gitd = Path("testdata/git-test/checkout")
	)
	defer gitd.RemoveAll()
	tests{
		{
			Expect: nil,
			Got:    gitd.Clone("https://github.com/mrosset/via-test"),
		},
		{
			Expect: nil,
			Got:    Checkout(gitd, "refs/heads/master"),
		},
	}.equals(t)
}

func TestCloneBranch(t *testing.T) {
	t.Parallel()
	var (
		gitd = Path("testdata/git-test/clone-branch")
	)
	defer gitd.RemoveAll()
	tests{
		{
			Name:   "CloneBranch",
			Expect: nil,
			Got: CloneBranch(
				gitd,
				"https://github.com/mrosset/via-test",
				"x86_64-via-linux-gnu-release",
			),
		},
	}.equals(t)

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
