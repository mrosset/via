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
		expect_one  = "glibc"
		expect_more = []string{"glibc", "glibc-arm"}
	)

	for i := 0; i <= 100; i++ {
		got := repo.Owns(tfile)
		if expect_one != got {
			t.Errorf(EXPECT_GOT_FMT, expect_one, got)
		}
		got = inverse.Owns(tfile)
		if expect_one != got {
			t.Errorf(EXPECT_GOT_FMT, expect_one, got)
		}
	}

	if got := repo.Owners(tfile); !reflect.DeepEqual(got, expect_more) {
		t.Errorf(EXPECT_GOT_FMT, expect_more, got)
	}

}
