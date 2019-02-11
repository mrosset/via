package via

import (
	"os"
	"path/filepath"
)

// Cache provides a type for working with build caching
type Cache string

// Dir returns a directory path joined to the end of the Cache path
func (c Cache) Dir(dir string) string {
	return filepath.Join(string(c), dir)
}

// Packages returns the cache directory containing build
// packages. This directory is used to install the built package
// before they are packaged into tarballs
func (c Cache) Packages() string {
	return filepath.Join(string(c), "pkg")
}

// Sources returns the directory cache that contains the source
// tarballs downloaded from upstream. This can also contain git
// repositories though only portions of git support currently exists
func (c Cache) Sources() string {
	return filepath.Join(string(c), "src")
}

// Builds returns the directory cache that contains plan builds. This
// is where out of source tree builds are built. Not all build systems
// support plans support out of source tree builds
func (c Cache) Builds() string {
	return filepath.Join(string(c), "bld")
}

// Stages returns the stages directory cache. This directory is used
// to cache decompressed source trees.
func (c Cache) Stages() string {
	return filepath.Join(string(c), "stg")
}

// String provides a stringer interface to convert Cache back to a
// string. The resulting string has environment variables expanded.
func (c Cache) String() string {
	return string(c.Expand())
}

// Expand returns a Cache that has had environmental variables
// expanded. Generally this is used for allowing us to specify
// $HOME/.cache/via in our Config file. though this might change and
// be removed at some point.
func (c Cache) Expand() Cache {
	return Cache(
		os.ExpandEnv(string(c)),
	)
}

// Init creates each cache directory ensuring it exists
func (c Cache) Init() {
	for _, path := range []string{c.Packages(), c.Sources(), c.Builds(), c.Stages()} {
		fatal(os.MkdirAll(path, 0755))
	}
}
