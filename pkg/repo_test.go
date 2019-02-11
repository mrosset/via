package via

import (
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
