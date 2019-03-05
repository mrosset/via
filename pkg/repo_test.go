package via

import (
	. "github.com/mrosset/via/pkg/test"
	"testing"
)

func TestRepoFilePaths(t *testing.T) {
	Tests{
		{
			Expect: Path("testdata/plans/repo.json"),
			Got:    testConfig.Repo.File(testConfig),
		},
		{
			Expect: Path("testdata/plans/files.json"),
			Got:    testConfig.Repo.FilesFile(testConfig),
		},
	}.Equals(t)
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
		Test{
			Expect: expectOne,
			Got:    repo.Owns(tfile),
		}.Equals(t)
		Test{
			Expect: expectOne,
			Got:    inverse.Owns(tfile),
		}.Equals(t)
	}

	Test{
		Expect: expectMore,
		Got:    repo.Owners(tfile),
	}.Equals(t)
}

func TestRepoCreate(t *testing.T) {
	Tests{
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
	}.Equals(t)
}

func TestRepo_Exists(t *testing.T) {
	Tests{
		{
			Name:   "ensure",
			Expect: nil,
			Got:    Repo{"testdata/repo"}.MkdirAll(),
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
	}.Equals(t)
}
