package via

import (
	"archive/tar"
	"fmt"
	"io"
	"path"
)

type Manifest struct {
	Plan  *Plan
	Files []*tar.Header
}

func TarManifest(plan *Plan) (man *Manifest, err error) {
	man = &Manifest{Plan: plan}
	pfile := path.Join(config.Repo, plan.PackageFile())
	tr, err := NewTarGzReader(pfile)
	if err != nil {
		return nil, err
	}
	defer tr.Close()
	for {
		hdr, err := tr.Tr.Next()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if hdr == nil {
			break
		}
		fmt.Println(hdr.Name)
	}
	return man, nil
}

func ReadManifest(file string) (man *Manifest, err error) {
	return man, nil
}

func SaveManifest(plan *Plan) (err error) {
	return nil
}
