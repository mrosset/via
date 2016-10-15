package via

import (
	"github.com/str1ngs/util/file"
	"os"
	"testing"
)

func TestClone(t *testing.T) {
	var (
		expected = []string{
			"testdata/via",
			"testdata/via/Makefile",
			"testdata/via/plans",
			"testdata/via/plans/config.json",
		}
	)
	defer os.RemoveAll("testdata/via")
	err := clone("testdata", "https://github.com/mrosset/via")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range expected {
		if !file.Exists(p) {
			t.Errorf("path %s does not exist.", p)
		}
	}
}
