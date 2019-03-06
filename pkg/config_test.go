package via

import (
	"encoding/json"
	"fmt"
	. "github.com/mrosset/test"
	"os"
	"reflect"
	"sort"
	"testing"
)

var (
	wd, _      = os.Getwd()
	testConfig = &Config{
		Root:  Path(wd).Join("testdata/root"),
		DB:    Path("/var/lib/via/db").ToDB(),
		Repo:  Repo{"testdata/repo"},
		Cache: Path(wd).Join("testdata/cache").ToCache(),
		Plans: Plans{"testdata/plans"},
		OS:    "linux",
		Arch:  "x86_64",

		Binary:  "http://localhost:8080/ipfs/",
		Threads: 8,
		IpfsAPI: "localhost:5001",
		Env: map[string]string{
			"PATH":    "/bin:/usr/bin",
			"LDFLAGS": "",
			"PREFIX":  "/opt/via",
		},
		Prefix: "/opt/via",
	}
)

// init
func init() {
	testConfig.Cache = Cache{
		Path(wd).Join("testdata/cache"),
	}

	for i, e := range testConfig.Env {
		os.Setenv(i, os.ExpandEnv(e))
	}
}

type TestConfig struct {
	Flags []string
}

func TestConfig_Unmarshal(t *testing.T) {
	var (
		data   = []byte(`{"Flags":["beta","alpha"]}`)
		config = new(ConfigJSON)
	)
	Tests{
		{
			Expect: nil,
			Got:    json.Unmarshal(data, config),
		},
		{
			Expect: Flags{"alpha", "beta"},
			Got:    config.Flags,
		},
	}.Equals(t)
}

func TestConfig_Marshal(t *testing.T) {
	config := ConfigJSON{
		Flags: []string{"beta", "alpha"},
	}
	b, err := config.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	Test{
		Expect: []byte(`{"Branch":"","Identity":"","Arch":"","OS":"","Root":"","PlansRepo":"","Threads":0,"IpfsAPI":"","Cache":"","DB":"","Plans":"","Repo":"","Binary":"","Prefix":"","Flags":["alpha","beta"],"Env":null,"Remove":null,"PostInstall":null}`),
		Got:    b,
	}.Equals(t)
}

// FIXME: this needs to run offline
func testConfigGetenv(t *testing.T) {
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

	Test{
		Expect: "https://bitbucket.org/strings/publish/raw/x86_64-via-linux-gnu/repo",
		Got:    c.Expand().Binary,
	}.Equals(t)

	Test{
		Expect: "/usr/local/via/include",
		Got:    c.Expand().Env["C_INCLUDE_PATH"],
	}.Equals(t)

	Test{
		Expect: "--build=x86_64-via-linux-gnu",
		Got:    c.Expand().Flags[0],
	}.Equals(t)

	Test{
		Expect: "-O2 -pipe",
		Got:    c.Expand().Env["CXXFLAGS"],
	}.Equals(t)

}

func TestDB_Installed(t *testing.T) {
	Tests{
		{
			Expect: Path(wd).Join("testdata/root/var/lib/via/db/installed"),
			Got:    testConfig.DB.Installed(testConfig),
		},
	}.Equals(t)

}
