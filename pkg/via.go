package via

import (
	"compress/gzip"
	"fmt"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/console"
	"github.com/mrosset/util/file"
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
	cache Cache

	client = new(http.Client)
	debug  = false
	deps   = false
	elog   = log.New(os.Stderr, "", log.Lshortfile)
	expand = os.ExpandEnv
	lfmt   = "%-20.20s %v\n"

	update  = false
	verbose = false
)

// Verbose sets the global verbosity level
//
// FIXME: this should be set via PlanContext
func Verbose(b bool) {
	verbose = b
}

// Update set if a plan should update after building
//
// FIXME: document what this actually does
func Update(b bool) {
	update = b
}

// Debug sets the global debugging level
func Debug(b bool) {
	debug = b
}

// DownloadSrc downloads the PlanContext Plans upstream source
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
		wget(ctx.Cache.Sources(), eurl)
	case "http", "https":
		return gurl.Download(ctx.Cache.Sources(), eurl)
	case "git":
		spath := filepath.Join(ctx.Cache.Sources(), ctx.Plan.Name)
		if err := Clone(spath, "https"+eurl[3:]); err != nil {
			elog.Println(err)
			return err
		}
	default:
		return fmt.Errorf("%s: URL scheme is not supported", u.Scheme)
	}
	return nil
}

// Stage the downloaded source in via's cache directory the stage only
// happens once unless BuilInStage is used
func Stage(ctx *PlanContext) (err error) {
	if ctx.Plan.Url == "" || file.Exists(ctx.StageDir()) {
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
		fmt.Println(ctx.SourcePath())
		fmt.Println(ctx.StageDir())
		if err := Clone(ctx.StageDir(), ctx.SourcePath()); err != nil {
			return err
		}
		goto patch
	}
	switch filepath.Ext(ctx.Plan.SourceFile()) {
	case ".zip":
		unzip(ctx.Cache.Stages(), ctx.SourcePath())
	default:
		GNUUntar(ctx.Cache.Stages(), ctx.SourcePath())
	}
patch:
	fmt.Printf(lfmt, "patch", ctx.Plan.NameVersion())
	return doCommands(&ctx.Config, join(ctx.Cache.Stages(), ctx.Plan.stageDir()), ctx.Plan.Patch)
}

// Build calls each shell command in the plans Build field.
func Build(ctx *PlanContext) (err error) {
	var (
		plan  = ctx.Plan
		build = plan.Build
	)
	for _, p := range plan.BuildDepends {
		if IsInstalled(&ctx.Config, p) {
			continue
		}
		dp, err := NewPlan(&ctx.Config, p)
		if err != nil {
			return err
		}
		if err := NewInstaller(&ctx.Config, dp).Install(); err != nil {
			return err
		}
	}
	// FIXME: flags should not be merged should have a ConfigFLags
	// and PlanFlags environment variable
	flags := append(ctx.Config.Flags, plan.Flags...)
	os.MkdirAll(ctx.BuildDir(), 0755)
	// Parent plan Build is run first this plans is added at the end.
	if plan.Inherit != "" {
		parent, _ := NewPlan(&ctx.Config, plan.Inherit)
		build = append(parent.Build, plan.Build...)
		flags = append(flags, parent.Flags...)
	}
	// FIXME: this should be set within exec.Cmd
	os.Setenv("SRCDIR", ctx.StageDir())
	os.Setenv("Flags", expand(flags.String()))
	err = doCommands(&ctx.Config, ctx.BuildDir(), build)
	if err != nil {
		return fmt.Errorf("%s in %s", err.Error(), ctx.BuildDir())
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
			Env:    config.SanitizeEnv(),
		}
		if verbose {
			cmd.Stdout = os.Stdout
		}
		if err := cmd.Run(); err != nil {
			elog.Printf("%s: %s\n", j, err)
			return err
		}
	}
	return nil
}

