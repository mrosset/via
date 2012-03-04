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
	"util"
	"util/file"
)

var (
	client  = new(http.Client)
	Verbose = false
)

type BuildFnc func(*Plan) error

func DownloadSrc(plan *Plan) (err error) {
	if file.Exists(path.Join(config.Sources(), plan.File)) {
		return nil
	}
	info("DownloadSrc", plan.Url())
	return gurl.Download(plan.Url(), config.Sources())
}

func Stage(plan *Plan) (err error) {
	info("Stage", plan.File)
	fd, err := os.Open(path.Join(config.Sources(), plan.File))
	util.CheckFatal(err)
	defer fd.Close()
	gz, err := gzip.NewReader(fd)
	util.CheckFatal(err)
	return Untar(gz, config.Stages())
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
	configure := path.Join(config.Stages(), plan.NameVersion(), "configure")
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
	info("Install", plan.NameVersion())
	return util.Run("make", config.GetBuildDir(plan.NameVersion()), "install", "DESTDIR="+config.GetPackageDir(plan.NameVersion()))
}

func Package(plan *Plan) (err error) {
	info("Package", plan.NameVersion())
	dirfile := path.Join(config.GetPackageDir(plan.NameVersion()), config.Prefix, "share", "info", "dir")
	if file.Exists(dirfile) {
		err := os.Remove(dirfile)
		if err != nil {
			return err
		}
	}
	fd, err := os.Create(plan.PackageFile())
	if err != nil {
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return Tar(gz, config.GetPackageDir(plan.NameVersion()))
}

func Install(name string) (err error) {
	plan, err := ReadPlan(name)
	if err != nil {
		return err
	}
	info("Installing", plan.NameVersion())
	err = CheckSig(plan.PackageFile())
	if err != nil {
		return err
	}
	fd, err := os.Open(plan.PackageFile())
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
	info("Removing", plan.NameVersion())
	err = RmTar(plan.PackageFile(), config.Root)
	if err != nil {
		return err
	}
	return
}

func info(prefix string, msg string) {
	fmt.Printf("%-20s %s\n", prefix, msg)
}
