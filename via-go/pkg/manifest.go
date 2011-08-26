package via

import (
	"bytes"
	"compress/gzip"
	"io"
	"json"
	"os"
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

func (t Manifest) Save(path string) (err os.Error) {
	fd, err := os.Create(path)
	if err != nil {
		return
	}
	defer fd.Close()
	gw, err := gzip.NewWriter(fd)
	if err != nil {
		return
	}
	defer gw.Close()
	b, err := json.Marshal(t)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, b, "", "\t")
	if err != nil {
		return
	}
	io.Copy(gw, buf)
	return
}
