package via

import "C"
import (
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"github.com/mrosset/via/pkg"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"path/filepath"
)

const (
	PACKAGE = "org.golang.todo.via"
	SRCPATH = "/data/data/org.golang.todo.libviaexample/files/src/via"
	GITURL  = "https://github.com/mrosset/plans"
)

var (
	cloneOptions = &git.CloneOptions{
		URL:           GITURL,
		Depth:         1,
		Tags:          git.NoTags,
		ReferenceName: "refs/heads/aarch64-via-linux-gnu",
	}
	srcpath = SRCPATH
)

func Hello() string {
	return "Via - A systems package manager"
}

func getConfig() (*via.Config, error) {
	var (
		config = &via.Config{}
	)
	if !file.Exists(srcpath) {
		log.Printf("cloning %s -> %s", GITURL, srcpath)
		if _, err := git.PlainClone(srcpath, false, cloneOptions); err != nil {
			return nil, err
		}
	}
	log.Printf("reading %s", filepath.Join(srcpath, "config.json"))
	if err := json.Read(config, filepath.Join(srcpath, "config.json")); err != nil {
		return nil, err
	}
	return config, nil
}

func GetBranch() (string, error) {
	config, err := getConfig()
	if err != nil {
		return "", err
	}
	return config.Branch, nil
}
