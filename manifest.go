package via

import (
	"debug/elf"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Manifest struct {
	Plan *Plan
}

func CreateManifest(dir string, plan *Plan) (err error) {
	mfile := filepath.Join(dir, "manifest.json.gz")
	man := Manifest{Plan: plan}
	files := []string{}
	if file.Exists(mfile) {
		err := os.Remove(mfile)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	walkFn := func(path string, info os.FileInfo, err error) error {
		if path == dir {
			return nil
		}
		spath := path[len(dir)+1:]
		stat, err := os.Lstat(path)
		if err != nil {
			log.Println(err, path)
			return err
		}
		if !stat.IsDir() {
			files = append(files, spath)
		}
		return nil
	}
	err = filepath.Walk(dir, walkFn)
	if err != nil {
		log.Println(err)
		return err
	}
	plan.Depends = Depends(plan.Name, dir, files)
	plan.Files = files
	plan.Save()
	return json.WriteGzJson(&man, mfile)
}

func Depends(pname, base string, files []string) []string {
	deps := []string{}
	for _, j := range files {
		d := depends(join(base, j))
		for _, k := range d {
			o := owns(k)
			if contains(deps, o) || pname == o {
				continue
			}
			deps = append(deps, o)
		}
	}
	return deps
}

func contains(sl []string, s string) bool {
	for _, j := range sl {
		if j == s {
			return true
		}
	}
	return false
}

func depends(file string) []string {
	f, err := elf.Open(file)
	if err != nil {
		return nil
	}
	im, err := f.ImportedLibraries()
	if err != nil {
		return nil
	}
	return im
}

func owns(file string) string {
	e, err := filepath.Glob(join(config.Plans, "*.json"))
	if err != nil {
		goto xerr
	}
	for _, j := range e {
		p, err := ReadPath(j)
		if err != nil {
			goto xerr
		}
		for _, f := range p.Files {
			if filepath.Base(f) == file {
				return p.Name
			}
		}
	}

	return ""
xerr:
	log.Println(err)
	return ""
}

func ReadManifest(name string) (man *Manifest, err error) {
	man = new(Manifest)
	err = json.Read(man, path.Join(config.DB.Installed(), name, "manifest.json"))
	if err != nil {
		return
	}
	return
}
