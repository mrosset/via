package via

import (
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	cache   Cache
	viapath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/mrosset/via")
	cfile   = filepath.Join(viapath, "plans/config.json")
	viaUrl  = "https://github.com/mrosset/via"
	planUrl = "https://github.com/mrosset/plans"
	config  = new(Config)
)

func init() {
	if os.Getenv("GOPATH") == "" {
		elog.Fatal("GOPATH must be set")
	}
	// TODO rework this to error and suggest user use 'via init'
	if !file.Exists(viapath) {
		elog.Println("cloning plans")
		if err := Clone(viapath, viaUrl); err != nil {
			elog.Fatal(err)
		}
	}
	pdir := filepath.Dir(cfile)
	if !file.Exists(pdir) {
		elog.Println("cloning plans")
		err := Clone(pdir, planUrl)
		if err != nil {
			elog.Fatal(err)
		}
	}
}

func init() {
	err := json.Read(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}
	// TODO: provide Lint for master config
	sort.Strings([]string(config.Flags))
	sort.Strings(config.Remove)
	err = json.Write(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}

	config = config.Expand()

	// if err := CheckLink(); err != nil {
	//	elog.Fatal(err)
	// }

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
	Threads   int
	IpfsApi   string
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
