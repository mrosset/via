package via

import (
	"github.com/mrosset/util/file"
	"os"
	gpath "path"
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
	return Path(gpath.Join(spaths...))
}

// like filepath.Join but joins the strings to this Path then returns a Path
func (path Path) JoinS(paths ...string) Path {
	paths = append([]string{string(path)}, paths...)
	return Path(gpath.Join(paths...))
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

// Returns a UNIX path. On windows this will convert c:/ to /c/
// Genereally this is only needed when passing a PATH to a msys2 tool
func (path Path) ToUnix() string {
	spath := path.String()
	if len(spath) == 0 {
		return ""
	}
	if spath[0] == '/' {
		return spath
	}
	spath = toUnixSlash(spath)
	switch spath[:2] {
		case "c:", "C:":
		return gpath.Join("/c",spath[2:])
	}
	return spath
}

// Returns a string replaces all \ with /
func toUnixSlash(path string) string {
	var (
		npath = ""
	)
	for _, c := range path {
		if c == '\\' {
			npath = npath + string('/')
		} else {
			npath = npath+string(c)
		}
	}
	return npath
}

// Provides stringer interface. This expands and returns this Path as a string
// this method should always be used when converting to a string.
func (path Path) String() string {
	return os.ExpandEnv(string(path))
}
