package main

import (
	"os"
	"testing"
)

var clean = os.Args

func testBuild(t *testing.T) {
	args := []string{"build", "ccache"}
	os.Args = append(clean, args...)
	main()
}

func TestLint(t *testing.T) {
	args := []string{"lint"}
	os.Args = append(clean, args...)
	main()
}
