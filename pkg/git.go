package via

import (
        "fmt"
        "gopkg.in/src-d/go-git.v4"
        "io/ioutil"
        "os"
        "strings"
)

// Clone remote URL into directory.
func Clone(dir Path, url string) error {
        _, err := git.PlainClone(dir.String(), false, &git.CloneOptions{
                URL:               url,
                Progress:          os.Stdout,
                RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
        })
        return err
}

// Branch returns the currently checked out branch for a .git directory
func Branch(path Path) (string, error) {
        var (
                head = path.Join(".git/HEAD")
                dir  = path.Base()
                sub  = path.Join("..", ".git", "modules", dir.String(), "HEAD")
        )
        if sub.Exists() {
                head = sub
        }
        b, err := ioutil.ReadFile(head.String())
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
