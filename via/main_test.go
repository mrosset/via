package main

import (
	"fmt"
	"os"
	"testing"
)

func TestDiff(t *testing.T) {
	os.Args = append(os.Args, "diff", "bash")
	main()
}

func TestGroup(t *testing.T) {
	os.Args = append([]string{}, "build", "devel")
	fmt.Println(os.Args)
	main()
}
