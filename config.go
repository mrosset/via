package via

import (
	"os"
	"path"
	"util"
	"util/file"
	"util/json"
)

var (
	home   = os.Getenv("HOME")
	config = ReadConfig()
	checkf = util.CheckFatal

	// dir aliases
	builds      = config.Cache.Dir("builds")
	stages      = config.Cache.Dir("stages")
	packages    = config.Cache.Dir("packages")
	sources     = config.Cache.Dir("sources")
	installed   = config.DB.Dir("installed")
	plans       = config.Home.Dir("plans")
	repo        = config.Home.Dir("repo")
	config_dirs = []Tree{}
)

type Config struct {
	Identity  string
	Root      string
	PlansRepo string
	Cache     *Tree
	DB        *Tree
	Home      *Tree
}

func init() {
	checkf(os.Setenv("CC", "gcc"))
	cfile := path.Join(home, ".via.json")
	checkf(json.Read(&config, cfile))
	checkf(json.Write(&config, cfile))
	config_dirs = append(config_dirs, builds, stages, sources, installed, plans, repo)
	createDirs()
}

func ReadConfig() *Config {
	cfile := path.Join(home, ".via.json")
	c := new(Config)
	checkf(json.Read(&c, cfile))
	return c
}

func createDirs() {
	for _, d := range config_dirs {
		if !file.Exists(d.String()) {
			checkf(os.MkdirAll(d.String(), 0755))
		}
	}
}

type Tree string

func (t Tree) Dir(n string) Tree {
	return Tree(path.Join(string(t), n))
}

func (t Tree) String() string {
	return string(t)
}

func (t Tree) File(n string) string {
	return path.Join(string(t), n)
}
