package via

import (
	"fmt"
	"os"
	"path/filepath"
)

type Repo struct {
	Arch      string
	Manifests map[string]*Manifest
}

func (this *Repo) AddManifest(man *Manifest) {
	this.Manifests[man.Meta.Name] = man
}

func LoadRepo(arch string) (rep *Repo, err os.Error) {
	file := fmt.Sprintf("%s/%s/repo-%s.json.gz", repo, arch, arch)
	rep = new(Repo)
	err = ReadGzFile(rep, file)
	if err != nil {
		return rep, err
	}
	return rep, err
}

func (this *Repo) Save(arch string) (err os.Error) {
	file := fmt.Sprintf("%s/%s/repo-%s.json.gz", repo, arch, arch)
	err = WriteGzFile(this, file)
	return err
}

func UpdateRepo(arch string) (err os.Error) {
	dir := filepath.Join(repo, arch)

	files, err := filepath.Glob(dir + "/*" + PackExt)
	if err != nil {
		return err
	}
	r := &Repo{arch, make(map[string]*Manifest)}
	for _, f := range files {
		man, err := UnpackManifest(f)
		if err != nil {
			return err
		}
		r.AddManifest(man)
	}
	err = r.Save(arch)
	return err
}

func uploadRepo(arch string) (err os.Error) {
	file := fmt.Sprintf("%s/%s/repo-%s.json.gz", repo, arch, arch)
	return upload(file)
}
