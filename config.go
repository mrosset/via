package via

import (
	"os"
	"path"
	"util"
	"util/json"
)

var (
	home           = os.Getenv("HOME")
	config         = &Config{}
	checkf         = util.CheckFatal
	cache          *Cache
	db             DB
	config_default = &Config{
		Repo: path.Join(home, ".via.json"),
	}
)

type Config struct {
	Identity  string
	Repo      string
	Root      string
	Plans     string
	PlansRepo string
	Cache     *Cache
	DB        DB
}

func init() {
	checkf(os.Setenv("CC", "gcc"))
	cfile := path.Join(os.Getenv("HOME"), ".via.json")
	checkf(json.Read(&config, cfile))
}

func (c *Config) StageDir(name string) string {
	return path.Join(c.Cache.Stages(), name)
}

func (c *Config) BuildDir(name string) string {
	return path.Join(c.Cache.Builds(), name)
}

func (c *Config) PackageDir(name string) string {
	return path.Join(c.Cache.Packages(), name)
}

type Cache string

func (c Cache) Builds() string {
	return path.Join(string(c), "builds")
}

func (c Cache) Stages() string {
	return path.Join(string(c), "stages")
}

func (c Cache) Sources() string {
	return path.Join(string(c), "sources")
}

func (c Cache) Packages() string {
	return path.Join(string(c), "packages")
}

type DB string

func (d DB) Installed() string {
	return path.Join(string(d), "installed")
}
