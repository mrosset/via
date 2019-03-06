package via

import (
	. "github.com/mrosset/test"
	"testing"
)

func TestClone(t *testing.T) {
	t.Parallel()
	var (
		gitd = Path("testdata/git-test/clone")
	)
	defer gitd.RemoveAll()
	Tests{
		{
			Expect: nil,
			Got:    gitd.Clone("https://github.com/mrosset/via-test"),
		},
		{
			Expect: true,
			Got:    gitd.Join("README.md").Exists(),
		},
	}.Equals(t)
}

func TestCheckout(t *testing.T) {
	t.Parallel()
	var (
		gitd = Path("testdata/git-test/checkout")
	)
	defer gitd.RemoveAll()
	Tests{
		{
			Expect: nil,
			Got:    gitd.Clone("https://github.com/mrosset/via-test"),
		},
		{
			Expect: nil,
			Got:    Checkout(gitd, "refs/heads/master"),
		},
	}.Equals(t)
}

func TestCloneBranch(t *testing.T) {
	t.Parallel()
	var (
		gitd   = Path("testdata/git-test/clone-branch")
		origin = "https://github.com/mrosset/via-test"
	)
	defer gitd.RemoveAll()
	Tests{
		{
			Name:   "CloneBranch",
			Expect: nil,
			Got: CloneBranch(
				gitd,
				origin,
				"x86_64-via-linux-gnu-release",
			),
		},
	}.Equals(t)
	got, err := Branch(gitd)
	Tests{
		{
			Name:   "Clone Branch error",
			Expect: nil,
			Got:    err,
		},
		{
			Name:   "Test branch",
			Expect: "x86_64-via-linux-gnu-release",
			Got:    got,
		},
	}.Equals(t)

}
