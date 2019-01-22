package via

import (
	"fmt"
	"github.com/mrosset/util/json"
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
	err := json.Read(&files, join(config.Plans, "files.json"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

func RepoCreate() error {
	var (
		repo  = []string{}
		files = map[string][]string{}
		rfile = join(config.Plans, "repo.json")
		ffile = join(config.Plans, "files.json")
	)
	e, err := PlanFiles()
	if err != nil {
		return err
	}
	for _, j := range e {
		p, err := ReadPath(config, j)
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
	return json.Write(files, ffile)
}
