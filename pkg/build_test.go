// +build online

package via

import (
	. "github.com/mrosset/via/pkg/test"
	"testing"
)

func TestNewBuildContext(t *testing.T) {
	var (
		bc = NewBuildContext(testConfig, testPlan)
	)
	Tests{
		{
			Expect: Path(wd).Join("testdata/cache/builds/hello-2.9"),
			Got:    bc.BuildDir,
		},
		{
			Expect: Path(wd).Join("testdata/cache/stages/hello-2.9"),
			Got:    bc.StageDir,
		},
		{
			Expect: Path(wd).Join("testdata/cache/packages/hello-2.9"),
			Got:    bc.PackageDir,
		},
		{
			Expect: Path(wd).Join("testdata/cache/sources/hello-2.9.tar.gz"),
			Got:    bc.SourcePath,
		},
	}.Equals(t)
}

func TestBuilder_Download(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	builder.Cache.Init()
	Tests{
		{
			Expect: nil,
			Got:    builder.Download(),
		},
		{
			Expect: true,
			Got:    builder.Context.SourcePath.Exists(),
		},
	}.Equals(t)

}

func TestBuilder_Stage(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	Tests{
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
	}.Equals(t)
}

func TestBuilder_Build(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	Tests{
		{
			Expect: nil,
			Got:    builder.Build(),
		},
	}.Equals(t)
}

func TestBuilder_Package(t *testing.T) {
	Verbose(true)
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	Tests{
		{
			Expect: nil,
			Got:    builder.Package(builder.Context.BuildDir),
		},
		{
			Name:   "Install exists",
			Expect: true,
			Got:    Path("testdata/cache/packages/hello-2.9/opt/via/bin/hello").Exists(),
		},
		{
			Expect: nil,
			Got:    builder.Context.PackageDir.RemoveAll(),
		},
		{
			Expect: false,
			Got:    builder.Context.PackageDir.Exists(),
		},
	}.Equals(t)
}

func TestBuilder_BuildSteps(t *testing.T) {
	var (
		builder = NewBuilder(testConfig, testPlan)
	)
	Test{
		Expect: nil,
		Got:    builder.BuildSteps(),
	}.Equals(t)
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
	Tests{
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
			Expect: "--cflag1 --cflag2",
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
	}.Equals(t)
}
