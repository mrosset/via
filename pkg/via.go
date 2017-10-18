package via

import (
	"compress/gzip"
	"fmt"
	"github.com/fatih/color"
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
	client         = new(http.Client)
	verbose        = false
	elog           = log.New(os.Stderr, "", log.Lshortfile)
	lfmt           = "%-20.20s %v\n"
	debug          = false
	expand         = os.ExpandEnv
	update         = false
	deps           = false
	INSTALL_PREFIX = Path("$HOME/via")
	PREFIX         = Path("/tmp/via")
	blue           = color.New(color.FgBlue).SprintFunc()
)

func init() {
	if !Symlinked() {
		INSTALL_PREFIX.MkDirAll(0700)
		err := INSTALL_PREFIX.Symlink(PREFIX)
		if err != nil {
			elog.Fatal(err)
		}
	}
	if !Symlinked() {
		elog.Fatalf("could not setup symlink %s to %s", INSTALL_PREFIX, PREFIX)
	}
}

func Symlinked() bool {
	return PREFIX.Exists() && INSTALL_PREFIX.Exists()
}

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

func DownloadSrc(plan *Plan) (err error) {
	if plan.SourceCid != "" {
		return nil
	}
	if file.Exists(plan.SourcePath()) && !update {
		return nil
	}
	fmt.Printf(lfmt, "download", plan.NameVersion())
	eurl := plan.Expand().Url
	u, err := url.Parse(eurl)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "ftp":
		wget(cache.Sources(), eurl)
	case "http", "https":
		return gurl.Download(cache.Sources(), eurl)
	default:
		return fmt.Errorf("%s URL scheme is not supported")
	}
	return nil
}

// Stages the downloaded source in via's cache directory. The stage only happens
// once unless BuilInStage is used
func Stage(plan *Plan) (err error) {
	if plan.Url == "" || file.Exists(plan.GetStageDir()) {
		// nothing to stage
		return nil
	}
	fmt.Printf(lfmt, "stage", plan.NameVersion())
	if plan.SourceCid != "" {
		return IpfsGet(Path(cache.Stages()), plan.SourceCid)
	}
	defer os.RemoveAll(plan.GetStageDir())
	switch filepath.Ext(plan.SourceFile()) {
	case ".zip":
		unzip(cache.Stages(), plan.SourcePath())
	default:
		s := Path(plan.SourcePath()).ToUnix()
		err := GNUUntar(cache.Stages(), s)
		if err != nil {
			return err
		}
	}
	if plan.SourceCid == "" {
		cid, err := IpfsAdd(Path(cache.Stages()).JoinS(plan.stageDir()), false)
		if err != nil {
			return err
		}
		plan.SourceCid = cid
		plan.Save()
		return Stage(plan)
	}
	return
}

// Calls each shell command in the plans Build field.
func Build(plan *Plan) (err error) {
	var (
		build = plan.Build
	)
	// if err = config.CheckBranches(); err != nil {
	//	return (err)
	// }
	if file.Exists(plan.PackagePath()) {
		fmt.Printf("FIXME: (short flags)  package %s exists building anyways.\n", plan.PackagePath())
	}
	for _, p := range plan.BuildDepends {
		if IsInstalled(p) {
			continue
		}
		if err := Install(p); err != nil {
			return err
		}
	}
	flags := append(config.Flags, plan.Flags...)
	os.MkdirAll(plan.BuildDir(), 0755)
	// Parent plan Build is run first this plans is added at the end.
	if plan.Inherit != "" {
		parent, _ := NewPlan(plan.Inherit)
		build = append(parent.Build, plan.Build...)
		flags = append(flags, parent.Flags...)
	}
	os.Setenv("SRCDIR", Path(plan.GetStageDir()).ToUnix())
	os.Setenv("Flags", expand(flags.String()))
	err = doCommands(plan.BuildDir(), build)
	if err != nil {
		return fmt.Errorf("%s in %s", err.Error(), plan.BuildDir())
	}
	return nil
}

func doCommands(dir string, cmds []string) (err error) {
	fmt.Println(dir)
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

func Package(dir string, plan *Plan) (err error) {
	var (
		pack = plan.Package
	)

	if err := config.CheckBranches(); err != nil {
		return (err)
	}
	pdir := join(cache.Packages(), plan.NameVersion())
	if file.Exists(pdir) {
		err := os.RemoveAll(pdir)
		if err != nil {
			return err
		}
	}

	if os.Mkdir(pdir, 0755) != nil {
		elog.Println(err)
		return err
	}
	os.Setenv("PKGDIR", Path(pdir).ToUnix())
	if plan.Inherit != "" {
		parent, _ := NewPlan(plan.Inherit)
		pack = append(parent.Package, plan.Package...)
	}

	if err := doCommands(dir, pack); err != nil {
		elog.Println(err)
		return err
	}
	for _, j := range plan.SubPackages {
		sub, err := NewPlan(j)
		sub.Version = plan.Version
		if err != nil {
			return err
		}
		if err = Package(dir, sub); err != nil {
			return err
		}
	}

	if err := CreatePackage(plan); err != nil {
		return (err)
	}
	cid, err := IpfsAdd(Path(plan.PackagePath()), false)
	if err != nil {
		return err
	}
	plan.Cid = cid
	return plan.Save()
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
func SyncHashs() {
	panic("not implimented")
	// plans, _ := GetPlans()
	// for _, p := range plans {
	//	if file.Exists(p.PackagePath()) {
	//		p.Oid, _ = file.Sha256sum(p.PackagePath())
	//		p.Save()
	//		log.Println(p.Oid, p.Name)
	//	}
	// }
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
	for _, d := range append(plan.AutoDepends, plan.ManualDepends...) {
		if IsInstalled(d) {
			continue
		}
		err := Install(d)
		if err != nil {
			return err
		}
	}
	db := filepath.Join(config.DB.Installed(), plan.Name)
	if file.Exists(db) {
		return fmt.Errorf("%s is already installed", name)
	}
	pfile := plan.PackagePath()
	if !file.Exists(pfile) {
		//return errors.New(fmt.Sprintf("%s does not exist", pfile))
		ddir := join(config.Repo, "repo")
		os.MkdirAll(ddir, 0755)
		err := gurl.Download(ddir, config.Binary+"/"+plan.PackageFile())
		if err != nil {
			elog.Println(pfile)
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
	cid, err := IpfsAdd(Path(plan.PackagePath()), true)
	if err != nil {
		return (err)
	}
	if cid != plan.Cid {
		return fmt.Errorf("%s Plans CID does not match tarballs got %s", plan.NameVersion(), cid)
	}
	man, err := ReadPackManifest(pfile)
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

func BuildDeps(plan *Plan) (err error) {
	deps := append(plan.AutoDepends, plan.ManualDepends...)
	for _, d := range deps {
		if IsInstalled(d) {
			continue
		}
		p, _ := NewPlan(d)
		if file.Exists(p.PackagePath()) {
			err := Install(p.Name)
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
	err = BuildSteps(plan)
	if err != nil {
		return err
	}
	return Install(plan.Name)
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
	if err := Stage(plan); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, blue("build"), plan.NameVersion())
	if err := Build(plan); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, blue("package"), plan.NameVersion())
	if err := Package(plan.BuildDir(), plan); err != nil {
		elog.Println(err)
		return err
	}
	return nil
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
