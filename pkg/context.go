package via

import (
	"fmt"
)

// ViaContext helps to tie config, cache and plan fields together
type ViaContext struct {
	Plan   *Plan
	Config Config
}

func NewViaContext(config Config, plan *Plan) *ViaContext {
	return &ViaContext{Config: config, Plan: plan}
}

func (c ViaContext) PackageFile() string {
	if c.Plan.Cid == "" {
		return fmt.Sprintf("%s-%s-%s.tar.gz", c.Plan.NameVersion(), c.Config.OS, c.Config.Arch)
	}
	return fmt.Sprintf("%s.tar.gz", c.Plan.Cid)
}

func (c ViaContext) PackagePath() string {
	return join(c.Config.Repo, c.PackageFile())
}
