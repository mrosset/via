package via

import (
	"os"
	"path"
)

type Cache string

func (c Cache) Packages() string {
	return path.Join(string(c), "pkg")
}

func (c Cache) Sources() string {
	return path.Join(string(c), "src")
}

func (c Cache) Builds() string {
	return path.Join(string(c), "bld")
}

func (c Cache) Stages() string {
	return path.Join(string(c), "stg")
}

func (c Cache) String() string {
	return os.ExpandEnv(string(c))
}

func (c Cache) Init() {
	for _, path := range []string{c.Packages(), c.Sources(), c.Builds(), c.Stages()} {
		fatal(os.MkdirAll(path, 0755))
	}
}
