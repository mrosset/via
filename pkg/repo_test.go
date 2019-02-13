package via

import (
	"os"
	"reflect"
	"testing"
)

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
			t.Errorf(EXPECT_GOT_FMT, expectOne, got)
		}
		got = inverse.Owns(tfile)
		if expectOne != got {
			t.Errorf(EXPECT_GOT_FMT, expectOne, got)
		}
	}

	if got := repo.Owners(tfile); !reflect.DeepEqual(got, expectMore) {
		t.Errorf(EXPECT_GOT_FMT, expectMore, got)
	}

}

func TestRepo_Exists(t *testing.T) {
	tests := []struct {
		name string
		r    Repo
		want bool
	}{
		{
			"",
			Repo("testdata/repo"),
			true,
		},
		{
			"",
			Repo("testdata/false"),
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
	tests := []struct {
		name string
		r    Repo
		want string
	}{
		{
			"",
			Repo("$VIA_TEST/repo"),
			"testdata/repo",
		},
		// TODO: Add test cases.
	}
	os.Setenv("VIA_TEST", "testdata")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Expand(); got != tt.want {
				t.Errorf("Repo.Expand() = %v, want %v", got, tt.want)
			}
		})
	}
}
