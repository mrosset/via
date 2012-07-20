package via

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/str1ngs/gurl"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/file/magic"
	"github.com/str1ngs/util/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	client  = new(http.Client)
	Verbose = false
	elog    = log.New(os.Stderr, "", log.Lshortfile)
	lfmt    = "%-20.20s %v\n"
)

func DownloadSrc(plan *Plan) (err error) {
	sfile := path.Join(cache.Srcs(), path.Base(plan.Url))
	if file.Exists(sfile) {
		return nil
	}
	return gurl.Download(cache.Srcs(), plan.Url)
}

func Stage(plan *Plan) (err error) {
	sdir := join(cache.Stages(), plan.stageDir())
	if file.Exists(sdir) {
		return nil
	}
	path := join(cache.Srcs(), path.Base(plan.Url))
	r, err := magic.GetReader(path)
	if err != nil {
		return err
	}
	err = Untar(cache.Stages(), r)
	if err != nil {
		return err
	}
	return
}

func Build(plan *Plan) (err error) {
	pfile := join(config.Repo, plan.PackageFile())
	if file.Exists(pfile) {
		fmt.Printf("FIXME: (short flags)  package %s exists building anyways.\n", plan.PackageFile())
	}
	flags := config.Flags
	if plan.Flags != nil {
		flags = append(flags, plan.Flags...)
	}
	os.Setenv("SRCDIR", join(cache.Stages(), plan.stageDir()))
	os.Setenv("Flags", flags.String())
	bdir := join(cache.Builds(), plan.NameVersion())
	if plan.BuildInStage {
		bdir = join(cache.Stages(), plan.NameVersion())
	}
	if !file.Exists(bdir) {
		os.Mkdir(bdir, 0755)
	}
	return doCommands(bdir, plan.Build)
}

func doCommands(dir string, cmds []string) (err error) {
	for _, j := range cmds {
		s := os.ExpandEnv(j)
		buf := new(bytes.Buffer)
		buf.WriteString(s)
		cmd := exec.Command("sh")
		cmd.Dir = dir
		cmd.Stdin = buf
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func Package(plan *Plan) (err error) {
	bdir := join(cache.Builds(), plan.NameVersion())
	pdir := join(cache.Pkgs(), plan.NameVersion())

	if plan.BuildInStage {
		bdir = join(cache.Stages(), plan.NameVersion())
	}
	if file.Exists(pdir) {
		err := os.RemoveAll(pdir)
		if err != nil {
			return err
		}
	}
	err = os.Mkdir(pdir, 0755)
	if err != nil {
		elog.Println(err)
		return err
	}
	os.Setenv("PKGDIR", pdir)
	err = doCommands(bdir, plan.Package)
	if err != nil {
		elog.Println(err)
		return err
	}
	for _, f := range config.CleanFiles {
		fmt.Println("cleaning", f)
		f = join(pdir, f)
		if file.Exists(f) {
			err := os.RemoveAll(f)
			if err != nil {
				elog.Println(err)
				return err
			}
		}
	}
	err = CreatePackage(plan)
	if err != nil {
		elog.Println(err)
		return err
	}
	return Sign(plan)
}

func CreatePackage(plan *Plan) (err error) {
	pfile := join(config.Repo, plan.PackageFile())
	fd, err := os.Create(pfile)
	if err != nil {
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return Tarball(gz, plan)
}

func Install(name string) (err error) {
	plan, err := ReadPlan(name)
	if err != nil {
		return
	}
	for _, d := range plan.Depends {
		if IsInstalled(d) {
			continue
		}
		err := Install(d)
		if err != nil {
			return err
		}
	}
	fmt.Printf(lfmt, "installing", plan.NameVersion())
	db := path.Join(config.DB.Installed(), plan.Name)
	if file.Exists(db) {
		return fmt.Errorf("%s is already installed", name)
	}
	pfile := path.Join(config.Repo, plan.PackageFile())
	err = CheckSig(pfile)
	if err != nil {
		return
	}
	man, err := ReadPackManifest(pfile)
	if err != nil {
		return err
	}
	errs := conflicts(man)
	if len(errs) > 0 {
		for _, e := range errs {
			elog.Println(e)
		}
		return errs[0]
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
	err = Untar(config.Root, gz)
	if err != nil {
		return err
	}
	err = os.MkdirAll(db, 0755)
	if err != nil {
		elog.Println(err)
		return err
	}
	return json.Write(man, path.Join(db, "manifest.json"))
}

func Remove(name string) (err error) {
	if !IsInstalled(name) {
		err = fmt.Errorf("%s is not installed.", name)
		elog.Println(err)
		return err
	}
	man, err := ReadManifest(name)
	if err != nil {
		return err
	}
	for _, f := range man.Files {
		fpath := join(config.Root, f)
		err = os.Remove(fpath)
		if err != nil {
			elog.Println(len(f), f, "not exist")
		}
		if file.Exists(fpath) {
			elog.Println(f, "not removed")
		}
	}
	return os.RemoveAll(path.Join(config.DB.Installed(), name))
}

type Step struct {
	Name string
	Run  func(*Plan) error
}

type Steps []Step

func (st Steps) Run(p *Plan) (err error) {
	for _, s := range st {
		fmt.Printf(lfmt, s.Name, p.NameVersion())
		if err = s.Run(p); err != nil {
			return
		}
	}
	return nil
}

// Run all of the functions required to build a package
func BuildSteps(plan *Plan) (err error) {
	if file.Exists(plan.PackageFile()) {
		return fmt.Errorf("package %s exists", plan.PackageFile())
	}
	steps := Steps{
		Step{"download", DownloadSrc},
		Step{"stage", Stage},
		Step{"build", Build},
		Step{"package", Package},
	}
	return steps.Run(plan)
}

// Creates a new plan from a given Url
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

func IsInstalled(name string) bool {
	return file.Exists(join(config.DB.Installed(), name))
}

func Lint() (err error) {
	e, err := PlanFiles()
	if err != nil {
		return err
	}
	for _, j := range e {
		plan, err := ReadPath(j)
		if err != nil {
			err = fmt.Errorf("%s %s", j, err)
			elog.Println(err)
			return err
		}
		fmt.Printf(lfmt, "lint", plan.NameVersion())
		sort.Strings(plan.Flags)
		err = plan.Save()
		if err != nil {
			elog.Println(err)
			return err
		}
	}
	return nil
}

func Clean(name string) error {
	plan, err := ReadPlan(name)
	if err != nil {
		return err
	}
	bdir := join(cache.Builds(), plan.NameVersion())
	if !file.Exists(bdir) {
		err = fmt.Errorf("%s: does not exist", bdir)
		elog.Println(err)
		return err
	}
	return os.RemoveAll(bdir)
}

func List() {
	e, _ := PlanFiles()
	for _, j := range e {
		file := path.Base(j)
		name := strings.Split(file, ".")[0]
		fmt.Println(name)
	}
}

func PlanFiles() ([]string, error) {
	return filepath.Glob(join(config.Plans, "*.json"))
}

func conflicts(man *Plan) (errs []error) {
	for _, f := range man.Files {
		if file.Exists(join(config.Root, f)) {
			errs = append(errs, fmt.Errorf("%s already exists.", f))
		}
	}
	return errs
}
