package via

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

func TestConfigGetenv(t *testing.T) {
	var (
		expect = []string{
			fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
			fmt.Sprintf("TERM=%s", os.Getenv("TERM")),
			"PKGDIR=",
			"SRCDIR=",
			"Flags=",
			"PATH=/bin:/usr/bin",
			"LDFLAGS=",
			"PREFIX=/opt/via",
		}
		got = testConfig.Getenv()
	)
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expect %s got %s", expect, got)
	}
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
