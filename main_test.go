package main

import (
	"os"
	"testing"
)

var (
	args = os.Args
)

func testBuild(t *testing.T) {
	os.Args = append(args, "build", "make")
	main()
}

func TestDiff(t *testing.T) {
	os.Args = append(args, "diff", "bash")
	main()
}

func TestList(t *testing.T) {
	os.Args = append(args, "list", "glibc")
	main()
}
