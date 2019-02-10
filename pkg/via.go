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
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	cache   Cache
	cfile   = filepath.Join(viapath, "plans/config.json")
	client  = new(http.Client)
	config  = new(Config)
	debug   = false
	deps    = false
	elog    = log.New(os.Stderr, "", log.Lshortfile)
	expand  = os.ExpandEnv
	lfmt    = "%-20.20s %v\n"
	planUrl = "https://github.com/mrosset/plans"
	update  = false
	verbose = false
	viaUrl  = "https://github.com/mrosset/via"
	viapath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/mrosset/via")
)

func init() {
	if os.Getenv("GOPATH") == "" {
		elog.Fatal("GOPATH must be set")
	}
	// TODO rework this to error and suggest user use 'via init'
	if !file.Exists(viapath) {
		elog.Println("cloning plans")
		if err := Clone(viapath, viaUrl); err != nil {
			elog.Fatal(err)
		}
	}
	pdir := filepath.Dir(cfile)
	if !file.Exists(pdir) {
		elog.Println("cloning plans")
		err := Clone(pdir, planUrl)
		if err != nil {
			elog.Fatal(err)
		}
	}
}

func init() {
	err := json.Read(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}
	// TODO: provide Lint for master config
	sort.Strings([]string(config.Flags))
	sort.Strings(config.Remove)
	err = json.Write(&config, cfile)
	if err != nil {
		elog.Fatal(err)
	}

	config = config.Expand()

	// if err := CheckLink(); err != nil {
	//	elog.Fatal(err)
	// }

	cache = Cache(os.ExpandEnv(string(config.Cache)))
	cache.Init()
	config.Plans = os.ExpandEnv(config.Plans)
	config.Repo = os.ExpandEnv(config.Repo)
	err = os.MkdirAll(config.Repo, 0755)
	if err != nil {
		elog.Fatal(err)
	}
	for i, j := range config.Env {
		os.Setenv(i, os.ExpandEnv(j))
	}
	for i, j := range config.Env {
		os.Setenv(i, os.ExpandEnv(j))
	}
}

func Root(path string) {
	config.Root = path
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

func DownloadSrc(ctx *PlanContext) (err error) {
	if file.Exists(ctx.SourcePath()) && !update {
		return nil
	}
	fmt.Printf(lfmt, "download", ctx.Plan.NameVersion())
	eurl := ctx.Plan.Expand().Url
	u, err := url.Parse(eurl)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "ftp":
		wget(cache.Sources(), eurl)
	case "http", "https":
		return gurl.Download(cache.Sources(), eurl)
	case "git":
		spath := filepath.Join(cache.Sources(), ctx.Plan.Name)
		if err := Clone(spath, "https"+eurl[3:]); err != nil {
			elog.Println(err)
			return err
		}
	default:
		return fmt.Errorf("%s: URL scheme is not supported", u.Scheme)
	}
	return nil
}

// Stages the downloaded source in via's cache directory
// the stage only happens once unless BuilInStage is used
func Stage(ctx *PlanContext) (err error) {
	if ctx.Plan.Url == "" || file.Exists(ctx.Plan.GetStageDir()) {
		// nothing to stage
		return nil
	}
	fmt.Printf(lfmt, "stage", ctx.Plan.NameVersion())
	u, err := url.Parse(ctx.Plan.Expand().Url)
	if err != nil {
		elog.Println(err)
		return err
	}
	//FIXME: move this down to switch statement so avoid goto
	if u.Scheme == "git" {
		fmt.Println(ctx.Plan.SourcePath())
		fmt.Println(ctx.Plan.GetStageDir())
		if err := Clone(ctx.Plan.GetStageDir(), ctx.Plan.SourcePath()); err != nil {
			return err
		}
		goto patch
	}
	switch filepath.Ext(ctx.Plan.SourceFile()) {
	case ".zip":
		unzip(cache.Stages(), ctx.Plan.SourcePath())
	default:
		GNUUntar(cache.Stages(), ctx.Plan.SourcePath())
	}
patch:
	fmt.Printf(lfmt, "patch", ctx.Plan.NameVersion())
	return doCommands(config, join(cache.Stages(), ctx.Plan.stageDir()), ctx.Plan.Patch)
}

// Calls each shell command in the plans Build field.
func Build(config *Config, plan *Plan) (err error) {
	var (
		build = plan.Build
	)
	for _, p := range plan.BuildDepends {
		if IsInstalled(config, p) {
			continue
		}
		dp, err := NewPlan(config, p)
		if err != nil {
			return err
		}
		if err := NewInstaller(config, dp).Install(); err != nil {
			return err
		}
	}
	// FIXME: flags should not be merged should have a ConfigFLags
	// and PlanFlags environment variable
	flags := append(config.Flags, plan.Flags...)
	os.MkdirAll(plan.BuildDir(), 0755)
	// Parent plan Build is run first this plans is added at the end.
	if plan.Inherit != "" {
		parent, _ := NewPlan(config, plan.Inherit)
		build = append(parent.Build, plan.Build...)
		flags = append(flags, parent.Flags...)
	}
	// FIXME: this should be set within exec.Cmd
	os.Setenv("SRCDIR", plan.GetStageDir())
	os.Setenv("Flags", expand(flags.String()))
	err = doCommands(config, plan.BuildDir(), build)
	if err != nil {
		return fmt.Errorf("%s in %s", err.Error(), plan.BuildDir())
	}
	return nil
}

