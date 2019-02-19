package via

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
func (c Cache) Packages() Path {
        return c.Join("packages")
}

// Sources returns the directory cache that contains the source
// tarballs downloaded from upstream. This can also contain git
// repositories though only portions of git support currently exists
func (c Cache) Sources() Path {
        return c.Join("sources")
}

// Builds returns the directory cache that contains plan builds. This
// is where out of source tree builds are built. Not all build systems
// support plans support out of source tree builds
func (c Cache) Builds() Path {
        return c.Join("builds")
}

// Stages returns the stages directory cache. This directory is used
// to cache decompressed source trees.
func (c Cache) Stages() Path {
        return c.Join("stages")
}

// Init creates each cache directory ensuring it exists
func (c Cache) Init() {
        for _, path := range []Path{c.Packages(), c.Sources(), c.Builds(), c.Stages()} {
                fatal(path.Ensure())
        }
}
