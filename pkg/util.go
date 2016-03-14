package via

// This file contains private functions reused throughout the Api

import (
	"fmt"
	"github.com/str1ngs/util/file"
	"os"
	"path"
	"path/filepath"
)

// function aliases
var (
	join   = path.Join
	base   = path.Base
	exists = file.Exists
)

func baseDir(p string) string {
	return base(path.Dir(p))
}

// Check if a string slice contains a
// string
func contains(sl []string, s string) bool {
	for _, j := range sl {
		if expand(j) == expand(s) {
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

func cpDir(s, d string) error {
	pdir := path.Dir(s)
	fn := func(p string, fi os.FileInfo, err error) error {
		spath := p[len(pdir)+1:]
		dpath := join(d, spath)
		if file.Exists(dpath) {
			return nil
		}
		if fi.IsDir() {
			return os.Mkdir(dpath, fi.Mode())
		}
		fd, err := os.OpenFile(dpath, os.O_CREATE|os.O_WRONLY, fi.Mode())
		if err != nil {
			elog.Println(err)
			return err
		}
		defer fd.Close()
		return file.Copy(fd, p)
	}
	return filepath.Walk(s, fn)
}

func printSlice(s []string) {
	for _, j := range s {
		fmt.Println(j)
	}
}

// for debugging only
func stop() {
	fmt.Println("STOP")
	os.Exit(1)
}
