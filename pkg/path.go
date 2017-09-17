package via

import (
	"github.com/mrosset/util/file"
	"os"
	"path/filepath"
)

type Path string

// Join's this path with subsequent Paths
func (path Path) Join(paths ...Path) Path {
	var (
		spaths = []string{}
	)
	spaths = append(spaths, path.String())
	for _, i := range paths {
		spaths = append(spaths, i.String())
	}
	return Path(filepath.Join(spaths...))
}

// like filepath.Join but joins the strings to this Path then returns a Path
func (path Path) JoinS(paths ...string) Path {
	paths = append([]string{string(path)}, paths...)
	return Path(filepath.Join(paths...))
}

// Returns true of this Path exists
func (path Path) Exists() bool {
	return file.Exists(path.String())
}

// Provides stringer interface. This expands and returns this Path as a string
// this method should always be used when converting to a string.
func (path Path) String() string {
	return os.ExpandEnv(string(path))
}
