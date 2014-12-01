package via

import (
	"os"
	"path"
)

type Cache string

func (c Cache) Pkgs() string {
	return path.Join(string(c), "pkg")
}

func (c Cache) Srcs() string {
	return path.Join(string(c), "src")
}

func (c Cache) Builds() string {
	return path.Join(string(c), "bld")
}

func (c Cache) Stages() string {
	return path.Join(string(c), "stg")
}

func (c Cache) Init() {
	for _, path := range []string{c.Pkgs(), c.Srcs(), c.Builds(), c.Stages()} {
		fatal(os.MkdirAll(path, 0755))
	}
}
