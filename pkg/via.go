package via

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/str1ngs/gurl"
	"github.com/str1ngs/util/console"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	client  = new(http.Client)
	verbose = false
	elog    = log.New(os.Stderr, "", log.Lshortfile)
	lfmt    = "%-20.20s %v\n"
	debug   = false
	expand  = os.ExpandEnv
)

func Root(s string) {
	config.Root = s
}

func Verbose(b bool) {
	verbose = b
}
func Debug(b bool) {
	debug = b
}
func DownloadSrc(plan *Plan) (err error) {
	if file.Exists(plan.SourcePath()) {
		return nil
	}
	fmt.Printf(lfmt, "download", plan.NameVersion())
	return gurl.Download(cache.Sources(), plan.Url)
}

// Stages the downloaded source via's cache directory
// the stage only happens once unless BuilInStage is used
func Stage(plan *Plan) (err error) {
	if plan.Url == "" || file.Exists(plan.GetStageDir()) {
		// nothing to stage
		return nil
	}
	err = GNUUntar(cache.Stages(), plan.SourceFile())
	if err != nil {
		elog.Println(err)
		return err
	}
	if err := doCommands(join(cache.Stages(), plan.stageDir()), plan.Patch); err != nil {
		return err
	}
	return
}

// Calls each shell command in the plans Build field.
func Build(plan *Plan) (err error) {
	var (
		build = plan.Build
	)
	if file.Exists(plan.PackagePath()) {
		fmt.Printf("FIXME: (short flags)  package %s exists building anyways.\n", plan.PackagePath())
	}
	flags := config.Flags
	if plan.Flags != nil {
		flags = append(flags, plan.Flags...)
	}
	if !file.Exists(plan.BuildDir()) {
		os.MkdirAll(plan.BuildDir(), 0755)
	}
	// Parent plan's Build is run first this plans is added at the end.
	if plan.Inherit != "" {
		parent, _ := NewPlan(plan.Inherit)
		build = append(parent.Build, plan.Build...)
		flags = append(flags, parent.Flags...)
	}
	os.Setenv("SRCDIR", plan.GetStageDir())
	os.Setenv("Flags", expand(flags.String()))
	return doCommands(plan.BuildDir(), build)
}

func doCommands(dir string, cmds []string) (err error) {
	for i, j := range cmds {
		j := expand(j)
		if debug {
			elog.Println(i, j)
		}
		cmd := exec.Command("sh", "-c", j)
		cmd.Dir = dir
		cmd.Stdin = os.Stdin
		if verbose {
			cmd.Stdout = os.Stdout
		}
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			elog.Printf("%s: %s\n", j, err)
			return err
		}
	}
	return nil
}

