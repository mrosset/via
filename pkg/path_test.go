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
        tests{
                {
                        Expect: "testdata",
                        Got:    Path("testdata").Expand(),
                },
        }.equals(t)
}

func TestPath_Join(t *testing.T) {
        test{
                Expect: Path("testdata/join"),
                Got:    Path("testdata").Join("join"),
        }.equals(t)
}
