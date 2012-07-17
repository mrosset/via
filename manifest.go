package via

import (
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
	plan.Files = files
	plan.Save()
	return json.WriteGzJson(&man, mfile)
}

func ReadManifest(name string) (man *Manifest, err error) {
	man = new(Manifest)
	err = json.Read(man, path.Join(config.DB.Installed(), name, "manifest.json"))
	if err != nil {
		return
	}
	return
}
