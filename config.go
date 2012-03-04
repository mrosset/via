package via

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"util"
	"util/file"
)

var (
	config = &Config{cache: "cache", Prefix: "/usr/local", Root: "/", plans: "plans", Installed: "installed",Identity:"Mike Rosset <mike.rosset@gmail.com>"}
)

type Config struct {
	OS        string
	Arch      string
	cache     string
	Prefix    string
	plans     string
	Root      string
	Installed string
	Identity  string
}

func InitConfig() {
	dirs := []string{config.Cache(), config.Sources(), config.Builds(), config.Packages(), config.Stages(), config.Plans()}
	for _, d := range dirs {
		if !file.Exists(d) {
			fmt.Printf("%-20s %s\n", "creating", d)
			err := os.Mkdir(d, 0775)
			util.CheckFatal(err)
			continue
		}
		wd, err := os.Getwd()
		util.CheckFatal(err)
		rel, err := filepath.Rel(wd, d)
		util.CheckFatal(err)
		if Verbose {
			fmt.Printf("%-20s %s\n", rel, "OK")
		}
	}
}

func (c *Config) Cache() string {
	dir, err := filepath.Abs(c.cache)
	util.CheckFatal(err)
	return dir
}

func (c *Config) Plans() string {
	dir, err := filepath.Abs(c.plans)
	util.CheckFatal(err)
	return dir
}

func (c *Config) Sources() string {
	return path.Join(c.Cache(), "sources")
}

func (c *Config) Builds() string {
	return path.Join(c.Cache(), "builds")
}

func (c *Config) Stages() string {
	return path.Join(c.Cache(), "stages")
}

func (c *Config) Packages() string {
	return path.Join(c.Cache(), "packages")
}

func (c *Config) GetStageDir(name string) string {
	return path.Join(c.Stages(), name)
}

func (c *Config) GetBuildDir(name string) string {
	return path.Join(c.Builds(), name)
}

func (c *Config) GetPackageDir(name string) string {
	return path.Join(c.Packages(), name)
}
