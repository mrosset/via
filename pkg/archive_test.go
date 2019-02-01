package via

import (
	"compress/gzip"
	"fmt"
	"github.com/mrosset/util/file"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// FIXME: create a test to verify we can extract long file names.
func TestLongNames(t *testing.T) {
	var (
		longName = fmt.Sprintf("L%sng", strings.Repeat("o", 98))
		archPath = filepath.Join("testdata", "archive")
		longPath = filepath.Join(archPath, longName)
		gzipPath = filepath.Join("testdata", "archive.tar.gz")
	)
	if len(longName) != 101 {
		t.Fatalf("Expect longName length of 101 got %d", len(longName))
	}
	if !file.Exists(archPath) {
		os.Mkdir(archPath, 0700)
	}

	fd, err := os.Create(longPath)
	if err != nil {
		t.Fatal(err)
	}
	fd.Close()

	fd, err = os.Create(gzipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	err = archive(gz, "testdata/archive")
	if err != nil {
		t.Error(err)
	}
}
