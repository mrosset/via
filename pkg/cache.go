package via

import (
	"os"
)

// Cache provides a type for working with build caching
type Cache struct {
	Path
}

// NewCache returns a new initialized Cache
func NewCache(path string) Cache {
	return Cache{
		Path: Path(path),
	}
}

// Packages returns the cache directory containing build
// packages. This directory is used to install the built package
// before they are packaged into tarballs
func (c Cache) Packages() string {
	return c.Join("pkg")
}

// Sources returns the directory cache that contains the source
// tarballs downloaded from upstream. This can also contain git
// repositories though only portions of git support currently exists
func (c Cache) Sources() string {
	return c.Join("src")
}

// Builds returns the directory cache that contains plan builds. This
// is where out of source tree builds are built. Not all build systems
// support plans support out of source tree builds
func (c Cache) Builds() string {
	return c.Join("bld")
}

// Stages returns the stages directory cache. This directory is used
// to cache decompressed source trees.
func (c Cache) Stages() string {
	return c.Join("stg")
}

// Init creates each cache directory ensuring it exists
func (c Cache) Init() {
	for _, path := range []string{c.Packages(), c.Sources(), c.Builds(), c.Stages()} {
		fatal(os.MkdirAll(path, 0755))
	}
}
