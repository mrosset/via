package via

import (
	"path"
)

var (
	home     = "/home/strings/via"
	plans    = path.Join(home, "plans")
	repo     = path.Join(home, "repo")
	cache    = path.Join(home, "cache")
	packages = path.Join(cache, "packages")
)

func GetRepo() string {
	return repo
}
