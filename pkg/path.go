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

// Symlinks this path to 'new' Path
func (path Path) Symlink(new Path) error {
	elog.Println("symlink", path, new)
	return os.Symlink(path.String(), new.String())
}

// Makes directory's in this Path with 'mode'
func (path Path) MkDirAll(mode os.FileMode) error {
	return os.MkdirAll(path.String(), mode)
}

// Returns true if this Path exists
func (path Path) Exists() bool {
	return file.Exists(path.String())
}

// Provides stringer interface. This expands and returns this Path as a string
// this method should always be used when converting to a string.
func (path Path) String() string {
	return os.ExpandEnv(string(path))
}
