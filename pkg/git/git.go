package git

import (
	"errors"
	"fmt"
	"github.com/str1ngs/util/file"
	"io/ioutil"
	"path"
	"strings"
)

var (
	join = path.Join
)

const (
	git_head = "%s/.git/HEAD"
)

func Branch(path string) (string, error) {
	head := fmt.Sprintf(git_head, path)
	if !file.Exists(head) {
		return "", errors.New(".git not found in " + path)
	}
	b, err := ioutil.ReadFile(head)
	if err != nil {
		return "", err
	}
	in := strings.Split(string(b), "/")
	branch := in[len(in)-1]
	branch = strings.Trim(branch, "\n\r")
	return branch, nil
}
