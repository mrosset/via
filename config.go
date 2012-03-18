package via

import (
	"os"
	"path"
	"util"
	"util/file"
	"util/json"
)

var (
	config = &Config{}
	checkf = util.CheckFatal
)

type Config struct {
	Arch     string
	Identity string
	OS       string
	Prefix   string
	Repo     string
	Root     string
	Cache    Cache
	Plans    string
	DB       string
}

func init() {
	checkf(os.Setenv("CC", "gcc"))
	cfile := path.Join(os.Getenv("HOME"), ".via.json")
	checkf(json.Read(&config, cfile))
	checkf(config.Cache.Create())
	for _, dir := range []string{config.Repo, config.DB} {
		if !file.Exists(dir) {
			info("mkdir", dir)
			checkf(os.MkdirAll(dir, 0755))
		}
	}
}

func (c *Config) GetStageDir(name string) string {
	return path.Join(c.Cache.Stages(), name)
}

func (c *Config) GetBuildDir(name string) string {
	return path.Join(c.Cache.Builds(), name)
}

func (c *Config) GetPackageDir(name string) string {
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

func (c Cache) Create() error {
	paths := []string{
		string(c),
		string(c.Builds()),
		string(c.Stages()),
		string(c.Sources()),
		string(c.Packages()),
		string(c.Packages()),
	}
	for _, d := range paths {
		if !file.Exists(d) {
			info("mkdir", d)
			if err := os.Mkdir(d, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
