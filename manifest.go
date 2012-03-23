package via

import (
	"util/json"
)

type Manifest struct {
	Plan  *Plan
	Files []string
	Dirs  []string
}

func ReadManifest(name string) (man *Manifest, err error) {
	man = new(Manifest)
	err = json.Read(man, join(inst, name, "manifest.json"))
	if err != nil {
		return
	}
	return
}
