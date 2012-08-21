package via

import (
	"github.com/str1ngs/gurl"
	"github.com/str1ngs/util/json"
	"os"
	"path"
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

// TODO: replace this with git?
func PlanSync() error {
	pdir := config.DB.Plans()
	if !exists(pdir) {
		err := os.MkdirAll(pdir, 0755)
		if err != nil {
			return err
		}
	}
	local := join(pdir, "repo.json")
	remote := config.PlansRepo + "/repo.json"
	err := gurl.Download(pdir, remote)
	if err != nil {
		return err
	}
	repo := []string{}
	if err = json.Read(&repo, local); err != nil {
		return err
	}
	for _, j := range repo {
		rurl := config.PlansRepo + "/" + j
		dir := join(pdir, path.Dir(j))
		if !exists(dir) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		}
		if err = gurl.Download(dir, rurl); err != nil {
			return err
		}
	}
	return nil
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
