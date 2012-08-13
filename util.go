package via

// This file contains private functions reused throughout the Api

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// Varies aliases
var (
	join = path.Join
)

// Check if a string slice contains a
// string
func contains(sl []string, s string) bool {
	for _, j := range sl {
		if j == s {
			return true
		}
	}
	return false
}

// Checks if a dir path is empty
func isEmpty(p string) bool {
	e, err := filepath.Glob(join(p, "*"))
	if err != nil {
		elog.Println(err)
		return false
	}
	return len(e) == 0
}

// Walks a path and prints each dir/file
func walkPath(p string) error {
	fn := func(path string, fi os.FileInfo, err error) error {
		if path == p {
			return nil
		}
		fmt.Println(path)
		return nil
	}
	return filepath.Walk(p, fn)
}
