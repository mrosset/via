package via

import (
	"os"
	"testing"
)

func TestCacheString(t *testing.T) {
	var (
		cache  = Cache("$VIA_CACHE_TEST/testdata/cache")
		expect = "./testdata/cache"
	)
	os.Setenv("VIA_CACHE_TEST", ".")
	got := string(cache.String())
	if expect != got {
		t.Errorf(EXPECT_GOT_FMT, "", expect, got)
	}
}
