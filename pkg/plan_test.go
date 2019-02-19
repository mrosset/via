package via

import (
        "encoding/json"
        mjson "github.com/mrosset/util/json"
        "testing"
)

var (
        testPlan = &Plan{
                Name:          "hello",
                Version:       "2.9",
                Url:           "http://mirrors.kernel.org/gnu/hello/hello-{{.Version}}.tar.gz",
                ManualDepends: []string{"libgomp"},
                BuildInStage:  false,
                Build:         []string{"touch hello"},
                Package:       []string{"install -m755 -D hello $PKGDIR/$PREFIX/bin/hello"},
                Files:         []string{"a.out"},
                Group:         "core",
        }
)

func TestPlanDepends(t *testing.T) {
        test{
                Expect: []string{"libgomp"},
                Got:    testPlan.Depends(),
        }.equals(t)
}

func TestPlanExpand(t *testing.T) {
        test{
                Expect: "http://mirrors.kernel.org/gnu/hello/hello-2.9.tar.gz",
                Got:    testPlan.Expand().Url,
        }.equals(t)

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
        )
        test{
                Got:    PackagePath(testConfig, plan),
                Expect: "testdata/repo/QmdmdqJZ5NuyiiEYhjsPfEHU87PYHXSNrFLM34misaZBo4.tar.gz",
        }.equals(t)
}

func TestPlanJSON_Encode(t *testing.T) {
        var (
                jplan = PlanJSON{
                        Flags: []string{"beta", "alpha"},
                }
                file = "testdata/plan.json"
        )
        test{
                Expect: nil,
                Got:    mjson.Write(jplan, file),
        }.equals(t)
}

func TestPlanJSON_MarshalJSON(t *testing.T) {
        var (
                jplan = PlanJSON{
                        SubPackages:   []string{"beta", "alpha"},
                        Flags:         []string{"beta", "alpha"},
                        Remove:        []string{"beta", "alpha"},
                        AutoDepends:   []string{"beta", "alpha"},
                        ManualDepends: []string{"beta", "alpha"},
                        BuildDepends:  []string{"beta", "alpha"},
                }
                plan   Plan
                expect = []string{"alpha", "beta"}
        )
        got, err := jplan.MarshalJSON()
        if err != nil {
                t.Fatal(err)
        }
        if err = json.Unmarshal(got, &plan); err != nil {
                t.Fatal(err)
        }
        tests{
                {
                        Expect: expect,
                        Got:    plan.SubPackages,
                },
                {
                        Expect: Flags(expect),
                        Got:    plan.Flags,
                },
                {
                        Expect: expect,
                        Got:    plan.AutoDepends,
                },
                {
                        Expect: expect,
                        Got:    plan.ManualDepends,
                },
                {
                        Expect: expect,
                        Got:    plan.BuildDepends,
                },
        }.equals(t)
}
