package via

import (
	"compress/gzip"
	"fmt"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"os"
	"path/filepath"
)

type Installer struct {
	config *Config
	plan   *Plan
}

func NewInstaller(config *Config, plan *Plan) *Installer {
	return &Installer{
		config: config,
		plan:   plan,
	}
}

func Download(config *Config, plan *Plan) error {
	var (
		url   = config.Binary + "/" + plan.Cid
		pfile = plan.PackagePath()
	)
	if isDocker() {
		url = "http://172.17.0.1:8080/ipfs/" + plan.Cid
	}
	if file.Exists(pfile) {
		return nil
	}
	if !file.Exists(config.Repo) {
		if err := os.MkdirAll(config.Repo, 0775); err != nil {
			return err
		}
	}
	return gurl.NameDownload(config.Repo, url, plan.PackageFile())
}

func (i Installer) VerifyCid() error {
	cid, err := HashOnly(i.config, Path(i.plan.PackagePath()))
	if err != nil {
		return err
	}
	if cid != i.plan.Cid {
		return fmt.Errorf("%s Plans CID does not match tarballs got %s", i.plan.NameVersion(), cid)
	}
	return nil
}

func (i Installer) Install() error {
	var (
		name = i.plan.Name
	)
	if IsInstalled(i.config, name) {
		err := Remove(i.config, name)
		if err != nil {
			return err
		}
	}
	db := filepath.Join(i.config.DB.Installed(i.config), name)
	if file.Exists(db) {
		return fmt.Errorf("%s is already installed", name)
	}
	if err := Download(i.config, i.plan); err != nil {
		return err
	}
	man, err := ReadPackManifest(i.plan.PackagePath())
	if err != nil {
		elog.Println(err)
		return err
	}
	errs := conflicts(i.config, man)
	if len(errs) > 0 {
		//return errs[0]
		for _, e := range errs {
			elog.Println(e)
		}
	}
	fd, err := os.Open(i.plan.PackagePath())
	if err != nil {
		elog.Println(err)
		return err
	}
	defer fd.Close()
	gz, err := gzip.NewReader(fd)
	if err != nil {
		elog.Println(err)
		return err
	}
	defer gz.Close()
	os.MkdirAll(i.config.Root, 0755)
	if err = Untar(i.config.Root, gz); err != nil {
		elog.Println(err)
		return err
	}
	if err = os.MkdirAll(db, 0755); err != nil {
		elog.Println(db, err)
		return err
	}
	man.Cid = i.plan.Cid
	err = json.Write(man, join(db, "manifest.json"))
	if err != nil {
		elog.Println(db, err)
		return err
	}
	return nil
}
