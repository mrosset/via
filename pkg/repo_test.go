package via

import (
        "testing"
)

func TestRepoFilePaths(t *testing.T) {
        tests{
                {
                        Expect: Path("testdata/plans/repo.json"),
                        Got:    testConfig.Repo.File(testConfig),
                },
                {
                        Expect: Path("testdata/plans/files.json"),
                        Got:    testConfig.Repo.FilesFile(testConfig),
                },
        }.equals(t)
}

func TestRepoFilesOwns(t *testing.T) {
        var (
                tfile = "libc.so"
                repo  = RepoFiles{
                        "glibc":     []string{tfile},
                        "glibc-arm": []string{tfile},
                }
                inverse = RepoFiles{
                        "glibc-arm": []string{tfile},
                        "glibc":     []string{tfile},
                }
                expectOne  = "glibc"
                expectMore = []string{"glibc", "glibc-arm"}
        )

        for i := 0; i <= 100; i++ {
                test{
                        Expect: expectOne,
                        Got:    repo.Owns(tfile),
                }.equals(t)
                test{
                        Expect: expectOne,
                        Got:    inverse.Owns(tfile),
                }.equals(t)
        }

        test{
                Expect: expectMore,
                Got:    repo.Owners(tfile),
        }.equals(t)
}

func TestRepoCreate(t *testing.T) {
        tests{
                {
                        Expect: nil,
                        Got:    RepoCreate(testConfig),
                },
                {
                        Name:   "files.json",
                        Expect: true,
                        Got:    Path("testdata/plans/files.json").Exists(),
                },
                {
                        Name:   "repo.json",
                        Expect: true,
                        Got:    Path("testdata/plans/repo.json").Exists(),
                },
        }.equals(t)
}

func TestRepo_Exists(t *testing.T) {
        tests{
                {
                        Name:   "ensure",
                        Expect: nil,
                        Got:    Repo{"testdata/repo"}.Ensure(),
                },
                {
                        Name:   "exists",
                        Expect: true,
                        Got:    Repo{"testdata/repo"}.Exists(),
                },
                {
                        Name:   "fail",
                        Expect: false,
                        Got:    Repo{"testdata/false"}.Exists(),
                },
        }.equals(t)
}
