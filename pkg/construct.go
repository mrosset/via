package via

import (
	"compress/gzip"
	"fmt"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/file"
	"net/url"
	"os"
	"path/filepath"
)

// Construct contains everthing needed to build and install a plan. While it's the
// motar which pieces together a Plan and a Config
type Construct struct {
	Plan   *Plan
	Config *Config
	Cache  Cache
}

func (c *Construct) BuildPath() string {
	if c.Plan.BuildInStage {
		return filepath.Join(c.Cache.Stages(), c.Plan.stageDir())
	}
	return filepath.Join(c.Cache.Builds(), c.Plan.NameVersion())
}

func (c *Construct) PlanStagePath() string {
	return filepath.Join(c.Cache.Stages(), c.Plan.stageDir())
}

func (c *Construct) PackageFileName() string {
	return fmt.Sprintf("%s-%s-%s.tar.gz", c.Plan.NameVersion(), c.Config.OS, c.Config.Arch)
}

func (c *Construct) PackageFilePath() string {
	return filepath.Join(c.Config.Repo, "repo", c.PackageFileName())
}

func (c *Construct) PackageDirPath() string {
	return filepath.Join(c.Cache.Packages(), c.Plan.NameVersion())
}

func (c *Construct) PlanSourcePath() string {
	return filepath.Join(c.Cache.Sources(), filepath.Base(c.Plan.Expand().Url))
}

func NewConstruct(config *Config, plan *Plan) *Construct {
	return &Construct{Config: config, Plan: plan, Cache: config.Cache}
}

// Run all of the functions required to build a package
func (c *Construct) BuildSteps() (err error) {
	if file.Exists(c.PackageFilePath()) {
		elog.Printf("package %s exists", c.PackageFileName())
	}
	if err := c.DownloadSource(); err != nil {
		elog.Println(err)
		return err
	}
	if err := c.Stage(); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "build", c.Plan.NameVersion())
	if err := c.Build(); err != nil {
		elog.Println(err)
		return err
	}
	fmt.Printf(lfmt, "package", c.Plan.NameVersion())
	if err := c.Package(); err != nil {
		elog.Println(err)
		return err
	}
	return nil
}

// Stages the downloaded source in via's cache directory
// the stage only happens once unless BuilInStage is used
func (c *Construct) Stage() (err error) {
	if c.Plan.Url == "" || file.Exists(c.PlanStagePath()) {
		// nothing to stage
		return nil
	}
	fmt.Printf(lfmt, "stage", c.Plan.NameVersion())
	switch filepath.Ext(c.Plan.SourceFileName()) {
	case ".zip":
		unzip(c.Cache.Stages(), c.PlanSourcePath())
	default:
		GNUUntar(c.Cache.Stages(), c.PlanSourcePath())
	}
	fmt.Printf(lfmt, "patch", c.Plan.NameVersion())
	if err := doCommands(c.PlanStagePath(), c.Plan.Patch); err != nil {
		return err
	}
	return
}

// Calls each shell command in the plans Build field.
func (c *Construct) Build() (err error) {
	var (
		build = c.Plan.Build
	)
	if err = config.CheckBranches(); err != nil {
		return (err)
	}
	if file.Exists(c.PackageFilePath()) {
		fmt.Printf("FIXME: (short flags)  package %s exists building anyways.\n", c.PackageFilePath())
	}
	flags := append(config.Flags, c.Plan.Flags...)
	os.MkdirAll(c.BuildPath(), 0755)
	// Parent plan Build is run first this plans is added at the end.
	if c.Plan.Inherit != "" {
		parent, _ := FindPlan(c.Config, c.Plan.Inherit)
		build = append(parent.Build, c.Plan.Build...)
		flags = append(flags, parent.Flags...)
	}
	os.Setenv("SRCDIR", c.PlanStagePath())
	os.Setenv("Flags", expand(flags.String()))
	err = doCommands(c.BuildPath(), build)
	if err != nil {
		return fmt.Errorf("%s in %s", err.Error(), c.BuildPath())
	}
	return nil
}

func (c *Construct) Package() (err error) {
	var (
		pack = c.Plan.Package
	)
	if err = c.Config.CheckBranches(); err != nil {
		return (err)
	}
	if file.Exists(c.PackageDirPath()) {
		err := os.RemoveAll(c.PackageDirPath())
		if err != nil {
			return err
		}
	}
	err = os.Mkdir(c.PackageDirPath(), 0755)
	if err != nil {
		elog.Println(err)
		return err
	}
	os.Setenv("PKGDIR", c.PackageDirPath())
	if c.Plan.Inherit != "" {
		parent, _ := FindPlan(c.Config, c.Plan.Inherit)
		pack = append(parent.Package, c.Plan.Package...)
	}
	err = doCommands(c.BuildPath(), pack)
	if err != nil {
		return err
	}
	for _, j := range c.Plan.SubPackages {
		sub, err := FindPlan(c.Config, j)
		if err != nil {
			return err
		}
		if err = NewConstruct(c.Config, sub).Package(); err != nil {
			return err
		}
	}
	err = c.GzipPackageDir()
	if err != nil {
		return (err)
	}
	c.Plan.Oid, err = file.Sha256sum(c.PackageFilePath())
	if err != nil {
		return (err)
	}
	return c.Plan.Save(config)
	/*
		err = CreatePackage(plan)
		if err != nil {
			return err
		}
		return Sign(plan)
	*/
}

func (construct *Construct) GzipPackageDir() (err error) {
	os.MkdirAll(filepath.Dir(construct.PackageFilePath()), 0755)
	fd, err := os.Create(construct.PackageFilePath())
	if err != nil {
		elog.Println(err)
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return Tarball(gz, construct.Plan)
}

func (c *Construct) DownloadSource() (err error) {
	if file.Exists(c.PlanSourcePath()) && !update {
		return nil
	}
	fmt.Printf(lfmt, "download", c.Plan.NameVersion())
	eurl := c.Plan.Expand().Url
	u, err := url.Parse(eurl)
	if err != nil {
		elog.Println(err)
		return err
	}
	switch u.Scheme {
	case "ftp":
		wget(c.Cache.Sources(), eurl)
	case "http", "https":
		return gurl.Download(c.Cache.Sources(), eurl)
	default:
		return fmt.Errorf("%s URL scheme is not supported")
	}
	return nil
}
