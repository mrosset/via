package main

import (
	"os"
	"testing"
)

var aclean = os.Args

func testBuild(t *testing.T) {
	args := []string{"build", "ccache"}
	os.Args = append(aclean, args...)
	main()
}

func TestLint(t *testing.T) {
	args := []string{"lint"}
	os.Args = append(aclean, args...)
	main()
}
