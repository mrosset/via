package via

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
)

var (
	client  = new(http.Client)
	debug   = false
	deps    = false
	elog    = log.New(os.Stderr, "", log.Lshortfile)
	expand  = os.ExpandEnv
	lfmt    = "%-20.20s %v\n"
	update  = false
	verbose = false
)

// Verbose sets the global verbosity level
//
// FIXME: this should be set via Builder or Installer
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
		fpath := config.Root.Join(f)
		if err := os.Remove(fpath.String()); err != nil {
			elog.Println(err)
		}
	}
	return config.DB.Installed(config).Join(name).RemoveAll()
}

// BuildDeps build's all of a packages dependency's are build if they
// are not built already
func BuildDeps(config *Config, plan *Plan) (err error) {
	for _, d := range plan.Depends() {
		if IsInstalled(config, d) {
			continue
		}
		p, _ := NewPlan(config, d)
		if p.IsRebuilt {
			b := NewBatch(config, os.Stdout)
			b.Walk(p)
			if err := b.Install(); err != nil {
				return err[0]
			}
			continue
		}
		fmt.Println("building", d, "for", plan.NameVersion())
		err := BuildDeps(config, p)
		if err != nil {
			elog.Println(err)
			return err
		}
	}
	bd := NewBuilder(config, plan)
	if err := bd.BuildSteps(); err != nil {
		return err
	}
	b := NewBatch(config, os.Stdout)
	b.Walk(plan)
	if err := b.Install(); err != nil {
		return err[0]
	}
	return nil
}

var (
	rexName   = regexp.MustCompile("[A-Za-z]+")
	rexTruple = regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+")
	rexDouble = regexp.MustCompile("[0-9]+.[0-9]+")
)

// Create a new plan from a given Url
func Create(config *Config, url, group string) (err error) {
	var (
		xfile   = path.Base(url)
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

	pfile := PlanFilePath(config, plan)
	if pfile.Exists() {
		return fmt.Errorf("%s already exists", pfile)
	}
	return WritePlan(config, plan)
}

// IsInstalled returns true if a plan is installed
func IsInstalled(config *Config, name string) bool {
	return config.DB.Installed(config).Join(name).Exists()
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Clean the Plans build directory
func Clean(config *Config, plan *Plan) error {
	var (
		cache = config.Cache
	)
	fmt.Printf(lfmt, "clean", plan.NameVersion())
	err := cache.Builds().Join(plan.NameVersion()).RemoveAll()
	if err != nil {
		return err
	}
	if plan.BuildInStage {
		return cache.Stages().Join(plan.stageDir()).RemoveAll()
	}
	return nil
}

func conflicts(config *Config, man *Plan) (errs []error) {
	for _, f := range man.Files {
		if config.Root.Join(f).Exists() {
			errs = append(errs, fmt.Errorf("%s already exists", f))
		}
	}
	return errs
}
