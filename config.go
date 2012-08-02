package via

import (
	"os"
	"path"
	"strings"
)

var (
	cache  Cache
	home   = os.Getenv("HOME")
	cfile  = path.Join(home, "via.json")
	config = &Config{
		Arch:      "amd64",
		OS:        "linux",
		Cache:     "$HOME/via/cache",
		DB:        "/usr/local/via/var/db/via",
		Identity:  "test user <test@test.com>",
		Plans:     "$HOME/via/plans",
		PlansRepo: "https://code.google.com/p/via.plans",
		Repo:      "$HOME/via/repo",
		Root:      "/",
		Flags: []string{
			"--disable-multilib",
			"--disable-dependency-tracking",
			"--disable-nls",
			"--with-shared",
			"--prefix=$PREFIX",
			"-q",
		},
		Env: map[string]string{
			"CC": "clang",
			"PREFIX":      "/usr/local/via",
			"MAKEFLAGS":   "-j3 -sw",
			"CFLAGS":      "-pipe -O2",
		},
		Remove: []string{
		},
	}
	join = path.Join
)

func init() {
	/*
		if !file.Exists(cfile) {
			elog.Println("WARNING no config was found writing new one to", cfile)
			elog.Println("please review it.")
			err := json.Write(&config, cfile)
			if err != nil {
				elog.Fatal(err)
			}
			return
		}
		config = &Config{}
		err := json.Read(&config, cfile)
		if err != nil {
			elog.Fatal(err)
		}
		// TODO: provide Lint for master config
		err = json.Write(&config, cfile)
		if err != nil {
			elog.Fatal(err)
		}
	*/
	cache = Cache(os.ExpandEnv(string(config.Cache)))
	config.Plans = os.ExpandEnv(config.Plans)
	config.Repo = os.ExpandEnv(config.Repo)
	for i, j := range config.Env {
		os.Setenv(i, j)
	}
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
