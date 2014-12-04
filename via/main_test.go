package main

import (
	"os"
	"testing"
)

func TestDiff(t *testing.T) {
	os.Args = append(os.Args, "diff", "musl", "bash")
	main()
}
