package via

import (
	"compress/gzip"
	"fmt"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/console"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	update  = false
	deps    = false
)

func Root(s string) {
	config.Root = s
}

func Verbose(b bool) {
	verbose = b
}

func Deps(b bool) {
	deps = b
}

func Update(b bool) {
	update = b
}

func Debug(b bool) {
	debug = b
}

func doCommands(dir string, cmds []string) (err error) {
	for i, j := range cmds {
		if debug {
			elog.Println(i, j)
		}
		cmd := exec.Command("bash", "-c", j)
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

// Updates each plans Oid to the Oid of the tarball in publish git repo
// this function should never be used in production. It's used for making sure
// the plans Oid match the git repo's Oid
func SyncHashs(cons *Construct) {
	plans, _ := GetPlans()
	for _, p := range plans {
		if file.Exists(cons.PackageFilePath()) {
			p.Oid, _ = file.Sha256sum(cons.PackageFilePath())
			p.Save(cons.Config)
			log.Println(p.Oid, p.Name)
		}
	}
}

func Install(config *Config, name string) (err error) {
	plan, err := FindPlan(config, name)
	if err != nil {
		elog.Println(name, err)
		return
	}
	cons := NewConstruct(config, plan)
	fmt.Printf(lfmt, "installing", plan.Name)
	if IsInstalled(name) {
		fmt.Printf("FIXME: (short flags) package %s installed upgrading anyways.\n", plan.NameVersion())
		err := Remove(name)
		if err != nil {
			return err
		}
	}
	for _, d := range append(plan.AutoDepends, plan.ManualDepends...) {
		if IsInstalled(d) {
			continue
		}
		err := Install(config, d)
		if err != nil {
			return err
		}
	}
	db := filepath.Join(config.DB.Installed(), plan.Name)
	if file.Exists(db) {
		return fmt.Errorf("%s is already installed", name)
	}
	if !file.Exists(cons.PackageFilePath()) {
		//return errors.New(fmt.Sprintf("%s does not exist", pfile))
		ddir := join(config.Repo, "repo")
		os.MkdirAll(ddir, 0755)
		err := gurl.Download(ddir, config.Binary+"/"+plan.PackageFile())
		if err != nil {
			elog.Println(cons.PackageFilePath())
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
	sha, err := file.Sha256sum(cons.PackageFilePath())
	if err != nil {
		return (err)
	}
	if sha != plan.Oid {
		return fmt.Errorf("%s Plans OID does not match tarballs got %s", plan.NameVersion(), sha)
	}
	man, err := ReadPackManifest(cons.PackageFilePath())
	if err != nil {
		return err
	}
	errs := conflicts(man)
	if len(errs) > 0 {
		//return errs[0]
		for _, e := range errs {
			elog.Println(e)
		}
	}
	fd, err := os.Open(cons.PackageFilePath())
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

func BuildDeps(plan *Plan) (err error) {
	cons := NewConstruct(config, plan)
	deps := append(plan.AutoDepends, plan.ManualDepends...)
	for _, d := range deps {
		if IsInstalled(d) {
			continue
		}
		p, _ := FindPlan(cons.Config, d)
		if file.Exists(cons.PackageFilePath()) {
			err := Install(cons.Config, p.Name)
			if err != nil {
				return err
			}
			continue
		}
		fmt.Println("building", d, "for", plan.NameVersion())
		err := BuildDeps(p)
		if err != nil {
			elog.Println(err)
			return err
		}
	}
	err = cons.BuildSteps()
	if err != nil {
		return err
	}
	return Install(cons.Config, plan.Name)
}

var (
	rexName   = regexp.MustCompile("[A-Za-z]+")
	rexTruple = regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+")
	rexDouble = regexp.MustCompile("[0-9]+.[0-9]+")
)

// Creates a new plan from a given Url
func Create(url, group string) (err error) {
	var (
		xfile   = filepath.Base(url)
		name    = rexName.FindString(xfile)
		triple  = rexTruple.FindString(xfile)
		double  = rexDouble.FindString(xfile)
		version string
	)
	switch {
	case triple != "":
		version = triple
	case double != "":
		version = double
	default:
		return fmt.Errorf("regex fail for %s", xfile)
	}
	plan := &Plan{Name: name, Version: version, Url: url, Group: group}
	plan.Inherit = "gnu"
	if file.Exists(plan.Path(config)) {
		return fmt.Errorf("%s already exists", plan.Path(config))
	}
	return plan.Save(config)
}

func IsInstalled(name string) bool {
	return file.Exists(join(config.DB.Installed(), name))
}

func refactor(plan *Plan) {
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
			console.Println("lint", plan.Name, plan.Version)
		}
		sort.Strings(plan.SubPackages)
		sort.Strings(plan.Flags)
		sort.Strings(plan.Remove)
		sort.Strings(plan.AutoDepends)
		sort.Strings(plan.ManualDepends)
		sort.Strings(plan.BuildDepends)
		refactor(plan)
		err = plan.Save(config)
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
	plan, err := FindPlan(config, name)
	if err != nil {
		return err
	}
	fmt.Printf(lfmt, "clean", plan.NameVersion())
	dir := join(cache.Builds(), plan.NameVersion())
	if err = os.RemoveAll(dir); err != nil {
		return err
	}

	if plan.BuildInStage {
		dir = join(cache.Stages(), plan.stageDir())
		return os.RemoveAll(dir)
	}
	return nil
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

// Setup Dynamic linker
func CheckLink() error {
	real := fmt.Sprintf(RUNTIME_LINKER, filepath.Join(config.Root, config.Prefix))
	ldir := filepath.Dir(config.Linker)

	if !file.Exists(real) {
		elog.Printf("%s real linker does not exist", real)
	}

	os.MkdirAll(ldir, 0755)

	elog.Printf("linking\t %s\t %s", config.Linker, real)
	return os.Symlink(real, config.Linker)
}

func GetConfig() *Config {
	return config
}
