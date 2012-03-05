package via

import (
	"archive/tar"
)

type Manifest struct {
	Plan  *Plan
	Files []*tar.Header
}

func TarManifest(plan *Plan) (man *Manifest, err error) {
	man = &Manifest{Plan: plan}
	return man, nil
}

func ReadManifest(file string) (man *Manifest, err error) {
	return man, nil
}

func SaveManifest(plan *Plan) (err error) {
	return nil
}
