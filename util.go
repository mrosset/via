package via

// This file contains private functions reused throughout the Api

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// function aliases
var (
	join = path.Join
	base = path.Base
)

func baseDir(p string) string {
	return base(path.Dir(p))
}

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

// Walks a path and returns a slice of dir/files
func walkPath(p string) ([]string, error) {
	s := []string{}
	fn := func(path string, fi os.FileInfo, err error) error {
		s = append(s, path)
		return nil
	}
	err := filepath.Walk(p, fn)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func printSlice(s []string) {
	for _, j := range s {
		fmt.Println(j)
	}
}