// Package calls each shell command in Plans package field
func Package(ctx *PlanContext, bdir string) (err error) {
	var (
		plan = ctx.Plan
		pack = plan.Package
	)
	// Remove plans Cid it's assumed we'll be creating a new one
	plan.Cid = ""
	defer os.Remove(plan.PackagePath())
	pdir := join(ctx.Cache.Packages(), plan.NameVersion())
	if bdir == "" {
		bdir = join(ctx.Cache.Builds(), plan.NameVersion())
	}
	if plan.BuildInStage {
		bdir = join(ctx.Cache.Stages(), plan.stageDir())
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
		parent, _ := NewPlan(&ctx.Config, plan.Inherit)
		pack = append(parent.Package, plan.Package...)
	}
	err = doCommands(&ctx.Config, bdir, pack)
	if err != nil {
		return err
	}
	for _, j := range plan.SubPackages {
		sub, err := NewPlanContextByName(&ctx.Config, j)
		if err != nil {
			return err
		}
		if err = Package(sub, bdir); err != nil {
			return err
		}
	}
	err = CreatePackage(ctx)
	if err != nil {
		elog.Println(err)
		return (err)
	}
	plan.Cid, err = IpfsAdd(&ctx.Config, Path(plan.PackagePath()))
	if err != nil {
		return err
	}
	plan.IsRebuilt = true
	return ctx.WritePlan()
	/*
		err = CreatePackage(plan)
		if err != nil {
			return err
		}
		return Sign(plan)
	*/
}

// CreatePackage Walks a Plans package directory and creates a tarball
func CreatePackage(ctx *PlanContext) (err error) {
	var (
		plan  = ctx.Plan
		pfile = plan.PackagePath()
	)
	os.MkdirAll(filepath.Dir(pfile), 0755)
	fd, err := os.Create(pfile)
	if err != nil {
		elog.Println(err)
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return Tarball(ctx, gz)
}

// PostInstall calls each of the Plans PostInstall commands
func PostInstall(config *Config, plan *Plan) (err error) {
	return doCommands(config, "/", append(plan.PostInstall, config.PostInstall...))
}

// Remove a plan from the system
func Remove(config *Config, name string) (err error) {
	if !IsInstalled(config, name) {
		err = fmt.Errorf("%s is not installed", name)
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

// BuildSteps runs all of the functions required to build a package
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
	if err := Build(ctx); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "package", ctx.Plan.NameVersion())
	if err := Package(ctx, ""); err != nil {
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

// Create a new plan from a given Url
func Create(config *Config, url, group string) (err error) {
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
	ctx := NewPlanContext(config, plan)
	if file.Exists(ctx.PlanPath()) {
		return fmt.Errorf("%s already exists", ctx.PlanPath())
	}
	return ctx.WritePlan()
}

// IsInstalled returns true if a plan is installed
func IsInstalled(config *Config, name string) bool {
	return file.Exists(join(config.DB.Installed(config), name))
}

// Lint walks all plans and formats it sorting fields
//
// FIXME: this should be renamed to Format and a new Lint function
// created. Lint function should have no side effects just look for
// known style isses. For example we can check that each upstream URL
// is using https and not http
func Lint(config *Config) (err error) {
	e, err := PlanFiles(config)
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
		ctx := NewPlanContext(config, plan)
		if err := ctx.WritePlan(); err != nil {
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

// Clean the PlanContext build directory
func Clean(ctx *PlanContext) error {
	var (
		plan  = ctx.Plan
		cache = ctx.Cache
	)
	fmt.Printf(lfmt, "clean", plan.NameVersion())
	dir := join(cache.Builds(), plan.NameVersion())
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	if plan.BuildInStage {
		dir = join(cache.Stages(), plan.stageDir())
		return os.RemoveAll(dir)
	}
	return nil
}

func conflicts(config *Config, man *Plan) (errs []error) {
	for _, f := range man.Files {
		fpath := join(config.Root, f)
		if file.Exists(fpath) {
			errs = append(errs, fmt.Errorf("%s already exists", f))
		}
	}
	return errs
}
