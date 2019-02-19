// +build online

package via

import (
        "github.com/mrosset/util/file"
        "path/filepath"
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
                        Expect: Path(filepath.Join(wd, "testdata/cache/sources/hello-2.9.tar.gz")),
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
                        Got:    Path(filepath.Join(wd, "testdata/cache/stages/hello-2.9")).Exists(),
                },
                {
                        Expect: true,
                        Got:    file.Exists("testdata/cache/stages/hello-2.9/configure"),
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
