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
		Arch:      "arm",
		OS:        "linux",
		Cache:     "$HOME/via/cache",
		DB:        "/data/data/gnuoid/var/db/via",
		Identity:  "test user <test@test.com>",
		Plans:     "$HOME/via/plans",
		PlansRepo: "https://code.google.com/p/via.plans",
		Repo:      "$HOME/via/repo",
		Root:      "/",
		Flags: []string{
			"--host=arm-linux-gnueabi",
			"--prefix=/data/data/gnuoid",
			"-q",
		},
	}
	join = path.Join
)

func init() {
	os.Setenv("CC", "arm-linux-gnueabi-gcc -pipe -O2")
	os.Setenv("PATH", os.Getenv("PATH")+":/opt/tools/bin")
	os.Setenv("MAKEFLAGS", "-j3  -sw")
	os.Setenv("LDFLAGS", "-Wl,-rpath -Wl,/data/data/gnuoid/lib")
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
	// TODO: provide Lint for master config
	err = json.Write(&config, cfile)
	if err != nil {
		log.Fatal(err)
	}
	cache = Cache(os.ExpandEnv(string(config.Cache)))
	config.Plans = os.ExpandEnv(config.Plans)
	config.Repo = os.ExpandEnv(config.Repo)
}

type Config struct {
	Identity  string
	Arch      string
	OS        string
	Root      string
	PlansRepo string

	// Paths
	Cache Cache
	DB    DB
	Plans string
	Repo  string

	// Toolchain
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