func doCommands(config *Config, dir string, cmds []string) (err error) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		return err
	}
	for _, j := range cmds {
		cmd := &exec.Cmd{
			Path:   bash,
			Args:   []string{"bash", "-c", j},
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
			Dir:    dir,
			Env:    config.Getenv(),
		}
		if verbose {
			cmd.Stdout = os.Stdout
		}
		if debug {
			fmt.Println(config.Getenv())
			fmt.Println(os.ExpandEnv(j))
		}
		err = cmd.Run()
		if err != nil {
			elog.Printf("%s: %s\n", j, err)
			return err
		}
	}
	return nil
}

func Package(config *Config, bdir string, plan *Plan) (err error) {
	var (
		pack = plan.Package
	)
	// Remove plans Cid it's assumed we'll be creating a new one
	plan.Cid = ""
	defer os.Remove(plan.PackagePath())
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
		parent, _ := NewPlan(config, plan.Inherit)
		pack = append(parent.Package, plan.Package...)
	}
	err = doCommands(config, bdir, pack)
	if err != nil {
		return err
	}
	for _, j := range plan.SubPackages {
		sub, err := NewPlan(config, j)
		if err != nil {
			return err
		}
		if err = Package(config, bdir, sub); err != nil {
			return err
		}
	}
	err = CreatePackage(config, plan)
	if err != nil {
		elog.Println(err)
		return (err)
	}
	plan.Cid, err = IpfsAdd(config, Path(plan.PackagePath()))
	if err != nil {
		return err
	}
	plan.IsRebuilt = true
	return plan.Save()
	/*
		err = CreatePackage(plan)
		if err != nil {
			return err
		}
		return Sign(plan)
	*/
}

func CreatePackage(config *Config, plan *Plan) (err error) {
	pfile := plan.PackagePath()
	os.MkdirAll(filepath.Dir(pfile), 0755)
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

// Updates each plans Oid to the Oid of the tarball in publish git repo
// this function should never be used in production. It's used for making sure
// the plans Oid match the git repo's Oid
func SyncHashs(config *Config) {
	plans, _ := GetPlans()
	for _, p := range plans {
		if file.Exists(p.PackagePath()) {
			p.Cid, _ = HashOnly(config, Path(p.PackagePath()))
			p.Save()
			log.Println(p.Cid, p.Name)
		}
	}
}

func PostInstall(config *Config, plan *Plan) (err error) {
	return doCommands(config, "/", append(plan.PostInstall, config.PostInstall...))
}

func Remove(config *Config, name string) (err error) {
	if !IsInstalled(config, name) {
		err = fmt.Errorf("%s is not installed.", name)
		return err
	}
	man, err := ReadManifest(config, name)
	if err != nil {
		elog.Println(err)
		return err
	}
	for _, f := range man.Files {
		fpath := join(config.Root, f)
		if err := os.Remove(fpath); err != nil {
			elog.Println(err)
		}
	}
	return os.RemoveAll(join(config.DB.Installed(config), name))
}

// func BuildDeps(config *Config, plan *Plan) (err error) {
//	for _, d := range plan.Depends() {
//		if IsInstalled(config, d) {
//			continue
//		}
//		p, _ := NewPlan(config, d)
//		if file.Exists(p.PackagePath()) {
//			if err := NewInstaller(config, plan).Install(); err != nil {
//				return err
//			}
//			continue
//		}
//		fmt.Println("building", d, "for", plan.NameVersion())
//		err := BuildDeps(config, p)
//		if err != nil {
//			elog.Println(err)
//			return err
//		}
//	}
//	err = BuildSteps(config, plan)
//	if err != nil {
//		return err
//	}
//	return NewInstaller(config, plan).Install()
// }

// Run all of the functions required to build a package
func BuildSteps(ctx *PlanContext) (err error) {
	if file.Exists(ctx.PackageFile()) {
		return fmt.Errorf("package %s exists", ctx.PackageFile())
	}
	if err := DownloadSrc(ctx); err != nil {
		elog.Println(err)
		return err
	}
	if err := Stage(ctx); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "build", ctx.Plan.NameVersion())
	if err := Build(&ctx.Config, ctx.Plan); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "package", ctx.Plan.NameVersion())
	if err := Package(&ctx.Config, "", ctx.Plan); err != nil {
		elog.Println(err)
		return err
	}
	return RepoCreate(&ctx.Config)
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
	if file.Exists(plan.Path()) {
		return fmt.Errorf("%s already exists", plan.Path())
	}
	return plan.Save()
}

func IsInstalled(config *Config, name string) bool {
	return file.Exists(join(config.DB.Installed(config), name))
}

func refactor(plan *Plan) {
	if len(plan.SubPackages) > 0 {
		for _, j := range plan.SubPackages {
			s, _ := NewPlan(config, j)
			if s.Version == plan.Version {
				continue
			}
			s.Version = plan.Version
			s.Save()
		}
	}
}

func Lint() (err error) {
	e, err := PlanFiles()
	if err != nil {
		return err
	}
	for _, j := range e {
		plan, err := ReadPath(config, j)
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
			console.Println("lint", plan.Name, plan.Version, plan.IsRebuilt)
		}
		sort.Strings(plan.SubPackages)
		sort.Strings(plan.Flags)
		sort.Strings(plan.Remove)
		sort.Strings(plan.AutoDepends)
		sort.Strings(plan.ManualDepends)
		sort.Strings(plan.BuildDepends)
		refactor(plan)
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
	plan, err := NewPlan(config, name)
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

func conflicts(config *Config, man *Plan) (errs []error) {
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
