package via

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	wd, _      = os.Getwd()
	testConfig = &Config{
		Root:    "testdata/root",
		Repo:    "testdata/repo",
		Cache:   "testdata/cache",
		DB:      "",
		Binary:  "http://localhost:8080/ipfs/",
		Threads: 8,
		IpfsApi: "localhost:5001",
		Env: map[string]string{
			"PATH":    "/bin:/usr/bin",
			"LDFLAGS": "",
			"PREFIX":  "/opt/via",
		},
	}
)

// init
func init() {
	testConfig.DB = DB(filepath.Join(wd, "testdata/root/db"))
	for i, e := range testConfig.Env {
		os.Setenv(i, os.ExpandEnv(e))
	}
	cache = Cache(filepath.Join(wd, string(testConfig.Cache)))
	cache.Init()
}

func TestConfigExpand(t *testing.T) {
	var (
		c = &Config{
			Branch: "x86_64-via-linux-gnu",
			Prefix: "/usr/local/via",
			Binary: "https://bitbucket.org/strings/publish/raw/{{.Branch}}/repo",
			Env: map[string]string{
				"C_INCLUDE_PATH": "{{.Prefix}}/include",
			},
		}
		expect = "https://bitbucket.org/strings/publish/raw/x86_64-via-linux-gnu/repo"
		got    = ""
	)

	got = c.Expand().Binary
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}

	expect = "/usr/local/via/include"
	got = c.Expand().Env["C_INCLUDE_PATH"]
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestConfig(t *testing.T) {
	if config == nil {
		t.Errorf("config is nil")
	}
}
