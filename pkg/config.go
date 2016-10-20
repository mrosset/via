package via

import (
	"errors"
	"fmt"
	"github.com/mrosset/via/pkg/git"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	// Toolchain
	Flags Flags

	Env         map[string]string
	Remove      []string
	PostInstall []string
}

const (
	ERR_BRANCH_MISMATCH = "Branches do not Match"
)

// Checks that the plan branch and the publish repo branch match the configured
// branch
func (c Config) CheckBranches() error {
	if c.PlanBranch() != config.Branch {
		msg := fmt.Sprintf("(%s) (%s) %s", c.PlanBranch(), config.Branch, ERR_BRANCH_MISMATCH)
		return (errors.New(msg))
	}
	if c.RepoBranch() != config.Branch {
		msg := fmt.Sprintf("(%s) (%s) %s", c.RepoBranch(), config.Branch, ERR_BRANCH_MISMATCH)
		return (errors.New(msg))
	}
	return nil
}

// Returns the checked out branch for plans directory
func (c Config) RepoBranch() string {
	p := filepath.Join(c.Repo, "/../.git")
	b, err := git.Branch(p)
	if err != nil {
		elog.Fatalf("%s %s", c.Repo, err)
	}
	return b
}

// Returns the checked out branch for plans directory
func (c Config) PlanBranch() string {
	p := filepath.Join(c.Plans, "../.git/modules/plans")
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
