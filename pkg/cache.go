package via

import (
	"os"
	"path/filepath"
)

type Cache string

func (c Cache) Packages() string {
	return filepath.Join(string(c), "pkg")
}

func (c Cache) Sources() string {
	return filepath.Join(string(c), "src")
}

func (c Cache) Builds() string {
	return filepath.Join(string(c), "bld")
}

func (c Cache) Stages() string {
	return filepath.Join(string(c), "stg")
}

func (c Cache) String() string {
	return os.ExpandEnv(string(c))
}

func (c Cache) Expand() Cache {
	return Cache(
		os.ExpandEnv(string(c)),
	)
}

func (c Cache) Init() {
	for _, path := range []string{c.Packages(), c.Sources(), c.Builds(), c.Stages()} {
		fatal(os.MkdirAll(path, 0755))
	}
}
