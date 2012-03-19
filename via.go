package via

import (
	"compress/gzip"
	"errors"
	"gurl/pkg"
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

func DownloadSrc(plan *Plan) (err error) {
	sfile := sources.File(path.Base(plan.Url))
	if file.Exists(sfile) {
		return nil
	}
	defer util.Println()
	return gurl.Download(plan.Url, sources.String())
}

func Stage(plan *Plan) (err error) {
	dir := stages.File(plan.NameVersion())
	if file.Exists(dir) {
		return nil
	}
	sfile := sources.File(path.Base(plan.Url))
	r, err := magic.GetReader(sfile)
	if err != nil {
		return err
	}
	_, err = Untar(r, stages.String())
	return
}

func GnuBuild(plan *Plan) (err error) {
	bdir := builds.File(plan.NameVersion())
	sdir := stages.File(plan.NameVersion())
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
	configure := stages.Dir(plan.NameVersion()).File("configure")
	switch {
	case file.Exists(configure):
		return GnuBuild(plan)
	default:
		log.Fatal(errors.New("could not determine build type"))
	}
	return
}

func MakeInstall(plan *Plan) (err error) {
	pdir := packages.File(plan.NameVersion())
	bdir := builds.File(plan.NameVersion())
	return util.Run("make", bdir, "install", "DESTDIR="+pdir)
}

func CreatePackage(plan *Plan) (err error) {
	pfile := repo.File(plan.PackageFile())
	fd, err := os.Create(pfile)
	if err != nil {
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return Package(gz, plan)
}

func Install(plan *Plan) (err error) {
	pfile := repo.File(plan.PackageFile())
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
	db := installed.File(plan.Name)
	err = os.MkdirAll(db, 0755)
	if err != nil {
		return
	}
	return json.Write(man, path.Join(db, "manifest.json"))
}

func Remove(plan *Plan) (err error) {
	man, err := ReadManifest(plan)
	if err != nil {
		return err
	}
	for _, f := range man.Files {
		err = os.Remove(path.Join(config.Root, f))
		if err != nil {
			return err
		}
	}
	return os.RemoveAll(installed.File(plan.Name))
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
