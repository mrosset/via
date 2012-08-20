package via

import (
	"compress/gzip"
	"fmt"
	"github.com/str1ngs/gurl"
	"github.com/str1ngs/util/console"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/file/magic"
	"github.com/str1ngs/util/human"
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
	sfile := path.Join(cache.Srcs(), path.Base(plan.Expand("Url")))
	if file.Exists(sfile) {
		return nil
	}
	return gurl.Download(cache.Srcs(), plan.Expand("Url"))
}

// Stages the downloaded source via's cache directory
// the stage only happens once unless BuilInStage is used
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

// Calls each shell command in the plans Build field.
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
	os.Setenv("Flags", expand(flags.String()))
	bdir := join(cache.Builds(), plan.NameVersion())
	if plan.BuildInStage {
		bdir = join(cache.Stages(), plan.stageDir())
	}
	if !file.Exists(bdir) {
		os.Mkdir(bdir, 0755)
	}
	return doCommands(bdir, plan.Build)
}

func doCommands(dir string, cmds []string) (err error) {
	for _, j := range cmds {
		j := expand(j)
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
	pdir := join(cache.Pkgs(), plan.NameVersion())
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
	err = doCommands(bdir, plan.Package)
	if err != nil {
		return err
	}
	for _, j := range plan.SubPackages {
		sub, err := FindPlan(j)
		if err != nil {
			return err
		}
		if err = Package(bdir, sub); err != nil {
			return err
		}
	}
	err = CreatePackage(plan)
	if err != nil {
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
	plan, err := FindPlan(name)
	if err != nil {
		return
	}
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
	return json.Write(man, join(db, "manifest.json"))
}

func PostInstall(plan *Plan) (err error) {
	return doCommands("/", plan.PostInstall)
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
		fmt.Printf(lfmt, "download", plan.NameVersion())
		return err
	}
	if err := Stage(plan); err != nil {
		fmt.Printf(lfmt, "stage", plan.NameVersion())
		return err
	}
	if err := Build(plan); err != nil {
		fmt.Printf(lfmt, "build", plan.NameVersion())
		return err
	}
	if err := Package("", plan); err != nil {
		fmt.Printf(lfmt, "package", plan.NameVersion())
		return err
	}
	return nil
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
		// If Group is empty, we can set it
		if plan.Group == "" {
			plan.Group = baseDir(j)
		}
		if verbose {
			console.Println("lint", plan.NameVersion())
		}
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

func Clean(name string) error {
	plan, err := FindPlan(name)
	if err != nil {
		return err
	}
	bdir := join(cache.Builds(), plan.NameVersion())
	if plan.BuildInStage {
		bdir = join(cache.Stages(), plan.stageDir())
	}
	if !file.Exists(bdir) {
		err = fmt.Errorf("%s: does not exist", bdir)
		elog.Println(err)
	}
	return os.RemoveAll(bdir)
}

func Search() {
	e, _ := PlanFiles()
	for _, j := range e {
		plan, err := ReadPath(j)
		if err != nil {
			elog.Println(err)
		}
		fmt.Printf(lfmt, plan.Name, human.ByteSize(plan.Size))
	}
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
