package via

import (
	"fmt"
	"github.com/mrosset/util/json"
	"path/filepath"
)

// PlanContext helps to tie config, cache and plan fields together
type PlanContext struct {
	Plan   *Plan
	Config Config
	Cache  Cache
	Update bool
	Debug  bool
	Verbos bool
}

// Returns a new VieaContext
func NewPlanContext(config *Config, plan *Plan) *PlanContext {
	return &PlanContext{
		Config: *config,
		Plan:   plan,
		Cache:  config.Cache,
	}
}

func NewPlanContextByName(config *Config, name string) (*PlanContext, error) {
	plan, err := NewPlan(config, name)
	if err != nil {
		return nil, err
	}
	return &PlanContext{
		Config: *config,
		Cache:  config.Cache,
		Plan:   plan,
	}, nil
}

func (c PlanContext) PlanPath() string {
	return filepath.Join(c.Config.Plans, c.Plan.Group, c.Plan.Name+".json")
}

func (c PlanContext) BuildDir() string {
	bdir := join(c.Cache.Builds(), c.Plan.NameVersion())
	if c.Plan.BuildInStage {
		bdir = join(c.Cache.Stages(), c.Plan.stageDir())
	}
	return bdir
}

func (c PlanContext) WritePlan() error {
	return json.Write(c.Plan, c.PlanPath())
}

func (c PlanContext) StageDir() string {
	return join(c.Cache.Stages(), c.Plan.stageDir())
}

func (c PlanContext) SourcePath() string {
	s := filepath.Join(c.Cache.Sources(), filepath.Base(c.Plan.Expand().Url))
	return s
}

func (c PlanContext) PackageFile() string {
	if c.Plan.Cid == "" {
		return fmt.Sprintf("%s-%s-%s.tar.gz", c.Plan.NameVersion(), c.Config.OS, c.Config.Arch)
	}
	return fmt.Sprintf("%s.tar.gz", c.Plan.Cid)
}

func (c PlanContext) PackagePath() string {
	return join(c.Config.Repo, c.PackageFile())
}
