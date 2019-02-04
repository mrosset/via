package via

import (
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testPlan = &Plan{
		Name:          "hello",
		Version:       "2.9",
		Url:           "http://mirrors.kernel.org/gnu/hello/hello-{{.Version}}.tar.gz",
		ManualDepends: []string{"libgomp"},
		BuildInStage:  false,
		Build:         []string{"touch a.out"},
		Package:       []string{"install -m775 -D a.out $PKGDIR/$PREFIX/bin/a.out"},
		Files:         []string{"a.out"},
		Group:         "core",
		config:        testConfig,
	}
)

func TestPlanDepends(t *testing.T) {
	var (
		expect = []string{"libgomp"}
		got    = testPlan.Depends()
	)
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestPlanExpand(t *testing.T) {
	var (
		expect = "http://mirrors.kernel.org/gnu/hello/hello-2.9.tar.gz"
		got    = testPlan.Expand().Url
	)
	if expect != got {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestFindPlan(t *testing.T) {
	var (
		expect = &Plan{
			Name: "sed",
			Url:  "http://mirrors.kernel.org/gnu/sed/sed-{{.Version}}.tar.xz",
		}
	)
	got, err := NewPlan(config, "sed")
	if err != nil {
		t.Fatal(err)
	}
	if expect.Name != got.Name || got.Url != expect.Url {
		t.Errorf("expected %s got %s", expect.Url, got.Url)
	}
}

func TestBuildDir(t *testing.T) {
	var (
		expect = filepath.Join(wd, "testdata/cache/bld/hello-2.9")
		got    = testPlan.BuildDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}

func TestStageDir(t *testing.T) {
	var (
		expect = filepath.Join(wd, "testdata/cache/stg/hello-2.9")
		got    = testPlan.GetStageDir()
	)
	if got != expect {
		t.Errorf("expect '%s' -> got '%s'", expect, got)
	}
}

func TestPlanPackagePath(t *testing.T) {
	var (
		plan = &Plan{
			Name:    "hello",
			Version: "2.9",
			Cid:     "QmdmdqJZ5NuyiiEYhjsPfEHU87PYHXSNrFLM34misaZBo4",
			config:  testConfig,
		}
		got    = plan.PackagePath()
		expect = "testdata/repo/QmdmdqJZ5NuyiiEYhjsPfEHU87PYHXSNrFLM34misaZBo4.tar.gz"
	)
	if got != expect {
		t.Errorf("expect '%s' -> got %s", expect, got)
	}
}
