package via

import (
	"fmt"
	"github.com/mrosset/util/file"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Clones remote URL into directory
// name is the reference name to clone.
// e.g reference name ref/heads/master
func Clone(dir, url string) error {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               url,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	return err
}

func Name(path string) {
}

// Returns the currently checked out branch for a Git directory
func Branch(path string) (string, error) {
	// path, err := filepath.Abs(path)
	// if err != nil {
	//	return "", err
	// }
	var (
		head = join(path, ".git/HEAD")
		dir  = filepath.Base(path)
		sub  = join(path, "../.git/modules", dir, "HEAD")
	)
	if file.Exists(sub) {
		head = sub
	}
	b, err := ioutil.ReadFile(head)
	if err != nil {
		return "", err
	}
	in := strings.Split(string(b), "/")
	branch := in[len(in)-1]
	branch = strings.Trim(branch, "\n\r")
	if branch == "" {
		return "", fmt.Errorf("No branch found")
	}
	return branch, nil
}
