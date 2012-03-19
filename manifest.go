package via

import (
	"util/json"
)

type Manifest struct {
	Plan  *Plan
	Files []string
	Dirs  []string
}

func ReadManifest(plan *Plan) (man *Manifest, err error) {
	man = new(Manifest)
	err = json.Read(man, installed.Dir(plan.Name).File("manifest.json"))
	if err != nil {
		return
	}
	return
}
