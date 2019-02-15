package via

import (
	"github.com/mrosset/util/file"
	"os"
	"reflect"
	"testing"
)

func TestRepoFilePaths(t *testing.T) {
	tests{
		{
			Expect: "testdata/plans/repo.json",
			Got:    testConfig.Repo.File(testConfig),
		},
		{
			Expect: "testdata/plans/files.json",
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
		got := repo.Owns(tfile)
		if expectOne != got {
			t.Errorf(EXPECT_GOT_FMT, "", expectOne, got)
		}
		got = inverse.Owns(tfile)
		if expectOne != got {
			t.Errorf(EXPECT_GOT_FMT, "", expectOne, got)
		}
	}

	if got := repo.Owners(tfile); !reflect.DeepEqual(got, expectMore) {
		t.Errorf(EXPECT_GOT_FMT, "", expectMore, got)
	}

}

func TestRepoCreate(t *testing.T) {
	tests{
		{
			Expect: nil,
			Got:    RepoCreate(testConfig),
		},
		{
			Label:  "files.json",
			Expect: true,
			Got:    file.Exists("testdata/plans/files.json"),
		},
		{
			Label:  "repo.json",
			Expect: true,
			Got:    file.Exists("testdata/plans/repo.json"),
		},
	}.equals(t)
}

func TestRepo_Exists(t *testing.T) {
	tests := []struct {
		name string
		r    Repo
		want bool
	}{
		{
			"",
			Repo{"testdata/repo"},
			true,
		},
		{
			"",
			Repo{"testdata/false"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Exists(); got != tt.want {
				t.Errorf("Repo.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Expand(t *testing.T) {
	os.Setenv("VIA_TEST_DATA", "testdata")
	tests{
		{
			Label:  "repo expand",
			Expect: "testdata/repo",
			Got:    Repo{"$VIA_TEST_DATA/repo"}.Expand(),
		},
	}.equals(t)
}
