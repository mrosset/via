package main

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
)

var (
	client = new(http.Client)
)

type BuildFnc func(*Plan) error

func init() {
	util.Verbose = *verbose
}

func DownloadSrc(plan *Plan) (err error) {
	if util.FileExists(path.Join(config.Sources(), plan.File)) {
		return nil
	}
	info("download", plan.Url())
	return gurl.Download(plan.Url(), config.Sources())
}

func Stage(plan *Plan) (err error) {
	info("stage", plan.File)
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
	if !util.FileExists(bdir) {
		info("creating", bdir)
		err = os.Mkdir(bdir, 0775)
		if err != nil {
			return err
		}
	}
	err = util.Run(sdir+"/configure", bdir, "--config-cache", "--prefix="+config.Root)
	if err != nil {
		return err
	}

	return util.Run("make", bdir)
}

func Build(plan *Plan) (err error) {
	configure := path.Join(config.Stages(), plan.NameVersion(), "configure")
	switch {
	case util.FileExists(configure):
		info("GnuBuild", plan.NameVersion())
		return GnuBuild(plan)
	default:
		log.Fatal(errors.New("could not determine build type"))
	}
	return
}

func Install(plan *Plan) (err error) {
	info("installing", plan.NameVersion())
	return util.Run("make", config.GetBuildDir(plan.NameVersion()), "install")
}

/*
func Package(name string) (err error) {
	info("packaging", name)
	walkFn := func(path string, info os.FileInfo, err error) error {
		spath := strings.Replace(path, config.GetPackageDir(name)+"/", "", -1)
		fmt.Println(spath)
		return nil
	}
	return filepath.Walk(config.GetPackageDir(name), walkFn)
}
*/

func info(prefix string, msg string) {
	fmt.Printf("%-20s %s\n", prefix, msg)
}
