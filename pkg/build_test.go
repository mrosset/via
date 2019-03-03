// +build online

package via

import (
	"testing"
)

func TestBuilder_PathMethods(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	tests{
		{
			Expect: Path(wd).Join("testdata/cache/stages/hello-2.9"),
			Got:    builder.StageDir(),
		},
		{
			Expect: Path(wd).Join("testdata/cache/sources/hello-2.9.tar.gz"),
			Got:    builder.SourcePath(),
		},
	}.equals(t)
}

func TestBuilder_Download(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	builder.Cache.Init()
	tests{
		{
			Expect: nil,
			Got:    builder.Download(),
		},
		{
			Expect: true,
			Got:    builder.SourcePath().Exists(),
		},
	}.equals(t)

}

func TestBuilder_Stage(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	tests{
		{
			Expect: nil,
			Got:    builder.Stage(),
		},
		{
			Expect: true,
			Got:    Path(wd).Join("testdata/cache/stages/hello-2.9").Exists(),
		},
		{
			Expect: true,
			Got:    Path("testdata/cache/stages/hello-2.9/configure").Exists(),
		},
	}.equals(t)
}

func TestBuilder_Build(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	tests{
		{
			Expect: nil,
			Got:    builder.Build(),
		},
	}.equals(t)
}

func TestBuilder_Package(t *testing.T) {
	Verbose(true)
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	tests{
		{
			Expect: nil,
			Got:    builder.Package(),
		},
		{
			Name:   "Install exists",
			Expect: true,
			Got:    Path("testdata/cache/packages/hello-2.9/opt/via/bin/hello").Exists(),
		},
		{
			Expect: nil,
			Got:    builder.PackageDir().RemoveAll(),
		},
		{
			Expect: false,
			Got:    builder.PackageDir().Exists(),
		},
	}.equals(t)
}

func TestBuilder_BuildSteps(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	test{
		Expect: nil,
		Got:    builder.BuildSteps(),
	}.equals(t)
}

func TestBuilder_Expand(t *testing.T) {
	var (
		config = &Config{
			Prefix: "/opt/via",
			Cache:  Cache{"testdata/cache"},
			Flags:  []string{"--cflag1", "--cflag2"},
		}
		plan = &Plan{
			Name:    "test",
			Version: "1.0.0",
			Flags:   []string{"--pflag1", "--pflag2"},
		}
		builder = NewBuilder(config, plan)
	)
	tests{
		{
			Expect: "/opt/via",
			Got:    ExpandCommand("$PREFIX", builder),
		},
		{
			Expect: "testdata/cache/stages/test-1.0.0",
			Got:    ExpandCommand("$SRCDIR", builder),
		},
		{
			Expect: "testdata/cache/packages/test-1.0.0",
			Got:    ExpandCommand("$PKGDIR", builder),
		},
		{
			Expect: "--cflag1 --cflag2 --pflag1 --pflag2",
			Got:    ExpandCommand("$Flags", builder),
		},
		{
			Expect: "--pflag1 --pflag2",
			Got:    ExpandCommand("$PlanFlags", builder),
		},
		{
			Expect: "testdata/cache/packages/test-1.0.0//opt/via",
			Got:    ExpandCommand("$PKGDIR/$PREFIX", builder),
		},
	}.equals(t)
}
