package via

import (
	"fmt"
)

// Context helps to tie config, cache and plan fields together
type Context struct {
	Plan   *Plan
	Config Config
}

func NewContext(config Config, plan *Plan) *Context {
	return &Context{Config: config, Plan: plan}
}

func (c Context) PackageFile() string {
	if c.Plan.Cid == "" {
		return fmt.Sprintf("%s-%s-%s.tar.gz", c.Plan.NameVersion(), c.Config.OS, c.Config.Arch)
	}
	return fmt.Sprintf("%s.tar.gz", c.Plan.Cid)
}

func (c Context) PackagePath() string {
	return join(c.Config.Repo, c.PackageFile())
}
