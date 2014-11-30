package via

import (
	"fmt"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"os"
	"os/exec"
)

type RepoFiles map[string][]string

func (rf *RepoFiles) Owns(file string) string {
	for pack, files := range *rf {
		for _, f := range files {
			if file == base(f) {
				return pack
			}
		}
	}
	fmt.Println("warning: can not resolve", file)
	return ""
}

func ReadRepoFiles() (RepoFiles, error) {
	files := RepoFiles{}
	err := json.ReadGz(&files, join(config.Plans, "files.json.gz"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

func PlanSync() error {
	dir := config.Plans
	arg := "pull"
	if !file.Exists(dir) {
		arg = "clone"
		dir = dir + "/../"
	}
	git := exec.Command("git", arg, config.PlansRepo)
	git.Dir = "/home/strings/via"
	git.Stdin = os.Stdin
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	elog.Println("syncing", config.PlansRepo)
	err := git.Run()
	if err != nil {
		return err
	}
	return RepoCreate()
}

func RepoCreate() error {
	var (
		repo  = []string{}
		files = map[string][]string{}
		rfile = join(config.Plans, "repo.json")
		ffile = join(config.Plans, "files.json.gz")
	)
	e, err := PlanFiles()
	if err != nil {
		return err
	}
	for _, j := range e {
		p, err := ReadPath(j)
		if err != nil {
			return err
		}
		repo = append(repo, join(p.Group, p.Name+".json"))
		files[p.Name] = p.Files
	}
	err = json.Write(repo, rfile)
	if err != nil {
		return err
	}
	return json.WriteGz(files, ffile)
}
