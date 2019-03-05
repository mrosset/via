package via

import (
	. "github.com/mrosset/via/pkg/test"
	"os"
	"testing"
)

func TestPath_String(t *testing.T) {
	Tests{
		{
			Expect: "testdata/plans/core/hello.json",
			Got:    Path("testdata/plans/core/hello.json").String(),
		},
	}.Equals(t)
}

func TestPath_Expand(t *testing.T) {
	os.Setenv("_PATH", "testdata")
	Test{
		Expect: Path("testdata"),
		Got:    Path("$_PATH").Expand(),
	}.Equals(t)
}

func TestPath_Join(t *testing.T) {
	Test{
		Expect: Path("testdata/join"),
		Got:    Path("testdata").Join("join"),
	}.Equals(t)
}

func TestPath_NewPath(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			Test{
				Expect: os.ErrNotExist,
				Got:    r.(error),
			}.Equals(t)
		}
	}()
	Tests{
		{
			Expect: Path("testdata/repo"),
			Got:    NewPath("testdata", "repo"),
		},
		{
			Expect: nil,
		},
	}.Equals(t)
	NewPath("testdata", "fail")
}

func TestPath_ToPath(t *testing.T) {
	Test{
		Expect: Path("testdata/plans"),
		Got:    testConfig.Plans.ToPath(),
	}.Equals(t)
}

func TestPath_Glob(t *testing.T) {
	var (
		top = Path("testdata/glob")
		one = top.Join("one")
		two = top.Join("two")
	)
	top.MkdirAll()
	one.Touch()
	two.Touch()
	defer top.RemoveAll()
	glob, err := top.Glob()

	Tests{
		{
			Expect: nil,
			Got:    err,
		},
		{
			Expect: []Path{"testdata/glob/one", "testdata/glob/two"},
			Got:    glob,
		},
	}.Equals(t)
}
