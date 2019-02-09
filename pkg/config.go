package via

import (
	"fmt"
	"github.com/mrosset/util/json"
	"os"
	"path/filepath"
	"strings"
)

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

func (c Config) Getenv() []string {
	var (
		keep = []string{"HOME", "TERM", "PKGDIR", "SRCDIR", "Flags"}
		env  = []string{}
	)
	for _, e := range keep {
		env = append(env, fmt.Sprintf("%s=%s", e, os.Getenv(e)))
	}
	for i, v := range c.Env {
		env = append(env, fmt.Sprintf("%s=%s", i, os.ExpandEnv(v)))
	}
	return env
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

func (d DB) Installed(config *Config) string {
	return join(config.Root, string(d), "installed")
}

func (d DB) Plans(config *Config) string {
	return join(config.Root, string(d), "plans")
}

func (d DB) InstalledFiles(config *Config) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(d.Installed(config), "*", "*.json"))
	if err != nil {
		return nil, err
	}
	return files, nil
}
