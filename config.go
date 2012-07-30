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
			//"--target=arm-linux-gnueabi",
			"--host=arm-linux-gnueabi",
			"--disable-multilib",
			"--disable-dependency-tracking",
			"--disable-nls",
			"--with-shared",
			"--prefix=$PREFIX",
			"-q",
		},
		Env: map[string]string{
			"PREFIX":      "/data/data/gnuoid",
			"MAKEFLAGS":   "-j3 -sw",
			"CFLAGS":      "-pipe -O2 -static-libgcc -I/data/data/gnuoid/include",
			"LD_RUN_PATH": "/data/data/gnuoid/lib",
			"LDFLAGS":     "-L/data/data/gnuoid/lib",
		},
		Remove: []string{
			"data/data/gnuoid/share/doc",
			"data/data/gnuoid/share/info",
			"data/data/gnuoid/share/man",
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
