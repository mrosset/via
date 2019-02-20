package via

import (
	"compress/gzip"
	"fmt"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"os"
)

// Installer provides Installer type
type Installer struct {
	config *Config
	plan   *Plan
}

// NewInstaller returns a new Installer that has been initialized
func NewInstaller(config *Config, plan *Plan) *Installer {
	return &Installer{
		config: config,
		plan:   plan,
	}
}

// Download gets the plans binary tarball package from ipfs http
// gateway. If run in a docker instance it will use a local docker ip.
//
// FIXME: now that we have have contain namespaces we probably don't
// need docker logic here. And this will probably produce corner cases
// down the road.
func Download(config *Config, plan *Plan) error {
	var (
		url   = config.Binary + "/" + plan.Cid
		pfile = PackagePath(config, plan)
	)
	if isDocker() {
		url = "http://172.17.0.1:8080/ipfs/" + plan.Cid
	}
	if file.Exists(pfile) {
		return nil
	}
	config.Repo.Ensure()
	return gurl.NameDownload(config.Repo.String(), url, PackageFile(config, plan))
}

// VerifyCid verifies that the download tarball matches the plans Cid
func (i Installer) VerifyCid() error {
	cid, err := HashOnly(i.config, PackagePath(i.config, i.plan))
	if err != nil {
		return err
	}
	if cid != i.plan.Cid {
		return fmt.Errorf("%s Plans CID does not match tarballs got %s", i.plan.NameVersion(), cid)
	}
	return nil
}

// Install method does the final installation of decompressing and
// extracting the tarball. The manifest which is essentially the
// plan's json file is then stored in the DB installed directory.  It
// also updates the manifest's Cid resulting in byte to byte parity
// with the manifest.json and plan.json files.
func (i Installer) Install() error {
	var (
		name  = i.plan.Name
		pfile = PackagePath(i.config, i.plan)
	)
	if err := i.VerifyCid(); err != nil {
		return err
	}
	if IsInstalled(i.config, name) {
		if err := Remove(i.config, name); err != nil {
			return err
		}
	}
	db := i.config.DB.Installed(i.config).Join(name)
	if db.Exists() {
		return fmt.Errorf("%s is already installed", name)
	}
	if err := Download(i.config, i.plan); err != nil {
		return err
	}
	man, err := ReadPackManifest(pfile)
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
	fd, err := os.Open(pfile)
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
	i.config.Root.Ensure()
	if err = Untar(i.config.Root, gz); err != nil {
		elog.Println(err)
		return err
	}
	db.Ensure()
	man.Cid = i.plan.Cid
	err = json.Write(man, db.Join("manifest.json").String())
	if err != nil {
		elog.Println(db, err)
		return err
	}
	return nil
}