func Package(bdir string, plan *Plan) (err error) {
	var (
		pack = plan.Package
	)
	pdir := join(cache.Packages(), plan.NameVersion())
	if bdir == "" {
		bdir = join(cache.Builds(), plan.NameVersion())
	}
	if plan.BuildInStage {
		bdir = join(cache.Stages(), plan.stageDir())
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
	if plan.Inherit != "" {
		parent, _ := NewPlan(plan.Inherit)
		pack = parent.Package
		pack = append(pack, plan.Package...)
	}
	err = doCommands(bdir, pack)
	if err != nil {
		return err
	}
	for _, j := range plan.SubPackages {
		sub, err := NewPlan(j)
		if err != nil {
			return err
		}
		if err = Package(bdir, sub); err != nil {
			return err
		}
	}
	return CreatePackage(plan)
	/*
		err = CreatePackage(plan)
		if err != nil {
			return err
		}
		return Sign(plan)
	*/
}

func CreatePackage(plan *Plan) (err error) {
	pfile := plan.PackagePath()
	os.MkdirAll(path.Dir(pfile), 0755)
	fd, err := os.Create(pfile)
	if err != nil {
		elog.Println(err)
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return Tarball(gz, plan)
}

func Install(name string) (err error) {
	plan, err := NewPlan(name)
	if err != nil {
		elog.Println(name, err)
		return
	}
	fmt.Printf(lfmt, "installing", plan.Name)
	if IsInstalled(name) {
		fmt.Printf("FIXME: (short flags) package %s installed upgrading anyways.\n", plan.NameVersion())
		err := Remove(name)
		if err != nil {
			return err
		}
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
	db := path.Join(config.DB.Installed(), plan.Name)
	if file.Exists(db) {
		return fmt.Errorf("%s is already installed", name)
	}
	pfile := plan.PackagePath()
	if !file.Exists(pfile) {
		//return errors.New(fmt.Sprintf("%s does not exist", pfile))
		err := gurl.Download(config.Repo+"/master", config.Binary+"/"+plan.PackageFile())
		if err != nil {
			log.Fatal(err)
		}
		//fatal(gurl.Download(config.Repo, config.Binary+"/"+plan.PackageFile()+".sig"))
	}
	/*
		err = CheckSig(pfile)
		if err != nil {
			return
		}
	*/
	man, err := ReadPackManifest(pfile)
	if err != nil {
		return err
	}
	errs := conflicts(man)
	if len(errs) > 0 {
		for _, e := range errs {
			elog.Println(e)
		}
		//return errs[0]
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
	err = json.Write(man, join(db, "manifest.json"))
	if err != nil {
		return err
	}
	return PostInstall(plan)
}

func PostInstall(plan *Plan) (err error) {
	return doCommands("/", append(plan.PostInstall, config.PostInstall...))
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
			elog.Println(f, err)
		}
	}

	return os.RemoveAll(join(config.DB.Installed(), name))
}

// Run all of the functions required to build a package
func BuildSteps(plan *Plan) (err error) {
	if file.Exists(plan.PackageFile()) {
		return fmt.Errorf("package %s exists", plan.PackageFile())
	}
	if err := DownloadSrc(plan); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "stage", plan.NameVersion())
	if err := Stage(plan); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "build", plan.NameVersion())
	if err := Build(plan); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "package", plan.NameVersion())
	if err := Package("", plan); err != nil {
		elog.Println(err)
		return err
	}
	return nil
}

var (
	rexName   = regexp.MustCompile("[a-z]+")
	rexTruple = regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+")
	rexDouble = regexp.MustCompile("[0-9]+.[0-9]+")
)

// Creates a new plan from a given Url
func Create(url string) (err error) {
	var (
		xfile   = path.Base(url)
		name    = rexName.FindString(xfile)
		truple  = rexTruple.FindString(xfile)
		double  = rexDouble.FindString(xfile)
		version string
	)
	switch {
	case truple != "":
		version = truple
	case double != "":
		version = double
	}
	plan := &Plan{Name: name, Version: version, Url: url, Group: "core"}
	if file.Exists(plan.Path()) {
		return errors.New(fmt.Sprintf("%s already exists", plan.Path()))
	}
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
		// If Group is empty, we can set it
		if plan.Group == "" {
			plan.Group = baseDir(j)
		}
		if verbose {
			console.Println("lint", plan.NameVersion(), plan.Package)
		}
		sort.Strings(plan.SubPackages)
		sort.Strings(plan.Flags)
		sort.Strings(plan.Remove)
		sort.Strings(plan.Depends)
		err = plan.Save()
		if err != nil {
			elog.Println(err)
			return err
		}
	}
	console.Flush()
	return nil
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func Clean(name string) error {
	plan, err := NewPlan(name)
	if err != nil {
		return err
	}
	fmt.Printf(lfmt, "clean", plan.NameVersion())
	dir := join(cache.Builds(), plan.NameVersion())
	if err = os.RemoveAll(dir); err != nil {
		return err
	}

	dir = join(cache.Stages(), plan.stageDir())
	return os.RemoveAll(dir)
}

func PlanFiles() ([]string, error) {
	return filepath.Glob(join(config.Plans, "*", "*.json"))
}

func conflicts(man *Plan) (errs []error) {
	for _, f := range man.Files {
		fpath := join(config.Root, f)
		if file.Exists(fpath) {
			errs = append(errs, fmt.Errorf("%s already exists.", f))
		}
	}
	return errs
}

func GetConfig() *Config {
	return config
}
