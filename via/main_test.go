package main

import (
	"os"
	"testing"
)

var (
	args = os.Args
)

func TestDiff(t *testing.T) {
	os.Args = append(args, "diff", "bash")
	main()
}

func TestGroup(t *testing.T) {
	os.Args = append(args, "build", "sed")
	main()
}

func TestList(t *testing.T) {
	os.Args = append(args, "list", "glibc")
	main()
}
