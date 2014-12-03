package main

import (
	"os"
	"testing"
)

func TestLint(t *testing.T) {
	os.Args = append(os.Args, "diff", "musl", "bash")
	main()
}
