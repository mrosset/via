package via

import (
	"encoding/json"
	"fmt"
	mjson "github.com/mrosset/util/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Config represents via configuration type
type Config struct {
	Branch    string
	Identity  string
	Arch      string
	OS        string
	Root      Path
	PlansRepo string
	Threads   int
	IpfsAPI   string
	// Paths
	Cache  Cache
	DB     DB
	Plans  Plans
	Repo   Repo
	Binary string
	Prefix Path

	// Toolchain
	Flags Flags

	Env         Env
	Remove      []string
	PostInstall []string

	// Internal Fields
	template *Config
}

// ConfigJSON provides json Marshal and Unmarshal interface for Config
type ConfigJSON Config

// ToConfig returns Config type
func (j ConfigJSON) ToConfig() *Config {
	config := Config(j)
	return &config
}

// UnmarshalJSON provides Unmarshal interface
func (j *ConfigJSON) UnmarshalJSON(data []byte) error {
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	sort.Strings(c.Flags)
	*j = ConfigJSON(c)
	return nil
}

// MarshalJSON provides Marshal interface
func (j *ConfigJSON) MarshalJSON() ([]byte, error) {
	sort.Strings(j.Flags)
	return json.Marshal(*j)
}

// ReadConfig reads config path and returns a new initialized Config
func ReadConfig(path Path) (*Config, error) {
	var jconfig ConfigJSON
	if err := mjson.Read(&jconfig, path.String()); err != nil {
		return nil, err
	}

	// TODO: create a marshal command to sort these fields?
	// sort.Strings([]string(config.Flags))
	// sort.Strings(config.Remove)

	if err := mjson.Write(&jconfig, path.String()); err != nil {
		return nil, err
	}

	config := jconfig.ToConfig().Expand()
	config.Cache = Cache{
		Path(config.Cache.Expand()),
	}
	config.Cache.Init()
	config.Plans = Plans{
		Path(config.Plans.Expand()),
	}
	config.Repo = Repo{
		Path(config.Repo.Expand()),
	}

	for i, j := range config.Env {
		os.Setenv(i, os.ExpandEnv(j))
	}
	return config, nil
}

// SanitizeEnv returns an os.Environ() environment string slice that
// keeps only white listed environment variables. This ensures when we
// exec command calls nothing leaks from system environment
func (c Config) SanitizeEnv() []string {
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

// Expand returns a Config that have had its fields parsed through
// go's template engine. Basically this allows for self referencing
// json. For example we use this to reduce repetition for things like
// the Prefix field. We can then reuse {{.Prefix}} to represent that
// field in other parts of the config file
func (c *Config) Expand() *Config {
	if c.template != nil {
		return c.template
	}
	o := new(Config)
	err := mjson.Parse(o, c)
	if err != nil {
		panic(err)
	}
	c.template = o
	return c.template
}

// Flags provides a string slice type for working with flags
type Flags []string

// Join joins flags into a string separated with a space
func (f Flags) Join() string {
	return strings.Join(f, " ")
}

// DB provides string type for working with DB installed path
type DB struct {
	Path
}

// Installed returns the path string of the installed directory
func (d DB) Installed(config *Config) Path {
	return config.Root.Join(d.String(), "installed")
}

// InstalledFiles returns all of the json manifests for each install
// package
func (d DB) InstalledFiles(config *Config) ([]string, error) {
	files, err := filepath.Glob(
		d.Installed(config).Join("*", "*.json").String(),
	)
	if err != nil {
		return nil, err
	}
	return files, nil
}
