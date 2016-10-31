package via

import (
	"errors"
	"fmt"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"github.com/mrosset/via/pkg/git"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	ERR_BRANCH_MISMATCH = "Branches do not Match"
)

var (
	cache  Cache
	gopath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/mrosset/via")
	cfile  = filepath.Join(gopath, "plans/config.json")
	viaUrl = "https://github.com/mrosset/via"
	config = new(Config)
)

func init() {
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
	err = os.MkdirAll(config.Repo, 0755)
	if err != nil {
		elog.Fatal(err)
	}
	for i, j := range config.Env {
		os.Setenv(i, os.ExpandEnv(j))
	}
	for i, j := range config.Env {
		os.Setenv(i, os.ExpandEnv(j))
	}
}

type Config struct {
	Branch    string
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
	Prefix string

	// Toolchain
	Flags Flags

	Env         map[string]string
	Remove      []string
	PostInstall []string

	// Internal Fields
	template *Config
}

func (c *Config) Expand() *Config {
	if c.template != nil {
		return c.template
	}
	o := new(Config)
	err := json.Parse(o, c)
	if err != nil {
		panic(err)
	}
	c.template = o
	return o
}

// Checks all branches match the Config branch
func (c Config) CheckBranches() error {
	if c.PlanBranch() != config.Branch || c.RepoBranch() != config.Branch {
		msg := fmt.Sprintf("%s: %s %s", ERR_BRANCH_MISMATCH, config.Branch)
		return (errors.New(msg))
	}
	return nil
}

// Returns the checked out branch for repo directory
func (c Config) RepoBranch() string {
	b, err := git.Branch(c.Repo)
	if err != nil {
		elog.Fatalf("%s %s", c.Repo, err)
	}
	return b
}

// Returns the checked out branch for plans directory
func (c Config) PlanBranch() string {
	p := filepath.Join(c.Plans)
	b, err := git.Branch(p)
	if err != nil {
		elog.Fatal(err)
	}
	return b
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
