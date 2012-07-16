package via

import (
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"log"
	"os"
	"path"
	"strings"
)

var (
	cache  Cache
	home   = os.Getenv("HOME")
	cfile  = path.Join(home, "via.json")
	config = &Config{
		Identity:  "test user <test@test.com>",
		Root:      "/",
		PlansRepo: "https://code.google.com/p/via.plans",
		Cache:     "$HOME/via/cache",
		DB:        "/usr/local/via",
		Plans:     "$HOME/via/plans",
		Repo:      "$HOME/via/repo",
		Flags: []string{
			"--host=arm-linux-gnueabi",
			"--prefix=/data/data/gnuoid",
			"-q",
		},
	}
	join = path.Join
)

func init() {
	os.Setenv("CC", "arm-linux-gnueabi-gcc")
	os.Setenv("PATH", os.Getenv("PATH")+":/opt/tools/bin")
	if !file.Exists(cfile) {
		err := json.Write(&config, cfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := json.Read(&config, cfile)
	if err != nil {
		log.Fatal(err)
	}
	cache = Cache(os.ExpandEnv(string(config.Cache)))
	config.Plans = os.ExpandEnv(config.Plans)
	config.Repo = os.ExpandEnv(config.Repo)
}

type Config struct {
	Identity  string
	Root      string
	PlansRepo string

	// Paths
	Cache Cache
	DB    DB
	Plans string
	Repo  string

	// Gnu
	Flags Flags
}

type Flags []string

func (f Flags) String() string {
	return strings.Join(f, " ")
}

type DB string

func (d DB) Installed() string {
	return path.Join(string(d), "installed")
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
