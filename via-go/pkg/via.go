package via

import (
	"fmt"
	"path"
)

var (
	home     = "/home/strings/via"
	plans    = path.Join(home, "plans")
	repo     = path.Join(home, "repo")
	cache    = path.Join(home, "cache")
	packages = path.Join(cache, "packages")
)

const (
	PackExt = "tar.gz"
)

func PkgFile(plan *Plan, arch string) string {
	return fmt.Sprintf("%s-%s-%s.%s", plan.Name, plan.Version, arch, PackExt)
}

func PkgAbsFile(plan *Plan, arch string) string {
	return fmt.Sprintf("%s/%s/%s", repo, arch, PkgFile(plan, arch))
}
