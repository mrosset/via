package via

import (
	"compress/gzip"
	"errors"
	"fmt"
	"gurl/pkg"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"util"
	"util/file"
	"util/file/magic"
)

var (
	client  = new(http.Client)
	Verbose = false
)

type BuildFnc func(*Plan) error

func DownloadSrc(plan *Plan) (err error) {
	sfile := path.Join(config.Cache.Sources(), path.Base(plan.Url))
	if file.Exists(sfile) {
		return nil
	}
	info("DownloadSrc", plan.Url)
	defer fmt.Println()
	return gurl.Download(plan.Url, config.Cache.Sources())
}

func Stage(plan *Plan) (err error) {
	if file.Exists(config.GetStageDir(plan.NameVersion())) {
		info("Stage", "skipping")
		return nil
	}
	info("Stage", path.Base(plan.Url))
	path := path.Join(config.Cache.Sources(), path.Base(plan.Url))
	r, err := magic.GetReader(path)
	if err != nil {
		return err
	}
	return Untar(r, config.Cache.Stages())
}

func GnuBuild(plan *Plan) (err error) {
	bdir := config.GetBuildDir(plan.NameVersion())
	sdir := config.GetStageDir(plan.NameVersion())
	if !file.Exists(bdir) {
		err = os.Mkdir(bdir, 0775)
		if err != nil {
			return err
		}
	}
	err = util.Run(sdir+"/configure", bdir, "--config-cache", "--prefix="+config.Prefix)
	if err != nil {
		return err
	}

	return util.Run("make", bdir)
}

func Build(plan *Plan) (err error) {
	configure := path.Join(config.Cache.Stages(), plan.NameVersion(), "configure")
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
	return util.Run("make", config.GetBuildDir(plan.NameVersion()), "install", "DESTDIR="+config.GetPackageDir(plan.NameVersion()))
}

func CreatePackage(plan *Plan) (err error) {
	info("Package", plan.NameVersion())
	dirfile := path.Join(config.GetPackageDir(plan.NameVersion()), config.Prefix, "share", "info", "dir")
	if file.Exists(dirfile) {
		err := os.Remove(dirfile)
		if err != nil {
			return err
		}
	}
	pfile := path.Join(config.Repo, plan.PackageFile())
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
		return err
	}
	info("Installing", plan.NameVersion())
	pfile := path.Join(config.Repo, plan.PackageFile())
	err = CheckSig(pfile)
	if err != nil {
		return err
	}
	fd, err := os.Open(pfile)
	if err != nil {
		return err
	}
	defer fd.Close()
	gz, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}
	defer gz.Close()
	err = Untar(gz, config.Root)
	if err != nil {
		return err
	}
	return nil
}

func Remove(name string) (err error) {
	plan, err := ReadPlan(name)
	if err != nil {
		return err
	}
	pfile := path.Join(config.Repo, plan.PackageFile())
	info("Removing", plan.NameVersion())
	err = RmTar(pfile, config.Root)
	if err != nil {
		return err
	}
	return
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
	err = plan.Save()
	if err != nil {
		return err
	}
	return BuildSteps(plan)
}

func info(prefix string, msg string) {
	fmt.Printf("%-20s %s\n", prefix, msg)
}
