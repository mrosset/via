package via

import (
	"fmt"
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
		return filepath.Join(c.Cache.Stages(), c.Plan.NameVersion())
	}
	return filepath.Join(c.Cache.Stages(), c.Plan.stageDir())
}

func (c *Construct) StagesPath() string {
	return filepath.Join(c.Cache.Stages(), c.Plan.stageDir())
}

func (c *Construct) PackageFilePath() string {
	return fmt.Sprintf("%s-%s-%s.tar.gz", c.Plan.NameVersion(), c.Config.OS, c.Config.Arch)
}

func (c *Construct) PackagePath() string {
	return filepath.Join(c.Config.Repo, "repo", c.PackageFilePath())
}

func (c *Construct) SourcePath() string {
	return filepath.Join(c.Cache.Sources(), filepath.Base(c.Plan.Expand().Url))
}

func NewConstruct(config *Config, plan *Plan) *Construct {
	return &Construct{Config: config, Plan: plan}
}
