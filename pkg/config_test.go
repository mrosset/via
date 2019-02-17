package via

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

var (
	wd, _      = os.Getwd()
	testConfig = &Config{
		Root:    "testdata/root",
		Repo:    Repo{"testdata/repo"},
		Cache:   Cache{"testdata/cache"},
		Plans:   Plans{"testdata/plans"},
		OS:      "linux",
		Arch:    "x86_64",
		DB:      "",
		Binary:  "http://localhost:8080/ipfs/",
		Threads: 8,
		IpfsAPI: "localhost:5001",
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
		got = testConfig.SanitizeEnv()
	)
	sort.Strings(got)
	sort.Strings(expect)
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expect %s got %s", expect, got)
	}
}

func TestConfigExpand(t *testing.T) {
	var (
		c = &Config{
			Branch: "x86_64-via-linux-gnu",
			Arch:   "x86_64",
			OS:     "linux",
			Prefix: "/usr/local/via",
			Binary: "https://bitbucket.org/strings/publish/raw/{{.Branch}}/repo",
			Flags: []string{
				"--build={{.Arch}}-via-{{.OS}}-gnu",
			},
			Env: map[string]string{
				"CFLAGS":         "-O2 -pipe",
				"CXXFLAGS":       "{{.Env.CFLAGS}}",
				"C_INCLUDE_PATH": "{{.Prefix}}/include",
			},
		}
	)

	test{
		Expect: "https://bitbucket.org/strings/publish/raw/x86_64-via-linux-gnu/repo",
		Got:    c.Expand().Binary,
	}.equals(t)

	test{
		Expect: "/usr/local/via/include",
		Got:    c.Expand().Env["C_INCLUDE_PATH"],
	}.equals(t)

	test{
		Expect: "--build=x86_64-via-linux-gnu",
		Got:    c.Expand().Flags[0],
	}.equals(t)

	test{
		Expect: "-O2 -pipe",
		Got:    c.Expand().Env["CXXFLAGS"],
	}.equals(t)

}
