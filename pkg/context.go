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
	Update bool
	Debug  bool
	Verbos bool
}

// Returns a new VieaContext
func NewPlanContext(config *Config, plan *Plan) *PlanContext {
	return &PlanContext{Config: *config, Plan: plan}
}

func NewPlanContextByName(config *Config, name string) (*PlanContext, error) {
	plan, err := NewPlan(config, name)
	if err != nil {
		return nil, err
	}
	return &PlanContext{Config: *config, Plan: plan}, nil
}

func (c PlanContext) PlanPath() string {
	return filepath.Join(c.Config.Plans, c.Plan.Group, c.Plan.Name+".json")
}

func (c PlanContext) WritePlan() error {
	return json.Write(c.Plan, c.PlanPath())
}

func (c PlanContext) SourcePath() string {
	return c.Plan.SourcePath()
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
