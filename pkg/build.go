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

// BuildState type contains everything required to build a Plan
type BuildState struct {
	BuildDir    Path
	StageDir    Path
	PackageDir  Path
	SourcePath  Path
	PackagePath Path
	GlobalFlags Flags
	PlanFlags   Flags
	Patch       []string
	Build       []string
	Package     []string
}

// NewBuildContext returns a newly initialized BuildState
func NewBuildContext(config *Config, plan *Plan) BuildState {
	return BuildState{
		BuildDir:    buildDir(config, plan),
		StageDir:    stageDir(config, plan),
		PackageDir:  packageDir(config, plan),
		SourcePath:  sourcePath(config, plan),
		PackagePath: packagePath(config, plan),
		GlobalFlags: config.Flags,
		PlanFlags:   plan.Flags,
	}
}

// buildDir returns the path for the plans' build directory
func buildDir(config *Config, plan *Plan) Path {
	bdir := config.Cache.Builds().Join(plan.NameVersion())
	if plan.BuildInStage {
		bdir = config.Cache.Stages().Join(plan.stageDir())
	}
	return bdir
}

// stageDir returns the plans's stage directory
func stageDir(config *Config, plan *Plan) Path {
	return config.Cache.Stages().Join(plan.stageDir())
}

// packageDir return the full path for plan's package directory
func packageDir(config *Config, plan *Plan) Path {
	return config.Cache.Packages().Join(plan.NameVersion())
}

// sourcePath returns the full path of the Plans source file
//
// FIXME: don't use base of source URL for source filename
func sourcePath(config *Config, plan *Plan) Path {
	file := filepath.Base(plan.Expand().Url)
	return config.Cache.Sources().Join(file)
}

// packagePath returns the full path of the plans package file
func packagePath(config *Config, plan *Plan) Path {
	return config.Repo.Join(PackageFile(config, plan))
}

// Builder provides type for building a Plan
type Builder struct {
	Config  *Config
	Plan    *Plan
	Cache   Cache
	Context BuildState
}

// NewBuilder returns new Builder that has been initialized
func NewBuilder(config *Config, plan *Plan) Builder {
	return Builder{
		Config:  config,
		Plan:    plan,
		Cache:   config.Cache,
		Context: NewBuildContext(config, plan),
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
	if err := b.Package(b.Context.BuildDir); err != nil {
		return err
	}
	return RepoCreate(b.Config)
}

// Download Plans sources to Cache
func (b Builder) Download() error {
	if b.Context.SourcePath.Exists() {
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
	if b.Plan.Url == "" || b.Context.StageDir.Exists() {
		return nil
	}
	url, err := url.Parse(b.SourceURL())
	if err != nil {
		return err
	}
	if url.Scheme == "git" {
		if err := Clone(b.Context.StageDir, b.Context.SourcePath.String()); err != nil {
			return err
		}
	} else {
		switch b.Context.SourcePath.Ext() {
		case "zip":
			unzip(b.Cache.Stages(), b.Context.SourcePath.String())
		default:
			if err := GNUUntar(b.Cache.Stages(), b.Context.SourcePath.String()); err != nil {
				return err
			}
		}
	}
	return b.doCommands(b.Context.StageDir, b.Plan.Patch)
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
	b.Context.BuildDir.Ensure()
	// Parent plan Build is run first this plans is added at the end.
	if b.Plan.Inherit != "" {
		parent, err := NewPlan(b.Config, b.Plan.Inherit)
		if err != nil {
			return err
		}
		build = append(parent.Build, b.Plan.Build...)
		flags = append(flags, parent.Flags...)
	}
	if err := b.doCommands(b.Context.BuildDir, build); err != nil {
		return fmt.Errorf("%s in %s", err, b.Context.BuildDir)
	}
	return nil
}

// ExpandCommand expands input's variables using builders Expands
func ExpandCommand(input string, builder Builder) string {
	return os.Expand(input, builder.Expand)
}

// Expand shell like variables returning an expanded string. This
// conforms to os.Expand function mapper argument.
//
// Supported variables are.
// PREFIX
// SRCDIR
// PKGDIR
// Flags
// PlanFlags
// LDFLAGS
//
// FIXME: instead of using os.Expand using encoding/template. os.Expand
// is used so current plans do not break
func (b Builder) Expand(in string) string {
	switch in {
	case "PREFIX":
		return b.Config.Prefix.String()
	case "SRCDIR":
		return b.Context.StageDir.String()
	case "PKGDIR":
		return b.Context.PackageDir.String()
	case "Flags":
		return b.Config.Flags.Join()
	case "PlanFlags":
		return b.Plan.Flags.Join()
	case "LDFLAGS":
		return b.Config.Env["LDFLAGS"]
	}

	return ""
}

// Package the Plan
func (b Builder) Package(dir Path) error {
	// Remove plans Cid it's assumed we'll be creating a new one
	b.Plan.Cid = ""
	var (
		plan = b.Plan
		pack = plan.Package
		// TODO: use temp file here
		pfile = b.Context.PackagePath
	)

	b.Context.PackageDir.RemoveAll()
	b.Context.PackageDir.Ensure()

	if plan.Inherit != "" {
		parent, _ := NewPlan(b.Config, plan.Inherit)
		pack = append(parent.Package, plan.Package...)
	}
	// Run package commands
	if err := b.doCommands(dir, pack); err != nil {
		return err
	}
	// Package sub packages
	for _, sub := range plan.SubPackages {
		p, err := NewPlan(b.Config, sub)
		if err != nil {
			return err
		}
		if err := NewBuilder(b.Config, p).Package(dir); err != nil {
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
	b.Config.Repo.Ensure()

	fd, err := os.Create(b.Context.PackagePath.String())
	if err != nil {
		elog.Println(err)
		return err
	}
	defer fd.Close()
	gz := gzip.NewWriter(fd)
	defer gz.Close()
	return b.Tarball(gz)
}

// Tarball creates manifest and walks packageDir taring and
// compressing package files
func (b Builder) Tarball(wr io.Writer) (err error) {
	if err := CreateManifest(b.Config, b.Plan, b.Context.PackageDir.String()); err != nil {
		return err
	}
	return archive(wr, b.Context.PackageDir.String())
}

// SourceURL returns the Plans expanded Url
func (b Builder) SourceURL() string {
	return b.Plan.Expand().Url
}

// doCommands runs in cmd in dir Path
func (b Builder) doCommands(dir Path, cmds []string) (err error) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		return err
	}
	for _, j := range cmds {
		args := ExpandCommand(j, b)
		cmd := &exec.Cmd{
			Path:   bash,
			Args:   []string{"bash", "-c", args},
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
			Dir:    dir.String(),
			Env:    b.Config.SanitizeEnv(),
		}
		if verbose {
			cmd.Stdout = os.Stdout
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %s", j, err)
		}
	}
	return nil
}
