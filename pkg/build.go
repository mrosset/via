package via

import (
	"compress/gzip"
	"fmt"
	"github.com/mrosset/gurl"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

// Builder provides type for building a Plan
type Builder struct {
	Config *Config
	Plan   *Plan
	Cache  Cache
}

// NewBuilder returns new Builder that has been initialized
func NewBuilder(config *Config, plan *Plan) Builder {
	return Builder{
		Config: config,
		Plan:   plan,
		Cache:  config.Cache,
	}
}

// NewBuilderByName returns a new builder looking Plan by name
func NewBuilderByName(config *Config, name string) (Builder, error) {
	plan, err := NewPlan(config, name)
	if err != nil {
		return Builder{}, err
	}
	return NewBuilder(config, plan), nil
}

// BuildSteps calls all of the methods required to build a Plan
func (b Builder) BuildSteps() error {
	fmt.Printf(lfmt, "download", b.Plan.NameVersion())
	if err := b.Download(); err != nil {
		return err
	}
	fmt.Printf(lfmt, "stage", b.Plan.NameVersion())
	if err := b.Stage(); err != nil {
		return err
	}
	fmt.Printf(lfmt, "build", b.Plan.NameVersion())
	if err := b.Build(); err != nil {
		return err
	}
	fmt.Printf(lfmt, "package", b.Plan.NameVersion())
	if err := b.Package(); err != nil {
		return err
	}
	return RepoCreate(b.Config)
}

// Download Plans sources to Cache
func (b Builder) Download() error {
	if b.SourcePath().Exists() {
		return nil
	}
	url, err := url.Parse(b.SourceURL())
	if err != nil {
		return err
	}
	switch url.Scheme {
	case "ftp":
		wget(b.Cache.Sources(), b.SourceURL())
	case "http", "https":
		return gurl.Download(b.Cache.Sources().String(), b.SourceURL())
	case "git":
		spath := b.Cache.Sources().Join(b.Plan.Name)
		giturl := "https" + b.SourceURL()[3:]
		if err := Clone(spath, giturl); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s: URL scheme is not supported", url.Scheme)
	}
	return nil
}

// Stage the Plans source files into it's staging directory
func (b Builder) Stage() error {
	if b.Plan.Url == "" || b.StageDir().Exists() {
		return nil
	}
	url, err := url.Parse(b.SourceURL())
	if err != nil {
		return err
	}
	if url.Scheme == "git" {
		if err := Clone(b.StageDir(), b.SourcePath().String()); err != nil {
			return err
		}
	} else {
		switch b.SourcePath().Ext() {
		case "zip":
			unzip(b.Cache.Stages(), b.SourcePath().String())
		default:
			if err := GNUUntar(b.Cache.Stages(), b.SourcePath().String()); err != nil {
				return err
			}
		}
	}

	return nil
}

// Build runs the Plans Build section
func (b Builder) Build() error {
	var (
		build = b.Plan.Build
	)
	// install build deps
	for _, p := range b.Plan.BuildDepends {
		if IsInstalled(b.Config, p) {
			continue
		}
		dp, err := NewPlan(b.Config, p)
		if err != nil {
			return err
		}
		ins := NewInstaller(b.Config, dp)
		if err := ins.Install(); err != nil {
			return err
		}
	}
	flags := append(b.Config.Flags, b.Plan.Flags...)
	Path(b.BuildDir()).Ensure()
	// Parent plan Build is run first this plans is added at the end.
	if b.Plan.Inherit != "" {
		parent, _ := NewPlan(b.Config, b.Plan.Inherit)
		build = append(parent.Build, b.Plan.Build...)
		flags = append(flags, parent.Flags...)
	}
	// FIXME: this should be set within exec.Cmd
	os.Setenv("SRCDIR", b.StageDir().String())
	os.Setenv("Flags", expand(flags.Join()))
	if err := doCommands(b.Config, b.BuildDir(), build); err != nil {
		return fmt.Errorf("%s in %s", err, b.BuildDir())
	}
	return nil
}

// Package the Plan
func (b Builder) Package() error {
	// Remove plans Cid it's assumed we'll be creating a new one
	b.Plan.Cid = ""
	var (
		plan = b.Plan
		pack = plan.Package
		// TODO: use temp file here
		pfile = PackagePath(b.Config, plan)
	)

	b.PackageDir().RemoveAll()
	b.PackageDir().Ensure()

	os.Setenv("PKGDIR", b.PackageDir().String())
	if plan.Inherit != "" {
		parent, _ := NewPlan(b.Config, plan.Inherit)
		pack = append(parent.Package, plan.Package...)
	}
	// Run package commands
	if err := doCommands(b.Config, b.BuildDir(), pack); err != nil {
		return err
	}
	// Package sub packages
	for _, sub := range plan.SubPackages {
		p, err := NewPlan(b.Config, sub)
		if err != nil {
			return err
		}
		if err := NewBuilder(b.Config, p).Package(); err != nil {
			return err
		}
	}
	// If repo.json and files.json do not exist create them
	if !b.Config.Repo.FilesFile(b.Config).Exists() {
		if err := RepoCreate(b.Config); err != nil {
			return err
		}
	}
	// Create the tarball
	if err := b.CreatePackage(); err != nil {
		return err
	}

	var err error
	if b.Plan.Cid, err = IpfsAdd(b.Config, pfile); err != nil {
		return err
	}
	plan.IsRebuilt = true
	return WritePlan(b.Config, b.Plan)
}

// CreatePackage create Tarball package
func (b Builder) CreatePackage() error {
	var (
		pfile = PackagePath(b.Config, b.Plan)
	)

	b.Config.Repo.Ensure()

	fd, err := os.Create(pfile)
	if err != nil {
		elog.Println(err)
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return b.Tarball(gz)
}

// Tarball creates manifest and walks PackageDir taring and
// compressing package files
func (b Builder) Tarball(wr io.Writer) (err error) {
	if err := CreateManifest(b.Config, b.Plan, b.PackageDir().String()); err != nil {
		return err
	}
	if err != nil {
		elog.Println(err)
		return err
	}
	return archive(wr, b.PackageDir().String())
}

// SourceURL returns the Plans expanded Url
func (b Builder) SourceURL() string {
	return b.Plan.Expand().Url
}

// SourcePath returns the full path of the Plans source file
func (b Builder) SourcePath() Path {
	file := filepath.Base(b.Plan.Expand().Url)
	return b.Cache.Sources().Join(file)
}

// PackageDir return the full path for Builder's package directory
func (b Builder) PackageDir() Path {
	return b.Cache.Packages().Join(b.Plan.NameVersion())
}

// BuildDir returns the path for the Builder's build directory
func (b Builder) BuildDir() Path {
	bdir := b.Cache.Builds().Join(b.Plan.NameVersion())
	if b.Plan.BuildInStage {
		bdir = b.Cache.Stages().Join(b.Plan.stageDir())
	}
	return bdir
}

// StageDir returns the Builders stage directory
func (b Builder) StageDir() Path {
	return b.Cache.Stages().Join(b.Plan.stageDir())
}

// doCommands runs in cmd in dir Path
func doCommands(config *Config, dir Path, cmds []string) (err error) {
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
			Dir:    dir.String(),
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
