package via

import (
	"fmt"
	"io"
	"os"
)

const (
	manifestName = "manifest.json.gz"
)

const (
	EntryDir = iota
	EntryFile
	EntryLink
)

type FileEntry struct {
	Path      string
	EntryType int
}

type Manifest struct {
	Meta  *Plan
	Files []*FileEntry
}

func NewManifest(plan *Plan) *Manifest {
	return &Manifest{plan, []*FileEntry{}}
}

func (t *Manifest) AddEntry(file string, eType int) {
	t.Files = append(t.Files, &FileEntry{file, eType})
}

func ReadManifest(r io.Reader) (man *Manifest, err os.Error) {
	man = new(Manifest)
	err = ReadGzIo(man, r)
	return man, err
}

func UnPackManifest(file string) (man *Manifest, err os.Error) {
	tbr, err := NewTarBallReader(file)
	if err != nil {
		return nil, err
	}
	defer tbr.Close()
	for {
		hdr, err := tbr.tr.Next()
		if err == os.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Name == manifestName {
			man, err = ReadManifest(tbr.tr)
			return man, err
		}
	}
	return nil, fmt.Errorf("Could not find %s in %s", manifestName, file)
}
