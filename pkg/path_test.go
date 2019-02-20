package via

import (
	"os"
	"testing"
)

func TestPath_String(t *testing.T) {
	tests{
		{
			Expect: "testdata/plans/core/hello.json",
			Got:    Path("testdata/plans/core/hello.json").String(),
		},
	}.equals(t)
}

func TestPath_Expand(t *testing.T) {
	os.Setenv("_PATH", "testdata")
	test{
		Expect: Path("testdata"),
		Got:    Path("$_PATH").Expand(),
	}.equals(t)
}

func TestPath_Join(t *testing.T) {
	test{
		Expect: Path("testdata/join"),
		Got:    Path("testdata").Join("join"),
	}.equals(t)
}

func TestPath_NewPath(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			test{
				Expect: os.ErrNotExist,
				Got:    r.(error),
			}.equals(t)
		}
	}()
	tests{
		{
			Expect: Path("testdata/repo"),
			Got:    NewPath("testdata", "repo"),
		},
		{
			Expect: nil,
		},
	}.equals(t)
	NewPath("testdata", "fail")
}

func TestPath_ToPath(t *testing.T) {
	test{
		Expect: Path("testdata/plans"),
		Got:    testConfig.Plans.ToPath(),
	}.equals(t)
}
