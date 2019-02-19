package via

import (
        "github.com/mrosset/util/file"
        "testing"
)

func TestClone(t *testing.T) {
        t.Parallel()
        var (
                expect = "testdata/git/README"
                gitd   = Path("testdata/git")
        )
        defer gitd.RemoveAll()
        if err := Clone(gitd, "https://github.com/mrosset/gur"); err != nil {
                t.Fatal(err)
        }
        if !file.Exists(expect) {
                t.Errorf("exected %s but file does not exist", expect)
        }
        expect = "master"
        got, err := Branch(gitd)
        if err != nil {
                t.Fatal(err)
        }
        if expect != got {
                t.Logf("expect '%s' got '%s'", expect, got)
        }
}
