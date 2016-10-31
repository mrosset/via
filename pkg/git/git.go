package git

import (
	"fmt"
	"github.com/mrosset/util/file"
	"gopkg.in/src-d/go-git.v4"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	join = filepath.Join
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func Clone(dir, url string) error {
	r, err := git.NewFilesystemRepository(dir)
	if err != nil {
		log.Println(err)
		return err
	}

	err = r.Clone(&git.CloneOptions{
		URL: url,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return (err)
	}
	// ... retrieving the commit object
	commit, err := r.Commit(ref.Hash())
	if err != nil {
		return (err)
	}

	// ... we get all the files from the commit
	files, err := commit.Files()
	if err != nil {
		return (err)
	}

	writeFile := func(f *git.File) error {
		abs := filepath.Join(dir, f.Name)
		dir := filepath.Dir(abs)

		os.MkdirAll(dir, 0755)
		file, err := os.Create(abs)
		if err != nil {
			return err
		}

		defer file.Close()
		r, err := f.Reader()
		if err != nil {
			return err
		}

		defer r.Close()

		if err := file.Chmod(f.Mode); err != nil {
			return err
		}

		_, err = io.Copy(file, r)
		return err
	}
	return files.ForEach(writeFile)
}

func Name(path string) {
}

// Returns the currently checked out branch for a Git directory
func Branch(path string) (string, error) {
	// path, err := filepath.Abs(path)
	// if err != nil {
	// 	return "", err
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
