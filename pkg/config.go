package via

import (
	"bitbucket.org/strings/via/pkg/git"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	cache  Cache
	gopath = filepath.Join(os.Getenv("GOPATH"), "src/bitbucket.org/strings/via")
	cfile  = filepath.Join(gopath, "plans/config.json")
	viaUrl = "https://github.com/mrosset/via"
	config = new(Config)
)

func init() {
	sync := false
	// TODO rework this to error and suggest user use 'via init'
	if !file.Exists(gopath) {
		err := clone(gopath, viaUrl)
		if err != nil {
			elog.Fatal(err)
		}
	}
	err := json.Read(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}
	if sync {
		PlanSync()
		os.Remove("/tmp/config.json")
	}
	sort.Strings([]string(config.Flags))
	sort.Strings(config.Remove)
	// TODO: provide Lint for master config
	err = json.Write(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}
	cache = Cache(os.ExpandEnv(string(config.Cache)))
	cache.Init()
	config.Plans = os.ExpandEnv(config.Plans)
	config.Repo = os.ExpandEnv(config.Repo)
	os.MkdirAll(config.Repo, 0755)
	for i, j := range config.Env {
		os.Setenv(i, os.ExpandEnv(j))
	}
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

	Env         map[string]string
	Remove      []string
	PostInstall []string
}

func (c Config) Branch() (string, error) {
	p := filepath.Join(c.Plans, "../.git/modules/plans")
	return git.Branch(p)
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
