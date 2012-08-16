package via

import (
	"github.com/str1ngs/gurl"
	"github.com/str1ngs/util/json"
	"os"
	"path"
)

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
	r := []string{}
	e, err := PlanFiles()
	if err != nil {
		return err
	}
	for _, j := range e {
		p, err := ReadPath(j)
		if err != nil {
			return err
		}
		r = append(r, join(p.Group, p.Name+".json"))
	}
	return json.Write(r, join(config.Plans, "repo.json"))
}
