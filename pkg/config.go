package via

import (
	"github.com/str1ngs/util/json"
	"os"
	"path"
	"sort"
	"strings"
)

var (
	cache  Cache
	cfile  = path.Join(os.Getenv("HOME"), "via", "plans", "config.json")
	config = new(Config)
)

func init() {
	err := json.Read(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}

	sort.Strings([]string(config.Flags))
	sort.Strings(config.Remove)
	// TODO: provide Lint for master config
	err = json.Write(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}
	cache = Cache(os.ExpandEnv(string(config.Cache)))
	config.Plans = os.ExpandEnv(config.Plans)
	config.Repo = os.ExpandEnv(config.Repo)
	for i, j := range config.Env {
		os.Setenv(i, os.ExpandEnv(j))
	}
}

type Config struct {
	Identity  string
	Arch      string
	OS        string
	Root      string
	PlansRepo string

	// Paths
	Cache  Cache
	DB     DB
	Plans  string
	Repo   string
	Binary string

	// Toolchain
	Flags Flags

	Env    map[string]string
	Remove []string
}

type Flags []string

func (f Flags) String() string {
	return strings.Join(f, " ")
}

type DB string

func (d DB) Installed() string {
	return join(config.Root, string(d), "installed")
}

func (d DB) Plans() string {
	return join(config.Root, string(d), "plans")
}

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
