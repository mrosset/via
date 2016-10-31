package git

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

var (
	join = path.Join
)

func Branch(head string) (string, error) {
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
