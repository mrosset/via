package via

import (
	"github.com/mrosset/util/json"
	"path/filepath"
)

// PlanContext type ties config, cache and plan fields together
type PlanContext struct {
	Plan   *Plan
	Config Config
	Cache  Cache
	Update bool
	Debug  bool
	Verbos bool
}

// NewPlanContext creates an initializes a new PlanContext
func NewPlanContext(config *Config, plan *Plan) *PlanContext {
	return &PlanContext{
		Config: *config,
		Plan:   plan,
		Cache:  config.Cache,
	}
}

// NewPlanContextByName creates and initializes a new PlanContext this
// is like NewPlanContext but instead finds a new Plan by name
func NewPlanContextByName(config *Config, name string) (*PlanContext, error) {
	plan, err := NewPlan(config, name)
	if err != nil {
		return nil, err
	}
	return NewPlanContext(config, plan), nil
}

// PlanPath returns the full path of this contexts Plan's json file
func (c PlanContext) PlanPath() string {
	return filepath.Join(c.Config.Plans, c.Plan.Group, c.Plan.Name+".json")
}

// BuildDir returns the full path of this context Plan's build directory
func (c PlanContext) BuildDir() string {
	bdir := join(c.Cache.Builds(), c.Plan.NameVersion())
	if c.Plan.BuildInStage {
		bdir = join(c.Cache.Stages(), c.Plan.stageDir())
	}
	return bdir
}

// WritePlan saves the serialized go struct to it's json file. The
// json file is pretty formatted so to keep consistency
func (c PlanContext) WritePlan() error {
	return json.Write(c.Plan, c.PlanPath())
}

// StageDir returns the full path for the PlanContext staging
// directory
func (c PlanContext) StageDir() string {
	return join(c.Cache.Stages(), c.Plan.stageDir())
}

// SourcePath returns the full path for the PlanContext source file or
// directory
func (c PlanContext) SourcePath() string {
	s := filepath.Join(c.Cache.Sources(), filepath.Base(c.Plan.Expand().Url))
	return s
}
