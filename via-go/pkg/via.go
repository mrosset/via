package via

import (
	"fmt"
	"log"
	"path"
)

var (
	home     = "/home/strings/via"
	plans    = path.Join(home, "plans")
	repo     = path.Join(home, "repo")
	cache    = path.Join(home, "cache")
	packages = path.Join(cache, "packages")
)

func init() {
	log.SetPrefix("via: ")
	log.SetFlags(0)
}

const (
	PackExt = "tar.gz"
)

func GetRepo() string {
	return repo
}

func PkgFile(plan *Plan, arch string) string {
	return fmt.Sprintf("%s-%s-%s.%s", plan.Name, plan.Version, arch, PackExt)
}

func PkgAbsFile(plan *Plan, arch string) string {
	return fmt.Sprintf("%s/%s/%s", repo, arch, PkgFile(plan, arch))
}
