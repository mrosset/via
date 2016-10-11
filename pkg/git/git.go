package git

import (
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
	git_head = "%s/HEAD"
)

func Branch(path string) (string, error) {
	head := fmt.Sprintf(git_head, path)
	if !file.Exists(head) {
		return "", fmt.Errorf(".git not found in %s", path)
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
