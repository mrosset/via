package via

import (
	"compress/gzip"
	"errors"
	"fmt"
	"gurl"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"util"
	"util/file"
	"util/file/magic"
	"util/json"
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
	sfile := path.Join(srcs, path.Base(plan.Url))
	if file.Exists(sfile) {
		return nil
	}
	info("DownloadSrc", plan.Url)
	return gurl.Download(srcs, plan.Url)
}

func Stage(plan *Plan) (err error) {
	if file.Exists(path.Join(plan.NameVersion())) {
		info("Stage", "skipping")
		return nil
	}
	info("Stage", path.Base(plan.Url))
	path := path.Join(srcs, path.Base(plan.Url))
	r, err := magic.GetReader(path)
	if err != nil {
		return err
	}
	_, err = Untar(r, stgs)
	return
}

func GnuBuild(plan *Plan) (err error) {
	bdir := path.Join(blds, plan.NameVersion())
	sdir := path.Join(stgs, plan.NameVersion())
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
	configure := path.Join(stgs, plan.NameVersion(), "configure")
	switch {
	case file.Exists(configure):
		info("GnuBuild", plan.NameVersion())
		return GnuBuild(plan)
	default:
		log.Fatal(errors.New("could not determine build type"))
	}
	return
}

func MakeInstall(plan *Plan) (err error) {
	info("MakeInstall", plan.NameVersion())
	dir := path.Join(blds, plan.NameVersion())
	pdir := path.Join(pkgs, plan.NameVersion())
	return util.Run("make", dir, "install", "DESTDIR="+pdir)
}

func CreatePackage(plan *Plan) (err error) {
	info("Package", plan.NameVersion())
	pfile := path.Join(config.Repo, plan.PackageFile())
	fmt.Println(pfile)
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
	info("Installing", plan.NameVersion())
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
	db := path.Join(inst, plan.Name)
	err = os.MkdirAll(db, 0755)
	if err != nil {
		fmt.Println("*WARNING*", err)
	}
	return json.Write(man, path.Join(db, "manifest.json"))
}

func List(name string) (err error) {
	man, err := ReadManifest(name)
	if err != nil {
		return
	}
	for _, f := range man.Files {
		fmt.Println("file:", path.Join(config.Root, f))
	}
	return
}

func Remove(name string) (err error) {
	man, err := ReadManifest(name)
	if err != nil {
		return err
	}
	info("Removing", man.Plan.NameVersion())
	for _, f := range man.Files {
		err = os.Remove(path.Join(config.Root, f))
		if err != nil {
			return err
		}
	}
	return os.RemoveAll(path.Join(inst, name))
}

// libtorrent-0.13.0.tar.gz

func BuildSteps(plan *Plan) (err error) {
	if err := DownloadSrc(plan); err != nil {
		return err
	}
	if err := Stage(plan); err != nil {
		return err
	}
	if err := Build(plan); err != nil {
		return err
	}
	if err := MakeInstall(plan); err != nil {
		return err
	}
	if err := CreatePackage(plan); err != nil {
		return err
	}
	return Sign(plan)
}

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

func info(prefix string, msg string) {
	fmt.Printf("%-20s %s\n", prefix, msg)
}
