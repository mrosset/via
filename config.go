package via

import (
	"log"
	"os"
	"path"
	"util/file"
	"util/json"
)

var (
	config   *Config
	blds     string
	inst     string
	pkgs     string
	srcs     string
	stgs     string
	home     = os.Getenv("HOME")
	cfile    = path.Join(home, "via.json")
	defaults = &Config{
		Cache:     "/home/strings/via/cache",
		DB:        "/usr/local/via",
		Identity:  "test user <test@test.com>",
		Plans:     "/home/strings/via/plans",
		PlansRepo: "https://code.google.com/p/via.plans",
		Repo:      "/home/strings/via/repo",
		Root:      "/",
	}
	join = path.Join
)

func init() {
	if !file.Exists(cfile) {
		err := json.Write(&defaults, cfile)
		if err != nil {
			log.Fatal(err)
		}
		config = defaults
		return
	}
	err := json.Read(&config, cfile)
	if err != nil {
		log.Fatal(err)
	}

	blds = join(config.Cache, "builds")
	inst = join(config.DB, "installed")
	pkgs = join(config.Cache, "packages")
	srcs = join(config.Cache, "sources")
	stgs = join(config.Cache, "stages")
}

type Config struct {
	Cache     string
	DB        string
	Identity  string
	Plans     string
	PlansRepo string
	Repo      string
	Root      string
}
