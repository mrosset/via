package via

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/str1ngs/gurl"
	"github.com/str1ngs/util"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/file/magic"
	"github.com/str1ngs/util/json"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
)

var (
	client  = new(http.Client)
	Verbose = false
)

//type BuildFnc func(*Plan) error

func Init() (err error) {
	return nil
}

func DownloadSrc(plan *Plan) (err error) {
	sfile := path.Join(cache.Srcs(), path.Base(plan.Url))
	if file.Exists(sfile) {
		return nil
	}
	return gurl.Download(cache.Srcs(), plan.Url)
}

func Stage(plan *Plan) (err error) {
	if file.Exists(path.Join(plan.NameVersion())) {
		return nil
	}
	path := path.Join(cache.Srcs(), path.Base(plan.Url))
	r, err := magic.GetReader(path)
	if err != nil {
		return err
	}
	_, err = Untar(r, cache.Stages())
	return
}

func GnuBuild(plan *Plan) (err error) {
	bdir := path.Join(cache.Builds(), plan.NameVersion())
	sdir := path.Join(cache.Stages(), plan.NameVersion())
	if !file.Exists(bdir) {
		err = os.Mkdir(bdir, 0775)
		if err != nil {
			return err
		}
	}
	err = util.Run(sdir+"/configure", bdir, "--config-cache")
	if err != nil {
		return err
	}

	return util.Run("make", bdir)
}

func Build(plan *Plan) (err error) {
	configure := path.Join(cache.Stages(), plan.NameVersion(), "configure")
	switch {
	case file.Exists(configure):
		return GnuBuild(plan)
	default:
		log.Fatal(errors.New("could not determine build type"))
	}
	return
}

func MakeInstall(plan *Plan) (err error) {
	dir := path.Join(cache.Builds(), plan.NameVersion())
	pdir := path.Join(cache.Pkgs(), plan.NameVersion())
	return util.Run("make", dir, "install", "DESTDIR="+pdir)
}

func CreatePackage(plan *Plan) (err error) {
	pfile := path.Join(string(config.Repo), plan.PackageFile())
	fd, err := os.Create(pfile)
	if err != nil {
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return Package(gz, plan)
}

func Install(name string) (err error) {
	plan, err := ReadPlan(name)
	if err != nil {
		return
	}
	db := path.Join(config.DB.Installed(), plan.Name)
	if file.Exists(db) {
		return fmt.Errorf("%s is already installed", name)
	}
	pfile := path.Join(config.Repo, plan.PackageFile())
	err = CheckSig(pfile)
	if err != nil {
		return
	}
	fd, err := os.Open(pfile)
	if err != nil {
		return
	}
	defer fd.Close()
	gz, err := gzip.NewReader(fd)
	if err != nil {
		return
	}
	defer gz.Close()
	man, err := Untar(gz, config.Root)
	if err != nil {
		return err
	}
	err = os.MkdirAll(db, 0755)
	if err != nil {
		fmt.Println("*WARNING*", err)
	}
	return json.Write(man, path.Join(db, "manifest.json"))
}

func Remove(name string) (err error) {
	man, err := ReadManifest(name)
	if err != nil {
		return err
	}
	for _, f := range man.Files {
		err = os.Remove(path.Join(config.Root, f))
		if err != nil {
			return err
		}
	}
	return os.RemoveAll(path.Join(config.DB.Installed(), name))
}

type Step struct {
	Name string
	Run  func(*Plan) error
}

type Steps []Step

func (this Steps) Run(p *Plan) (err error) {
	for _, s := range this {
		fmt.Printf("%-20.20s %s\n", s.Name, p.NameVersion())
		if err = s.Run(p); err != nil {
			return
		}
	}
	return nil
}

// Run all of the functions required to build a package
func BuildSteps(plan *Plan) (err error) {
	steps := Steps{
		Step{"download", DownloadSrc},
		Step{"stage", Stage},
		Step{"build", Build},
		Step{"make install", MakeInstall},
		Step{"package", CreatePackage},
		Step{"sign", Sign},
	}
	return steps.Run(plan)
}

// Creates a new plan from a give Url
func Create(url string) (err error) {
	var (
		file    = path.Base(url)
		name    = regexp.MustCompile("[a-z]+").FindString(file)
		truple  = regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+").FindString(file)
		double  = regexp.MustCompile("[0-9]+.[0-9]+").FindString(file)
		version string
	)
	switch {
	case truple != "":
		version = truple
	case double != "":
		version = double
	}
	plan := &Plan{Name: name, Version: version, Url: url}
	return plan.Save()
}
