package main

import (
	"os"
	"testing"
)

var reset = os.Args

func TestLint(t *testing.T) {
	os.Args = append(os.Args, "show", "bash")
	main()
}

func TestSearch(t *testing.T) {
	os.Args = reset
	os.Args = append(os.Args, "search")
	main()
}
