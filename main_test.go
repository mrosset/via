package main

import (
	"os"
	"testing"
)

func TestHelp(t *testing.T) {
	os.Args = append([]string{}, "help")
	main()
}
