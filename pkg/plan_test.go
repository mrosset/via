package via

import (
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
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}
}

func TestFindPlan(t *testing.T) {
	var (
		expect = &Plan{
			Name: "sed",
			Url:  "http://mirrors.kernel.org/gnu/sed/sed-{{.Version}}.tar.xz",
		}
	)
	got, err := NewPlan(testConfig, "sed")
	if err != nil {
		t.Fatal(err)
	}
	if expect.Name != got.Name || got.Url != expect.Url {
		t.Errorf("expected %s got %s", expect.Url, got.Url)
	}
}

func TestPlanPackagePath(t *testing.T) {
	var (
		plan = &Plan{
			Name:    "hello",
			Group:   "core",
			Version: "2.9",
			Cid:     "QmdmdqJZ5NuyiiEYhjsPfEHU87PYHXSNrFLM34misaZBo4",
		}
		got    = PackagePath(testConfig, plan)
		expect = "testdata/repo/QmdmdqJZ5NuyiiEYhjsPfEHU87PYHXSNrFLM34misaZBo4.tar.gz"
	)
	if got != expect {
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}
}
